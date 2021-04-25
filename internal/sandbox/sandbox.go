package sandbox

type Sandbox interface {
	ID() string
	Run() (stdout, stderr string, err error)
	IsRunning() (bool, error)
	Kill() error
	Delete() error
}

type Provider interface {
	CreateSandbox(spec RunSpec) (Sandbox, error)
}
