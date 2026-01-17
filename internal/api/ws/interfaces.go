package ws

import (
	"github.com/ranna-go/ranna/internal/config"
	"github.com/ranna-go/ranna/pkg/models"
)

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
	KillAndCleanUp(id string) (bool, error)
}
