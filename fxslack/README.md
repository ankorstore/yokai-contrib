# Fx Slack Module

[![ci](https://github.com/ankorstore/yokai-contrib/actions/workflows/fxgcppubsub-ci.yml/badge.svg)](https://github.com/ankorstore/yokai-contrib/actions/workflows/fxgcppubsub-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai-contrib/fxgcppubsub)](https://goreportcard.com/report/github.com/ankorstore/yokai-contrib/fxgcppubsub)
[![codecov](https://codecov.io/gh/ankorstore/yokai-contrib/graph/badge.svg?token=ghUBlFsjhR&flag=fxgcppubsub)](https://app.codecov.io/gh/ankorstore/yokai-contrib/tree/main/fxgcppubsub)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai-contrib%2Ffxgcppubsub)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai-contrib/fxgcppubsub)](https://pkg.go.dev/github.com/ankorstore/yokai-contrib/fxgcppubsub)

> [Fx](https://uber-go.github.io/fx/) module for [GCP Pub/Sub](https://cloud.google.com/pubsub).

<!-- TOC -->
* [Overview](#overview)
* [Installation](#installation)
* [Configuration](#configuration)
* [Health Check](#health-check)
* [Testing](#testing)
<!-- TOC -->

## Overview

This module provides to your Fx application a [slack.Client](https://pkg.go.dev/github.com/slack-go/slack#Client),
that you can `inject` anywhere to interact with the Slack API.

## Installation

First install the module:

```shell
go get github.com/ankorstore/yokai-contrib/fxslack
```

Then activate it in your application bootstrapper:

```go
// internal/bootstrap.go
package internal

import (
	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai-contrib/fxslack"
)

var Bootstrapper = fxcore.NewBootstrapper().WithOptions(
	// load fxslack module
	fxslack.FxSlackbModule,
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
  slack:
    token: ${SLACK_TOKEN}  # Slack Token
```

## Testing

In `test` mode, this client is configured to interact with a [fake slack server](https://github.com/slack-go/slack/tree/master/slacktest).
