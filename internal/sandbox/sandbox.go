package sandbox

import "github.com/ranna-go/ranna/pkg/models"

type Sandbox interface {
	ID() string
	Run(bufferCap int) (res *models.ExecutionResponse, err error)
	IsRunning() (bool, error)
	Kill() error
	Delete() error
}

type Provider interface {
	Prepare(spec models.Spec) error
	CreateSandbox(spec RunSpec) (Sandbox, error)
}
