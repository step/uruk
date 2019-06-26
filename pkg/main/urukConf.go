package main

import (
	"flag"
)

var queueName string
var sourceMountPoint string
var containerSourcePath string
var numOfWorkers int

func init() {
	flag.StringVar(&queueName, "queue", "test", "Job `queue` where Uruk messages are queued")
	flag.StringVar(&sourceMountPoint, "source-mount-point", "/source", "`location` where source is mounted on Uruk")
	flag.StringVar(&containerSourcePath, "container-source-path", "source", "`path` (without leading /) where source is copied to on the container that runs the job")
	flag.IntVar(&numOfWorkers, "num-of-workers", 2, "`number` of workers that can run containers in parallel")
}
