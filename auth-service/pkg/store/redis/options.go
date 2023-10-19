package redis

import (
	"context"
	"net"

	"github.com/redis/go-redis/v9"
)

type Option func(opt *redis.Options)

func WithDialer(f func(ctx context.Context, network, addr string) (net.Conn, error)) Option {
	return func(Opt *redis.Options) {
		Opt.Dialer = f
	}
}

func WithUsername(username string) Option {
	return func(opt *redis.Options) {
		opt.Username = username
	}
}

func WithPassword(password string) Option {
	return func(opt *redis.Options) {
		opt.Password = password
	}
}

func WithDatabase(databaseNumber int) Option {
	return func(opt *redis.Options) {
		opt.DB = databaseNumber
	}
}
