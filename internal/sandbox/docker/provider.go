package docker

import (
	"path"
	"strings"

	dockerclient "github.com/fsouza/go-dockerclient"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/ranna/internal/config"
	"github.com/zekroTJA/ranna/internal/sandbox"
	"github.com/zekroTJA/ranna/internal/static"
	"github.com/zekroTJA/ranna/internal/util"
)

const (
	containerRootPath = "/var/tmp/exec"
)

type DockerSandboxProvider struct {
	cfg    config.Provider
	client *dockerclient.Client
}

func NewDockerSandboxProvider(ctn di.Container) (dsp *DockerSandboxProvider, err error) {
	dsp = &DockerSandboxProvider{}

	dsp.cfg = ctn.Get(static.DiConfigProvider).(config.Provider)

	dsp.client, err = dockerclient.NewClientFromEnv()
	if err != nil {
		return
	}

	return
}

func (dsp *DockerSandboxProvider) CreateSandbox(spec sandbox.RunSpec) (sbx sandbox.Sandbox, err error) {
	repo, tag := getImage(spec.Image)

	err = dsp.client.PullImage(dockerclient.PullImageOptions{
		Repository: repo,
		Tag:        tag,
	}, dockerclient.AuthConfiguration{})
	if err != nil {
		return
	}

	workingDir := path.Join(containerRootPath, spec.Subdir)
	cfg := &dockerclient.Config{
		Image:           repo + ":" + tag,
		WorkingDir:      workingDir,
		Entrypoint:      spec.GetEntrypoint(),
		Cmd:             spec.GetCommandWithArgs(),
		Env:             spec.GetEnv(),
		NetworkDisabled: true,
	}

	resources := dsp.cfg.Config().Resources
	if resources != nil {
		cfg.Memory, err = util.ParseMemoryStr(resources.Memory)
		if err != nil {
			return
		}
	}

	hostDir := spec.GetAssambledHostDir()
	container, err := dsp.client.CreateContainer(dockerclient.CreateContainerOptions{
		Config: cfg,
		HostConfig: &dockerclient.HostConfig{
			Binds: []string{hostDir + ":" + workingDir},
		},
	})

	sbx = &DockerSandbox{
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
