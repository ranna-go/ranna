package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	v1 "github.com/zekroTJA/ranna/internal/api/v1"
	"github.com/zekroTJA/ranna/internal/config"
	"github.com/zekroTJA/ranna/internal/static"
	"github.com/zekroTJA/ranna/pkg/models"
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

	r.app = fiber.New(fiber.Config{
		DisableStartupMessage: !cfg.Config().Debug,
		ServerHeader:          "ranna",
		ErrorHandler:          errorHandler,
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
