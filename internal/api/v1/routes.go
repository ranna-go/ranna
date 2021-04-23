package v1

import (
	"path"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/ranna/internal/config"
	"github.com/zekroTJA/ranna/internal/file"
	"github.com/zekroTJA/ranna/internal/sandbox"
	"github.com/zekroTJA/ranna/internal/spec"
	"github.com/zekroTJA/ranna/internal/static"
	"github.com/zekroTJA/ranna/internal/util"
)

var (
	errUnsupportredLanguage = fiber.NewError(fiber.StatusBadRequest, "unsupported language")
	errTimedOut             = fiber.NewError(fiber.StatusRequestTimeout, "code execution timed out")
)

type Router struct {
	sandbox sandbox.Provider
	spec    spec.Provider
	file    file.Provider
	cfg     config.Provider
}

func (r *Router) Setup(route fiber.Router, ctn di.Container) {
	r.sandbox = ctn.Get(static.DiSandboxProvider).(sandbox.Provider)
	r.spec = ctn.Get(static.DiSpecProvider).(spec.Provider)
	r.file = ctn.Get(static.DiFileProvider).(file.Provider)
	r.cfg = ctn.Get(static.DiConfigProvider).(config.Provider)

	route.Post("/exec", r.postExec)
}

func (r *Router) postExec(ctx *fiber.Ctx) (err error) {
	req := new(executionRequest)
	if err = ctx.BodyParser(req); err != nil {
		return
	}

	spc, ok := r.spec.Spec().Get(req.Language)
	if !ok {
		return errUnsupportredLanguage
	}

	spc.Subdir = "test1"
	spc.HostDir = r.cfg.Config().HostRootDir
	spc.Cmd = spc.FileName

	hostDir := spc.GetAssambledHostDir()
	if err = r.file.CreateDirectory(hostDir); err != nil {
		return
	}

	fileDir := path.Join(hostDir, spc.FileName)
	if err = r.file.CreateFileWithContent(fileDir, req.Code); err != nil {
		return
	}

	sbx, err := r.sandbox.CreateSandbox(spc)
	if err != nil {
		return
	}

	res := new(executionResponse)
	timedOut := util.RunBlockingWithTimeout(func() {
		res.StdOut, res.StdErr, err = sbx.Run()
	}, time.Duration(r.cfg.Config().ExecutionTimeoutSeconds)*time.Second)
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
	if err = r.file.DeleteDirectory(fileDir); err != nil {
		return
	}

	return ctx.JSON(res)
}