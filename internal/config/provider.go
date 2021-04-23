package config

type Provider interface {
	Load() error
	Get() *Config
}
