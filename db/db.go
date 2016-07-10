package db

import (
	"github.com/fzzy/radix/redis"
	"log"
	"os"
	"sync"
)

var instance *redis.Client
var once sync.Once

func GetInstance() *redis.Client {
	redisAddr := os.Getenv("REDIS_ADDR")

	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	once.Do(func() {
		client, err := redis.Dial("tcp", redisAddr)

		if err != nil {
			log.Fatal(err)
		}

		instance = client
	})

	return instance
}
