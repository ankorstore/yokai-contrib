package fxgcppubsub_test

import (
	"fmt"
	"testing"

	"github.com/ankorstore/yokai-contrib/fxgcppubsub"
	"github.com/stretchr/testify/assert"
)

func TestPrepareSchema(t *testing.T) {
	t.Parallel()

	result := fxgcppubsub.PrepareSchema(fxgcppubsub.PrepareSchemaParams{})

	assert.Equal(t, "fx.invokeOption", fmt.Sprintf("%T", result))
}

func TestPrepareTopic(t *testing.T) {
	t.Parallel()

	result := fxgcppubsub.PrepareTopic(fxgcppubsub.PrepareTopicParams{})

	assert.Equal(t, "fx.invokeOption", fmt.Sprintf("%T", result))
}

func TestPrepareTopicWithSchema(t *testing.T) {
	t.Parallel()

	result := fxgcppubsub.PrepareTopicWithSchema(fxgcppubsub.PrepareTopicWithSchemaParams{})

	assert.Equal(t, "fx.invokeOption", fmt.Sprintf("%T", result))
}

func TestPrepareTopicAndSubscription(t *testing.T) {
	t.Parallel()

	result := fxgcppubsub.PrepareTopicAndSubscription(fxgcppubsub.PrepareTopicAndSubscriptionParams{})

	assert.Equal(t, "fx.invokeOption", fmt.Sprintf("%T", result))
}

func TestPrepareTopicAndSubscriptionWithSchema(t *testing.T) {
	t.Parallel()

	result := fxgcppubsub.PrepareTopicAndSubscriptionWithSchema(fxgcppubsub.PrepareTopicAndSubscriptionWithSchemaParams{})

	assert.Equal(t, "fx.invokeOption", fmt.Sprintf("%T", result))
}
