package docker

import (
	"bytes"
	"fmt"

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

	fmt.Println(s.container.ID)
	err = s.client.StartContainer(s.container.ID, nil)
	if err != nil {
		return
	}

	waiter.Wait()
	stdout = buffStdout.String()
	stderr = buffStderr.String()
	return
}

func (s *DockerSandbox) Delete() error {
	return s.client.RemoveContainer(dockerclient.RemoveContainerOptions{
		ID: s.container.ID,
	})
}
