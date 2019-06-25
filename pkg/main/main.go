package main

import (
	"flag"

	u "github.com/step/uruk/pkg/uruk"
)

func main() {
	flag.Parse()
	dClient := getDockerClient()
	qClient := getRedisClient()

	uruk := u.Uruk{qClient, dClient}

	uruk.Start(queueName)
}
