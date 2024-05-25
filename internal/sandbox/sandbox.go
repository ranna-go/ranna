package sandbox

import "github.com/ranna-go/ranna/pkg/models"

// Sandbox defines an interface to control an encapsulated
// code execution environment.
type Sandbox interface {

	// ID returns a unique ID of the sandbox.
	ID() string

	// Run starts the execution of the sandbox
	// blocking and returns the execution response
	// information.
	//
	// bufferCap defines the maximum size of the
	// output stream buffers used to capture the
	// sandbox stdout and stderr streams.
	Run(cOut, cErr chan []byte, cClose chan bool) (err error)

	// IsRunning returns true if the sandbox is
	// still executing.
	IsRunning() (bool, error)

	// Kill stops the sandbox instantly without
	// taking care of the teardown of internal
	// processes.
	//
	// It's like plugging the cable. ;)
	Kill() error

	// Delete tears down the used resources
	// of the sandbox and deletes it.
	Delete() error
}

// Provider defines an interface to prepare the
// environment by the given specs and creating
// sandboxes by given spec.
type Provider interface {

	// Prepare runs necessary tasks to speed up first
	// time startups of sandboxes.
	//
	// This pulls images used in specs, for example.
	Prepare(spec models.Spec, force bool) error

	// CreateSandbox creates a new sandbox by given spec,
	// allocates necessary resources for the sandbox and
	// prepare it to be run.
	CreateSandbox(spec RunSpec) (Sandbox, error)

	// Info returns general information about the used
	// sandbox provider.
	Info() (*models.SandboxInfo, error)
}
