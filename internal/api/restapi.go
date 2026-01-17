package api

import (
	"errors"
	"strings"

	"github.com/zekrotja/rogu/log"

	"github.com/gofiber/fiber/v2"
	v1 "github.com/ranna-go/ranna/internal/api/v1"
	"github.com/ranna-go/ranna/pkg/models"
)

type RestAPI struct {
	bindAddress string
	app         *fiber.App
}

func NewRestAPI(cfg ConfigProvider, spec SpecProvider, manager SandboxManager) (t *RestAPI, err error) {

	t = &RestAPI{
		bindAddress: cfg.Config().API.BindAddress,
	}

	var trustedProxies []string
	if tp := cfg.Config().API.TrustedProxies; tp != "" {
		trustedProxies = strings.Split(tp, " ")
	}
	t.app = fiber.New(fiber.Config{
		DisableStartupMessage:   !cfg.Config().Debug,
		ServerHeader:            "ranna",
		ErrorHandler:            errorHandler,
		EnableTrustedProxyCheck: len(trustedProxies) != 0,
		TrustedProxies:          trustedProxies,
		ProxyHeader:             "X-Forwarded-For",
	})

	new(v1.Router).Setup(t.app.Group("/v1"), cfg, spec, manager)

	return
}

func (t *RestAPI) ListenAndServeBlocking() error {
	log.Info().Field("addr", t.bindAddress).Msg("Starting REST API ...")
	return t.app.Listen(t.bindAddress)
}

func errorHandler(ctx *fiber.Ctx, err error) error {
	var fErr *fiber.Error
	if errors.As(err, &fErr) {
		ctx.Status(fErr.Code)
		return ctx.JSON(&models.ErrorModel{
			Error: fErr.Message,
			Code:  fErr.Code,
		})
	}

	return errorHandler(ctx,
		fiber.NewError(fiber.StatusInternalServerError, err.Error()))
}
