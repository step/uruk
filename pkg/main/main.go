package main

import (
	"flag"

	"github.com/step/uruk/pkg/tarutils"

	u "github.com/step/uruk/pkg/uruk"
)

func main() {
	flag.Parse()
	dClient := getDockerClient()
	qClient := getRedisClient()
	tarable := tarutils.NewDefaultTarable(containerSourcePath)

	uruk := u.Uruk{
		QClient:             qClient,
		DClient:             dClient,
		Tarable:             tarable,
		SourceMountPoint:    sourceMountPoint,
		ContainerSourcePath: containerSourcePath,
		NumOfWorkers:        numOfWorkers,
	}

	uruk.Start(queueName)
}
