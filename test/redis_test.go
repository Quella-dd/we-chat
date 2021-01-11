package test

import (
	"github.com/go-redis/redis"
	"testing"
	"fmt"
)

func TestRedis(t *testing.T) {
	fmt.Println("golang连接redis")

	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		Password: "123456",
		DB: 0,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		t.Log("error has been happend")
	} else {
		t.Log("ping success:", pong)
	}
}