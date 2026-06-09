// Package recorder provides a concurrency-safe sink for capturing published
// events during an end-to-end run.
//
// A publisher decorator writes each event into a Sink, and the application's
// Boot returns Sink.Recorded as e2e.BootResult.Events. For apps on the shared
// gcppubsub.Publisher the pubsubmem helper provides that decorator (and the
// Sink); only an app on a different publisher needs to write its own.
package recorder

import (
	"sync"

	"github.com/ankorstore/yokai-contrib/fxe2e/e2e"
)

// Sink collects events recorded during a run. The zero value is ready to use.
type Sink struct {
	mu     sync.Mutex
	events []e2e.RecordedEvent
}

// Record appends a captured event.
func (s *Sink) Record(e e2e.RecordedEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = append(s.events, e)
}

// Recorded returns a copy of every event captured so far. Its signature matches
// e2e.BootResult.Events, so it can be passed there directly.
func (s *Sink) Recorded() []e2e.RecordedEvent {
	s.mu.Lock()
	defer s.mu.Unlock()

	return append([]e2e.RecordedEvent(nil), s.events...)
}
