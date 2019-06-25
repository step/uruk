package main

import (
	"log"

	"github.com/docker/docker/client"
	dc "github.com/docker/docker/client"
	rc "github.com/step/angmar/pkg/redisclient"
	u "github.com/step/uruk/pkg/uruk"
)

func NewDefaultClient(version string) *client.Client {
	dClient, err := dc.NewClientWithOpts(dc.WithVersion(version))
	if err != nil {
		log.Panic("Unable to create Docker Client")
	}
	return dClient
}

func main() {
	dClient := NewDefaultClient("1.37")
	qClient := rc.NewDefaultClient(rc.RedisConf{
		Address:  "localhost:6379",
		Password: "", // no password set
		Db:       2,  // use default DB
	})

	uruk := u.Uruk{qClient, dClient}

	uruk.Start("foo")
}
