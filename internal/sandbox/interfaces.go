package sandbox

import (
	"github.com/ranna-go/ranna/internal/config"
	"github.com/ranna-go/ranna/internal/spec"
)

type SpecProvider interface {
	Spec() *spec.SafeSpecMap
}

type FileProvider interface {
	CreateDirectory(path string) error
	CreateFileWithContent(path, content string) error
	DeleteDirectory(path string) error
}

type ConfigProvider interface {
	Config() *config.Config
}

type NamespaceProvider interface {
	Get() (string, error)
}
