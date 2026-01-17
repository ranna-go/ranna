package v1

import (
	"github.com/ranna-go/ranna/internal/config"
	"github.com/ranna-go/ranna/internal/sandbox"
	"github.com/ranna-go/ranna/internal/spec"
	"github.com/ranna-go/ranna/pkg/models"
)

type SpecProvider interface {
	Spec() *spec.SafeSpecMap
}

type ConfigProvider interface {
	Config() *config.Config
}

type SandboxManager interface {
	RunInSandbox(
		req *models.ExecutionRequest,
		cSpn chan string,
		cOut chan []byte,
		cErr chan []byte,
		cClose chan bool,
	) (err error)
	PrepareEnvironments(force bool) []error
	KillAndCleanUp(id string) (bool, error)
	Cleanup() []error
	GetProvider() sandbox.Provider
}
