package uruk

import (
	"fmt"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/step/saurontypes"
)

type ContainerCreationError struct {
	UrukMessage saurontypes.UrukMessage
	err         error
}

func (cce ContainerCreationError) Error() string {
	msg := cce.UrukMessage
	return fmt.Sprintf("Unable to create container %s for %s\n%s", msg.ImageName, msg.RepoLocation, cce.err.Error())
}

type CopyToContainerError struct {
	UrukMessage         saurontypes.UrukMessage
	SourceMountPoint    string
	ContainerSourcePath string
	err                 error
}

func (ctce CopyToContainerError) Error() string {
	msg := ctce.UrukMessage
	source := filepath.Join(ctce.SourceMountPoint, msg.RepoLocation)
	destination := filepath.Join("/", ctce.ContainerSourcePath)
	return fmt.Sprintf("Unable to copy from %s to %s:%s\n%s", source, msg.ImageName, destination, ctce.err.Error())
}

type StartContainerError struct {
	Response    container.ContainerCreateCreatedBody
	UrukMessage saurontypes.UrukMessage
	err         error
}

func (sce StartContainerError) Error() string {
	msg := sce.UrukMessage
	return fmt.Sprintf("Unable to start container of image %s with id %s for %s\nWarnings: %s\n%s",
		msg.ImageName,
		sce.Response.ID,
		msg.RepoLocation,
		sce.Response.Warnings,
		sce.err.Error())
}

type CopyFromContainerError struct {
	UrukMessage saurontypes.UrukMessage
	src         string
	ContainerId string
	err         error
}

func (cfce CopyFromContainerError) Error() string {
	msg := cfce.UrukMessage
	return fmt.Sprintf("Unable to copy %s from container of image %s with id %s\n%s",
		cfce.src,
		msg.ImageName,
		cfce.ContainerId,
		cfce.err.Error())
}
