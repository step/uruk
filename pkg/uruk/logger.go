package uruk

import (
	"fmt"
	"strings"

	"github.com/docker/docker/api/types/container"
)

func (u Uruk) logStart(qName string) {
	var builder strings.Builder
	builder.WriteString("Starting Uruk...\n")
	builder.WriteString("---\n")
	builder.WriteString(u.String() + "\n")
	builder.WriteString("Queue: " + qName + "\n")
	builder.WriteString("---\n")
	u.Logger.Println(builder.String())
}

func (u Uruk) logError(desc string, err error) {
	var builder strings.Builder
	builder.WriteString(desc + "\n")
	builder.WriteString("---\n")
	builder.WriteString(err.Error() + "\n")
	builder.WriteString("---\n")
	u.Logger.Println(builder.String())
}

func (u Uruk) logCreateContainer(imageName string) {
	message := fmt.Sprintf("Creating container for image %s", imageName)
	u.Logger.Println(message)
}

func (u Uruk) logStartContainer(containerId string) {
	message := fmt.Sprintf("Starting container %s", containerId)
	u.Logger.Println(message)
}

func (u Uruk) logContainerSuccessful(containerId string, status container.ContainerWaitOKBody) {
	message := fmt.Sprintf("Container %s ran successfully\n---\nStatus Code: %d\nErrors: %s\n---", containerId, status.StatusCode, status.Error)
	u.Logger.Println(message)
}

func (u Uruk) logWaitingForContainer(containerId string) {
	message := fmt.Sprintf("Waiting for container %s", containerId)
	u.Logger.Println(message)
}

func (u Uruk) logRemoveContainer(containerId string) {
	message := fmt.Sprintf("Removing container %s", containerId)
	u.Logger.Println(message)
}

func (u Uruk) logCopyToContainer(containerId, repoLocation string) {
	message := fmt.Sprintf("Copying %s to container %s:/", repoLocation, containerId)
	u.Logger.Println(message)
}