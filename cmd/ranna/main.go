package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/ranna/internal/sandbox/docker"
	"github.com/zekroTJA/ranna/internal/spec"
	"github.com/zekroTJA/ranna/internal/static"
)

func main() {
	godotenv.Load()

	specFile, ok := os.LookupEnv("RANA_SPECFILE")
	if !ok {
		specFile = "spec/spec.yaml"
	}

	diBuilder, _ := di.NewBuilder()

	diBuilder.Add(di.Def{
		Name: static.DiSpecProvider,
		Build: func(ctn di.Container) (interface{}, error) {
			return spec.NewFileProvider(specFile), nil
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

	ctn := diBuilder.Build()

	specInst := ctn.Get(static.DiSpec).(spec.SpecMap)
	fmt.Printf("%+v\n", specInst["go"])
}
