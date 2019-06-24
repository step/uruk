package main

import (
	"log"
	"time"

	"github.com/go-redis/redis"

	"github.com/docker/docker/client"
	dc "github.com/docker/docker/client"
	u "github.com/step/uruk/pkg/uruk"
)

type RedisClient struct {
	actualClient *redis.Client
}

func (r RedisClient) Enqueue(name, value string) error {
	r.actualClient.LPush(name, value)
	return nil
}

func (r RedisClient) Dequeue(name string) (string, error) {
	resp := r.actualClient.BRPop(time.Second*3, name)
	values, err := resp.Result()
	if err != nil {
		return "", err
	}
	return values[1], err
}

func (r RedisClient) SwitchQueue(src, dest string) (string, error) {
	return "", nil
}

func (r RedisClient) String() string {
	return r.actualClient.String()
}

func NewDefaultClient(version string) *client.Client {
	dClient, err := dc.NewClientWithOpts(dc.WithVersion("1.37"))
	if err != nil {
		log.Panic("Unable to create Docker Client")
	}
	return dClient
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       2,  // use default DB
	})
	dClient := NewDefaultClient("1.37")
	qClient := RedisClient{client}
	uruk := u.Uruk{qClient, dClient}

	// qClient.Enqueue("foo", "{\"imageName\":\"orc_sample\", \"repoLocation\":\"luciferankon/master/luciferankon-gauge-js-948694c\"}")
	uruk.Start("foo")
}
