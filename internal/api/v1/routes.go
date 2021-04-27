package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ranna-go/ranna/internal/config"
	"github.com/ranna-go/ranna/internal/sandbox"
	"github.com/ranna-go/ranna/internal/spec"
	"github.com/ranna-go/ranna/internal/static"
	"github.com/ranna-go/ranna/internal/util"
	"github.com/ranna-go/ranna/pkg/models"
	"github.com/sarulabs/di/v2"
)

var (
	errOutputLenExceeded = fiber.NewError(fiber.StatusBadRequest, "output len exceeded")
)

type Router struct {
	spec    spec.Provider
	cfg     config.Provider
	manager sandbox.Manager
}

func (r *Router) Setup(route fiber.Router, ctn di.Container) {
	r.cfg = ctn.Get(static.DiConfigProvider).(config.Provider)
	r.spec = ctn.Get(static.DiSpecProvider).(spec.Provider)
	r.manager = ctn.Get(static.DiSandboxManager).(sandbox.Manager)

	route.Use(r.optionsBypass)

	route.Get("/spec", r.getSpec)
	route.Post("/exec", r.postExec)
}

func (r *Router) optionsBypass(ctx *fiber.Ctx) error {
	if ctx.Method() == "OPTIONS" {
		return ctx.SendStatus(fiber.StatusOK)
	}
	return ctx.Next()
}

func (r *Router) getSpec(ctx *fiber.Ctx) (err error) {
	return ctx.JSON(r.spec.Spec())
}

func (r *Router) postExec(ctx *fiber.Ctx) (err error) {
	req := new(models.ExecutionRequest)
	if err = ctx.BodyParser(req); err != nil {
		return
	}

	res, err := r.manager.RunInSandbox(req)
	if err != nil {
		if sandbox.IsSystemError(err) {
			return err
		}
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
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
