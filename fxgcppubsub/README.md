# Yokai GCP Pub/Sub Module

[![ci](https://github.com/ankorstore/yokai-contrib/actions/workflows/fxgcppubsub-ci.yml/badge.svg)](https://github.com/ankorstore/yokai-contrib/actions/workflows/fxgcppubsub-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai-contrib/fxgcppubsub)](https://goreportcard.com/report/github.com/ankorstore/yokai-contrib/fxgcppubsub)
[![codecov](https://codecov.io/gh/ankorstore/yokai-contrib/graph/badge.svg?token=ghUBlFsjhR&flag=fxgcppubsub)](https://app.codecov.io/gh/ankorstore/yokai-contrib/tree/main/fxgcppubsub)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai-contrib%2Ffxgcppubsub)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai-contrib/fxgcppubsub)](https://pkg.go.dev/github.com/ankorstore/yokai-contrib/fxgcppubsub)

> [Yokai](https://github.com/ankorstore/yokai) module for [GCP Pub/Sub](https://cloud.google.com/pubsub).

<!-- TOC -->
* [Overview](#overview)
* [Installation](#installation)
* [Configuration](#configuration)
* [Publish](#publish)
  * [Raw message](#raw-message)
  * [Avro message](#avro-message)
  * [Protobuf message](#protobuf-message)
* [Subscribe](#subscribe)
  * [Raw message](#raw-message-1)
  * [Avro message](#avro-message-1)
  * [Protobuf message](#protobuf-message-1)
* [Health Check](#health-check)
* [Testing](#testing)
<!-- TOC -->

## Overview

This module provides to your [Yokai](https://github.com/ankorstore/yokai) application the possibility to `publish` and/or `subscribe` on a [GCP Pub/Sub](https://cloud.google.com/pubsub) instance.

It also provides the support of [Avro](https://avro.apache.org/) and [Protobuf](https://protobuf.dev/) schemas.

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
modules:
  gcppubsub:
    project:
      id: ${GCP_PROJECT_ID}  # GCP project id
    healthcheck:
      topics:                # list of topics to check for the topics probe
        - some-topic         # refers to projects/${GCP_PROJECT_ID}/topics/some-topic
      subscriptions:         # list of subscriptions to check for the subscriptions probe
        - some-subscription  # refers to projects/${GCP_PROJECT_ID}/subscriptions/some-subscription
```

## Publish

This module provides a high level [Publisher](publisher.go) that you can use to `publish` messages on a `topic`, that you can inject anywhere.

If the `topic` is associated to an `avro` or `protobuf` schema, the publisher will automatically handle the message encoding.

This module also provides a [pubsub.Client](https://pkg.go.dev/cloud.google.com/go/pubsub#Client), that you can use for [low level publishing](https://cloud.google.com/pubsub/docs/samples/pubsub-quickstart-publisher).

### Raw message

To publish a raw message (without associated schema) on a topic:

```go
// publish on projects/${GCP_PROJECT_ID}/topics/some-topic
res, err := publisher.Publish(context.Backgound(), "some-topic", "some message")
```

### Avro message

The publisher can accept any struct, and will automatically handle the avro encoding based on the following tags:

- `avro` to drive avro binary encoding (see the [underlying library documentation](https://github.com/hamba/avro) for more details)
- `json` to drive avro json encoding (see the [underlying library documentation](https://github.com/linkedin/goavro) for more details)

Considering this avro schema:

```avroschema
{
  "namespace": "Simple",
  "type": "record",
  "name": "Avro",
  "fields": [
    {
      "name": "StringField",
      "type": "string"
    },
    {
      "name": "FloatField",
      "type": "float"
    },
    {
      "name": "BooleanField",
      "type": "boolean"
    }
  ]
}
```

To publish a message on a topic associated to this avro schema:

```go
// struct with tags, representing the message
type SimpleRecord struct {
    StringField  string  `avro:"StringField" json:"StringField"`
    FloatField   float32 `avro:"FloatField" json:"FloatField"`
    BooleanField bool    `avro:"BooleanField" json:"BooleanField"`
}

// publish on projects/${GCP_PROJECT_ID}/topics/some-topic
res, err := publisher.Publish(context.Backgound(), "some-topic", &SimpleRecord{
    StringField:  "some string",
    FloatField:   12.34,
    BooleanField: true,
})
```

### Protobuf message

The publisher can accept any [proto.Message](https://github.com/golang/protobuf/blob/master/proto/proto.go), and will automatically handle the protobuf binary or json encoding.

Considering this protobuf schema:

```protobuf
syntax = "proto3";

package simple;

option go_package = "github.com/ankorstore/yokai-contrib/fxgcppubsub/testdata/proto";

message SimpleRecord {
    string string_field = 1;
    float float_field = 2;
    bool boolean_field = 3;
}
```

To publish a message on a topic associated to this protobuf schema:

```go
// generated protobuf stub proto.Message, representing the message
type SimpleRecord struct {
    state         protoimpl.MessageState
    sizeCache     protoimpl.SizeCache
    unknownFields protoimpl.UnknownFields
    
    StringField  string  `protobuf:"bytes,1,opt,name=string_field,json=stringField,proto3" json:"string_field,omitempty"`
    FloatField   float32 `protobuf:"fixed32,2,opt,name=float_field,json=floatField,proto3" json:"float_field,omitempty"`
    BooleanField bool    `protobuf:"varint,3,opt,name=boolean_field,json=booleanField,proto3" json:"boolean_field,omitempty"`
}

// publish on projects/${GCP_PROJECT_ID}/topics/some-topic
res, err := publisher.Publish(context.Backgound(), "some-topic", &SimpleRecord{
    StringField:  "test proto",
    FloatField:   56.78,
    BooleanField: false,
})
```

## Subscribe

This module provides a high level [Subscriber](subscriber.go) that you can use to `subscribe` messages from a `subscription`, that you can inject anywhere.

If the subscription's `topic` is associated to an `avro` or `protobuf` schema, the subscriber will offer a message from which you can handle the decoding.

This module also provides a [pubsub.Client](https://pkg.go.dev/cloud.google.com/go/pubsub#Client), that you can use for [low level subscribing](https://cloud.google.com/pubsub/docs/samples/pubsub-quickstart-subscriber).

### Raw message

To subscribe on a subscription receiving raw messages (without associated schema):

```go
// subscribe from projects/${GCP_PROJECT_ID}/subscriptions/some-subscription
err := subscriber.Subscribe(ctx, "some-subscription", func(ctx context.Context, m *message.Message) {
    fmt.Printf("%s", m.Data())
    
    m.Ack()
})
```

### Avro message

The subscriber message can be decoded into any struct with the following tags:

- `avro` to drive avro binary decoding (see the [underlying library documentation](https://github.com/hamba/avro) for more details)
- `json` to drive avro json decoding (see the [underlying library documentation](https://github.com/linkedin/goavro) for more details)

Considering this avro schema:

```avroschema
{
  "namespace": "Simple",
  "type": "record",
  "name": "Avro",
  "fields": [
    {
      "name": "StringField",
      "type": "string"
    },
    {
      "name": "FloatField",
      "type": "float"
    },
    {
      "name": "BooleanField",
      "type": "boolean"
    }
  ]
}
```

To subscribe from a subscription associated to this avro schema:

```go
// struct with tags, representing the message
type SimpleRecord struct {
    StringField  string  `avro:"StringField" json:"StringField"`
    FloatField   float32 `avro:"FloatField" json:"FloatField"`
    BooleanField bool    `avro:"BooleanField" json:"BooleanField"`
}

// subscribe from projects/${GCP_PROJECT_ID}/subscriptions/some-subscription
err := subscriber.Subscribe(ctx, "some-subscription", func(ctx context.Context, m *message.Message) {
    var rec SimpleRecord
    
    err = m.Decode(&rec)
    if err != nil {
        m.Nack()
    }
    
    fmt.Printf("%v", rec)
    
    m.Ack()
})
```

### Protobuf message

The subscriber message can be decoded into any [proto.Message](https://github.com/golang/protobuf/blob/master/proto/proto.go), for protobuf binary or json encoding.

Considering this protobuf schema:

```protobuf
syntax = "proto3";

package simple;

option go_package = "github.com/ankorstore/yokai-contrib/fxgcppubsub/testdata/proto";

message SimpleRecord {
  string string_field = 1;
  float float_field = 2;
  bool boolean_field = 3;
}
```

To subscribe from a subscription associated to this protobuf schema:

```go
// generated protobuf stub proto.Message, representing the message
type SimpleRecord struct {
    state         protoimpl.MessageState
    sizeCache     protoimpl.SizeCache
    unknownFields protoimpl.UnknownFields
    
    StringField  string  `protobuf:"bytes,1,opt,name=string_field,json=stringField,proto3" json:"string_field,omitempty"`
    FloatField   float32 `protobuf:"fixed32,2,opt,name=float_field,json=floatField,proto3" json:"float_field,omitempty"`
    BooleanField bool    `protobuf:"varint,3,opt,name=boolean_field,json=booleanField,proto3" json:"boolean_field,omitempty"`
}

// subscribe from projects/${GCP_PROJECT_ID}/subscriptions/some-subscription
err := subscriber.Subscribe(ctx, "some-subscription", func(ctx context.Context, m *message.Message) {
    var rec SimpleRecord
    
    err = m.Decode(&rec)
    if err != nil {
        m.Nack()
    }
    
    fmt.Printf("%v", rec)
    
    m.Ack()
})
```

## Health Check

This module provides ready to use health check probes, to be used by
the [Health Check](https://ankorstore.github.io/yokai/modules/fxhealthcheck/) module:

- [GcpPubSubTopicsProbe](healthcheck/topic.go): to check existence of the topics in `modules.gcppubsub.healthcheck.topics`
- [GcpPubSubSubscriptionsProbe](healthcheck/subscription.go): to check existence of topics in `modules.gcppubsub.healthcheck.subscriptions`

Considering the following configuration:

```yaml
# ./configs/config.yaml
app:
modules:
  gcppubsub:
    project:
      id: ${GCP_PROJECT_ID}   # GCP project id
    healthcheck:
      topics:                 # list of topics to check for the topics probe
        - some-topic          # refers to projects/${GCP_PROJECT_ID}/topics/some-topic
        - other-topic         # refers to projects/${GCP_PROJECT_ID}/topics/other-topic
      subscriptions:          # list of subscriptions to check for the subscriptions probe
        - some-subscription   # refers to projects/${GCP_PROJECT_ID}/subscriptions/some-subscription
        - other-subscription  # refers to projects/${GCP_PROJECT_ID}/subscriptions/other-subscription
```

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
		// register the GcpPubSubTopicsProbe for some-topic and other-topic
		fxhealthcheck.AsCheckerProbe(healthcheck.NewGcpPubSubTopicsProbe),
		// register the GcpPubSubSubscriptionsProbe for some-subscription and other-subscription
		fxhealthcheck.AsCheckerProbe(healthcheck.NewGcpPubSubSubscriptionsProbe),
		// ...
	)
}
```

Notes:

- if your application is interested only in `publishing`, activate the `GcpPubSubTopicsProbe` only
- if it is interested only in `subscribing`, activate the `GcpPubSubSubscriptionsProbe` only

## Testing

In `test` mode:
- the high level [Publisher](publisher.go) and [Subscriber](https://pkg.go.dev/cloud.google.com/go/pubsub@v1.40.0/pstest)
- and the low level [pubsub.Client](https://pkg.go.dev/cloud.google.com/go/pubsub#Client) and [pubsub.SchemaClient](https://pkg.go.dev/cloud.google.com/go/pubsub#SchemaClient)

are all configured to work against a [pstest.Server](https://pkg.go.dev/cloud.google.com/go/pubsub@v1.40.0/pstest), avoiding the need to spin up any `Pub/Sub` real (or emulator) instance, for better tests portability.

You can create `topics`, `subscriptions` and `schemas` locally only for your tests.

For example:

```go
// internal/example/example_test.go
package example_test

import (
	"context"
	"testing"
	"time"
	
	"github.com/ankorstore/yokai-contrib/fxgcppubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/message"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/reactor/ack"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/stretchr/testify/assert"
	"github.com/foo/bar/internal"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestPubSub(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("GCP_PROJECT_ID", "test-project")

	var publisher fxgcppubsub.Publisher
	var subscriber fxgcppubsub.Subscriber
	var supervisor ack.AckSupervisor

	ctx := context.Background()

	// test app
	internal.RunTest(
		t,
		// prepare test topic and subscription
		fxgcppubsub.PrepareTopicAndSubscription(fxgcppubsub.PrepareTopicAndSubscriptionParams{
			TopicID:        "test-topic",
			SubscriptionID: "test-subscription",
		}),
		fx.Populate(&publisher, &subscriber, &supervisor),
	)
	
	// publish to test-topic
	res, err := publisher.Publish(ctx, "test-topic", []byte("test data"))
	assert.NotNil(t, res)
	assert.NoError(t, err)

	sid, err := res.Get(ctx)
	assert.NotEmpty(t, sid)
	assert.NoError(t, err)
	
	waiter := supervisor.StartAckWaiter("test-subscription")
	
	// subscribe from test-subscription
	go subscriber.Subscribe(ctx, "test-subscription", func(ctx context.Context, m *message.Message) {
		assert.Equal(t, []byte("test data"), m.Data())

		m.Ack()
	})

	// waits for subscriber message ack
	_, err = waiter.WaitMaxDuration(ctx, time.Second)
	assert.NoError(t, err)
}
```

Notes:

- you can prepare the test `topics`, `subscriptions` and `schemas` using the [provided helpers](prepare.go)
- you can find tests involving `avro` and `protobuf` schemas in the module [test examples](module_test.go)
