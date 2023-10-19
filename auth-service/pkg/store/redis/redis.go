package redis

import "github.com/redis/go-redis/v9"

func New(address string) *redis.Client {
	cache := redis.NewClient(&redis.Options{
		Addr: address,
	})

	return cache
}
