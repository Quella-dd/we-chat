package database

import "github.com/go-redis/redis"

var option = &redis.Options{
	Addr:     "192.168.99.101:6379",
	Password: "123456",
	DB:       0,
}

func NewRedisClient() (*redis.Client, error) {
	client := redis.NewClient(option)

	if err := client.Ping().Err(); err != nil {
		return nil, err
	}
	return client, nil
}
