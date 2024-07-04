package fxgcppubsub_test

import (
	"fmt"
	"testing"

	"github.com/ankorstore/yokai-contrib/fxgcppubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/reactor/ack"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/reactor/log"
	"github.com/stretchr/testify/assert"
)

func TestAsPubSubTestServerReactor(t *testing.T) {
	t.Parallel()

	result := fxgcppubsub.AsPubSubTestServerReactor(log.NewLogReactor)

	assert.Equal(t, "fx.provideOption", fmt.Sprintf("%T", result))
}

func TestAsPubSubTestServerReactors(t *testing.T) {
	t.Parallel()

	result := fxgcppubsub.AsPubSubTestServerReactors(log.NewLogReactor, ack.NewAckReactor)

	assert.Equal(t, "fx.optionGroup", fmt.Sprintf("%T", result))
}
