package v1

import (
	"path"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/ranna/internal/config"
	"github.com/zekroTJA/ranna/internal/file"
	"github.com/zekroTJA/ranna/internal/namespace"
	"github.com/zekroTJA/ranna/internal/sandbox"
	"github.com/zekroTJA/ranna/internal/spec"
	"github.com/zekroTJA/ranna/internal/static"
	"github.com/zekroTJA/ranna/internal/util"
	"github.com/zekroTJA/ranna/pkg/models"
	"github.com/zekroTJA/ranna/pkg/timeout"
)

var (
	errUnsupportredLanguage = fiber.NewError(fiber.StatusBadRequest, "unsupported language")
	errTimedOut             = fiber.NewError(fiber.StatusRequestTimeout, "code execution timed out")
	errOutputLenExceeded    = fiber.NewError(fiber.StatusBadRequest, "output len exceeded")
)

type Router struct {
	sandbox sandbox.Provider
	spec    spec.Provider
	file    file.Provider
	cfg     config.Provider
	ns      namespace.Provider
}

func (r *Router) Setup(route fiber.Router, ctn di.Container) {
	r.sandbox = ctn.Get(static.DiSandboxProvider).(sandbox.Provider)
	r.spec = ctn.Get(static.DiSpecProvider).(spec.Provider)
	r.file = ctn.Get(static.DiFileProvider).(file.Provider)
	r.cfg = ctn.Get(static.DiConfigProvider).(config.Provider)
	r.ns = ctn.Get(static.DiNamespaceProvider).(namespace.Provider)

	route.Get("/spec", r.getSpec)
	route.Post("/exec", r.postExec)
}

func (r *Router) getSpec(ctx *fiber.Ctx) (err error) {
	return ctx.JSON(r.spec.Spec())
}

func (r *Router) postExec(ctx *fiber.Ctx) (err error) {
	req := new(models.ExecutionRequest)
	if err = ctx.BodyParser(req); err != nil {
		return
	}

	spc, ok := r.spec.Spec().Get(req.Language)
	if !ok {
		return errUnsupportredLanguage
	}

	runSpc := sandbox.RunSpec{Spec: spc}

	if runSpc.Subdir, err = r.ns.Get(); err != nil {
		return
	}

	runSpc.HostDir = r.cfg.Config().HostRootDir
	runSpc.Arguments = req.Arguments
	runSpc.Environment = req.Environment

	if runSpc.Cmd == "" {
		runSpc.Cmd = spc.FileName
	}

	hostDir := runSpc.GetAssambledHostDir()
	if err = r.file.CreateDirectory(hostDir); err != nil {
		return
	}

	fileDir := path.Join(hostDir, spc.FileName)
	if err = r.file.CreateFileWithContent(fileDir, req.Code); err != nil {
		return
	}

	sbx, err := r.sandbox.CreateSandbox(runSpc)
	if err != nil {
		return
	}

	res := new(models.ExecutionResponse)
	timedOut := timeout.RunBlockingWithTimeout(func() {
		res.StdOut, res.StdErr, err = sbx.Run()
	}, time.Duration(r.cfg.Config().Sandbox.TimeoutSeconds)*time.Second)
	if timedOut {
		if err = sbx.Kill(); err != nil {
			return
		}
		return errTimedOut
	}
	if err != nil {
		return
	}

	if err = sbx.Delete(); err != nil {
		return
	}
	if err = r.file.DeleteDirectory(hostDir); err != nil {
		return
	}

	if err = r.checkOutputLen(res.StdOut, res.StdErr); err != nil {
		return
	}

	return ctx.JSON(res)
}

func (r *Router) checkOutputLen(stdout, stderr string) (err error) {
	max, err := util.ParseMemoryStr(r.cfg.Config().API.MaxOutputLen)
	if err != nil {
		return
	}
	if int64(len(stdout))+int64(len(stderr)) > max {
		err = errOutputLenExceeded
	}
	return
}
