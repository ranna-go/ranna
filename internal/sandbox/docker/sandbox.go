package docker

import (
	"bytes"

	dockerclient "github.com/fsouza/go-dockerclient"
)

type DockerSandbox struct {
	client    *dockerclient.Client
	container *dockerclient.Container
}

func (s *DockerSandbox) Run() (stdout, stderr string, err error) {
	var buffStdout, buffStderr bytes.Buffer
	waiter, err := s.client.AttachToContainerNonBlocking(dockerclient.AttachToContainerOptions{
		Container:    s.container.ID,
		Stdout:       true,
		Stderr:       true,
		Stream:       true,
		OutputStream: &buffStdout,
		ErrorStream:  &buffStderr,
	})
	if err != nil {
		return
	}

	err = s.client.StartContainer(s.container.ID, nil)
	if err != nil {
		return
	}

	waiter.Wait()
	stdout = buffStdout.String()
	stderr = buffStderr.String()
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
