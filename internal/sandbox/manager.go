package sandbox

import (
	"errors"
	"fmt"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/ranna-go/ranna/pkg/models"
	"github.com/ranna-go/ranna/pkg/timeout"
	"github.com/sirupsen/logrus"
)

var (
	errUnsupportedLanguage        = errors.New("unsupported language spec")
	errNoInlineExpressionsSupport = errors.New("this spec has no support for inline expressions")
	errTimedOut                   = errors.New("code execution timed out")
)

// Manager is a higher level abstraction used to create and
// run sandboxes, prepare the environment for given specs and
// cleaning up running containers on teardown.
type Manager interface {

	// RunInSandbox tries to extract the desired spec
	// to be used defined by the req. Then, a new sandbox
	// is created with this spec and given runtime variables.
	// The sandbox is then started and the current go routine is
	// blocked until the execution is finished or timed out.
	//
	// On success returns an execution response.
	RunInSandbox(
		req *models.ExecutionRequest,
		cSpn chan string,
		cOut, cErr chan []byte,
		cClose chan bool,
	) (err error)

	// PrepareEnvironments prepares the sandbox environment for
	// faster first time creation of sandboxes.
	//
	// This pulls required images, for example.
	//
	// If force is true, the environment is being prepared even
	// though it has been already prepared before. This is useful
	// to perform image updates, for example.
	PrepareEnvironments(force bool) []error

	// KillAndCleanUp takes a sandbox ID and, if
	// existing, kills the running sandbox.
	KillAndCleanUp(id string) (bool, error)

	// Cleanup tries to kill and delete all running sandboxes.
	Cleanup() []error

	// GetProvider returns the utilized sandbox provider instance.
	GetProvider() Provider
}

// ManagerImpl is the standard implementation
// of Manager.
type ManagerImpl struct {
	sandbox Provider
	spec    SpecProvider
	file    FileProvider
	cfg     ConfigProvider
	ns      NamespaceProvider

	runningSandboxes *sync.Map
	isCleanup        bool
}

var _ Manager = (*ManagerImpl)(nil)

// sandboxWrapper wraps a sandbox instance and
// the used hostDir.
type sandboxWrapper struct {
	Sandbox
	hostDir string
}

var _ Sandbox = (*sandboxWrapper)(nil)

// SystemError wraps an error occuring with the
// environment or sandbox system.
type SystemError struct {
	error
}

// IsSystemError returns true when the passed
// error is type of SystemError.
func IsSystemError(err error) (ok bool) {
	var systemError SystemError
	return errors.As(err, &systemError)
}

// NewManager returns a new instance of ManagerImpl.
func NewManager(
	sandbox Provider,
	spec SpecProvider,
	file FileProvider,
	cfg ConfigProvider,
	ns NamespaceProvider,
) (m *ManagerImpl, err error) {
	m = &ManagerImpl{}

	m.sandbox = sandbox
	m.spec = spec
	m.file = file
	m.cfg = cfg
	m.ns = ns

	m.runningSandboxes = &sync.Map{}

	return
}

func (m *ManagerImpl) PrepareEnvironments(force bool) (errs []error) {
	errs = []error{}

	for _, spec := range m.spec.Spec().GetSnapshot() {
		if spec.Image == "" {
			continue
		}
		if err := m.sandbox.Prepare(*spec, force); err != nil {
			logrus.WithField("image", spec.Image).WithError(err).Error("failed preparing env")
			errs = append(errs, err)
		}
	}

	return
}

