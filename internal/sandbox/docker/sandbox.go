package docker

import (
	dockerclient "github.com/fsouza/go-dockerclient"
	"github.com/ranna-go/ranna/pkg/chanwriter"
)

// DockerSandbox implements Sandbox for
// Docker containers.
type DockerSandbox struct {
	client    *dockerclient.Client
	container *dockerclient.Container
}

func (s *DockerSandbox) ID() string {
	return s.container.ID
}

func (s *DockerSandbox) Run(cOut, cErr chan []byte, cClose chan bool) (err error) {
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

func (s *DockerSandbox) IsRunning() (ok bool, err error) {
	ctn, err := s.client.InspectContainer(s.container.ID)
	if err != nil {
		return
	}

	ok = ctn.State.Running
	return
}

func (s *DockerSandbox) Kill() error {
	return s.client.KillContainer(dockerclient.KillContainerOptions{
		ID: s.container.ID,
	})
}

func (s *DockerSandbox) Delete() error {
	return s.client.RemoveContainer(dockerclient.RemoveContainerOptions{
		ID: s.container.ID,
	})
}
