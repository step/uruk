package uruk

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/step/saurontypes"
	"github.com/step/uruk/pkg/tarutils"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

	q "github.com/step/angmar/pkg/queueclient"
	s "github.com/step/uruk/pkg/streamClient"
)

// Uruk is service which takes a UrukMessage and
// parses it to get the details of job and launches a
// container with the image specified for the job.
// It copies the source and data to the container and
// after container process is done running it copies the
// artifacts from the container and kills the container
type Uruk struct {
	QClient             q.QueueClient
	SClient             s.StreamClient
	DClient             *client.Client
	Tarable             tarutils.Tarable
	SourceMountPoint    string
	ContainerDataPath   string
	ContainerSourcePath string
	NumOfWorkers        int
	Logger              *log.Logger
}

func (u Uruk) String() string {
	var builder strings.Builder
	builder.WriteString(u.QClient.String() + "\n")
	builder.WriteString(fmt.Sprintf("%v\n", u.DClient))
	builder.WriteString("Source mounted at: " + u.SourceMountPoint + "\n")
	builder.WriteString("Copying source to container at: " + u.ContainerSourcePath + "\n")
	builder.WriteString("Copying data to container at: " + u.ContainerDataPath + "\n")
	builder.WriteString("Number of Workers: " + fmt.Sprintf("%d", u.NumOfWorkers))
	return builder.String()
}

func (u Uruk) executeJob(urukMessage saurontypes.UrukMessage) (rerr error) {
	resp, err := u.createContainer(urukMessage)
	if err != nil {
		return ContainerCreationError{urukMessage, err}
	}

	defer func() {
		if rerr == nil {
			rerr = u.removeContainer(context.Background(), resp.ID)
		}
	}()

	sourceLocation := filepath.Join(u.SourceMountPoint, urukMessage.RepoLocation)
	err = u.copyToContainer(resp.ID, sourceLocation, u.ContainerSourcePath)
	if err != nil {
		return CopyToContainerError{urukMessage, sourceLocation, u.ContainerSourcePath, err}
	}

	err = u.copyToContainer(resp.ID, urukMessage.DataPath, u.ContainerDataPath)
	if err != nil {
		return CopyToContainerError{urukMessage, urukMessage.DataPath, u.ContainerDataPath, err}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err = u.startContainer(ctx, resp.ID)
	if err != nil {
		return StartContainerError{resp, urukMessage, err}
	}

	okCh, errCh := u.DClient.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	u.logWaitingForContainer(resp.ID)
	select {
	case status := <-okCh:
		u.logContainerSuccessful(resp.ID, status)
		fileToCopy := "/results/result.json"
		err, content := u.copyFromContainer(resp.ID, fileToCopy)
		if err != nil {
			return CopyFromContainerError{urukMessage, fileToCopy, resp.ID, err}
		}

		details := make(map[string]string)

		details["job"] = urukMessage.Job
		details["result"] = content

		marshalledDetails, err := json.Marshal(details)   
    if err != nil {
        fmt.Println(err.Error())
        return
    }
     
    jsonStr := string(marshalledDetails)
		u.publishEvent(urukMessage, "job_complete", jsonStr)

	case err := <-errCh:
		return err
	case <-ctx.Done():
		u.logContainerTookTooLong(resp.ID)
		return u.killContainer(context.Background(), resp.ID)
	}

	// TODO: InformDB()
	return nil
}

func (u Uruk) publishEvent(message saurontypes.UrukMessage, eventType, details string) {
	event := saurontypes.Event{
		Source:    "uruk",
		Type:      eventType,
		FlowID:    message.FlowID,
		Timestamp: time.Now().String(),
		PusherID:  message.Pusher,
		Project:   message.Project,
		Details:   details,
	}
	u.SClient.Add(message.Stream, event.ConvertToEntry())
}

func worker(id int, u Uruk, messages <-chan saurontypes.UrukMessage) {
	for message := range messages {
		err := u.executeJob(message)
		if err != nil {
			u.logError("Error executing job\n"+message.String(), err)
		}
	}
}

// Start method should be called to start of uruk
// it takes queue name as parameter to determine which
// queue to listen to
func (u Uruk) Start(qName string) {
	u.logStart(qName)
	jobs := make(chan saurontypes.UrukMessage, 10)

	for index := 0; index < u.NumOfWorkers; index++ {
		go worker(index, u, jobs)
	}

	for {
		var urukMessage saurontypes.UrukMessage
		msg, err := u.QClient.Dequeue(qName)
		if err != nil {
			continue
		}
		err = json.Unmarshal([]byte(msg), &urukMessage)
		if err != nil {
			u.logError("Unable to unmarshall\n---\n"+msg+"---\n", err)
		}
		jobs <- urukMessage
	}
}
