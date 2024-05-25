package docker

import (
	dockerclient "github.com/fsouza/go-dockerclient"
	"github.com/ranna-go/ranna/internal/sandbox"
	"github.com/ranna-go/ranna/pkg/chanwriter"
)

// Sandbox implements Sandbox for
// Docker containers.
type Sandbox struct {
	client    *dockerclient.Client
	container *dockerclient.Container
}

var _ sandbox.Sandbox = (*Sandbox)(nil)

func (s *Sandbox) ID() string {
	return s.container.ID
}

func (s *Sandbox) Run(cOut, cErr chan []byte, cClose chan bool) (err error) {
	buffStdout := chanwriter.New(cOut)
	buffStderr := chanwriter.New(cErr)
	waiter, err := s.client.AttachToContainerNonBlocking(dockerclient.AttachToContainerOptions{
		Container:    s.container.ID,
		Stdout:       true,
		Stderr:       true,
		Stream:       true,
		OutputStream: buffStdout,
		ErrorStream:  buffStderr,
	})
	if err != nil {
		return
	}

	err = s.client.StartContainer(s.container.ID, nil)
	if err != nil {
		return
	}

	waiter.Wait()
	cClose <- true
	return
}

func (s *Sandbox) IsRunning() (ok bool, err error) {
	ctn, err := s.client.InspectContainerWithOptions(dockerclient.InspectContainerOptions{ID: s.container.ID})
	if err != nil {
		return
	}

	ok = ctn.State.Running
	return
}

func (s *Sandbox) Kill() error {
	return s.client.KillContainer(dockerclient.KillContainerOptions{
		ID: s.container.ID,
	})
}

func (s *Sandbox) Delete() error {
	return s.client.RemoveContainer(dockerclient.RemoveContainerOptions{
		ID: s.container.ID,
	})
}
