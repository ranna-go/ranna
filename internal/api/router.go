package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
)

type Router interface {
	Setup(route fiber.Route, ctn di.Container)
}
