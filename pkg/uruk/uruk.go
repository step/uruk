package uruk

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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

func getNewFilename(fileName, src string) string {
	return strings.Replace(fileName, src, "", -1)
}

func tarContents(src string, buffer io.Writer) {
	tarWriter := tar.NewWriter(buffer)
	defer tarWriter.Close()

	filepath.Walk(src, func(fileName string, fi os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return fmt.Errorf("Unable to fetch header for %s", fi.Name())
		}

		newFilename := getNewFilename(fileName, src)
		newFilename = filepath.Join("source", newFilename)
		header.Name = newFilename
		if err := tarWriter.WriteHeader(header); err != nil {
			return fmt.Errorf("Unable to write header for %s", fi.Name())
		}

		if !fi.Mode().IsRegular() {
			return nil
		}

		file, err := os.Open(fileName)
		if err != nil {
			return fmt.Errorf("Unable to open %s", fileName)
		}

		if _, err := io.Copy(tarWriter, file); err != nil {
			return fmt.Errorf("Unable to tar %s", fileName)
		}

		fmt.Println(" +", header.Name)
		if err := file.Close(); err != nil {
			return fmt.Errorf("Unable to close %s", fileName)
		}

		return nil
	})
}

func (u Uruk) CopyToContainer(containerId, repoLocation string) error {
	var buffer bytes.Buffer
	tarContents(repoLocation, &buffer)
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

func (u Uruk) executeJob(urukMessage UrukMessage) {
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
