// Package await provides poll-until helpers for effect-based waits. The engine
// uses Until + Count to implement a case's "await" spec (poll a COUNT(*) query
// until it reaches a target); the helpers are also usable directly from a custom
// Boot.
package await

import (
	"database/sql"
	"testing"
	"time"
)

// Until polls cond until it returns true or timeout elapses, sleeping interval
// between attempts. A cond error or a timeout fails the test.
func Until(tb testing.TB, timeout, interval time.Duration, desc string, cond func() (bool, error)) {
	tb.Helper()

	deadline := time.Now().Add(timeout)
	for {
		done, err := cond()
		if err != nil {
			tb.Fatalf("await: checking %s: %v", desc, err)
		}
		if done {
			return
		}
		if time.Now().After(deadline) {
			tb.Fatalf("await: timed out after %s waiting for %s", timeout, desc)
		}
		time.Sleep(interval)
	}
}

// Count builds a condition satisfied once the given COUNT(*) query returns want
// — e.g. Count(db, 0, "SELECT COUNT(*) FROM jobs WHERE status = 'running'").
func Count(db *sql.DB, want int, query string, args ...any) func() (bool, error) {
	return func() (bool, error) {
		var n int
		if err := db.QueryRow(query, args...).Scan(&n); err != nil {
			return false, err
		}

		return n == want, nil
	}
}

// NoRows is Count with a target of zero.
func NoRows(db *sql.DB, query string, args ...any) func() (bool, error) {
	return Count(db, 0, query, args...)
}
