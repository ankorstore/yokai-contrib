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

## Testing

This module provides a simple and consistent API for creating test containers. See the [test examples](module_test.go) for usage patterns.

The module automatically applies sensible defaults:
- If no `WaitingFor` strategy is provided and a `Port` is specified, it will wait for the port to be listening
- Exposed ports are automatically configured based on the `Port` and `ExposedPorts` fields
- Container names are used for error reporting
- Containers are started automatically and ready to use when the function returns
