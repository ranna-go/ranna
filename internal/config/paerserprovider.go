package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/ranna-go/paerser/env"
	"github.com/ranna-go/paerser/file"
)

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

func (p *Paerser) Config() *Config {
	return p.cfg
}

func (p *Paerser) Load() (err error) {
	cfg := defaults

	cfgFile := defaultConfigLoc
	if p.configFile != "" {
		cfgFile = p.configFile
	}
	if err = file.Decode(cfgFile, &cfg); err != nil && !os.IsNotExist(err) {
		return
	}

	godotenv.Load()
	if err = env.Decode(os.Environ(), "RANNA_", &cfg); err != nil {
		return
	}

	p.cfg = &cfg

	return
}
