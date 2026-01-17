package v1

import (
	"runtime"

	"github.com/zekrotja/rogu/log"

	"github.com/gofiber/fiber/v2"
	"github.com/ranna-go/ranna/internal/api/ws"
	"github.com/ranna-go/ranna/internal/sandbox"
	"github.com/ranna-go/ranna/internal/static"
	"github.com/ranna-go/ranna/internal/util"
	"github.com/ranna-go/ranna/pkg/cappedbuffer"
	"github.com/ranna-go/ranna/pkg/models"
)

var (
	errOutputLenExceeded = fiber.NewError(fiber.StatusBadRequest, "output len exceeded")
	errEmptyCode         = fiber.NewError(fiber.StatusBadRequest, "code is empty")
)

// Router
//
// @title ranna main API
// @version 1.0
// @description The ranna main REST API.
// @basepath /v1
type Router struct {
	spec            SpecProvider
	cfg             ConfigProvider
	manager         SandboxManager
	streamBufferCap int
}

func (t *Router) Setup(route fiber.Router,
	cfg ConfigProvider,
	spec SpecProvider,
	manager SandboxManager,
) {
	t.cfg = cfg
	t.spec = spec
	t.manager = manager

	sbc, err := util.ParseMemoryStr(t.cfg.Config().Sandbox.StreamBufferCap)
	if err != nil {
		log.Fatal().Err(err).Msg("Invalid value for stream buffer cap")
		return
	}
	t.streamBufferCap = int(sbc)

	route.Use(t.optionsBypass)

	route.Get("/spec", t.getSpec)
	route.Post("/exec", t.postExec)
	route.Get("/info", t.getInfo)
	route.Use("/ws", ws.Upgrade())
	route.Get("/ws", ws.Handler(cfg, manager))
}

func (t *Router) optionsBypass(ctx *fiber.Ctx) error {
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
func (t *Router) getInfo(ctx *fiber.Ctx) (err error) {
	sandboxInfo, err := t.manager.GetProvider().Info(ctx.Context())
	if err != nil {
		return err
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
func (t *Router) getSpec(ctx *fiber.Ctx) (err error) {
	return ctx.JSON(t.spec.Spec().GetSnapshot())
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
func (t *Router) postExec(ctx *fiber.Ctx) (err error) {
	req := new(models.ExecutionRequest)
	if err = ctx.BodyParser(req); err != nil {
		return err
	}

	if req.Code == "" {
		return errEmptyCode
	}

	cStdOut := make(chan []byte)
	cStdErr := make(chan []byte)

	stdOut := cappedbuffer.New([]byte{}, t.streamBufferCap)
	stdErr := cappedbuffer.New([]byte{}, t.streamBufferCap)
	cClose := make(chan struct{}, 1)

	go func() {
		for {
			select {
			case <-cClose:
				return
			case p := <-cStdOut:
				stdOut.Write(p)
			case p := <-cStdErr:
				stdErr.Write(p)
			}
		}
	}()
	defer func() {
		cClose <- struct{}{}
	}()

	execTime := util.MeasureTime(func() {
		err = t.manager.RunInSandbox(ctx.Context(), req, nil, cStdOut, cStdErr)
	})

	if err != nil {
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

	if err = t.checkOutputLen(res.StdOut, res.StdErr); err != nil {
		return
	}

	return ctx.JSON(res)
}

// --- UTIL ---

func (t *Router) checkOutputLen(stdout, stderr string) (err error) {
	maxOutLen, err := util.ParseMemoryStr(t.cfg.Config().API.MaxOutputLen)
	if err != nil {
		return err
	}

	if int64(len(stdout))+int64(len(stderr)) > maxOutLen {
		return errOutputLenExceeded
	}

	return nil
}
