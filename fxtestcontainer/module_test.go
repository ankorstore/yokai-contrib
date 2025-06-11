package fxtestcontainer_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai-contrib/fxtestcontainer"
	"github.com/ankorstore/yokai/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestFxTestContainerModule(t *testing.T) {
	app := fxtest.New(
		t,
		fx.NopLogger,
		fxtestcontainer.FxTestContainerModule,
	)

	app.RequireStart()
	assert.NoError(t, app.Err(), "failed to start the Fx application")

	app.RequireStop()
	assert.NoError(t, app.Err(), "failed to stop the Fx application")
}

func TestCreateGenericContainer_DirectUsage(t *testing.T) {
	ctx := context.Background()

	config := &fxtestcontainer.ContainerConfig{
		Name:  "direct-test",
		Image: "nginx:alpine",
		Port:  "80/tcp",
	}

	// Test direct usage
	container, err := fxtestcontainer.CreateGenericContainer(ctx, config)
	if err != nil {
		t.Skipf("Skipping test due to container creation failure (Docker may not be available): %v", err)
	}

	defer func() {
		if terminateErr := container.Terminate(ctx); terminateErr != nil {
			t.Logf("Failed to terminate container: %v", terminateErr)
		}
	}()

	assert.NotNil(t, container)
}

func TestCreateGenericContainer_EmptyImage(t *testing.T) {
	ctx := context.Background()

	config := &fxtestcontainer.ContainerConfig{
		Name:  "empty-image-test",
		Image: "", // Empty image should cause error
		Port:  "80/tcp",
	}

	// Test that empty image returns error
	container, err := fxtestcontainer.CreateGenericContainer(ctx, config)
	assert.Error(t, err)
	assert.Nil(t, container)
	assert.ErrorIs(t, err, fxtestcontainer.ErrContainerImageRequired)
}

func TestCreateGenericContainer_InvalidImage(t *testing.T) {
	ctx := context.Background()

	config := &fxtestcontainer.ContainerConfig{
		Name:  "invalid-image-test",
		Image: "nonexistent/invalid-image:impossible-tag",
		Port:  "80/tcp",
	}

	// Test that invalid image returns error
	container, err := fxtestcontainer.CreateGenericContainer(ctx, config)
	if err != nil {
		// This should fail due to invalid image, which covers the error path in container creation
		assert.Error(t, err)
		assert.Nil(t, container)
		assert.Contains(t, err.Error(), "failed to create container")
	} else {
		// If it somehow succeeds (maybe Docker pulls the image), clean up
		defer func() {
			if terminateErr := container.Terminate(ctx); terminateErr != nil {
				t.Logf("Failed to terminate container: %v", terminateErr)
			}
		}()
		t.Skip("Skipping test as invalid image was unexpectedly available")
	}
}

func TestCreateGenericContainerFromConfig(t *testing.T) {
	t.Setenv("APP_ENV", "test")

	ctx := context.Background()
	cfg, err := config.NewDefaultConfigFactory().Create(config.WithFilePaths("./testdata/config"))
	require.NoError(t, err)

	factory := fxtestcontainer.NewDefaultContainerConfigFactory(cfg)

	t.Run("test redis container creation from config", func(t *testing.T) {
		container, err := fxtestcontainer.CreateGenericContainerFromConfig(ctx, factory, "redis")
		if err != nil {
			t.Skipf("Skipping test due to container creation failure (Docker may not be available): %v", err)
		}

		defer func() {
			if terminateErr := container.Terminate(ctx); terminateErr != nil {
				t.Logf("Failed to terminate container: %v", terminateErr)
			}
		}()

		assert.NotNil(t, container)

		// Get the container endpoint
		endpoint, err := container.Endpoint(ctx, "")
		if err == nil {
			t.Logf("Redis available at: %s", endpoint)
		}
	})

	t.Run("test postgres container creation from config", func(t *testing.T) {
		container, err := fxtestcontainer.CreateGenericContainerFromConfig(ctx, factory, "postgres")
		if err != nil {
			t.Skipf("Skipping test due to container creation failure (Docker may not be available): %v", err)
		}

		defer func() {
			if terminateErr := container.Terminate(ctx); terminateErr != nil {
				t.Logf("Failed to terminate container: %v", terminateErr)
			}
		}()

		assert.NotNil(t, container)

		// Get the container endpoint
		endpoint, err := container.Endpoint(ctx, "")
		if err == nil {
			t.Logf("PostgreSQL available at: %s", endpoint)
		}
	})

	t.Run("test container creation from non-existent config", func(t *testing.T) {
		container, err := fxtestcontainer.CreateGenericContainerFromConfig(ctx, factory, "nonexistent")
		assert.Error(t, err)
		assert.Nil(t, container)
		assert.ErrorIs(t, err, fxtestcontainer.ErrContainerConfigNotFound)
	})
}

