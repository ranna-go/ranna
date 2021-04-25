package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/ranna/internal/api"
	"github.com/zekroTJA/ranna/internal/config"
	"github.com/zekroTJA/ranna/internal/file"
	"github.com/zekroTJA/ranna/internal/namespace"
	"github.com/zekroTJA/ranna/internal/sandbox"
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
			p := config.NewConfitaProvider()
			return p, p.Load()
		},
	})

	diBuilder.Add(di.Def{
		Name: static.DiSpecProvider,
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get(static.DiConfigProvider).(config.Provider)
			p := spec.NewFileProvider(cfg.Config().SpecFile)
			return p, p.Load()
		},
	})

	diBuilder.Add(di.Def{
		Name: static.DiSandboxProvider,
		Build: func(ctn di.Container) (interface{}, error) {
			return docker.NewDockerSandboxProvider(ctn)
		},
	})

	diBuilder.Add(di.Def{
		Name: static.DiSandboxManager,
		Build: func(ctn di.Container) (interface{}, error) {
			return sandbox.NewManager(ctn)
		},
		Close: func(obj interface{}) error {
			logrus.Info("cleaning up running sandboxes...")
			m := obj.(sandbox.Manager)
			m.TryCleanup()
			return nil
		},
	})

	diBuilder.Add(di.Def{
		Name: static.DiFileProvider,
		Build: func(ctn di.Container) (v interface{}, err error) {
			cfg := ctn.Get(static.DiConfigProvider).(config.Provider)
			if cfg.Config().Debug {
				v = file.NewDummyFileProvider()
			} else {
				v = file.NewLocalFileProvider()
			}
			return
		},
	})

	diBuilder.Add(di.Def{
		Name: static.DiAPI,
		Build: func(ctn di.Container) (interface{}, error) {
			return api.NewRestAPI(ctn)
		},
	})

	diBuilder.Add(di.Def{
		Name: static.DiNamespaceProvider,
		Build: func(ctn di.Container) (v interface{}, err error) {
			cfg := ctn.Get(static.DiConfigProvider).(config.Provider)
			if cfg.Config().Debug {
				v = namespace.NewDummyProvider("test1")
			} else {
				v = namespace.NewRandomProvider()
			}
			return
		},
	})

	ctn := diBuilder.Build()

	cfg := ctn.Get(static.DiConfigProvider).(config.Provider)
	logrus.SetLevel(logrus.Level(cfg.Config().Log.Level))
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: cfg.Config().Debug,
	})

	api := ctn.Get(static.DiAPI).(api.API)
	go api.ListenAndServeBlocking()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Tear down dependency instances
	ctn.DeleteWithSubContainers()
}
