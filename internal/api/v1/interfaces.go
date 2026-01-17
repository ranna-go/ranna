package v1

import (
	"context"

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
		ctx context.Context,
		req *models.ExecutionRequest,
		cSpn chan string,
		cOut chan []byte,
		cErr chan []byte,
	) (err error)
	PrepareEnvironments(ctx context.Context, force bool) []error
	KillAndCleanUp(ctx context.Context, id string) (bool, error)
	Cleanup(ctx context.Context) []error
	GetProvider() sandbox.Provider
}
