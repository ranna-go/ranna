package config

import (
	"os"
)

type EnvProvider struct {
	prefix string

	c *Config
}

func NewEnvProvider(prefix string) *EnvProvider {
	return &EnvProvider{
		prefix: prefix,
		c: &Config{
			API: &API{},
		},
	}
}

func (ep *EnvProvider) Load() (err error) {

	ep.c.SpecFile = ep.getString("SPECFILE", "spec/spec.yaml")
	ep.c.API.BindAddress = ep.getString("API_BINDADDRESS", ":8080")

	return
}

func (ep *EnvProvider) Get() *Config {
	return ep.c
}

func (ep *EnvProvider) getString(key, def string) (v string) {
	var ok bool
	if v, ok = os.LookupEnv(ep.prefix + key); !ok {
		v = def
	}
	return
}
