package sandbox

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/zekrotja/rogu"
	"github.com/zekrotja/rogu/log"

	"github.com/ranna-go/ranna/pkg/models"
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
		cOut chan []byte,
		cErr chan []byte,
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
	logger  rogu.Logger

	runningSandboxes *sync.Map
}

// sandboxWrapper wraps a sandbox instance and
// the used hostDir.
type sandboxWrapper struct {
	Sandbox
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
) (t *ManagerImpl, err error) {
	t = &ManagerImpl{}

	t.sandbox = sandbox
	t.spec = spec
	t.file = file
	t.cfg = cfg
	t.ns = ns
	t.logger = log.Tagged("Manager")

	t.runningSandboxes = &sync.Map{}

	return t, nil
}

func (t *ManagerImpl) PrepareEnvironments(ctx context.Context, force bool) (errs []error) {
	errs = []error{}

	for _, spec := range t.spec.Spec().GetSnapshot() {
		if spec.Image == "" {
			continue
		}
		if err := t.sandbox.Prepare(ctx, *spec, force); err != nil {
			t.logger.Error().Field("image", spec.Image).Err(err).Msg("failed preparing env")
			errs = append(errs, err)
		}
	}

	return errs
}

func (t *ManagerImpl) RunInSandbox(
	ctx context.Context,
	req *models.ExecutionRequest,
	cSpn chan string,
	cOut chan []byte,
	cErr chan []byte,
) (err error) {
	defer func() {
		if err != nil && IsSystemError(err) {
			t.logger.Error().
				Err(err).
				Field("spec", req.Language).
				Msg("sandbox run failed")
		}
	}()

	// Try to get spec from specified language
	spc, ok := t.spec.Spec().Get(req.Language)
	if !ok {
		return errUnsupportedLanguage
	}

	// Process the specified code if it is an inline expression
	if req.InlineExpression {
		// Check if the spec supports inline expressions
		if !spc.SupportsTemplating() {
			return errNoInlineExpressionsSupport
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
	if runSpc.Subdir, err = t.ns.Get(); err != nil {
		return SystemError{err}
	}

	// Set HostDir, Arguments and Environment Variables
	runSpc.HostDir = t.cfg.Config().HostRootDir
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
	if err = t.file.CreateDirectory(hostDir); err != nil {
		return SystemError{err}
	}

	// Create code snippet file in the host + sub-directory
	fileDir := path.Join(hostDir, spc.FileName)
	if err = t.file.CreateFileWithContent(fileDir, req.Code); err != nil {
		return SystemError{err}
	}

	// Create sandbox using RunSpec
	sbx, err := t.sandbox.CreateSandbox(ctx, runSpc)
	if err != nil {
		return SystemError{err}
	}
	if cSpn != nil {
		cSpn <- sbx.ID()
	}
	t.logger.Info().Fields("id", sbx.ID(), "spec", req.Language).Msg("created sandbox")

	// Store sandbox to track run state later
	wrapper := &sandboxWrapper{sbx, hostDir}
	t.runningSandboxes.Store(sbx.ID(), wrapper)

	timeout := time.Duration(t.cfg.Config().Sandbox.TimeoutSeconds) * time.Second
	runCtx, cancelRunCtx := context.WithTimeoutCause(ctx, timeout, errTimedOut)
	defer cancelRunCtx()

	err = sbx.Run(runCtx, cOut, cErr)
	defer func() {
		// Kill container if it is still running, delete the
		// container after as well as delete the snippet host
		// directory.
		if cErr := t.killAndCleanUp(ctx, wrapper); cErr != nil {
			err = SystemError{error: errors.Join(err, cErr)}
		}
		t.logger.Info().Fields("id", sbx.ID(), "spec", req.Language).Msg("sandbox cleaned up")
	}()
	if err != nil {
		if errors.Is(err, errTimedOut) {
			t.logger.Debug().Fields("id", sbx.ID(), "spec", req.Language).Msg("execution timed out")
			return err
		}
		return SystemError{err}
	}

	return err
}

func (t *ManagerImpl) KillAndCleanUp(ctx context.Context, id string) (ok bool, err error) {
	v, ok := t.runningSandboxes.Load(id)
	if !ok {
		return false, nil
	}

	sbx := v.(Sandbox)
	if ok, err = sbx.IsRunning(ctx); !ok || err != nil {
		return ok, err
	}

	if err = sbx.Kill(ctx); err != nil {
		return false, err
	}

	return true, nil
}

func (t *ManagerImpl) Cleanup(ctx context.Context) (errs []error) {
	errs = []error{}

	t.runningSandboxes.Range(func(key, value any) bool {
		w, ok := value.(*sandboxWrapper)
		if ok {
			t.logger.Info().Field("id", w.ID()).Msg("killing and cleaning up running container")
			if err := t.killAndCleanUp(ctx, w); err != nil {
				errs = append(errs, err)
			}
		}
		return true
	})

	return errs
}

func (t *ManagerImpl) GetProvider() Provider {
	return t.sandbox
}

func (t *ManagerImpl) killAndCleanUp(ctx context.Context, w *sandboxWrapper) (err error) {
	defer func() {
		if err != nil {
			t.logger.Error().
				Err(err).
				Field("id", w.ID()).
				Msg("failed cleaning up container")
		}
	}()

	t.logger.Debug().Field("id", w.ID()).Msg("calling killAndCleanUp")

	ok, err := w.IsRunning(ctx)
	if err != nil {
		return err
	}
	if ok {
		err = w.Kill(ctx)
	}
	if err = w.Delete(ctx); err != nil {
		return err
	}
	if err = t.file.DeleteDirectory(w.hostDir); err != nil {
		return err
	}
	t.runningSandboxes.Delete(w.ID())
	return nil
}
