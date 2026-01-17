package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/ranna-go/paerser/env"
	"github.com/ranna-go/paerser/file"
)

const envPrefix = "RANNA_"

const defaultConfigLoc = "./config.yaml"

type Paerser struct {
	cfg        *Config
	configFile string
}

func NewPaerser(configFile string) *Paerser {
	return &Paerser{
		configFile: configFile,
	}
}

func (t *Paerser) Config() *Config {
	return t.cfg
}

func (t *Paerser) Load() (err error) {
	cfg := defaults

	cfgFile := defaultConfigLoc
	if t.configFile != "" {
		cfgFile = t.configFile
	}
	if err = file.Decode(cfgFile, &cfg); err != nil && !os.IsNotExist(err) {
		return
	}

	godotenv.Load()
	if err = env.Decode(os.Environ(), envPrefix, &cfg); err != nil {
		return
	}

	t.cfg = &cfg

	return
}
