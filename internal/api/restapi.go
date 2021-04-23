package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	v1 "github.com/zekroTJA/ranna/internal/api/v1"
	"github.com/zekroTJA/ranna/internal/config"
	"github.com/zekroTJA/ranna/internal/static"
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
		ServerHeader: "ranna",
		ErrorHandler: func(c *fiber.Ctx, e error) error {
			return fiber.DefaultErrorHandler(c, e)
		},
	})

	new(v1.Router).Setup(r.app.Group("/v1"), ctn)

	return
}

func (r *RestAPI) ListenAndServeBlocking() error {
	return r.app.Listen(r.bindAddress)
}
