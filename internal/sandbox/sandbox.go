package sandbox

type Sandbox interface {
	Run() (stdout, stderr string, err error)
	Kill() error
	Delete() error
}

type Provider interface {
	CreateSandbox(spec RunSpec) (Sandbox, error)
}
