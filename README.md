# Yokai Contrib Modules

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go version](https://img.shields.io/badge/Go-≥1.20-blue)](https://go.dev/)
[![codecov](https://codecov.io/gh/ankorstore/yokai-contrib/graph/badge.svg?token=CxImMei31C)](https://codecov.io/gh/ankorstore/yokai-contrib)

> Contrib modules repository for the [Yokai](https://github.com/ankorstore/yokai) framework.

## Modules

| Module                             | Description                                                              |
|------------------------------------|--------------------------------------------------------------------------|
| [fxgcppubsub](fxgcppubsub)         | Module for [GCP Pub/Sub](https://cloud.google.com/pubsub)                |
| [fxgomysqlserver](fxgomysqlserver) | Module for [Go Mysql Server](https://github.com/dolthub/go-mysql-server) |
| [fxjsonapi](fxjsonapi)             | Module for [JSON API](https://github.com/google/jsonapi)                 |
| [fxslack](fxslack)                 | Module for [Slack](https://api.slack.com/)                               |
| [fxredis](fxredis)                 | Module for [Redis](https://redis.io/docs/connect/clients/go/)            |

## Contributing

This repository uses [release-please](https://github.com/googleapis/release-please) to automate Yokai's contrib modules release process.

> [!IMPORTANT]
> You must provide [atomic](https://en.wikipedia.org/wiki/Atomic_commit#Revision_control) and [conventional](https://www.conventionalcommits.org/en/v1.0.0/) commits, as the release process relies on them to determine the version to release and to generate the release notes.
