# Fx Elasticsearch Module

[![ci](https://github.com/ankorstore/yokai-contrib/actions/workflows/fxelasticsearch-ci.yml/badge.svg)](https://github.com/ankorstore/yokai-contrib/actions/workflows/fxelasticsearch-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai-contrib/fxelasticsearch)](https://goreportcard.com/report/github.com/ankorstore/yokai-contrib/fxelasticsearch)
[![codecov](https://codecov.io/gh/ankorstore/yokai-contrib/graph/badge.svg?token=ghUBlFsjhR&flag=fxelasticsearch)](https://app.codecov.io/gh/ankorstore/yokai-contrib/tree/main/fxelasticsearch)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai-contrib%2Ffxelasticsearch)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai-contrib/fxelasticsearch)](https://pkg.go.dev/github.com/ankorstore/yokai-contrib/fxelasticsearch)

> [Fx](https://uber-go.github.io/fx/) module for [Elasticsearch](https://www.elastic.co/elasticsearch/).

<!-- TOC -->
* [Overview](#overview)
* [Installation](#installation)
* [Configuration](#configuration)
* [Usage](#usage)
<!-- TOC -->

## Overview

This module provides to your Fx application an [elasticsearch.Client](https://pkg.go.dev/github.com/elastic/go-elasticsearch/v8),
that you can `inject` anywhere to interact with [Elasticsearch](https://www.elastic.co/elasticsearch/).

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

## Usage

Inject the Elasticsearch client in your services:

```go
package service

import (
    "context"
    "github.com/elastic/go-elasticsearch/v8"
    "github.com/elastic/go-elasticsearch/v8/esapi"
)

type MyService struct {
    es *elasticsearch.Client
}

func NewMyService(es *elasticsearch.Client) *MyService {
    return &MyService{es: es}
}

func (s *MyService) SearchDocuments(ctx context.Context, index string, query string) (*esapi.Response, error) {
    // Use the Elasticsearch client to search for documents
    return s.es.Search(
        s.es.Search.WithContext(ctx),
        s.es.Search.WithIndex(index),
        s.es.Search.WithBody(strings.NewReader(query)),
    )
}
``` 