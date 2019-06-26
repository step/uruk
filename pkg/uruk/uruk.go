package uruk

import (
	"bytes"
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
	QClient q.QueueClient
	DClient *client.Client
	Tarable tarutils.Tarable
}

func (u Uruk) CreateContainer(message saurontypes.UrukMessage) (container.ContainerCreateCreatedBody, error) {
	name := message.ImageName
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	return u.DClient.ContainerCreate(ctx, &container.Config{
		Image: name,
		Env:   []string{},
	}, nil, nil, "")
}

func (u Uruk) CopyToContainer(containerId, repoLocation string) error {
	var buffer bytes.Buffer
	tarutils.Tar(repoLocation, &buffer, u.Tarable)
	ctx := context.Background()
	return u.DClient.CopyToContainer(ctx, containerId, "/", &buffer, types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
	})
}

func (u Uruk) CopyFromContainer(containerId, src string) io.ReadCloser {
	readCloser, _, err := u.DClient.CopyFromContainer(context.Background(), containerId, src)
	fmt.Println("Error copying", err)
	return readCloser
}

func (u Uruk) StartContainer(ctx context.Context, containerId string) error {
	return u.DClient.ContainerStart(ctx, containerId, types.ContainerStartOptions{})
}

func (u Uruk) executeJob(urukMessage saurontypes.UrukMessage) {
	resp, err := u.CreateContainer(urukMessage)
	fmt.Println(resp, err)
	if err != nil {
		fmt.Println("error", err)
	}

	err = u.CopyToContainer(resp.ID, urukMessage.RepoLocation)
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err = u.StartContainer(ctx, resp.ID)
	if err != nil {
		fmt.Println(err)
	}

	okCh, errCh := u.DClient.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case status := <-okCh:
		fmt.Println("everything ok. container done", status.StatusCode)
		readCloser := u.CopyFromContainer(resp.ID, "/results/result.json")
		defer readCloser.Close()
		io.Copy(os.Stdout, readCloser)
	case <-errCh:
		fmt.Println("some weird error")
	case <-ctx.Done():
		u.DClient.ContainerKill(context.Background(), resp.ID, "SIGTERM")
	}
	u.DClient.ContainerRemove(context.Background(), resp.ID, types.ContainerRemoveOptions{Force: true})
	// TODO: InformDB()
}

func (o Uruk) Start(qName string) {
	for {
		var urukMessage saurontypes.UrukMessage
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
