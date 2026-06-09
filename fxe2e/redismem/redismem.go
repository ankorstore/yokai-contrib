// Package redismem provides an in-process, in-memory Redis for end-to-end tests.
//
// Yokai's fxredis module, in test mode, hands back a go-redis/redismock client —
// an expectation-based mock where every command must be scripted in advance.
// That forces each application to bypass its real Redis-backed code with
// hand-written stubs. This package removes that need: it starts a real (but
// in-memory) Redis server via miniredis and points the application's
// *redis.Client at it, so the genuine Redis-backed code (locks, leases, OAuth
// state, rate limiters, caches, ...) runs unchanged against an isolated store
// with no external Redis and no per-test command scripting.
//
// Wire it into a Boot function alongside the other fx options:
//
//	internal.RunE2ETest(tb,
//	    redismem.FxOption(tb),
//	    // ... other options
//	)
//
// The server is bound to tb's lifecycle and shut down automatically when the
// test ends. Each test gets a fresh, empty server.
package redismem

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

// Server starts a fresh in-process miniredis bound to tb's lifecycle and returns
// it together with an fx.Option that decorates the application's *redis.Client to
// talk to it. The returned server can be used to seed or inspect state directly
// (e.g. mr.Set(...), mr.Keys()).
func Server(tb testing.TB) (*miniredis.Miniredis, fx.Option) {
	tb.Helper()

	mr := miniredis.RunT(tb)

	opt := fx.Decorate(func(*redis.Client) *redis.Client {
		return redis.NewClient(&redis.Options{Addr: mr.Addr()})
	})

	return mr, opt
}

// FxOption is the common case: start an in-memory Redis and return only the
// fx.Option that points the application's *redis.Client at it.
func FxOption(tb testing.TB) fx.Option {
	tb.Helper()

	_, opt := Server(tb)

	return opt
}
