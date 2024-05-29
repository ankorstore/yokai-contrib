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
<!-- TOC -->

## Overview

TODO

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

TODO
