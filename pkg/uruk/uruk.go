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
	QClient          q.QueueClient
	DClient          *client.Client
	Tarable          tarutils.Tarable
	SourceMountPoint string
	NumOfWorkers     int
}

func (u Uruk) executeJob(urukMessage saurontypes.UrukMessage) {
	resp, err := u.createContainer(urukMessage)
	fmt.Println(resp, err)
	if err != nil {
		fmt.Println("error", err)
	}

	err = u.copyToContainer(resp.ID, urukMessage.RepoLocation)
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err = u.startContainer(ctx, resp.ID)
	if err != nil {
		fmt.Println(err)
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
}

func worker(id int, u Uruk, messages <-chan saurontypes.UrukMessage) {
	for message := range messages {
		u.executeJob(message)
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
