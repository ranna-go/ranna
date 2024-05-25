package docker

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	dockerclient "github.com/fsouza/go-dockerclient"
	"github.com/ranna-go/ranna/internal/sandbox"
	"github.com/ranna-go/ranna/internal/util"
	"github.com/ranna-go/ranna/pkg/models"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
)

const (
	containerRootPath = "/var/tmp/exec"
)

type Provider struct {
	cfg    ConfigProvider
	client *dockerclient.Client
}

var _ sandbox.Provider = (*Provider)(nil)

func NewProvider(cfg ConfigProvider) (dsp *Provider, err error) {
	dsp = &Provider{}

	dsp.cfg = cfg

	dsp.client, err = dockerclient.NewClientFromEnv()
	if err != nil {
		return
	}

	return
}

func (dsp *Provider) Info() (v *models.SandboxInfo, err error) {
	info, err := dsp.client.Info()
	if err != nil {
		return
	}
	v = &models.SandboxInfo{
		Type:    "docker",
		Version: info.ServerVersion,
	}
	return
}

func (dsp *Provider) Prepare(spec models.Spec, force bool) (err error) {
	repo, tag := getImage(spec.Image)

	if force {
		err = dockerclient.ErrNoSuchImage
	} else {
		_, err = dsp.client.InspectImage(repo + ":" + tag)
	}
	if errors.Is(err, dockerclient.ErrNoSuchImage) {
		logrus.WithFields(logrus.Fields{
			"repo": repo,
			"tag":  tag,
		}).Info("pull image")
		err = dsp.client.PullImage(dockerclient.PullImageOptions{
			Repository: repo,
			Tag:        tag,
			Registry:   spec.Registry,
		}, dockerclient.AuthConfiguration{})
	}

	return
}

func (dsp *Provider) CreateSandbox(spec sandbox.RunSpec) (sbx sandbox.Sandbox, err error) {
	repo, tag := getImage(spec.Image)

	err = dsp.Prepare(spec.Spec, false)
	if err != nil {
		return
	}

	workingDir := path.Join(containerRootPath, spec.Subdir)
	ctnCfg := &dockerclient.Config{
		Image:           repo + ":" + tag,
		WorkingDir:      workingDir,
		Entrypoint:      spec.GetEntrypoint(),
		Cmd:             spec.GetCommandWithArgs(),
		Env:             spec.GetEnv(),
		NetworkDisabled: !dsp.cfg.Config().Sandbox.EnableNetworking,
	}

	hostDir, err := filepath.Abs(spec.GetAssembledHostDir())
	if err != nil {
		return
	}

	hostCfg := &dockerclient.HostConfig{
		Binds:   []string{hostDir + ":" + workingDir},
		Runtime: dsp.cfg.Config().Sandbox.Runtime,
	}

	hostCfg.Memory, err = util.ParseMemoryStr(dsp.cfg.Config().Sandbox.Memory)
	if err != nil {
		return
	}
	hostCfg.MemorySwap = hostCfg.Memory

	container, err := dsp.client.CreateContainer(dockerclient.CreateContainerOptions{
		Config:     ctnCfg,
		HostConfig: hostCfg,
		Name:       fmt.Sprintf("ranna-%s-%s", spec.Language, xid.New().String()),
	})

	sbx = &Sandbox{
		client:    dsp.client,
		container: container,
	}

	return
}

func getImage(environmentDescriptor string) (repo, tag string) {
	split := strings.SplitN(environmentDescriptor, ":", 2)
	if len(split) == 1 {
		split = append(split, "latest")
	}

	return split[0], split[1]
}
