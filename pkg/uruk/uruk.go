package uruk

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/step/saurontypes"
	"github.com/step/uruk/pkg/tarutils"

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
		if err := u.copyFromContainer(resp.ID, "/results/result.json"); err != nil {
			return CopyFromContainerError{urukMessage, resp.ID, err}
		}

	case err := <-errCh:
		return err
	case <-ctx.Done():
		fmt.Println("timed out...")
		return u.killContainer(context.Background(), resp.ID)
	}

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
