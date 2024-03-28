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
		NewRedisClient,
	),
)

// FxRedisClientParam allows injection of the required dependencies in [NewRedisClient].
type FxRedisClientParam struct {
	fx.In
	LifeCycle fx.Lifecycle
	Config    *config.Config
}

// NewRedisClient returns a [redis.Client].
func NewRedisClient(p FxRedisClientParam) (*redis.Client, *redismock.ClientMock) {
	if p.Config.IsTestEnv() {
		return createMockClient()
	} else {
		client := createClient(p)

		return client, nil
	}
}

func createClient(p FxRedisClientParam) *redis.Client {
	opt, err := redis.ParseURL(p.Config.GetString("modules.redis.dsn"))
	if err != nil {
		panic(err)
	}

	return redis.NewClient(opt)
}

func createMockClient() (*redis.Client, *redismock.ClientMock) {
	client, clientMock := redismock.NewClientMock()

	return client, &clientMock
}
