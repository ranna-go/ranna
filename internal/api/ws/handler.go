package ws

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func Upgrade() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			if ip := c.IP(); ip != "" {
				c.Locals("ip", ip)
			}
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}
}

func Handler(cfg ConfigProvider, manager SandboxManager) fiber.Handler {
	rlm := NewRateLimitManager(cfg)
	return newSession(rlm, manager).Handler()
}
