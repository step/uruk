package uruk

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/step/saurontypes"
	"github.com/step/uruk/pkg/tarutils"
)

func (u Uruk) createContainer(message saurontypes.UrukMessage) (container.ContainerCreateCreatedBody, error) {
	name := message.ImageName
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	u.logCreateContainer(name)
	return u.DClient.ContainerCreate(ctx, &container.Config{
		Image: name,
		Env:   []string{},
	}, nil, nil, "")
}

func (u Uruk) copyToContainer(containerId, repoLocation string) error {
	var buffer bytes.Buffer
	location := filepath.Join(u.SourceMountPoint, repoLocation)
	u.logCopyToContainer(containerId, location)
	if err := tarutils.Tar(location, &buffer, u.Tarable); err != nil {
		return err
	}
	ctx := context.Background()
	return u.DClient.CopyToContainer(ctx, containerId, "/", &buffer, types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
	})
}

func (u Uruk) copyFromContainer(containerId, src string) (rerr error) {
	u.logCopyFromContainer(containerId, src)
	readCloser, _, err := u.DClient.CopyFromContainer(context.Background(), containerId, src)
	if err != nil {
		return err
	}

	defer func() {
		if rerr == nil {
			rerr = readCloser.Close()
		}
	}()

	if _, err := io.Copy(os.Stdout, readCloser); err != nil {
		return err
	}

	return nil
}

func (u Uruk) startContainer(ctx context.Context, containerId string) error {
	u.logStartContainer(containerId)
	return u.DClient.ContainerStart(ctx, containerId, types.ContainerStartOptions{})
}

func (u Uruk) killContainer(ctx context.Context, containerId string) error {
	u.logKillContainer(containerId)
	err := u.DClient.ContainerKill(ctx, containerId, "SIGTERM")
	if err != nil {
		return err
	}
	return nil
}

func (u Uruk) removeContainer(ctx context.Context, containerId string) error {
	u.logRemoveContainer(containerId)
	return u.DClient.ContainerRemove(context.Background(), containerId, types.ContainerRemoveOptions{Force: true})
}
