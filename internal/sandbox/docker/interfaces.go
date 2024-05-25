package docker

import "github.com/ranna-go/ranna/internal/config"

type ConfigProvider interface {
	Config() *config.Config
}
