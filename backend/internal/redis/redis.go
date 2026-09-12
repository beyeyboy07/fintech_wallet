package redis

import (
	"context"
	"strconv"
	"time"

	"fintech-wallet/backend/internal/config"

	goredis "github.com/redis/go-redis/v9"
)

func NewClient(cfg config.Config) *goredis.Client {
	db, _ := strconv.Atoi(cfg.RedisDB)

	client := goredis.NewClient(&goredis.Options{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
		Password: cfg.RedisPassword,
		DB:       db,
	})

	return client
}

func Ping(client *goredis.Client) error {
	return client.Ping(context.Background()).Err()
}

func Set(
	client *goredis.Client,
	key string,
	value string,
	expiration time.Duration,
) error {
	return client.Set(
		context.Background(),
		key,
		value,
		expiration,
	).Err()
}

func Get(
	client *goredis.Client,
	key string,
) (string, error) {
	return client.Get(
		context.Background(),
		key,
	).Result()
}
