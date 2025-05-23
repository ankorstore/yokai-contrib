# Fx Elasticsearch Module

[![ci](https://github.com/ankorstore/yokai-contrib/actions/workflows/fxelasticsearch-ci.yml/badge.svg)](https://github.com/ankorstore/yokai-contrib/actions/workflows/fxelasticsearch-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai-contrib/fxelasticsearch)](https://goreportcard.com/report/github.com/ankorstore/yokai-contrib/fxelasticsearch)
[![codecov](https://codecov.io/gh/ankorstore/yokai-contrib/graph/badge.svg?token=ghUBlFsjhR&flag=fxelasticsearch)](https://codecov.io/gh/ankorstore/yokai-contrib/fxelasticsearch)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai-contrib%2Ffxelasticsearch)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai-contrib/fxelasticsearch)](https://pkg.go.dev/github.com/ankorstore/yokai-contrib/fxelasticsearch)

> [Fx](https://uber-go.github.io/fx/) module for [Elasticsearch](https://github.com/elastic/go-elasticsearch).

<!-- TOC -->
* [Overview](#overview)
* [Installation](#installation)
* [Configuration](#configuration)
* [Testing](#testing)
<!-- TOC -->

## Overview

This module provides to your Fx application an [elasticsearch.Client](https://pkg.go.dev/github.com/elastic/go-elasticsearch/v8),
that you can `inject` anywhere to interact with [Elasticsearch](https://github.com/elastic/go-elasticsearch).

## Installation

Install the module:

```shell
go get github.com/ankorstore/yokai-contrib/fxelasticsearch
```

Then activate it in your application bootstrapper:

```go
// internal/bootstrap.go
package internal

import (
	"github.com/ankorstore/yokai-contrib/fxelasticsearch"
	"github.com/ankorstore/yokai/fxcore"
)

var Bootstrapper = fxcore.NewBootstrapper().WithOptions(
	// load modules
	fxelasticsearch.FxElasticsearchModule,
	// ...
)
```

## Configuration

Configuration reference:

```yaml
# ./configs/config.yaml
app:
  name: app
  env: dev
  version: 0.1.0
  debug: true
modules:
  elasticsearch:
    address: ${ELASTICSEARCH_ADDRESS}
    username: ${ELASTICSEARCH_USERNAME}
    password: ${ELASTICSEARCH_PASSWORD}
```

Notes:
- The `modules.elasticsearch.address` configuration key is mandatory
- The `modules.elasticsearch.username` and `modules.elasticsearch.password` configuration keys are optional
- See [Elasticsearch configuration](https://pkg.go.dev/github.com/elastic/go-elasticsearch/v8#Config) documentation for more details

## Testing

In `test` mode, an additional mock Elasticsearch client is provided with HTTP transport-level mocking capabilities.

### Automatic Test Environment Support

When `APP_ENV=test`, the module automatically provides a default mock Elasticsearch client that returns empty successful responses. This allows your application to start and run basic tests without any additional setup.

### Custom Mock Clients

For specific test scenarios, you can create custom mock clients with controlled responses:

```go
package service_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/ankorstore/yokai-contrib/fxelasticsearch"
	"github.com/stretchr/testify/assert"
)

func TestMyService_Search(t *testing.T) {
	// Define mock response
	mockResponse := `{
		"took": 5,
		"timed_out": false,
		"hits": {
			"total": {"value": 1},
			"hits": [
				{
					"_source": {"title": "Test Document", "content": "Test content"}
				}
			]
		}
	}`

	// Create mock Elasticsearch client
	esClient, err := fxelasticsearch.NewMockESClientWithSingleResponse(mockResponse, 200)
	assert.NoError(t, err)

	// Use the mock client in your service
	service := NewMyService(esClient)

	// Test your service methods that use Elasticsearch
	results, err := service.Search(context.Background(), "test-index", "test query")
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "Test Document", results[0]["title"])
}

### Using Injected Mock Client in Tests

You can also use the automatically provided mock client in Fx-based tests:

```go
func TestWithFxInjection(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	
	var esClient *elasticsearch.Client
	
	app := fxtest.New(
		t,
		fxconfig.FxConfigModule,
		fxelasticsearch.FxElasticsearchModule,
		fx.Populate(&esClient),
	)
	
	app.RequireStart()
	
	// The injected client is a mock that returns empty successful responses
	res, err := esClient.Search()
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	
	app.RequireStop()
}
```

See [example](module_test.go).