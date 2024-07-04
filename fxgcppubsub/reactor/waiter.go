package reactor

import (
	"context"
	"time"
)

// Waiter is a component that can wait (blocking call) until stopped.
type Waiter struct {
	done chan struct{}
	data any
	err  error
}

// NewWaiter returns a new Waiter instance.
func NewWaiter() *Waiter {
	return &Waiter{
		done: make(chan struct{}),
	}
}

// Done returns the Waiter channel.
func (w *Waiter) Done() <-chan struct{} {
	return w.done
}

// Stop stops the waiter.
func (w *Waiter) Stop(data any, err error) {
	w.data = data
	w.err = err

	close(w.done)
}

// Wait is making a blocking call until the Stop func is called, or until the provided context is cancelled.
func (w *Waiter) Wait(ctx context.Context) (any, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-w.Done():
			return w.data, w.err
		}
	}
}

// WaitMaxDuration is making a blocking call until the Stop func is called, or until the provided context is cancelled, for a maximum duration.
func (w *Waiter) WaitMaxDuration(ctx context.Context, duration time.Duration) (any, error) {
	ctx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()

	return w.Wait(ctx)
}
