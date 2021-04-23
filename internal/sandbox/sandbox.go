package sandbox

type Sandbox interface {
	Run() (stdout, stderr string, err error)
	Delete() error
}

type SandboxProvider interface {
	CreateSandbox(spec Spec) (Sandbox, error)
}
