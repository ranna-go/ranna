package api

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	v1 "github.com/ranna-go/ranna/internal/api/v1"
	"github.com/ranna-go/ranna/pkg/models"
	"github.com/sirupsen/logrus"
)

type RestAPI struct {
	bindAddress string
	app         *fiber.App
}

func NewRestAPI(cfg ConfigProvider, spec SpecProvider, manager SandboxManager) (r *RestAPI, err error) {

	r = &RestAPI{
		bindAddress: cfg.Config().API.BindAddress,
	}

	var trustedProxies []string
	if tp := cfg.Config().API.TrustedProxies; tp != "" {
		trustedProxies = strings.Split(tp, " ")
	}
	r.app = fiber.New(fiber.Config{
		DisableStartupMessage:   !cfg.Config().Debug,
		ServerHeader:            "ranna",
		ErrorHandler:            errorHandler,
		EnableTrustedProxyCheck: len(trustedProxies) != 0,
		TrustedProxies:          trustedProxies,
		ProxyHeader:             "X-Forwarded-For",
	})

	new(v1.Router).Setup(r.app.Group("/v1"), cfg, spec, manager)

	return
}

func (r *RestAPI) ListenAndServeBlocking() error {
	logrus.WithFields(logrus.Fields{"addr": r.bindAddress}).Info("Starting REST API ...")
	return r.app.Listen(r.bindAddress)
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
