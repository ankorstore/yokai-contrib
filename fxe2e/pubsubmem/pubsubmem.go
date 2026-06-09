// Package pubsubmem captures published events for end-to-end assertions.
//
// To assert what a journey published (5_job.json) without a real PubSub broker,
// the harness needs to intercept publishing. This package decorates the shared
// go-modules gcppubsub.Publisher: every Publish is recorded (schema namespace +
// the JSON of the payload) and then delegated to the real publisher, so the
// genuine publish path keeps running.
//
// It works for any Ankorstore app on go-modules/gcppubsub, so the per-app
// recording-publisher boilerplate disappears. Wire it into a Boot:
//
//	sink, opt := pubsubmem.Recorder()
//	internal.RunE2ETest(tb, opt, ... )
//	return e2e.BootResult{ Events: sink.Recorded, ... }
package pubsubmem

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/go-modules/gcppubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/topic"
	"github.com/ankorstore/yokai-contrib/fxe2e/e2e"
	"github.com/ankorstore/yokai-contrib/fxe2e/recorder"
	"go.uber.org/fx"
)

// Recorder returns an event sink and an fx.Option that decorates the app's
// gcppubsub.Publisher to record every published event into the sink. Pass the
// sink's Recorded method to e2e.BootResult.Events.
func Recorder() (*recorder.Sink, fx.Option) {
	sink := &recorder.Sink{}

	opt := fx.Decorate(func(inner gcppubsub.Publisher) gcppubsub.Publisher {
		return &recordingPublisher{inner: inner, sink: sink}
	})

	return sink, opt
}

// recordingPublisher records each published event, then delegates to the real
// publisher. It implements gcppubsub.Publisher.
type recordingPublisher struct {
	inner gcppubsub.Publisher
	sink  *recorder.Sink
}

var _ gcppubsub.Publisher = (*recordingPublisher)(nil)

func (p *recordingPublisher) Publish(
	ctx context.Context,
	schemaNamespace string,
	data any,
	options ...topic.PublishOption,
) (*pubsub.PublishResult, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		payload = json.RawMessage("null")
	}

	p.sink.Record(e2e.RecordedEvent{
		Schema:  schemaNamespace,
		Payload: payload,
	})

	return p.inner.Publish(ctx, schemaNamespace, data, options...)
}

func (p *recordingPublisher) Stop() {
	p.inner.Stop()
}
