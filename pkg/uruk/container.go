package uruk

import (
	"bytes"
	"fmt"
	"context"
	"archive/tar"
	"time"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/step/saurontypes"
	"github.com/step/uruk/pkg/tarutils"
	"github.com/step/angmar/pkg/testutils"
)


func UntarWithoutGz(reader io.Reader, extractor Extractor) (rerr error) {
	tarReader := tar.NewReader(reader)
	// For each header in the tar stream call the appropriate Extractor function
	for header, err := tarReader.Next(); err != io.EOF; header, err = tarReader.Next() {
		if err != nil {
			return TarHeaderError{err, "Untar"}
		}
		extract := extractor.ExtractFile
		if header.FileInfo().IsDir() {
			extract = extractor.ExtractDir
		}

		if err := extract(*header, tarReader); err != nil {
			return ExtractionError{header.Name, header.Mode, err, "Untar"}
		}
	}

	return nil
}

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

func (u Uruk) copyToContainer(containerId, src, dest string) error {
	var buffer bytes.Buffer
	u.logCopyToContainer(containerId, src, dest)
	if err := tarutils.Tar(src, &buffer, u.Tarable); err != nil {
		return err
	}
	ctx := context.Background()
	return u.DClient.CopyToContainer(ctx, containerId, dest, &buffer, types.CopyToContainerOptions{
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

	mapFiles := testutils.NewMapFiles();
	UntarWithoutGz(readCloser, &mapFiles)
	fmt.Println(mapFiles)

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
