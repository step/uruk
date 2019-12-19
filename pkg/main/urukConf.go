package main

import (
	"flag"
	"github.com/step/uruk/pkg/uruk"
	"io"
	"log"
	"os"
)

var queueName string
var sourceMountPoint string
var containerSourcePath string
var containerDataPath string
var numOfWorkers int

func init() {
	flag.StringVar(&queueName, "queue", "test", "Job `queue` where Uruk messages are queued")
	flag.StringVar(&sourceMountPoint, "source-mount-point", "/source", "`location` where source is mounted on Uruk")
	flag.StringVar(&containerSourcePath, "container-source-path", "source", "`path` where source is copied to on the container that runs the job")
	flag.StringVar(&containerDataPath, "container-data-path", "data", "`path` where data is copied to on the container that runs the job")
	flag.IntVar(&numOfWorkers, "num-of-workers", 2, "`number` of workers that can run containers in parallel")
}

func getLogger(file *os.File) uruk.UrukLogger {
	multiWriter := io.MultiWriter(file, os.Stdout)

	actualLogger := log.New(multiWriter, "--> ", log.LstdFlags)
	return uruk.UrukLogger{Logger: actualLogger}
}
