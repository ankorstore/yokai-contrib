package fxredis_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai-contrib/fxredis"
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestFxRedisModule(t *testing.T) {
	app := fxtest.New(
		t,
		fx.NopLogger,
		fxredis.FxRedisModule,
		fxconfig.FxConfigModule,
	)

	app.RequireStart()
	assert.NoError(t, app.Err(), "failed to start the Fx application")

	app.RequireStop()
	assert.NoError(t, app.Err(), "failed to stop the Fx application")
}

func TestFxRedisClient(t *testing.T) {
	t.Setenv("APP_ENV", config.AppEnvDev)
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("TOKEN", "my-token")
	t.Setenv("REDIS_USER", "my-user")
	t.Setenv("REDIS_PASSWORD", "my-password")
	t.Setenv("REDIS_HOST", "my-host")
	t.Setenv("REDIS_PORT", "1234")
	t.Setenv("REDIS_DB", "0")

	var conf *config.Config
	var client *redis.Client
	var mockClient redismock.ClientMock

	app := fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxredis.FxRedisModule,
		fx.Populate(&conf, &client, &mockClient),
	)

	app.RequireStart()
	assert.NoError(t, app.Err(), "failed to create redis.Client")
	assert.NotNil(t, client)
	assert.Nil(t, mockClient)
	assert.Equal(t, "my-user", client.Options().Username)
	assert.Equal(t, "my-password", client.Options().Password)
	assert.Equal(t, "my-host:1234", client.Options().Addr)
	assert.Equal(t, 0, client.Options().DB)

	app.RequireStop()
	assert.NoError(t, app.Err(), "failed to close redis.Client")
}

func TestFxRedisTestClient(t *testing.T) {
	t.Setenv("APP_ENV", config.AppEnvTest)
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("TOKEN", "my-token")

	var conf *config.Config
	var client *redis.Client
	var mockClient redismock.ClientMock

	app := fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxredis.FxRedisModule,
		fx.Populate(&conf, &client, &mockClient),
	)

	app.RequireStart()
	assert.NoError(t, app.Err(), "failed to create test redis.Client")
	assert.NotNil(t, client)
	assert.NotNil(t, mockClient)

	mockClient.ExpectSet("key", "value", 0).SetVal("OK")
	mockClient.ExpectGet("key2").SetVal("value2")

	client.Set(context.Background(), "key", "value", 0)
	assert.Equal(t, "value2", client.Get(context.Background(), "key2").Val())

	if err := mockClient.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}

	app.RequireStop()
	assert.NoError(t, app.Err(), "failed to close test redis.Client")
}
