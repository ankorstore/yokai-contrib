# Fx Google Pub/Sub Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/fxgcppubsub-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/fxgcppubsub-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/fxgcppubsub)](https://goreportcard.com/report/github.com/ankorstore/yokai/fxgcppubsub)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=fxgcppubsub)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/fxgcppubsub)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Ffxgcppubsub)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/fxgcppubsub)](https://pkg.go.dev/github.com/ankorstore/yokai/fxgcppubsub)

> [Fx](https://uber-go.github.io/fx/) module for [Google Pub/Sub](https://cloud.google.com/pubsub).

<!-- TOC -->
* [Overview](#overview)
* [Installation](#installation)
* [Configuration](#configuration)
* [Health Check](#health-check)
* [Testing](#testing)
<!-- TOC -->

## Overview

This module provides to your Fx application a [pubsub.Client](https://pkg.go.dev/cloud.google.com/go/pubsub#Client),
that you can `inject` anywhere to `publish` or `subscribe` on a `Pub/Sub` instance.

## Installation

First install the module:

```shell
go get github.com/ankorstore/yokai-contrib/fxgcppubsub
```

Then activate it in your application bootstrapper:

```go
// internal/bootstrap.go
package internal

import (
	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub"
)

var Bootstrapper = fxcore.NewBootstrapper().WithOptions(
	// load fxgcppubsub module
	fxgcppubsub.FxGcpPubSubModule,
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
  gcppubsub:
    project:
      id: ${GCP_PROJECT_ID}  # GCP project id
    healthcheck:
      topics:                # list of topics to check for the topics probe
        - topic1
        - topic2
      subscriptions:         # list of subscriptions to check for the subscriptions probe
        - subscription1
        - subscription2
```

## Health Check

This module provides ready to use health check probes, to be used by
the [fxhealthcheck](https://ankorstore.github.io/yokai/modules/fxhealthcheck/) module:

- [GcpPubSubTopicsProbe](healthcheck/topic.go): to check existence of the topics in `modules.gcppubsub.healthcheck.topics`
- [GcpPubSubSubscriptionsProbe](healthcheck/subscription.go): to check existence of topics in `modules.gcppubsub.healthcheck.subscriptions`

To activate those probes, you just need to register them:

```go
// internal/services.go
package internal

import (
	"github.com/ankorstore/yokai/fxhealthcheck"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/healthcheck"
	"go.uber.org/fx"
)

func ProvideServices() fx.Option {
	return fx.Options(
		// register the GcpPubSubTopicsProbe
		fxhealthcheck.AsCheckerProbe(healthcheck.NewGcpPubSubTopicsProbe),
		// register the GcpPubSubSubscriptionsProbe
		fxhealthcheck.AsCheckerProbe(healthcheck.NewGcpPubSubSubscriptionsProbe),
		// ...
	)
}
```

If your application is interested only in `publishing`, activate the `GcpPubSubTopicsProbe` only.

If it is interested only in `subscribing`, activate the `GcpPubSubSubscriptionsProbe` only.

## Testing

In `test` mode, this client is configured to [work with
a ptest.Server](module.go), avoiding the need to run any `Pub/Sub`
instance, for better tests portability.

```go
// internal/example/example_test.go
package example_test

import (
	"context"
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/foo/bar/internal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestExample(t *testing.T) {
	var client *pubsub.Client

	internal.RunTest(t, fx.Populate(&client))

	ctx := context.Background()
	
	// prepare test topic on test server
	topic, err := client.CreateTopic(ctx, "test-topic")
	assert.NoError(t, err)

	// public on test topic
	topic.Publish(ctx, &pubsub.Message{Data: []byte("test message")})
	
	// ...
}
```
