package sandbox

import (
	"errors"
	"path"
	"sync"
	"time"

	"github.com/ranna-go/ranna/internal/config"
	"github.com/ranna-go/ranna/internal/file"
	"github.com/ranna-go/ranna/internal/namespace"
	"github.com/ranna-go/ranna/internal/spec"
	"github.com/ranna-go/ranna/internal/static"
	"github.com/ranna-go/ranna/internal/util"
	"github.com/ranna-go/ranna/pkg/models"
	"github.com/ranna-go/ranna/pkg/timeout"
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
)

var (
	errUnsupportredLanguage = errors.New("unsupported language spec")
	errTimedOut             = errors.New("code execution timed out")
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
	// On success, an execution response is returned.
	RunInSandbox(req *models.ExecutionRequest) (res *models.ExecutionResponse, err error)

	// PrepareEnvironment prepares the sandbox environment for
	// faster first time creation of sandboxes.
	//
	// This pulls required images, for example.
	//
	// If force is true, the environment is being prepared even
	// though it has been already prepared before. This is useful
	// to perform image updates, for example.
	PrepareEnvironments(force bool) []error

	// Cleanup tries to kill and delete all running sandboxes.
	Cleanup() []error

	// GetProvider returns the utilized sandbox provider instance.
	GetProvider() Provider
}

// managerImpl is the standard implementation
// of Manager.
type managerImpl struct {
	sandbox Provider
	spec    spec.Provider
	file    file.Provider
	cfg     config.Provider
	ns      namespace.Provider

	streamBufferCap  int
	runningSandboxes *sync.Map
	isCleanup        bool
}

// sandboxWrapper wraps a sandbox instance and
// the used hostDir.
type sandboxWrapper struct {
	sbx     Sandbox
	hostDir string
}

// SystemError wraps an error occuring with the
// environment or sandbox system.
type SystemError struct {
	error
}

// IsSystemError returns true when the passed
// error is type of SystemError.
func IsSystemError(err error) (ok bool) {
	_, ok = err.(SystemError)
	return
}

// NewManager returns a new instance of managerImpl.
func NewManager(ctn di.Container) (m *managerImpl, err error) {
	m = &managerImpl{}

	m.sandbox = ctn.Get(static.DiSandboxProvider).(Provider)
	m.spec = ctn.Get(static.DiSpecProvider).(spec.Provider)
	m.file = ctn.Get(static.DiFileProvider).(file.Provider)
	m.cfg = ctn.Get(static.DiConfigProvider).(config.Provider)
	m.ns = ctn.Get(static.DiNamespaceProvider).(namespace.Provider)

	m.runningSandboxes = &sync.Map{}
	sbc, err := util.ParseMemoryStr(m.cfg.Config().Sandbox.StreamBufferCap)
	if err != nil {
		return
	}
	m.streamBufferCap = int(sbc)

	return
}

func (m *managerImpl) PrepareEnvironments(force bool) (errs []error) {
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

func (m *managerImpl) RunInSandbox(req *models.ExecutionRequest) (res *models.ExecutionResponse, err error) {
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
		err = errUnsupportredLanguage
		return
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

	// Create host directory + sub directory on the
	// Docker host.
	hostDir := runSpc.GetAssambledHostDir()
	if err = m.file.CreateDirectory(hostDir); err != nil {
		err = SystemError{err}
		return
	}

	// Create code snippet file in the host + sub directory
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
	logrus.WithFields(logrus.Fields{
		"id":   sbx.ID(),
		"spec": req.Language,
	}).Info("created sandbox")

	// Store sandbox to track run state later
	wrapper := &sandboxWrapper{sbx, hostDir}
	m.runningSandboxes.Store(sbx.ID(), wrapper)

	// Run sandbox blocking with timeout
	res = new(models.ExecutionResponse)
	timedOut := timeout.RunBlockingWithTimeout(func() {
		res, err = sbx.Run(m.streamBufferCap)
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

func (m *managerImpl) Cleanup() (errs []error) {
	m.isCleanup = true
	errs = []error{}

	m.runningSandboxes.Range(func(key, value interface{}) bool {
		w, ok := value.(*sandboxWrapper)
		if ok {
			logrus.WithField("id", w.sbx.ID()).Info("killing and cleaning up running container")
			if err := m.killAndCleanUp(w); err != nil {
				errs = append(errs, err)
			}
		}
		return true
	})

	return
}

func (m *managerImpl) GetProvider() Provider {
	return m.sandbox
}

func (m *managerImpl) killAndCleanUp(w *sandboxWrapper) (err error) {
	defer func() {
		if err != nil {
			logrus.
				WithError(err).
				WithField("id", w.sbx.ID()).
				Error("failed cleaning up container")
		}
	}()

	logrus.WithField("id", w.sbx.ID()).Debug("calling killAndCleanUp")

	ok, err := w.sbx.IsRunning()
	if err != nil {
		return
	}
	if ok {
		err = w.sbx.Kill()
	}
	if err = w.sbx.Delete(); err != nil {
		return
	}
	if err = m.file.DeleteDirectory(w.hostDir); err != nil {
		return
	}
	m.runningSandboxes.Delete(w.sbx.ID())
	return
}
