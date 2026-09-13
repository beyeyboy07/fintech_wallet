package service

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type RedisService struct {
	client *goredis.Client
}

func NewRedisService(
	client *goredis.Client,
) *RedisService {
	return &RedisService{
		client: client,
	}
}

func (s *RedisService) Set(
	key string,
	value string,
	expiration time.Duration,
) error {
	return s.client.Set(
		context.Background(),
		key,
		value,
		expiration,
	).Err()
}

func (s *RedisService) Get(
	key string,
) (string, error) {
	return s.client.Get(
		context.Background(),
		key,
	).Result()
}

func (s *RedisService) Delete(key string) error {
	return s.client.Del(context.Background(), key).Err()
}
