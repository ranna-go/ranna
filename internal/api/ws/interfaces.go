package ws

import (
	"context"

	"github.com/ranna-go/ranna/internal/config"
	"github.com/ranna-go/ranna/pkg/models"
)

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
	KillAndCleanUp(ctx context.Context, id string) (bool, error)
}
