package fxredis

import (
	"github.com/ankorstore/yokai/config"
	"github.com/redis/go-redis/v9"
)

var _ RedisClientFactory = (*DefaultRedisClientFactory)(nil)

type RedisClientFactory interface {
	Create() (*redis.Client, error)
}

type DefaultRedisClientFactory struct {
	config *config.Config
}

func NewDefaultRedisClientFactory(config *config.Config) *DefaultRedisClientFactory {
	return &DefaultRedisClientFactory{
		config: config,
	}
}

func (f *DefaultRedisClientFactory) Create() (*redis.Client, error) {
	opt, err := redis.ParseURL(f.config.GetString("modules.redis.dsn"))
	if err != nil {
		return nil, err
	}

	return redis.NewClient(opt), nil
}
