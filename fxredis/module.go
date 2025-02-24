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
		fx.Annotate(NewDefaultRedisClientFactory, fx.As(new(RedisClientFactory))),
		NewFxRedisClient,
	),
)

// FxRedisClientParam allows injection of the required dependencies in [NewRedisClient].
type FxRedisClientParam struct {
	fx.In
	LifeCycle fx.Lifecycle
	Config    *config.Config
	Factory   RedisClientFactory
}

// NewFxRedisClient returns a [redis.Client] and a [redismock.ClientMock] in test mode.
func NewFxRedisClient(p FxRedisClientParam) (*redis.Client, redismock.ClientMock, error) {
	if p.Config.IsTestEnv() {
		client, clientMock := redismock.NewClientMock()

		return client, clientMock, nil
	} else {
		client, err := p.Factory.Create()

		return client, nil, err
	}
}
