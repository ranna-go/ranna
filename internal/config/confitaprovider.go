package config

import (
	"context"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
	"github.com/heetch/confita/backend/flags"
)

type Provider struct {
	cfg *Config
}

func NewProvider() *Provider {
	return &Provider{
		cfg: &defaults,
	}
}

func (p *Provider) Load() error {
	loader := confita.NewLoader(
		env.NewBackend(),
		file.NewOptionalBackend("config.json"),
		file.NewOptionalBackend("config.yaml"),
		flags.NewBackend(),
	)
	return loader.Load(context.Background(), p.cfg)
}

func (p *Provider) Config() *Config {
	return p.cfg
}
