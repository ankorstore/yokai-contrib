# Fx Go MySQL Server Module

[![ci](https://github.com/ankorstore/yokai-contrib/actions/workflows/fxgomysqlserver-ci.yml/badge.svg)](https://github.com/ankorstore/yokai-contrib/actions/workflows/fxgomysqlserver-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai-contrib/fxgomysqlserver)](https://goreportcard.com/report/github.com/ankorstore/yokai-contrib/fxgomysqlserver)
[![codecov](https://codecov.io/gh/ankorstore/yokai-contrib/graph/badge.svg?token=ghUBlFsjhR&flag=fxgomysqlserver)](https://app.codecov.io/gh/ankorstore/yokai-contrib/tree/main/fxgomysqlserver)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai-contrib%2Ffxgomysqlserver)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai-contrib/fxgomysqlserver)](https://pkg.go.dev/github.com/ankorstore/yokai-contrib/fxgomysqlserver)

> [Fx](https://uber-go.github.io/fx/) module for [Go Mysql Server](https://github.com/dolthub/go-mysql-server).

<!-- TOC -->
* [Overview](#overview)
* [Installation](#installation)
* [Configuration](#configuration)
* [Usage](#usage)
  * [TCP transport](#tcp-transport)
  * [Memory transport](#memory-transport)
* [Testing](#testing)
<!-- TOC -->

## Overview

This module integrates an embedded [Go MySQL server](https://github.com/dolthub/go-mysql-server) into your Yokai application.

This is made for `development / testing purposes`, not to replace a real [MySQL server](https://www.mysql.com/) for production applications.

It can be configured to accept connections:

- via `TCP`
- via a `socket`
- via `memory`

Make sure to acknowledge the underlying vendor [limitations](https://github.com/dolthub/go-mysql-server?tab=readme-ov-file#limitations-of-the-in-memory-database-implementation).

## Installation

First install the module:

```shell
go get github.com/ankorstore/yokai-contrib/fxgomysqlserver
```

Then activate it in your application bootstrapper:

```go
// internal/bootstrap.go
package internal

import (
	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai-contrib/fxgomysqlserver"
)

var Bootstrapper = fxcore.NewBootstrapper().WithOptions(
	// load fxgomysqlserver module
	fxgomysqlserver.FxGoMySQLServerModule,
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
  gomysqlserver:
    config:
      transport: tcp            # transport to use: tcp or memory (tcp by default)
      user: user                # database user name (user by default)
      password: password        # database user password (password by default)
      host: localhost           # database host (localhost by default)
      port: 3306                # database port (3306 by default)
      database: db              # database name (db by default)
    log:
      enabled: true             # to enable server logs
    trace:
      enabled: true             # to enable server traces
```

## Usage

### TCP transport

You can configure the server to accept `TCP` connections:

```yaml
# ./configs/config.yaml
modules:
  gomysqlserver:
    config:
      transport: tcp
      user: user
      password: password
      host: localhost
      port: 3306
      database: db
```

And then connect with:

```go
import "database/sql"

db, _ := sql.Open("mysql", "user:password@tcp(localhost:3306)/db")
```

### Memory transport

You can configure the server to accept `memory` connections:

```yaml
# ./configs/config.yaml
modules:
  gomysqlserver:
    config:
      transport: memory
      user: user
      password: password
      database: db
```

And then connect with:

```go
import "database/sql"

db, _ := sql.Open("mysql", "user:password@memory(bufconn)/db")
```

## Testing

The `memory` transport avoids to open TCP ports or sockets, making it particularly lean and useful for `testing` purposes.

In you test configuration:

```yaml
# ./configs/config.test.yaml
modules:
  gomysqlserver:
    config:
      transport: memory
```

In you application bootstrapper:

```go
// internal/bootstrap.go
package internal

import (
	"testing"
	
	"github.com/ankorstore/yokai-contrib/fxgomysqlserver"
)

// RunTest starts the application in test mode, with an optional list of [fx.Option].
func RunTest(tb testing.TB, options ...fx.Option) {
	tb.Helper()

	// ...

	Bootstrapper.RunTestApp(
		tb,
		// enable and start the server
		fxgomysqlserver.FxGoMySQLServerModule,
		// apply per test options
		fx.Options(options...),
	)
}
```

You can check the [tests of this module](module_test.go) to get testing examples.