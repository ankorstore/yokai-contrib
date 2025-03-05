package ack_test

import (
	"context"
	"testing"
	"time"

	"github.com/ankorstore/yokai-contrib/fxgcppubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/reactor/ack"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/subscription"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestAckSupervisor(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "../../testdata/config")
	t.Setenv("GCP_PROJECT_ID", "test-project")

	var supervisor ack.AckSupervisor

	ctx := context.Background()

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxgcppubsub.FxGcpPubSubModule,
		fx.Supply(fx.Annotate(ctx, fx.As(new(context.Context)))),
		fx.Populate(&supervisor),
	).RequireStart().RequireStop()

	t.Run("wait for ack", func(t *testing.T) {
		waiter := supervisor.StartAckWaiter("test-subscription")

		go func() {
			time.Sleep(1 * time.Millisecond)

			supervisor.StopAckWaiter(
				subscription.NormalizeSubscriptionName("test-project", "test-subscription"),
				[]string{"test-id"},
				nil,
			)
		}()

		data, err := waiter.WaitMaxDuration(context.Background(), 5*time.Millisecond)
		assert.NoError(t, err)
		assert.Equal(t, []string{"test-id"}, data)
	})

	t.Run("wait for ack with error", func(t *testing.T) {
		waiter := supervisor.StartAckWaiter("test-subscription")

		go func() {
			time.Sleep(1 * time.Millisecond)

			supervisor.StopAckWaiter(
				subscription.NormalizeSubscriptionName("test-project", "test-subscription"),
				[]string{"test-id"},
				assert.AnError,
			)
		}()

		data, err := waiter.WaitMaxDuration(context.Background(), 5*time.Millisecond)
		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
		assert.Equal(t, []string{"test-id"}, data)
	})

	t.Run("wait for nack", func(t *testing.T) {
		waiter := supervisor.StartNackWaiter("test-subscription")

		go func() {
			time.Sleep(1 * time.Millisecond)

			supervisor.StopNackWaiter(
				subscription.NormalizeSubscriptionName("test-project", "test-subscription"),
				[]string{"test-id"},
				nil,
			)
		}()

		data, err := waiter.WaitMaxDuration(context.Background(), 5*time.Millisecond)
		assert.NoError(t, err)
		assert.Equal(t, []string{"test-id"}, data)
	})

	t.Run("wait for nack with error", func(t *testing.T) {
		waiter := supervisor.StartNackWaiter("test-subscription")

		go func() {
			time.Sleep(1 * time.Millisecond)

			supervisor.StopNackWaiter(
				subscription.NormalizeSubscriptionName("test-project", "test-subscription"),
				[]string{"test-id"},
				assert.AnError,
			)
		}()

		data, err := waiter.WaitMaxDuration(context.Background(), 5*time.Millisecond)
		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
		assert.Equal(t, []string{"test-id"}, data)
	})
}
