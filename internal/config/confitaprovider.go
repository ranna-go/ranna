package config

import (
	"context"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
	"github.com/heetch/confita/backend/flags"
)

type ConfitaProvider struct {
	cfg *Config
}

func NewConfitaProvider() *ConfitaProvider {
	return &ConfitaProvider{
		cfg: &defaults,
	}
}

func (p *ConfitaProvider) Load() error {
	loader := confita.NewLoader(
		env.NewBackend(),
		file.NewOptionalBackend("config.json"),
		file.NewOptionalBackend("config.yaml"),
		flags.NewBackend(),
	)
	return loader.Load(context.Background(), p.cfg)
}

func (p *ConfitaProvider) Config() *Config {
	return p.cfg
}