func TestFxTestContainerModule_WithFactory(t *testing.T) {
	t.Setenv("APP_ENV", "test")

	var factory fxtestcontainer.ContainerConfigFactory

	app := fxtest.New(
		t,
		fx.NopLogger,
		fx.Provide(
			func() (*config.Config, error) {
				return config.NewDefaultConfigFactory().Create(config.WithFilePaths("./testdata/config"))
			},
		),
		fxtestcontainer.FxTestContainerModule,
		fx.Populate(&factory),
	)

	app.RequireStart()
	assert.NoError(t, app.Err(), "failed to start the Fx application")

	// Test that factory is properly injected and working
	containerConfig, err := factory.Create("redis")
	assert.NoError(t, err)
	assert.Equal(t, "test-redis", containerConfig.Name)
	assert.Equal(t, "redis:alpine", containerConfig.Image)

	app.RequireStop()
	assert.NoError(t, app.Err(), "failed to stop the Fx application")
}

// Configuration tests

func TestDefaultContainerConfigFactory_Create(t *testing.T) {
	t.Setenv("APP_ENV", "test")

	cfg, err := config.NewDefaultConfigFactory().Create(config.WithFilePaths("./testdata/config"))
	require.NoError(t, err)

	factory := fxtestcontainer.NewDefaultContainerConfigFactory(cfg)

	t.Run("test redis config creation", func(t *testing.T) {
		containerConfig, err := factory.Create("redis")
		require.NoError(t, err)

		assert.Equal(t, "test-redis", containerConfig.Name)
		assert.Equal(t, "redis:alpine", containerConfig.Image)
		assert.Equal(t, "6379/tcp", containerConfig.Port)
		assert.Equal(t, "", containerConfig.Environment["REDIS_PASSWORD"])
		assert.Empty(t, containerConfig.ExposedPorts)
		assert.Empty(t, containerConfig.Cmd)
	})

	t.Run("test postgres config creation", func(t *testing.T) {
		containerConfig, err := factory.Create("postgres")
		require.NoError(t, err)

		assert.Equal(t, "test-postgres", containerConfig.Name)
		assert.Equal(t, "postgres:13", containerConfig.Image)
		assert.Equal(t, "5432/tcp", containerConfig.Port)
		assert.Equal(t, "testdb", containerConfig.Environment["POSTGRES_DB"])
		assert.Equal(t, "testuser", containerConfig.Environment["POSTGRES_USER"])
		assert.Equal(t, "testpass", containerConfig.Environment["POSTGRES_PASSWORD"])
		assert.Contains(t, containerConfig.ExposedPorts, "5432/tcp")
		assert.Contains(t, containerConfig.Cmd, "postgres")
		assert.Contains(t, containerConfig.Cmd, "-c")
		assert.Contains(t, containerConfig.Cmd, "log_statement=all")
	})

	t.Run("test non-existent config", func(t *testing.T) {
		containerConfig, err := factory.Create("nonexistent")
		assert.Error(t, err)
		assert.Nil(t, containerConfig)
		assert.ErrorIs(t, err, fxtestcontainer.ErrContainerConfigNotFound)
		assert.Contains(t, err.Error(), "nonexistent")
	})
}
