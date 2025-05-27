package fxtestcontainer

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/fx"
)

// ModuleName is the module name.
const ModuleName = "testcontainer"

const (
	// DefaultStartupTimeout is the default timeout for container startup.
	DefaultStartupTimeout = 120 * time.Second
	// DefaultPollInterval is the default poll interval for waiting strategies.
	DefaultPollInterval = 1 * time.Second
)

// ErrContainerImageRequired is returned when no container image is provided.
var ErrContainerImageRequired = errors.New("container image is required")

// FxTestContainerModule provides container creation functionality for Fx.
var FxTestContainerModule = fx.Module(ModuleName)

// CreateGenericContainer creates a testcontainer from the provided configuration.
func CreateGenericContainer(ctx context.Context, config *ContainerConfig) (testcontainers.Container, error) {
	if config.Image == "" {
		return nil, ErrContainerImageRequired
	}

	// Prepare exposed ports
	exposedPorts := config.ExposedPorts
	if config.Port != "" {
		exposedPorts = append(exposedPorts, config.Port)
	}

	// Set default wait strategy if none provided
	waitStrategy := config.WaitingFor
	if waitStrategy == nil && config.Port != "" {
		port := nat.Port(config.Port)
		waitStrategy = wait.ForListeningPort(port).
			WithStartupTimeout(DefaultStartupTimeout).
			WithPollInterval(DefaultPollInterval)
	}

	// Create container request
	req := testcontainers.ContainerRequest{
		Image:        config.Image,
		ExposedPorts: exposedPorts,
		Env:          config.Environment,
		WaitingFor:   waitStrategy,
		Cmd:          config.Cmd,
		Mounts:       config.Mounts,
	}

	// Start container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create container %s: %w", config.Name, err)
	}

	return container, nil
}

// ContainerConfig provides configuration for creating test containers.
type ContainerConfig struct {
	// Name is a unique identifier for the container
	Name string
	// Image specifies the Docker image to use
	Image string
	// Port specifies the main port to expose (convenience field)
	Port string
	// ExposedPorts lists the ports to expose
	ExposedPorts []string
	// Environment provides environment variables
	Environment map[string]string
	// WaitingFor specifies the wait strategy
	WaitingFor wait.Strategy
	// Cmd specifies the command to run in the container
	Cmd []string
	// Mounts specifies volume mounts
	Mounts testcontainers.ContainerMounts
}
