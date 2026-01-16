package docker

import (
	"context"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	"github.com/rs/xid"
	"github.com/zekrotja/rogu"
	"github.com/zekrotja/rogu/log"

	"github.com/ranna-go/ranna/internal/sandbox"
	"github.com/ranna-go/ranna/internal/util"
	"github.com/ranna-go/ranna/pkg/models"
)

const (
	containerRootPath = "/var/tmp/exec"
)

type Provider struct {
	cfg    ConfigProvider
	logger rogu.Logger
	client *client.Client
}

var _ sandbox.Provider = (*Provider)(nil)

func NewProvider(cfg ConfigProvider) (dsp *Provider, err error) {
	dsp = &Provider{}

	dsp.cfg = cfg
	dsp.logger = log.Tagged("Provider")

	dsp.client, err = client.New(client.FromEnv)
	if err != nil {
		return nil, err
	}

	return dsp, nil
}

func (dsp *Provider) Info() (v *models.SandboxInfo, err error) {
	info, err := dsp.client.Info(context.TODO(), client.InfoOptions{})
	if err != nil {
		return nil, err
	}
	v = &models.SandboxInfo{
		Type:    "docker",
		Version: info.Info.ServerVersion,
	}
	return v, nil
}

func (dsp *Provider) Prepare(spec models.Spec, force bool) (err error) {
	repo, tag := getImage(spec.Image)

	if !force {
		dsp.logger.Debug().Fields("image", spec.Image).Msg("inspecting image")
		_, err = dsp.client.ImageInspect(context.TODO(), spec.Image)
		if err == nil {
			return nil
		}
	}

	dsp.logger.Info().Fields("repo", repo, "tag", tag).Msg("pull image")
	resp, err := dsp.client.ImagePull(context.TODO(), spec.Image, client.ImagePullOptions{})
	if err != nil {
		return err
	}
	return resp.Wait(context.TODO())
}

func (dsp *Provider) CreateSandbox(spec sandbox.RunSpec) (sbx sandbox.Sandbox, err error) {
	repo, tag := getImage(spec.Image)

	err = dsp.Prepare(spec.Spec, false)
	if err != nil {
		return nil, err
	}

	workingDir := path.Join(containerRootPath, spec.Subdir)
	ctnCfg := &container.Config{
		Image:           repo + ":" + tag,
		WorkingDir:      workingDir,
		Entrypoint:      spec.GetEntrypoint(),
		Cmd:             spec.GetCommandWithArgs(),
		Env:             spec.GetEnv(),
		NetworkDisabled: !dsp.cfg.Config().Sandbox.EnableNetworking,
	}

	hostDir, err := filepath.Abs(spec.GetAssembledHostDir())
	if err != nil {
		return nil, err
	}

	hostCfg := &container.HostConfig{
		Binds:   []string{hostDir + ":" + workingDir},
		Runtime: dsp.cfg.Config().Sandbox.Runtime,
	}

	hostCfg.Memory, err = util.ParseMemoryStr(dsp.cfg.Config().Sandbox.Memory)
	if err != nil {
		return nil, err
	}
	hostCfg.MemorySwap = hostCfg.Memory

	container, err := dsp.client.ContainerCreate(context.TODO(), client.ContainerCreateOptions{
		Config:     ctnCfg,
		HostConfig: hostCfg,
		Name:       fmt.Sprintf("ranna-%s-%s", spec.Language, xid.New().String()),
	})
	dsp.logger.Debug().Fields("spec", spec.Image, "id", container.ID).Msg("container created")

	sbx = newSandbox(dsp.client, &container)

	return sbx, nil
}

func getImage(environmentDescriptor string) (repo, tag string) {
	split := strings.SplitN(environmentDescriptor, ":", 2)
	if len(split) == 1 {
		split = append(split, "latest")
	}

	return split[0], split[1]
}
