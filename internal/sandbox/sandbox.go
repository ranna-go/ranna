package sandbox

import "github.com/zekroTJA/ranna/internal/models"

type Sandbox interface {
	Run() (stdout, stderr string, err error)
	Delete() error
}

type Provider interface {
	CreateSandbox(spec models.Spec) (Sandbox, error)
}
