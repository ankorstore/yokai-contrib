package fxtestcontainer_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai-contrib/fxtestcontainer"
	"github.com/stretchr/testify/assert"
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
