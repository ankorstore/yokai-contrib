package fxredis_test

import (
	"testing"

	"github.com/ankorstore/yokai-contrib/fxredis"
	"github.com/ankorstore/yokai/config"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestDefaultRedisClientFactory(t *testing.T) {
	createConfig := func() (*config.Config, error) {
		return config.NewDefaultConfigFactory().Create(
			config.WithFilePaths("./testdata/config"),
		)
	}

	t.Run("create success", func(t *testing.T) {
		cfg, err := createConfig()
		assert.NoError(t, err)

		factory := fxredis.NewDefaultRedisClientFactory(cfg)

		client, err := factory.Create()
		assert.NoError(t, err)
		assert.IsType(t, &redis.Client{}, client)
	})

	t.Run("create error with invalid dsn", func(t *testing.T) {
		t.Setenv("MODULES_REDIS_DSN", "invalid")

		cfg, err := createConfig()
		assert.NoError(t, err)

		factory := fxredis.NewDefaultRedisClientFactory(cfg)

		_, err = factory.Create()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid URL scheme")
	})
}
