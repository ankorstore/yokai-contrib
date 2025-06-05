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

func TestDefaultContainerConfigFactory_Create_DefaultName(t *testing.T) {
	t.Setenv("APP_ENV", "test")

	cfg, err := config.NewDefaultConfigFactory().Create(config.WithFilePaths("./testdata/config"))
	require.NoError(t, err)

	factory := fxtestcontainer.NewDefaultContainerConfigFactory(cfg)

	containerConfig, err := factory.Create("noname")
	require.NoError(t, err)

	// Should use configKey as default name when name is not provided
	assert.Equal(t, "noname", containerConfig.Name)
	assert.Equal(t, "nginx:alpine", containerConfig.Image)
	assert.Equal(t, "80/tcp", containerConfig.Port)
}

func TestDefaultContainerConfigFactory_Create_CompleteConfig(t *testing.T) {
	t.Setenv("APP_ENV", "test")

	cfg, err := config.NewDefaultConfigFactory().Create(config.WithFilePaths("./testdata/config"))
	require.NoError(t, err)

	factory := fxtestcontainer.NewDefaultContainerConfigFactory(cfg)

	containerConfig, err := factory.Create("postgres")
	require.NoError(t, err)

	// Test all fields are properly populated
	assert.Equal(t, "test-postgres", containerConfig.Name)
	assert.Equal(t, "postgres:13", containerConfig.Image)
	assert.Equal(t, "5432/tcp", containerConfig.Port)

	// Test exposed ports
	expectedExposedPorts := []string{"5432/tcp"}
	assert.Equal(t, expectedExposedPorts, containerConfig.ExposedPorts)

	// Test environment variables (should be uppercase)
	expectedEnv := map[string]string{
		"POSTGRES_DB":       "testdb",
		"POSTGRES_USER":     "testuser",
		"POSTGRES_PASSWORD": "testpass",
	}
	assert.Equal(t, expectedEnv, containerConfig.Environment)

	// Test command
	expectedCmd := []string{"postgres", "-c", "log_statement=all"}
	assert.Equal(t, expectedCmd, containerConfig.Cmd)
}

func TestNewDefaultContainerConfigFactory(t *testing.T) {
	t.Setenv("APP_ENV", "test")

	cfg, err := config.NewDefaultConfigFactory().Create(config.WithFilePaths("./testdata/config"))
	require.NoError(t, err)

	factory := fxtestcontainer.NewDefaultContainerConfigFactory(cfg)

	assert.NotNil(t, factory)
	assert.IsType(t, &fxtestcontainer.DefaultContainerConfigFactory{}, factory)
}
