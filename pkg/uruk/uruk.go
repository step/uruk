package uruk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/step/saurontypes"
	"github.com/step/uruk/pkg/tarutils"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

	q "github.com/step/angmar/pkg/queueclient"
)

type Uruk struct {
	QClient             q.QueueClient
	DClient             *client.Client
	Tarable             tarutils.Tarable
	SourceMountPoint    string
	ContainerSourcePath string
	NumOfWorkers        int
}

func (u Uruk) executeJob(urukMessage saurontypes.UrukMessage) error {
	resp, err := u.createContainer(urukMessage)
	// log response and error
	if err != nil {
		return ContainerCreationError{urukMessage, err}
	}

	err = u.copyToContainer(resp.ID, urukMessage.RepoLocation)
	if err != nil {
		return CopyToContainerError{urukMessage, u.SourceMountPoint, u.ContainerSourcePath, err}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err = u.startContainer(ctx, resp.ID)
	if err != nil {
		return StartContainerError{resp, urukMessage, err}
	}

	okCh, errCh := u.DClient.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case status := <-okCh:
		fmt.Println("everything ok. container done", status.StatusCode)
		readCloser := u.copyFromContainer(resp.ID, "/results/result.json")
		defer readCloser.Close()
		io.Copy(os.Stdout, readCloser)
	case <-errCh:
		fmt.Println("some weird error")
	case <-ctx.Done():
		u.DClient.ContainerKill(context.Background(), resp.ID, "SIGTERM")
	}

	u.DClient.ContainerRemove(context.Background(), resp.ID, types.ContainerRemoveOptions{Force: true})
	// TODO: InformDB()
	return nil
}

func worker(id int, u Uruk, messages <-chan saurontypes.UrukMessage) {
	for message := range messages {
		err := u.executeJob(message)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func (u Uruk) Start(qName string) {
	jobs := make(chan saurontypes.UrukMessage, 10)

	for index := 0; index < u.NumOfWorkers; index++ {
		go worker(index, u, jobs)
	}

	for {
		var urukMessage saurontypes.UrukMessage
		msg, err := u.QClient.Dequeue(qName)
		if err != nil {
			fmt.Println("Unable to dequeue")
			continue
		}
		err = json.Unmarshal([]byte(msg), &urukMessage)
		if err != nil {
			fmt.Println("Unable to unmarshal", msg)
		}
		jobs <- urukMessage
	}
}
