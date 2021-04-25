package config

type API struct {
	BindAddress  string `config:"api.bindaddress,required" json:"bindaddress" yaml:"api"`
	MaxOutputLen string `config:"api.maxoutputlen" json:"maxoutputlen" yaml:"maxoutputlen"`
}

type Sandbox struct {
	Memory         string `config:"sandbox.memory" json:"memory" yaml:"memory"`
	TimeoutSeconds int    `config:"sandbox.timeoutseconds" json:"executiontimeoutseconds" yaml:"executiontimeoutseconds"`
}

type Config struct {
	Debug       bool    `config:"debug" json:"debug" yaml:"debug"`
	SpecFile    string  `config:"specfile" json:"specfile" yaml:"specfile"`
	HostRootDir string  `config:"hostrootdir" json:"hostrootdir" yaml:"hostrootdir"`
	API         API     `json:"api" yaml:"api"`
	Sandbox     Sandbox `json:"sandbox" yaml:"sandbox"`
}

var defaults = Config{
	Debug:       false,
	SpecFile:    "spec/spec.yaml",
	HostRootDir: "/var/opt/ranna",
	API: API{
		BindAddress:  ":8080",
		MaxOutputLen: "1M",
	},
	Sandbox: Sandbox{
		Memory:         "100M",
		TimeoutSeconds: 20,
	},
}
