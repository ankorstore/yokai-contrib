# Fx Clockwork Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/fxclockwork-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/fxclockwork-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/fxclockwork)](https://goreportcard.com/report/github.com/ankorstore/yokai/fxclockwork)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=fxclockwork)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/fxclockwork)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Ffxclockwork)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/fxclockwork)](https://pkg.go.dev/github.com/ankorstore/yokai/fxclockwork)

> [Fx](https://uber-go.github.io/fx/) module for [Clockwork](https://github.com/jonboulle/clockwork).
<!-- TOC -->
* [Overview](#overview)
* [Installation](#installation)
<!-- TOC -->

## Overview

This module provides to your Fx application a [Clockwork.Clock](https://github.com/jonboulle/clockwork),
that you can `inject` anywhere.

## Installation

First install the module:

```shell
go get github.com/ankorstore/yokai-contrib/fxclockwork
```

Then activate it in your application bootstrapper:

```go
// internal/bootstrap.go
package internal

import (
	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai-contrib/fxclockwork"
)

var Bootstrapper = fxcore.NewBootstrapper().WithOptions(
	// load fxclockwork module
	fxclockwork.FxClockworkModule,
	// ...
)
```