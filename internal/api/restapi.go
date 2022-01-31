package api

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	v1 "github.com/ranna-go/ranna/internal/api/v1"
	"github.com/ranna-go/ranna/internal/config"
	"github.com/ranna-go/ranna/internal/static"
	"github.com/ranna-go/ranna/pkg/models"
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
)

type RestAPI struct {
	bindAddress string
	app         *fiber.App
}

func NewRestAPI(ctn di.Container) (r *RestAPI, err error) {
	cfg := ctn.Get(static.DiConfigProvider).(config.Provider)

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

	new(v1.Router).Setup(r.app.Group("/v1"), ctn)

	return
}

func (r *RestAPI) ListenAndServeBlocking() error {
	logrus.WithFields(logrus.Fields{"addr": r.bindAddress}).Info("Starting REST API ...")
	return r.app.Listen(r.bindAddress)
}

func errorHandler(ctx *fiber.Ctx, err error) error {
	if fErr, ok := err.(*fiber.Error); ok {
		ctx.Status(fErr.Code)
		return ctx.JSON(&models.ErrorModel{
			Error: fErr.Message,
			Code:  fErr.Code,
		})
	}

	return errorHandler(ctx,
		fiber.NewError(fiber.StatusInternalServerError, err.Error()))
}
