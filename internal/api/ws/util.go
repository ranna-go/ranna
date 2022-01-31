package ws

import (
	"strings"

	"github.com/gofiber/websocket/v2"
)

func getAddr(c *websocket.Conn) string {
	if ip, ok := c.Locals("ip").(string); ok && ip != "" {
		return ip
	}
	addr := c.RemoteAddr().String()
	return addr[:strings.LastIndex(addr, ":")]
}
