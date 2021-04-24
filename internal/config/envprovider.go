package config

import (
	"os"
	"strconv"
	"strings"
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

	ep.c.Debug = ep.getBool("DEBUG", false)
	ep.c.SpecFile = ep.getString("SPECFILE", "spec/spec.yaml")
	ep.c.HostRootDir = ep.getString("HOSTROOTDIR", "/var/opt/ranna")
	ep.c.API.BindAddress = ep.getString("API_BINDADDRESS", ":8080")

	ep.c.ExecutionTimeoutSeconds, err = ep.getInt("EXECUTIONTIMEOUTSECONDS", 20)
	if err != nil {
		return
	}

	return
}

func (ep *EnvProvider) Config() *Config {
	return ep.c
}

func (ep *EnvProvider) getString(key, def string) (v string) {
	var ok bool
	if v, ok = os.LookupEnv(ep.prefix + key); !ok {
		v = def
	}
	return
}

func (ep *EnvProvider) getBool(key string, def bool) (v bool) {
	defStr := ""
	if def {
		defStr = "true"
	}

	vStr := strings.ToLower(ep.getString(key, defStr))
	return vStr == "true" || vStr == "1"
}

func (ep *EnvProvider) getInt(key string, def int) (v int, err error) {
	vStr := ep.getString(key, "")
	if vStr == "" {
		v = def
		return
	}

	v, err = strconv.Atoi(vStr)
	return
}
