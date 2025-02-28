# Yokai JSON API Module

[![ci](https://github.com/ankorstore/yokai-contrib/actions/workflows/fxjsonapi-ci.yml/badge.svg)](https://github.com/ankorstore/yokai-contrib/actions/workflows/fxjsonapi-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai-contrib/fxjsonapi)](https://goreportcard.com/report/github.com/ankorstore/yokai-contrib/fxjsonapi)
[![codecov](https://codecov.io/gh/ankorstore/yokai-contrib/graph/badge.svg?token=ghUBlFsjhR&flag=fxjsonapi)](https://app.codecov.io/gh/ankorstore/yokai-contrib/tree/main/fxjsonapi)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai-contrib%2Ffxjsonapi)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai-contrib/fxjsonapi)](https://pkg.go.dev/github.com/ankorstore/yokai-contrib/fxjsonapi)

> [Yokai](https://github.com/ankorstore/yokai) module for [JSON API](https://jsonapi.org/), based on [google/jsonapi](https://github.com/google/jsonapi).

<!-- TOC -->
* [Overview](#overview)
* [Installation](#installation)
* [Configuration](#configuration)
* [Processing](#processing)
  * [Request processing](#request-processing)
  * [Response processing](#response-processing)
* [Error handling](#error-handling)
* [Testing](#testing)
<!-- TOC -->

## Overview

This module provides to your [Yokai](https://github.com/ankorstore/yokai) application a [Processor](processor.go), that you can `inject` in your HTTP handlers to process JSON API requests and responses.

It also provides automatic [error handling](error.go), compliant with the [JSON API specifications](https://jsonapi.org/).

## Installation

Install the module:

```shell
go get github.com/ankorstore/yokai-contrib/fxjsonapi
```

Then activate it in your application bootstrapper:

```go
// internal/bootstrap.go
package internal

import (
	"github.com/ankorstore/yokai-contrib/fxjsonapi"
	"github.com/ankorstore/yokai/fxcore"
)

var Bootstrapper = fxcore.NewBootstrapper().WithOptions(
	// load modules
	fxjsonapi.FxJSONAPIModule,
	// ...
)
```

## Configuration

Configuration reference:

```yaml
# ./configs/config.yaml
modules:
  jsonapi:
    log:
      enabled: true # to automatically log JSON API processing, disabled by default
    trace:
      enabled: true # to automatically trace JSON API processing, disabled by default
```

## Processing

### Request processing

You can use the provided [Processor](processor.go) to automatically process a JSON API request:

```go
package handler

import (
	"net/http"

	"github.com/ankorstore/yokai-contrib/fxjsonapi"
	"github.com/google/jsonapi"
	"github.com/labstack/echo/v4"
)

type Foo struct {
	ID   int    `jsonapi:"primary,foo"`
	Name string `jsonapi:"attr,name"`
	Bar  *Bar   `jsonapi:"relation,bar"`
}

func (f Foo) JSONAPIMeta() *jsonapi.Meta {
	return &jsonapi.Meta{
		"some": "foo meta",
	}
}

type Bar struct {
	ID   int    `jsonapi:"primary,bar"`
	Name string `jsonapi:"attr,name"`
}

func (b Bar) JSONAPIMeta() *jsonapi.Meta {
	return &jsonapi.Meta{
		"some": "bar meta",
	}
}

type JSONAPIHandler struct {
	processor fxjsonapi.Processor
}

func NewJSONAPIHandler(processor fxjsonapi.Processor) *JSONAPIHandler {
	return &JSONAPIHandler{
		processor: processor,
	}
}

func (h *JSONAPIHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		foo := Foo{}

		// unmarshall JSON API request payload in foo
		err := h.processor.ProcessRequest(
			// echo context
			c,
			// pointer to the struct to unmarshall
			&foo,
			// optionally override module config for logging
			fxjsonapi.WithLog(true),
			// optionally override module config for tracing
			fxjsonapi.WithTrace(true),
			)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, foo)
	}
}
```

Notes about `ProcessRequest()`:

- if the request payload does not respect the [JSON API specifications](https://jsonapi.org/), a `400` error will be automatically returned
- if the request `Content-Type` headers is not `application/vnd.api+json`, a `415` error will be automatically returned

### Response processing

You can use the provided [Processor](processor.go) to automatically process a JSON API response:

```go
package handler

import (
	"net/http"

	"github.com/ankorstore/yokai-contrib/fxjsonapi"
	"github.com/ankorstore/yokai-contrib/fxjsonapi/testdata/model"
	"github.com/labstack/echo/v4"
)

type Foo struct {
	ID   int    `jsonapi:"primary,foo"`
	Name string `jsonapi:"attr,name"`
	Bar  *Bar   `jsonapi:"relation,bar"`
}

func (f Foo) JSONAPIMeta() *jsonapi.Meta {
	return &jsonapi.Meta{
		"some": "foo meta",
	}
}

type Bar struct {
	ID   int    `jsonapi:"primary,bar"`
	Name string `jsonapi:"attr,name"`
}

func (b Bar) JSONAPIMeta() *jsonapi.Meta {
	return &jsonapi.Meta{
		"some": "bar meta",
	}
}

type JSONAPIHandler struct {
	processor fxjsonapi.Processor
}

func NewJSONAPIHandler(processor fxjsonapi.Processor) *JSONAPIHandler {
	return &JSONAPIHandler{
		processor: processor,
	}
}

func (h *JSONAPIHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		foo := Foo{
			ID:   123,
			Name: "foo",
			Bar: &Bar{
				ID:   456,
				Name: "bar",
			},
		}

		return h.processor.ProcessResponse(
			// echo context
			c,
			// HTTP status code
			http.StatusOK,
			// pointer to the struct to marshall
			&foo,
			// optionally pass metadata to the JSON API response
			fxjsonapi.WithMetadata(map[string]interface{}{
				"some": "response meta",
			}),
			// optionally remove the included from the JSON API response (enabled by default)
			fxjsonapi.WithIncluded(false),
			// optionally override module config for logging
			fxjsonapi.WithLog(true),
			// optionally override module config for tracing
			fxjsonapi.WithTrace(true),
		)
	}
}
```

Notes about `ProcessResponse()`:

- you can pass a pointer or a slice of pointers to marshall as JSON API
- `application/vnd.api+json` will be automatically added to the response `Content-Type` header

## Error handling

This module automatically enables the [ErrorHandler](error.go), to convert errors bubbling up in JSON API format.

It handles:

- [JSON API errors](https://github.com/google/jsonapi/blob/master/errors.go) errors (automatically sets a 500 status code)
- [validation](https://ankorstore.github.io/yokai/modules/fxvalidator/) errors (automatically sets a 400 status code)
- [HTTP](https://echo.labstack.com/docs/error-handling) errors (automatically sets the status code of the error)
- or any generic error (automatically sets a 500 status code)

## Testing

This module provides a [ProcessorMock](fxjsonapitest/mock.go) for mocking [Processor](processor.go), see [usage example](fxjsonapitest/mock_test.go).
