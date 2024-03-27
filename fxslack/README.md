# Fx Slack Module

[![ci](https://github.com/ankorstore/yokai-contrib/actions/workflows/fxslack-ci.yml/badge.svg)](https://github.com/ankorstore/yokai-contrib/actions/workflows/fxslack-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai-contrib/fxslack)](https://goreportcard.com/report/github.com/ankorstore/yokai-contrib/fxslack)
[![codecov](https://codecov.io/gh/ankorstore/yokai-contrib/graph/badge.svg?token=ghUBlFsjhR&flag=fxslack)](https://app.codecov.io/gh/ankorstore/yokai-contrib/tree/main/fxslack)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai-contrib%2Ffxslack)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai-contrib/fxslack)](https://pkg.go.dev/github.com/ankorstore/yokai-contrib/fxslack)

> [Fx](https://uber-go.github.io/fx/) module for [Slack](https://api.slack.com).

<!-- TOC -->
* [Overview](#overview)
* [Installation](#installation)
* [Configuration](#configuration)
* [Testing](#testing)
<!-- TOC -->

## Overview

This module provides to your Fx application a [slack.Client](https://pkg.go.dev/github.com/slack-go/slack#Client),
that you can `inject` anywhere to interact with the [Slack API](https://api.slack.com).

## Installation

This module requires [fxhttpclient](https://github.com/ankorstore/yokai/tree/main/fxhttpclient).

Install the modules:

```shell
go get github.com/ankorstore/yokai/fxhttpclient
go get github.com/ankorstore/yokai-contrib/fxslack
```

Then activate them in your application bootstrapper:

```go
// internal/bootstrap.go
package internal

import (
	"github.com/ankorstore/yokai-contrib/fxslack"
	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai/fxhttpclient"
)

var Bootstrapper = fxcore.NewBootstrapper().WithOptions(
	// load modules
	fxhttpclient.FxHttpClientModule,
	fxslack.FxSlackModule,
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
    auth_token: ${SLACK_AUTH_TOKEN}  # Slack Auth Token
    app_level_token: ${SLACK_APP_LEVEL_TOKEN} # Slack App level Token
```

## Testing

In `test` mode, this client is configured to interact with a [fake slack server](https://github.com/slack-go/slack/tree/master/slacktest).

See [example](module_test.go).
