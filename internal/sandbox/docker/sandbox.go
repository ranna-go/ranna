package docker

import (
	"context"

	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/client"
	"github.com/ranna-go/ranna/internal/sandbox"
	"github.com/ranna-go/ranna/pkg/chanwriter"
	"github.com/zekrotja/rogu"
	"github.com/zekrotja/rogu/log"
)

// Sandbox implements Sandbox for
// Docker containers.
type Sandbox struct {
	logger    rogu.Logger
	client    *client.Client
	container *client.ContainerCreateResult
}

var _ sandbox.Sandbox = (*Sandbox)(nil)

func newSandbox(client *client.Client, container *client.ContainerCreateResult) *Sandbox {
	return &Sandbox{
		logger:    log.Tagged("Sandbox"),
		client:    client,
		container: container,
	}
}

func (s *Sandbox) ID() string {
	return s.container.ID
}

func (s *Sandbox) Run(cOut, cErr chan []byte, cClose chan bool) (err error) {
	buffStdout := chanwriter.New(cOut)
	buffStderr := chanwriter.New(cErr)
	res, err := s.client.ContainerAttach(context.TODO(), s.container.ID, client.ContainerAttachOptions{
		Stdout: true,
		Stderr: true,
		Stream: true,
	})
	if err != nil {
		return err
	}
	defer res.Close()
	s.logger.Debug().Fields("id", s.container.ID).Msg("container attached")

	cErrStdCopy := make(chan error)

	go func() {
		_, err := stdcopy.StdCopy(buffStdout, buffStderr, res.Reader)
		if err != nil {
			s.logger.Error().Err(err).Msg("failed copying stdin/stdout")
			cErrStdCopy <- err
		}
	}()

	_, err = s.client.ContainerStart(context.TODO(), s.container.ID, client.ContainerStartOptions{})
	if err != nil {
		return err
	}
	s.logger.Debug().Fields("id", s.container.ID).Msg("container started")

	wait := s.client.ContainerWait(context.TODO(), s.container.ID, client.ContainerWaitOptions{})
	select {
	case err = <-wait.Error:
	case err = <-cErrStdCopy:
	case <-wait.Result:
	}

	s.logger.Debug().Fields("id", s.container.ID).Msg("container finished")

	cClose <- true
	return err
}

func (s *Sandbox) IsRunning() (ok bool, err error) {
	ctn, err := s.client.ContainerInspect(context.TODO(), s.container.ID, client.ContainerInspectOptions{})
	if err != nil {
		return
	}

	ok = ctn.Container.State.Running
	return
}

func (s *Sandbox) Kill() error {
	_, err := s.client.ContainerKill(context.TODO(), s.container.ID, client.ContainerKillOptions{})
	return err
}

func (s *Sandbox) Delete() error {
	_, err := s.client.ContainerRemove(context.TODO(), s.container.ID, client.ContainerRemoveOptions{})
	return err
}
