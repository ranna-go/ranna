package docker

import (
	"context"

	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/client"
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

func newSandbox(client *client.Client, container *client.ContainerCreateResult) *Sandbox {
	return &Sandbox{
		logger:    log.Tagged("Sandbox"),
		client:    client,
		container: container,
	}
}

func (t *Sandbox) ID() string {
	return t.container.ID
}

func (t *Sandbox) Run(cOut, cErr chan []byte, cClose chan bool) (err error) {
	buffStdout := chanwriter.New(cOut)
	buffStderr := chanwriter.New(cErr)
	res, err := t.client.ContainerAttach(context.TODO(), t.container.ID, client.ContainerAttachOptions{
		Stdout: true,
		Stderr: true,
		Stream: true,
	})
	if err != nil {
		return err
	}
	defer res.Close()
	t.logger.Debug().Fields("id", t.container.ID).Msg("container attached")

	cErrStdCopy := make(chan error)

	go func() {
		_, err := stdcopy.StdCopy(buffStdout, buffStderr, res.Reader)
		if err != nil {
			t.logger.Error().Err(err).Msg("failed copying stdin/stdout")
			cErrStdCopy <- err
		}
	}()

	_, err = t.client.ContainerStart(context.TODO(), t.container.ID, client.ContainerStartOptions{})
	if err != nil {
		return err
	}
	t.logger.Debug().Fields("id", t.container.ID).Msg("container started")

	wait := t.client.ContainerWait(context.TODO(), t.container.ID, client.ContainerWaitOptions{})
	select {
	case err = <-wait.Error:
	case err = <-cErrStdCopy:
	case <-wait.Result:
	}

	t.logger.Debug().Fields("id", t.container.ID).Msg("container finished")

	cClose <- true
	return err
}

func (t *Sandbox) IsRunning() (ok bool, err error) {
	ctn, err := t.client.ContainerInspect(context.TODO(), t.container.ID, client.ContainerInspectOptions{})
	if err != nil {
		return
	}

	ok = ctn.Container.State.Running
	return
}

func (t *Sandbox) Kill() error {
	_, err := t.client.ContainerKill(context.TODO(), t.container.ID, client.ContainerKillOptions{})
	return err
}

func (t *Sandbox) Delete() error {
	_, err := t.client.ContainerRemove(context.TODO(), t.container.ID, client.ContainerRemoveOptions{})
	return err
}
