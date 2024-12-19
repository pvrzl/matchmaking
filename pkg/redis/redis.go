package redis

import (
	"context"

	redis "github.com/go-redis/redis/v8"
)

type Config struct {
	Address  string
	Password string
}

func NewRedis(params Config) (*redis.Client, error) {
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:     params.Address,
		Password: params.Password,
	})

	err := client.Ping(ctx).Err()
	return client, err
}
