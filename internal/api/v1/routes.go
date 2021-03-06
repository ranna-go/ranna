package v1

import (
	"runtime"

	"github.com/gofiber/fiber/v2"
	"github.com/ranna-go/ranna/internal/api/ws"
	"github.com/ranna-go/ranna/internal/config"
	"github.com/ranna-go/ranna/internal/sandbox"
	"github.com/ranna-go/ranna/internal/spec"
	"github.com/ranna-go/ranna/internal/static"
	"github.com/ranna-go/ranna/internal/util"
	"github.com/ranna-go/ranna/pkg/cappedbuffer"
	"github.com/ranna-go/ranna/pkg/models"
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
)

var (
	errOutputLenExceeded = fiber.NewError(fiber.StatusBadRequest, "output len exceeded")
	errEmptyCode         = fiber.NewError(fiber.StatusBadRequest, "code is empty")
)

// @title ranna main API
// @version 1.0
// @description The ranna main REST API.
// @basepath /v1
type Router struct {
	spec            spec.Provider
	cfg             config.Provider
	manager         sandbox.Manager
	streamBufferCap int
}

func (r *Router) Setup(route fiber.Router, ctn di.Container) {
	r.cfg = ctn.Get(static.DiConfigProvider).(config.Provider)
	r.spec = ctn.Get(static.DiSpecProvider).(spec.Provider)
	r.manager = ctn.Get(static.DiSandboxManager).(sandbox.Manager)

	sbc, err := util.ParseMemoryStr(r.cfg.Config().Sandbox.StreamBufferCap)
	if err != nil {
		logrus.WithError(err).Fatal("Invalid value for stream buffer cap")
		return
	}
	r.streamBufferCap = int(sbc)

	route.Use(r.optionsBypass)

	route.Get("/spec", r.getSpec)
	route.Post("/exec", r.postExec)
	route.Get("/info", r.getInfo)
	route.Use("/ws", ws.Upgrade())
	route.Get("/ws", ws.Handler(ctn))
}

func (r *Router) optionsBypass(ctx *fiber.Ctx) error {
	if ctx.Method() == "OPTIONS" {
		return ctx.SendStatus(fiber.StatusOK)
	}
	return ctx.Next()
}

// @summary Get System Info
// @description Returns general system and version information.
// @produce json
// @success 200 {object} models.ExecutionResponse
// @failure 500 {object} models.ErrorModel
// @router /info [get]
func (r *Router) getInfo(ctx *fiber.Ctx) (err error) {
	sandboxInfo, err := r.manager.GetProvider().Info()
	if err != nil {
		return
	}
	info := &models.SystemInfo{
		SandboxInfo: sandboxInfo,
		Version:     static.Version,
		BuildDate:   static.BuildDate,
		GoVersion:   runtime.Version(),
	}
	return ctx.JSON(info)
}

// @summary Get Spec Map
// @description Returns the available spec map.
// @produce json
// @success 200 {object} models.SpecMap
// @router /spec [get]
func (r *Router) getSpec(ctx *fiber.Ctx) (err error) {
	return ctx.JSON(r.spec.Spec().GetSnapshot())
}

// @summary Get Spec Map
// @description Returns the available spec map.
// @accept json
// @produce json
// @param payload body models.ExecutionRequest true "The execution payload"
// @success 200 {object} models.ExecutionResponse
// @failure 400 {object} models.ErrorModel
// @failure 500 {object} models.ErrorModel
// @router /exec [post]
func (r *Router) postExec(ctx *fiber.Ctx) (err error) {
	req := new(models.ExecutionRequest)
	if err = ctx.BodyParser(req); err != nil {
		return
	}

	if req.Code == "" {
		return errEmptyCode
	}

	cStdOut := make(chan []byte)
	cStdErr := make(chan []byte)
	cStop := make(chan bool, 1)

	stdOut := cappedbuffer.New([]byte{}, r.streamBufferCap)
	stdErr := cappedbuffer.New([]byte{}, r.streamBufferCap)

	go func() {
		for {
			select {
			case <-cStop:
				return
			case p := <-cStdOut:
				stdOut.Write(p)
			case p := <-cStdErr:
				stdErr.Write(p)
			}
		}
	}()

	execTime := util.MeasureTime(func() {
		err = r.manager.RunInSandbox(req, nil, cStdOut, cStdErr, cStop)
	})

	if err != nil {
		cStop <- false
		if sandbox.IsSystemError(err) {
			return err
		}
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	res := &models.ExecutionResponse{
		StdOut:     stdOut.String(),
		StdErr:     stdErr.String(),
		ExecTimeMS: int(execTime.Milliseconds()),
	}

	if err = r.checkOutputLen(res.StdOut, res.StdErr); err != nil {
		return
	}

	return ctx.JSON(res)
}

// --- UTIL ---

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
