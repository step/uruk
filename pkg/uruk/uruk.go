package uruk

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

	q "github.com/step/angmar/pkg/queueclient"
)

type Uruk struct {
	QClient q.QueueClient
	DClient *client.Client
}

type UrukMessage struct {
	ImageName    string
	RepoLocation string
}

func (u Uruk) CreateContainer(message UrukMessage) (container.ContainerCreateCreatedBody, error) {
	name := message.ImageName
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	return u.DClient.ContainerCreate(ctx, &container.Config{
		Image: name,
		Env:   []string{},
	}, nil, nil, "")
}

func (u Uruk) executeJob(urukMessage UrukMessage) {
	resp, err := u.CreateContainer(urukMessage)
	fmt.Println(resp, err)
	if err != nil {
		fmt.Println("error", err)
	}

	// CopyToContainer()
	// StartContainer()
	// WaitForContainerToStop()
	// CopyFromContainer()
	// RemoveContainer()
	// InformDB()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := u.DClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
	okCh, errCh := u.DClient.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case <-okCh:
		fmt.Println("everything ok. container done")
	case <-errCh:
		fmt.Println("some weird error")
	case <-ctx.Done():
		u.DClient.ContainerKill(context.Background(), resp.ID, "SIGTERM")
	}
}

func (o Uruk) Start(qName string) {
	for {
		var urukMessage UrukMessage
		msg, err := o.QClient.Dequeue(qName)
		if err != nil {
			fmt.Println("Unable to dequeue")
			continue
		}
		err = json.Unmarshal([]byte(msg), &urukMessage)
		if err != nil {
			fmt.Println("Unable to unmarshal", msg)
		}
		o.executeJob(urukMessage)
	}
}
