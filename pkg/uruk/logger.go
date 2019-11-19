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

func (u Uruk) logStartContainer(containerID string) {
	message := fmt.Sprintf("Starting container %s", containerID)
	u.Logger.Println(message)
}

func (u Uruk) logContainerSuccessful(containerID string, status container.ContainerWaitOKBody) {
	message := fmt.Sprintf("Container %s ran successfully\n---\nStatus Code: %d\nErrors: %s\n---", containerID, status.StatusCode, status.Error)
	u.Logger.Println(message)
}

func (u Uruk) logWaitingForContainer(containerID string) {
	message := fmt.Sprintf("Waiting for container %s", containerID)
	u.Logger.Println(message)
}

func (u Uruk) logRemoveContainer(containerID string) {
	message := fmt.Sprintf("Removing container %s", containerID)
	u.Logger.Println(message)
}

func (u Uruk) logCopyToContainer(containerID, src, dest string) {
	message := fmt.Sprintf("Copying %s to container %s:%s", src, containerID, dest)
	u.Logger.Println(message)
}

func (u Uruk) logCopyFromContainer(containerID, src string) {
	message := fmt.Sprintf("Copying from container %s:%s", containerID, src)
	u.Logger.Println(message)
}

func (u Uruk) logContainerTookTooLong(containerID string) {
	message := fmt.Sprintf("Container %s did not finish execution in allotted time", containerID)
	u.Logger.Println(message)
}

func (u Uruk) logKillContainer(containerID string) {
	message := fmt.Sprintf("Killing container %s", containerID)
	u.Logger.Println(message)
}
