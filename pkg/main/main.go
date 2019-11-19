package main

import (
	"flag"
	"log"
	"os"

	"github.com/step/uruk/pkg/tarutils"

	u "github.com/step/uruk/pkg/uruk"
)

func main() {
	flag.Parse()
	dClient := getDockerClient()
	qClient := getRedisClient()
	tarable := tarutils.NewDefaultTarable("")

	logger := log.New(os.Stdout, "uruk ", log.LstdFlags)
	uruk := u.Uruk{
		QClient:             qClient,
		DClient:             dClient,
		SClient:             qClient,
		Tarable:             tarable,
		SourceMountPoint:    sourceMountPoint,
		ContainerSourcePath: containerSourcePath,
		ContainerDataPath:   containerDataPath,
		NumOfWorkers:        numOfWorkers,
		Logger:              logger,
	}

	uruk.Start(queueName)
}
