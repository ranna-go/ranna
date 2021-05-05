package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/ranna-go/ranna/internal/api"
	"github.com/ranna-go/ranna/internal/config"
	"github.com/ranna-go/ranna/internal/file"
	"github.com/ranna-go/ranna/internal/namespace"
	"github.com/ranna-go/ranna/internal/sandbox"
	"github.com/ranna-go/ranna/internal/sandbox/docker"
	"github.com/ranna-go/ranna/internal/scheduler"
	"github.com/ranna-go/ranna/internal/spec"
	"github.com/ranna-go/ranna/internal/static"
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
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
			specFile := cfg.Config().SpecFile
			var p spec.Provider
			if strings.HasPrefix(specFile, "https://") || strings.HasPrefix(specFile, "http://") {
				p = spec.NewHttpProvider(specFile)
			} else {
				p = spec.NewFileProvider(specFile)
			}
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
			logrus.Info("cleaning up running sandboxes ...")
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

	diBuilder.Add(di.Def{
		Name: static.DiScheduler,
		Build: func(ctn di.Container) (interface{}, error) {
			sched := scheduler.NewCronScheduler()
			sched.Start()
			return sched, nil
		},
		Close: func(obj interface{}) error {
			sched := obj.(scheduler.Scheduler)
			sched.Stop()
			return nil
		},
	})

	ctn := diBuilder.Build()

	cfg := ctn.Get(static.DiConfigProvider).(config.Provider)
	logrus.SetLevel(logrus.Level(cfg.Config().Log.Level))
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: cfg.Config().Debug,
	})

	if !cfg.Config().SkipStartupPrep {
		mgr := ctn.Get(static.DiSandboxManager).(sandbox.Manager)
		logrus.Info("Prepare spec environments ...")
		mgr.PrepareEnvironments(true)
	} else {
		logrus.Warn("Skipping spec preparation on startup")
	}

	if err := scheduleTasks(ctn); err != nil {
		logrus.WithError(err).Fatal("failed scheduling job")
	}

	api := ctn.Get(static.DiAPI).(api.API)
	go api.ListenAndServeBlocking()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Tear down dependency instances
	ctn.DeleteWithSubContainers()
}

func scheduleTasks(ctn di.Container) (err error) {
	cfg := ctn.Get(static.DiConfigProvider).(config.Provider)
	sched := ctn.Get(static.DiScheduler).(scheduler.Scheduler)
	mgr := ctn.Get(static.DiSandboxManager).(sandbox.Manager)

	schedule := func(name, spec string) (err error) {
		if spec != "" {
			logrus.WithField("name", name).WithField("spec", spec).Info("Scheduling job")
			_, err = sched.Schedule(spec, func() {
				logrus.Info("Updating spec environments ...")
				mgr.PrepareEnvironments(true)
			})
		}
		return
	}

	spec := cfg.Config().Scheduler.UpdateImages
	if err = schedule("update spec environments", spec); err != nil {
		return
	}

	return
}
