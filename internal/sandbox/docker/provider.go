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

func NewProvider(cfg ConfigProvider) (t *Provider, err error) {
	t = &Provider{}

	t.cfg = cfg
	t.logger = log.Tagged("Provider")

	t.client, err = client.New(client.FromEnv)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (t *Provider) Info(ctx context.Context) (v *models.SandboxInfo, err error) {
	info, err := t.client.Info(ctx, client.InfoOptions{})
	if err != nil {
		return nil, err
	}
	v = &models.SandboxInfo{
		Type:    "docker",
		Version: info.Info.ServerVersion,
	}
	return v, nil
}

func (t *Provider) Prepare(ctx context.Context, spec models.Spec, force bool) (err error) {
	repo, tag := getImage(spec.Image)

	if !force {
		t.logger.Debug().Fields("image", spec.Image).Msg("inspecting image")
		_, err = t.client.ImageInspect(ctx, spec.Image)
		if err == nil {
			return nil
		}
	}

	t.logger.Info().Fields("repo", repo, "tag", tag).Msg("pull image")
	resp, err := t.client.ImagePull(ctx, spec.Image, client.ImagePullOptions{})
	if err != nil {
		return err
	}
	return resp.Wait(ctx)
}

func (t *Provider) CreateSandbox(ctx context.Context, spec sandbox.RunSpec) (sbx sandbox.Sandbox, err error) {
	repo, tag := getImage(spec.Image)

	err = t.Prepare(ctx, spec.Spec, false)
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
		NetworkDisabled: !t.cfg.Config().Sandbox.EnableNetworking,
	}

	hostDir, err := filepath.Abs(spec.GetAssembledHostDir())
	if err != nil {
		return nil, err
	}

	hostCfg := &container.HostConfig{
		Binds:   []string{hostDir + ":" + workingDir},
		Runtime: t.cfg.Config().Sandbox.Runtime,
	}

	hostCfg.Memory, err = util.ParseMemoryStr(t.cfg.Config().Sandbox.Memory)
	if err != nil {
		return nil, err
	}
	hostCfg.MemorySwap = hostCfg.Memory

	container, err := t.client.ContainerCreate(ctx, client.ContainerCreateOptions{
		Config:     ctnCfg,
		HostConfig: hostCfg,
		Name:       fmt.Sprintf("ranna-%s-%s", spec.Language, xid.New().String()),
	})
	t.logger.Debug().Fields("spec", spec.Image, "id", container.ID).Msg("container created")

	sbx = newSandbox(t.client, &container)

	return sbx, nil
}

func getImage(environmentDescriptor string) (repo, tag string) {
	split := strings.SplitN(environmentDescriptor, ":", 2)
	if len(split) == 1 {
		split = append(split, "latest")
	}

	return split[0], split[1]
}
