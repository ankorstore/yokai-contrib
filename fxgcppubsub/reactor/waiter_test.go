package reactor_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ankorstore/yokai-contrib/fxgcppubsub/reactor"
	"github.com/stretchr/testify/assert"
)

func TestWaiter(t *testing.T) {
	t.Parallel()

	t.Run("wait until waiter completion", func(t *testing.T) {
		t.Parallel()

		waiter := reactor.NewWaiter()

		start := time.Now()

		go func(w *reactor.Waiter) {
			time.Sleep(1 * time.Millisecond)

			w.Stop("test", nil)
		}(waiter)

		data, err := waiter.Wait(context.Background())
		latency := time.Since(start)

		assert.NoError(t, err)
		assert.Equal(t, "test", data)
		assert.GreaterOrEqual(t, latency, 1*time.Millisecond)
	})

	t.Run("wait until waiter completion with error", func(t *testing.T) {
		t.Parallel()

		waiter := reactor.NewWaiter()

		start := time.Now()

		go func(w *reactor.Waiter) {
			time.Sleep(1 * time.Millisecond)

			w.Stop("test", fmt.Errorf("test error"))
		}(waiter)

		data, err := waiter.Wait(context.Background())
		latency := time.Since(start)

		assert.Error(t, err)
		assert.Equal(t, "test error", err.Error())
		assert.Equal(t, "test", data)
		assert.GreaterOrEqual(t, latency, 1*time.Millisecond)
	})

	t.Run("wait until context cancellation", func(t *testing.T) {
		t.Parallel()

		waiter := reactor.NewWaiter()

		start := time.Now()

		go func(w *reactor.Waiter) {
			time.Sleep(2 * time.Millisecond)

			w.Stop("test", nil)
		}(waiter)

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()

		data, err := waiter.Wait(ctx)
		latency := time.Since(start)

		assert.Error(t, err)
		assert.Nil(t, data)
		assert.Equal(t, "context deadline exceeded", err.Error())
		assert.GreaterOrEqual(t, latency, 1*time.Millisecond)
	})

	t.Run("waitMaxDuration until waiter completion", func(t *testing.T) {
		t.Parallel()

		waiter := reactor.NewWaiter()

		start := time.Now()

		go func(w *reactor.Waiter) {
			time.Sleep(1 * time.Millisecond)

			w.Stop("test", nil)
		}(waiter)

		data, err := waiter.WaitMaxDuration(context.Background(), 2*time.Millisecond)
		latency := time.Since(start)

		assert.NoError(t, err)
		assert.Equal(t, "test", data)
		assert.GreaterOrEqual(t, latency, 1*time.Millisecond)
	})

	t.Run("waitMaxDuration until max duration reached", func(t *testing.T) {
		t.Parallel()

		waiter := reactor.NewWaiter()

		start := time.Now()

		go func(w *reactor.Waiter) {
			time.Sleep(2 * time.Millisecond)

			w.Stop("test", nil)
		}(waiter)

		data, err := waiter.WaitMaxDuration(context.Background(), 1*time.Millisecond)
		latency := time.Since(start)

		assert.Error(t, err)
		assert.Nil(t, data)
		assert.Equal(t, "context deadline exceeded", err.Error())
		assert.GreaterOrEqual(t, latency, 1*time.Millisecond)
	})
}
