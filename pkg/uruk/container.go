package uruk

import (
	"archive/tar"
	"bytes"
	"context"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/step/angmar/pkg/testutils"
	"github.com/step/saurontypes"
	"github.com/step/uruk/pkg/tarutils"
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

func (u Uruk) copyToContainer(containerID, src, dest string) error {
	var buffer bytes.Buffer
	u.logCopyToContainer(containerID, src, dest)
	if err := tarutils.Tar(src, &buffer, u.Tarable); err != nil {
		return err
	}
	ctx := context.Background()
	return u.DClient.CopyToContainer(ctx, containerID, dest, &buffer, types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
	})
}

func (u Uruk) copyFromContainer(containerID, src string) (rerr error) {
	u.logCopyFromContainer(containerID, src)
	readCloser, _, err := u.DClient.CopyFromContainer(context.Background(), containerID, src)
	if err != nil {
		return err
	}

	defer func() {
		if rerr == nil {
			rerr = readCloser.Close()
		}
	}()

	mapFiles := testutils.NewMapFiles()
	UntarWithoutGz(readCloser, &mapFiles)

	return nil
}

func (u Uruk) startContainer(ctx context.Context, containerID string) error {
	u.logStartContainer(containerID)
	return u.DClient.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
}

func (u Uruk) killContainer(ctx context.Context, containerID string) error {
	u.logKillContainer(containerID)
	err := u.DClient.ContainerKill(ctx, containerID, "SIGTERM")
	if err != nil {
		return err
	}
	u.DClient.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{})
	return nil
}

func (u Uruk) removeContainer(ctx context.Context, containerID string) error {
	u.logRemoveContainer(containerID)
	return u.DClient.ContainerRemove(context.Background(), containerID, types.ContainerRemoveOptions{Force: true})
}
