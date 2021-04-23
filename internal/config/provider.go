package config

type Provider interface {
	Load() error
	Config() *Config
}
