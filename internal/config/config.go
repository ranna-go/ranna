package config

import "github.com/sirupsen/logrus"

type Log struct {
	Level int `config:"log.level" json:"level" yaml:"level"`
}

type API struct {
	BindAddress  string `config:"api.bindaddress,required" json:"bindaddress" yaml:"api"`
	MaxOutputLen string `config:"api.maxoutputlen" json:"maxoutputlen" yaml:"maxoutputlen"`
}

type Sandbox struct {
	Runtime          string `config:"sandbox.runtime" json:"runtime" yaml:"runtime"`
	EnableNetworking bool   `config:"sandbox.enablenetworking" json:"enablenetworking" yaml:"enablenetworking"`
	Memory           string `config:"sandbox.memory" json:"memory" yaml:"memory"`
	TimeoutSeconds   int    `config:"sandbox.timeoutseconds" json:"executiontimeoutseconds" yaml:"executiontimeoutseconds"`
	StreamBufferCap  string `config:"sandbox.streambuffercap" json:"streambuffercap" yaml:"streambuffercap"`
}

type Scheduler struct {
	UpdateImages string `config:"scheduler.updateimages" json:"updateimages" yaml:"updateimages"`
	UpdateSpecs  string `config:"scheduler.updatespecs" json:"updatespecs" yaml:"updatespecs"`
}

type Config struct {
	Debug           bool   `config:"debug" json:"debug" yaml:"debug"`
	SpecFile        string `config:"specfile" json:"specfile" yaml:"specfile"`
	HostRootDir     string `config:"hostrootdir" json:"hostrootdir" yaml:"hostrootdir"`
	SkipStartupPrep bool   `config:"skipstartupprep" json:"skipstartupprep" yaml:"skipstartupprep"`

	Log       Log       `json:"log" yaml:"log"`
	API       API       `json:"api" yaml:"api"`
	Sandbox   Sandbox   `json:"sandbox" yaml:"sandbox"`
	Scheduler Scheduler `json:"scheduler" yaml:"scheduler"`
}

var defaults = Config{
	Debug:           false,
	SpecFile:        "spec/spec.yaml",
	HostRootDir:     "/var/opt/ranna",
	SkipStartupPrep: false,

	Log: Log{
		Level: int(logrus.InfoLevel),
	},
	API: API{
		BindAddress:  ":8080",
		MaxOutputLen: "1M",
	},
	Sandbox: Sandbox{
		Runtime:          "",
		Memory:           "100M",
		TimeoutSeconds:   20,
		StreamBufferCap:  "50M",
		EnableNetworking: false,
	},
	Scheduler: Scheduler{
		UpdateImages: "0 3 * * *",
		UpdateSpecs:  "",
	},
}
