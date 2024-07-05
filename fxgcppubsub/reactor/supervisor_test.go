package reactor_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ankorstore/yokai-contrib/fxgcppubsub/reactor"
	"github.com/stretchr/testify/assert"
)

func TestDefaultWaiterSupervisor(t *testing.T) {
	t.Parallel()

	t.Run("supervise", func(t *testing.T) {
		t.Parallel()

		supervisor := reactor.NewDefaultWaiterSupervisor()

		waiter := supervisor.StartWaiter("target")

		start := time.Now()

		go func(s reactor.WaiterSupervisor) {
			time.Sleep(1 * time.Millisecond)

			s.StopWaiter("target", "test", nil)
		}(supervisor)

		data, err := waiter.Wait(context.Background())
		latency := time.Since(start)

		assert.NoError(t, err)
		assert.Equal(t, "test", data)
		assert.GreaterOrEqual(t, latency, 1*time.Millisecond)
	})

	t.Run("supervise with error", func(t *testing.T) {
		t.Parallel()

		supervisor := reactor.NewDefaultWaiterSupervisor()

		waiter := supervisor.StartWaiter("target")

		start := time.Now()

		go func(s reactor.WaiterSupervisor) {
			time.Sleep(1 * time.Millisecond)

			s.StopWaiter("target", "test", fmt.Errorf("test error"))
		}(supervisor)

		data, err := waiter.Wait(context.Background())
		latency := time.Since(start)

		assert.Error(t, err)
		assert.Equal(t, "test error", err.Error())
		assert.Equal(t, "test", data)
		assert.GreaterOrEqual(t, latency, 1*time.Millisecond)
	})
}
