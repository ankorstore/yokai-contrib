package fxtestcontainer_test

import (
	"testing"

	"github.com/ankorstore/yokai-contrib/fxtestcontainer"
	"github.com/ankorstore/yokai/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultContainerConfigFactory_Create_Success(t *testing.T) {
	t.Setenv("APP_ENV", "test")

	cfg, err := config.NewDefaultConfigFactory().Create(config.WithFilePaths("./testdata/config"))
	require.NoError(t, err)

	factory := fxtestcontainer.NewDefaultContainerConfigFactory(cfg)

	containerConfig, err := factory.Create("redis")
	require.NoError(t, err)

	assert.Equal(t, "test-redis", containerConfig.Name)
	assert.Equal(t, "redis:alpine", containerConfig.Image)
	assert.Equal(t, "6379/tcp", containerConfig.Port)
}

func TestDefaultContainerConfigFactory_Create_NotFound(t *testing.T) {
	t.Setenv("APP_ENV", "test")

	cfg, err := config.NewDefaultConfigFactory().Create(config.WithFilePaths("./testdata/config"))
	require.NoError(t, err)

	factory := fxtestcontainer.NewDefaultContainerConfigFactory(cfg)

	containerConfig, err := factory.Create("nonexistent")
	assert.Error(t, err)
	assert.Nil(t, containerConfig)
	assert.ErrorIs(t, err, fxtestcontainer.ErrContainerConfigNotFound)
}
