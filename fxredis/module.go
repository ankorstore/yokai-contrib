package fxredis

import (
	"github.com/ankorstore/yokai/config"
	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

// ModuleName is the module name.
const ModuleName = "redis"

// FxRedis is the [Fx] redis module.
//
// [Fx]: https://github.com/uber-go/fx
var FxRedisModule = fx.Module(
	ModuleName,
	fx.Provide(
		NewFxRedisClient,
	),
)

// FxRedisClientParam allows injection of the required dependencies in [NewRedisClient].
type FxRedisClientParam struct {
	fx.In
	LifeCycle fx.Lifecycle
	Config    *config.Config
}

// NewFxRedisClient returns a [redis.Client].
func NewFxRedisClient(p FxRedisClientParam) (*redis.Client, redismock.ClientMock, error) {
	if p.Config.IsTestEnv() {
		return createMockClient()
	} else {
		client, err := createClient(p)

		return client, nil, err
	}
}

func createClient(p FxRedisClientParam) (*redis.Client, error) {
	opt, err := redis.ParseURL(p.Config.GetString("modules.redis.dsn"))
	if err != nil {
		return nil, err
	}

	return redis.NewClient(opt), nil
}

func createMockClient() (*redis.Client, redismock.ClientMock, error) {
	client, clientMock := redismock.NewClientMock()

	return client, clientMock, nil
}
