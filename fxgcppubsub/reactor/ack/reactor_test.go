package ack_test

import (
	"context"
	"testing"
	"time"

	"cloud.google.com/go/pubsub/apiv1/pubsubpb"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/reactor/ack"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/subscription"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestAckReactor(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "../../testdata/config")
	t.Setenv("GCP_PROJECT_ID", "test-project")

	var sup ack.AckSupervisor

	ctx := context.Background()

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxgcppubsub.FxGcpPubSubModule,
		fx.Supply(fx.Annotate(ctx, fx.As(new(context.Context)))),
		fx.Populate(&sup),
	).RequireStart().RequireStop()

	react := ack.NewAckReactor(sup)

	t.Run("func names", func(t *testing.T) {
		assert.Equal(
			t,
			[]string{
				"Acknowledge",
				"ModifyAckDeadline",
			},
			react.FuncNames(),
		)
	})

	t.Run("react to ack", func(t *testing.T) {
		req := &pubsubpb.AcknowledgeRequest{
			Subscription: subscription.NormalizeSubscriptionName("test-project", "test-subscription"),
			AckIds:       []string{"test-id"},
		}

		waiter := sup.StartAckWaiter("test-subscription")

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

	t.Run("react to nack", func(t *testing.T) {
		req := &pubsubpb.ModifyAckDeadlineRequest{
			Subscription: subscription.NormalizeSubscriptionName("test-project", "test-subscription"),
			AckIds:       []string{"test-id"},
		}

		waiter := sup.StartNackWaiter("test-subscription")

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
