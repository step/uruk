package main

import (
	"flag"
)

var queueName string

func init() {
	flag.StringVar(&queueName, "queue", "test", "Job `queue` where Uruk messages are queued")
}