func (m *ManagerImpl) RunInSandbox(
	req *models.ExecutionRequest,
	cSpn chan string,
	cOut, cErr chan []byte,
	cClose chan bool,
) (err error) {
	defer func() {
		if err != nil && IsSystemError(err) {
			logrus.
				WithError(err).
				WithFields(logrus.Fields{
					"spec": req.Language,
				}).
				Error("sandbox run failed")
		}
	}()

	// Try to get spec from specified language
	spc, ok := m.spec.Spec().Get(req.Language)
	if !ok {
		err = errUnsupportedLanguage
		return
	}

	// Process the specified code if it is an inline expression
	if req.InlineExpression {
		// Check if the spec supports inline expressions
		if !spc.SupportsTemplating() {
			err = errNoInlineExpressionsSupport
			return
		}

		code := spc.Inline.Template

		if spc.Inline.ImportRegexCompiled != nil {
			// Extract the imports from the source string
			imports := spc.Inline.ImportRegexCompiled.FindAllString(req.Code, -1)
			req.Code = spc.Inline.ImportRegexCompiled.ReplaceAllString(req.Code, "")
			code = strings.ReplaceAll(code, "$${IMPORTS}", strings.Join(imports, "\n"))
			fmt.Println(imports, spc.Inline.ImportRegex)
		}

		// Wrap the code to execute using the specified template
		code = strings.ReplaceAll(code, "$${CODE}", req.Code)
		req.Code = code
		fmt.Println(code)
	}

	// Wrap in RunSpec
	runSpc := RunSpec{Spec: spc}

	// Get namespace as subdir
	if runSpc.Subdir, err = m.ns.Get(); err != nil {
		err = SystemError{err}
		return
	}

	// Set HostDir, Arguments and Environment Variables
	runSpc.HostDir = m.cfg.Config().HostRootDir
	runSpc.Arguments = req.Arguments
	runSpc.Environment = req.Environment

	// If command is not specified, set file name as
	// command.
	if runSpc.Cmd == "" {
		runSpc.Cmd = spc.FileName
	}

	// Create host directory + sub-directory on the
	// Docker host.
	hostDir := runSpc.GetAssembledHostDir()
	if err = m.file.CreateDirectory(hostDir); err != nil {
		err = SystemError{err}
		return
	}

	// Create code snippet file in the host + sub-directory
	fileDir := path.Join(hostDir, spc.FileName)
	if err = m.file.CreateFileWithContent(fileDir, req.Code); err != nil {
		err = SystemError{err}
		return
	}

	// Create sandbox using RunSpec
	sbx, err := m.sandbox.CreateSandbox(runSpc)
	if err != nil {
		err = SystemError{err}
		return
	}
	if cSpn != nil {
		cSpn <- sbx.ID()
	}
	logrus.WithFields(logrus.Fields{
		"id":   sbx.ID(),
		"spec": req.Language,
	}).Info("created sandbox")

	// Store sandbox to track run state later
	wrapper := &sandboxWrapper{sbx, hostDir}
	m.runningSandboxes.Store(sbx.ID(), wrapper)

	// Run sandbox blocking with timeout
	timedOut := timeout.RunBlockingWithTimeout(func() {
		err = sbx.Run(cOut, cErr, cClose)
	}, time.Duration(m.cfg.Config().Sandbox.TimeoutSeconds)*time.Second)

	if err != nil {
		err = SystemError{err}
		return
	}
	// When manager is in 'cleanup mode' and running containers
	// are killed, skip here and don't try to clean up container
	// again.
	if m.isCleanup {
		return
	}
	// Kill container if it is still running, delete the
	// container after as well as delete the snippet host
	// directory.
	if err = m.killAndCleanUp(wrapper); err != nil {
		return
	}
	if timedOut {
		err = errTimedOut
	}
	logrus.WithFields(logrus.Fields{
		"id":   sbx.ID(),
		"spec": req.Language,
	}).Info("sandbox cleaned up")

	return
}

func (m *ManagerImpl) KillAndCleanUp(id string) (ok bool, err error) {
	v, ok := m.runningSandboxes.Load(id)
	if !ok {
		return
	}

	sbx := v.(Sandbox)
	if ok, err = sbx.IsRunning(); !ok || err != nil {
		return
	}

	err = sbx.Kill()
	ok = true
	return
}

func (m *ManagerImpl) Cleanup() (errs []error) {
	m.isCleanup = true
	errs = []error{}

	m.runningSandboxes.Range(func(key, value interface{}) bool {
		w, ok := value.(*sandboxWrapper)
		if ok {
			logrus.WithField("id", w.ID()).Info("killing and cleaning up running container")
			if err := m.killAndCleanUp(w); err != nil {
				errs = append(errs, err)
			}
		}
		return true
	})

	return
}

func (m *ManagerImpl) GetProvider() Provider {
	return m.sandbox
}

func (m *ManagerImpl) killAndCleanUp(w *sandboxWrapper) (err error) {
	defer func() {
		if err != nil {
			logrus.
				WithError(err).
				WithField("id", w.ID()).
				Error("failed cleaning up container")
		}
	}()

	logrus.WithField("id", w.ID()).Debug("calling killAndCleanUp")

	ok, err := w.IsRunning()
	if err != nil {
		return
	}
	if ok {
		err = w.Kill()
	}
	if err = w.Delete(); err != nil {
		return
	}
	if err = m.file.DeleteDirectory(w.hostDir); err != nil {
		return
	}
	m.runningSandboxes.Delete(w.ID())
	return
}
