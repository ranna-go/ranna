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
			API:     API{},
			Sandbox: Sandbox{},
		},
	}
}

func (t *EnvProvider) Load() (err error) {

	t.c.Debug = t.getBool("DEBUG", false)
	t.c.SpecFile = t.getString("SPECFILE", "spec/spec.yaml")
	t.c.HostRootDir = t.getString("HOSTROOTDIR", "/var/opt/ranna")
	t.c.Sandbox.TimeoutSeconds, err = t.getInt("EXECUTIONTIMEOUTSECONDS", 20)
	if err != nil {
		return
	}

	t.c.API.BindAddress = t.getString("API_BINDADDRESS", ":8080")

	t.c.Sandbox.Memory = t.getString("RESOURCES_MEMORY", "100M")

	return
}

func (t *EnvProvider) Config() *Config {
	return t.c
}

func (t *EnvProvider) getString(key, def string) (v string) {
	var ok bool
	if v, ok = os.LookupEnv(t.prefix + key); !ok {
		v = def
	}
	return
}

func (t *EnvProvider) getBool(key string, def bool) (v bool) {
	defStr := ""
	if def {
		defStr = "true"
	}

	vStr := strings.ToLower(t.getString(key, defStr))
	return vStr == "true" || vStr == "1"
}

func (t *EnvProvider) getInt(key string, def int) (v int, err error) {
	vStr := t.getString(key, "")
	if vStr == "" {
		v = def
		return
	}

	v, err = strconv.Atoi(vStr)
	return
}
