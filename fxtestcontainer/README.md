# Yokai Test Container Module

[![ci](https://github.com/ankorstore/yokai-contrib/actions/workflows/fxtestcontainer-ci.yml/badge.svg)](https://github.com/ankorstore/yokai-contrib/actions/workflows/fxtestcontainer-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai-contrib/fxtestcontainer)](https://goreportcard.com/report/github.com/ankorstore/yokai-contrib/fxtestcontainer)
[![codecov](https://codecov.io/gh/ankorstore/yokai-contrib/graph/badge.svg?token=ghUBlFsjhR&flag=fxtestcontainer)](https://app.codecov.io/gh/ankorstore/yokai-contrib/tree/main/fxtestcontainer)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai-contrib%2Ffxtestcontainer)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai-contrib/fxtestcontainer)](https://pkg.go.dev/github.com/ankorstore/yokai-contrib/fxtestcontainer)

> [Yokai](https://github.com/ankorstore/yokai) module for [Testcontainers](https://github.com/testcontainers/testcontainers-go) integration.

<!-- TOC -->
* [Overview](#overview)
* [Installation](#installation)
* [Usage](#usage)
  * [Basic Usage](#basic-usage)
  * [Advanced Configuration](#advanced-configuration)
  * [Configuration-Based Usage](#configuration-based-usage)
* [Testing](#testing)
<!-- TOC -->

## Overview

This module provides a simple and consistent API for creating test containers with [Testcontainers](https://github.com/testcontainers/testcontainers-go) in your [Yokai](https://github.com/ankorstore/yokai) applications.

The module offers the `CreateGenericContainer` function that handles container creation with sensible defaults and flexible configuration options.

## Installation

Install the module:

```shell
go get github.com/ankorstore/yokai-contrib/fxtestcontainer
```

Then activate it in your application bootstrapper:

```go
// internal/bootstrap.go
package internal

import (
	"github.com/ankorstore/yokai-contrib/fxtestcontainer"
	"github.com/ankorstore/yokai/fxcore"
)

var Bootstrapper = fxcore.NewBootstrapper().WithOptions(
	// load modules
	fxtestcontainer.FxTestContainerModule,
	// ...
)
```

## Usage

### Basic Usage

Use the `CreateGenericContainer` function to create test containers:

```go
package service_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai-contrib/fxtestcontainer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMyService_WithRedis(t *testing.T) {
	ctx := context.Background()

	// Define container configuration
	config := &fxtestcontainer.ContainerConfig{
		Name:  "test-redis",
		Image: "redis:alpine",
		Port:  "6379/tcp",
		Environment: map[string]string{
			"REDIS_PASSWORD": "",
		},
	}

	// Create test container
	container, err := fxtestcontainer.CreateGenericContainer(ctx, config)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	// Get container endpoint
	endpoint, err := container.Endpoint(ctx, "")
	require.NoError(t, err)

	// Use the container in your service
	service := NewMyService(endpoint)
	
	// Test your service methods
	err = service.Set(ctx, "key", "value")
	assert.NoError(t, err)
	
	value, err := service.Get(ctx, "key")
	assert.NoError(t, err)
	assert.Equal(t, "value", value)
}
```

### Advanced Configuration

The `ContainerConfig` struct provides flexible configuration options:

```go
type ContainerConfig struct {
	// Name is a unique identifier for the container
	Name string
	// Image specifies the Docker image to use
	Image string
	// Port specifies the main port to expose (convenience field)
	Port string
	// ExposedPorts lists additional ports to expose
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
```

Example with custom wait strategy and multiple configuration options:

```go
func TestWithPostgres(t *testing.T) {
	ctx := context.Background()

	config := &fxtestcontainer.ContainerConfig{
		Name:  "test-postgres",
		Image: "postgres:13",
		Port:  "5432/tcp",
		Environment: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections"),
	}

	container, err := fxtestcontainer.CreateGenericContainer(ctx, config)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	// Get connection details
	host, err := container.Host(ctx)
	require.NoError(t, err)

	port, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err)

	// Use in your tests
	dsn := fmt.Sprintf("postgres://testuser:testpass@%s:%s/testdb?sslmode=disable", 
		host, port.Port())
	
	// Connect to database and run tests...
}
```

### Configuration-Based Usage

The module also supports configuration-based container creation using YAML files.

First, create a `config.test.yaml` file in your test configuration directory:

```yaml
app:
  env: test

modules:
  testcontainer:
    containers:
      # Redis configuration
      redis:
        name: "test-redis"
        image: "redis:alpine"
        port: "6379/tcp"
        environment:
          REDIS_PASSWORD: ""
      
      # PostgreSQL configuration
      postgres:
        name: "test-postgres"
        image: "postgres:13"
        port: "5432/tcp"
        environment:
          POSTGRES_DB: "testdb"
          POSTGRES_USER: "testuser"
          POSTGRES_PASSWORD: "testpass"
        exposed_ports:
          - "5432/tcp"
        cmd:
          - "postgres"
          - "-c"
          - "log_statement=all"
      
      # MySQL configuration
      mysql:
        name: "test-mysql"
        image: "mysql:8.0"
        port: "3306/tcp"
        environment:
          MYSQL_ROOT_PASSWORD: "rootpass"
          MYSQL_DATABASE: "testdb"
          MYSQL_USER: "testuser"
          MYSQL_PASSWORD: "testpass"
      
      # Elasticsearch configuration
      elasticsearch:
        name: "test-elasticsearch"
        image: "elasticsearch:8.11.0"
        port: "9200/tcp"
        exposed_ports:
          - "9200/tcp"
          - "9300/tcp"
        environment:
          DISCOVERY_TYPE: "single-node"
          XPACK_SECURITY_ENABLED: "false"
          ES_JAVA_OPTS: "-Xms512m -Xmx512m"
```

> **Note:** The configuration above shows examples of different container types you can configure. You don't need to include all of them - just add the containers you actually need for your tests.

> **Important:** Viper's map-based access methods return lowercase keys, but this module automatically converts environment variable keys back to uppercase for Docker containers. Your tests can use standard uppercase keys (e.g., `POSTGRES_PASSWORD`) when accessing the `Environment` map.

Then use the configuration-based approach in your tests:

```go
package service_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai-contrib/fxtestcontainer"
	"github.com/ankorstore/yokai/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMyService_WithConfigBasedContainers(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	
	ctx := context.Background()
	
	// Load configuration
	cfg, err := config.NewDefaultConfigFactory().Create(config.WithFilePaths("./path/to/config"))
	require.NoError(t, err)
	
	// Create factory
	factory := fxtestcontainer.NewDefaultContainerConfigFactory(cfg)
	
	// Create Redis container from config
	redisContainer, err := fxtestcontainer.CreateGenericContainerFromConfig(ctx, factory, "redis")
	require.NoError(t, err)
	defer redisContainer.Terminate(ctx)
	
	// Create PostgreSQL container from config
	postgresContainer, err := fxtestcontainer.CreateGenericContainerFromConfig(ctx, factory, "postgres")
	require.NoError(t, err)
	defer postgresContainer.Terminate(ctx)
	
	// Get connection details
	redisEndpoint, err := redisContainer.Endpoint(ctx, "")
	require.NoError(t, err)
	
	postgresHost, err := postgresContainer.Host(ctx)
	require.NoError(t, err)
	postgresPort, err := postgresContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)
	
	// Initialize your service with container endpoints
	service := NewMyService(redisEndpoint, postgresHost, postgresPort.Port())
	
	// Run your tests
	err = service.Process(ctx, "test-data")
	assert.NoError(t, err)
}
```

You can also use the module with Fx dependency injection:

```go
func TestWithFxInjection(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	
	var factory fxtestcontainer.ContainerConfigFactory
	
	app := fxtest.New(
		t,
		fx.NopLogger,
		fx.Provide(
			func() (*config.Config, error) {
				return config.NewDefaultConfigFactory().Create(config.WithFilePaths("./path/to/config"))
			},
		),
		fxtestcontainer.FxTestContainerModule,
		fx.Populate(&factory),
	)
	
	app.RequireStart()
	defer app.RequireStop()
	
	// Use the injected factory
	container, err := fxtestcontainer.CreateGenericContainerFromConfig(context.Background(), factory, "redis")
	require.NoError(t, err)
	defer container.Terminate(context.Background())
	
	// Test your service...
}
```

## Testing

This module provides a simple and consistent API for creating test containers. See the [test examples](module_test.go) for usage patterns.

The module automatically applies sensible defaults:
- If no `WaitingFor` strategy is provided and a `Port` is specified, it will wait for the port to be listening
- Exposed ports are automatically configured based on the `Port` and `ExposedPorts` fields
- Container names are used for error reporting
- Containers are started automatically and ready to use when the function returns
