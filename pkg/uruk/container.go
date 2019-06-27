package uruk

import (
	"bytes"
	"context"
	"fmt"
	"io"
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
	return u.DClient.ContainerCreate(ctx, &container.Config{
		Image: name,
		Env:   []string{},
	}, nil, nil, "")
}

func (u Uruk) copyToContainer(containerId, repoLocation string) error {
	var buffer bytes.Buffer
	location := filepath.Join(u.SourceMountPoint, repoLocation)
	if err := tarutils.Tar(location, &buffer, u.Tarable); err != nil {
		return err
	}
	ctx := context.Background()
	return u.DClient.CopyToContainer(ctx, containerId, "/", &buffer, types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
	})
}

func (u Uruk) copyFromContainer(containerId, src string) io.ReadCloser {
	readCloser, _, err := u.DClient.CopyFromContainer(context.Background(), containerId, src)
	fmt.Println("Error copying", err)
	return readCloser
}

func (u Uruk) startContainer(ctx context.Context, containerId string) error {
	return u.DClient.ContainerStart(ctx, containerId, types.ContainerStartOptions{})
}
