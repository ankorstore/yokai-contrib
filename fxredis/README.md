# Fx Redis Module

[![ci](https://github.com/ankorstore/yokai-contrib/actions/workflows/fxredis-ci.yml/badge.svg)](https://github.com/ankorstore/yokai-contrib/actions/workflows/fxredis-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai-contrib/fxredis)](https://goreportcard.com/report/github.com/ankorstore/yokai-contrib/fxredis)
[![codecov](https://codecov.io/gh/ankorstore/yokai-contrib/graph/badge.svg?token=ghUBlFsjhR&flag=fxredis)](https://app.codecov.io/gh/ankorstore/yokai-contrib/tree/main/fxredis)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai-contrib%2Ffxredis)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai-contrib/fxredis)](https://pkg.go.dev/github.com/ankorstore/yokai-contrib/fxredis)

> [Fx](https://uber-go.github.io/fx/) module for [Redis](https://redis.io/docs/connect/clients/go/).

<!-- TOC -->
* [Overview](#overview)
* [Installation](#installation)
* [Configuration](#configuration)
* [Testing](#testing)
<!-- TOC -->

## Overview

This module provides to your Fx application a [redis.Client](https://pkg.go.dev/github.com/go-redis/redis/v9#Client),
that you can `inject` anywhere to interact with [Redis](https://redis.io/docs/connect/clients/go/).

## Installation

Install the module:

```shell
go get github.com/ankorstore/yokai-contrib/fxredis
```

Then activate them in your application bootstrapper:

```go
// internal/bootstrap.go
package internal

import (
	"github.com/ankorstore/yokai-contrib/fxredis"
	"github.com/ankorstore/yokai/fxcore"
)

var Bootstrapper = fxcore.NewBootstrapper().WithOptions(
	// load modules
	fxredis.FxRedisModule,
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
  redis:
    dsn: redis://${REDIS_USER}:${REDIS_PASSWORD}@${REDIS_HOST}:${REDIS_PORT}/${REDIS_DB}
```

## Testing

In `test` mode, an additional [redismock.ClientMock](https://pkg.go.dev/github.com/go-redis/redismock/v9#ClientMock) is provided.

See [example](module_test.go).
