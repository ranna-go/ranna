package config

type API struct {
	BindAddress string `json:"bindaddress" yaml:"api"`
}

type Config struct {
	SpecFile string `json:"specfile" yaml:"specfile"`
	API      *API   `json:"api" yaml:"api"`
}
