package main

import (
	"github.com/joho/godotenv"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/ranna/internal/api"
	"github.com/zekroTJA/ranna/internal/config"
	"github.com/zekroTJA/ranna/internal/sandbox/docker"
	"github.com/zekroTJA/ranna/internal/spec"
	"github.com/zekroTJA/ranna/internal/static"
)

func main() {
	godotenv.Load()

	diBuilder, _ := di.NewBuilder()

	diBuilder.Add(di.Def{
		Name: static.DiConfigProvider,
		Build: func(ctn di.Container) (interface{}, error) {
			p := config.NewEnvProvider("RANNA_")
			err := p.Load()
			return p, err
		},
	})

	diBuilder.Add(di.Def{
		Name: static.DiSpecProvider,
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get(static.DiConfigProvider).(config.Provider)
			return spec.NewFileProvider(cfg.Get().SpecFile), nil
		},
	})

	diBuilder.Add(di.Def{
		Name: static.DiSpec,
		Build: func(ctn di.Container) (interface{}, error) {
			provider := ctn.Get(static.DiSpecProvider).(spec.Provider)
			return provider.Load()
		},
	})

	diBuilder.Add(di.Def{
		Name: static.DiSandboxProvider,
		Build: func(ctn di.Container) (interface{}, error) {
			return docker.NewDockerSandboxProvider()
		},
	})

	diBuilder.Add(di.Def{
		Name: static.DiAPI,
		Build: func(ctn di.Container) (interface{}, error) {
			return api.NewRestAPI(ctn)
		},
	})

	ctn := diBuilder.Build()

	api := ctn.Get(static.DiAPI).(api.API)
	api.ListenAndServeBlocking()
}
