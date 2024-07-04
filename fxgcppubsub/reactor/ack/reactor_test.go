package ack_test

import (
	"context"
	"testing"
	"time"

	"cloud.google.com/go/pubsub/apiv1/pubsubpb"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/reactor"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/reactor/ack"
	"github.com/stretchr/testify/assert"
)

func TestAckReactor(t *testing.T) {
	t.Parallel()

	supervisor := reactor.NewWaiterSupervisor()

	react := ack.NewAckReactor(supervisor)

	t.Run("func names", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, []string{"Acknowledge"}, react.FuncNames())
	})

	t.Run("react", func(t *testing.T) {
		t.Parallel()

		req := &pubsubpb.AcknowledgeRequest{
			Subscription: "test-subscription",
			AckIds:       []string{"test-id"},
		}

		waiter := supervisor.StartWaiter("test-subscription")

		go func() {
			rHandled, rRet, rErr := react.React(req)

			assert.False(t, rHandled)
			assert.Nil(t, rRet)
			assert.NoError(t, rErr)
		}()

		data, err := waiter.WaitMaxDuration(context.Background(), 1*time.Millisecond)
		assert.NoError(t, err)
		assert.Equal(t, []string{"test-id"}, data)
	})
}
