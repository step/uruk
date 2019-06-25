package main

import (
	"flag"
	"log"

	"github.com/docker/docker/client"
	dc "github.com/docker/docker/client"
)

var dockerVersion string

func init() {
	flag.StringVar(&dockerVersion, "docker-version", "1.37", "sets the docker `version` to use to communicate with the docker daemon")
}

func getDockerClient() *client.Client {
	dClient, err := dc.NewClientWithOpts(dc.WithVersion(dockerVersion))
	if err != nil {
		log.Panic("Unable to create Docker Client")
	}
	return dClient
}
