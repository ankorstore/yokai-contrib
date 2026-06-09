package pubsubmem

import (
	"context"
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/topic"
	"github.com/ankorstore/yokai-contrib/fxe2e/recorder"
	"github.com/onsi/gomega"
)

// fakePublisher is a stand-in gcppubsub.Publisher that records the calls it
// receives, so a test can assert the decorator delegated to it.
type fakePublisher struct {
	calls   int
	schemas []string
	stopped bool
}

func (f *fakePublisher) Publish(_ context.Context, schemaNamespace string, _ any, _ ...topic.PublishOption) (*pubsub.PublishResult, error) {
	f.calls++
	f.schemas = append(f.schemas, schemaNamespace)

	return nil, nil
}

func (f *fakePublisher) Stop() { f.stopped = true }

func TestRecordingPublisher_RecordsAndDelegates(t *testing.T) {
	g := gomega.NewWithT(t)

	inner := &fakePublisher{}
	sink := &recorder.Sink{}
	pub := &recordingPublisher{inner: inner, sink: sink}

	type orderPlaced struct {
		OrderUUID  string `json:"orderUuid"`
		TotalCents int    `json:"totalCents"`
	}
	_, err := pub.Publish(context.Background(), "order.placed.v1",
		orderPlaced{OrderUUID: "abc", TotalCents: 3750})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// Recorded for assertion: schema namespace + JSON of the payload.
	recorded := sink.Recorded()
	g.Expect(recorded).To(gomega.HaveLen(1))
	g.Expect(recorded[0].Schema).To(gomega.Equal("order.placed.v1"))
	g.Expect(string(recorded[0].Payload)).To(gomega.MatchJSON(`{"orderUuid":"abc","totalCents":3750}`))

	// Still delegated to the real publisher.
	g.Expect(inner.calls).To(gomega.Equal(1))
	g.Expect(inner.schemas).To(gomega.ConsistOf("order.placed.v1"))

	pub.Stop()
	g.Expect(inner.stopped).To(gomega.BeTrue())
}
