package ack

import (
	"cloud.google.com/go/pubsub/apiv1/pubsubpb"
)

// AckReactor is a pub/sub test server reactor for subscriptions message acks.
type AckReactor struct {
	supervisor AckSupervisor
}

// NewAckReactor returns a new AckReactor instance.
func NewAckReactor(supervisor AckSupervisor) *AckReactor {
	return &AckReactor{
		supervisor: supervisor,
	}
}

// FuncNames returns the list of function names this reactor will react to.
func (r *AckReactor) FuncNames() []string {
	return []string{
		"Acknowledge",
		"ModifyAckDeadline",
	}
}

// React is the reactor logic.
func (r *AckReactor) React(req any) (bool, any, error) {
	if ackReq, ok := req.(*pubsubpb.AcknowledgeRequest); ok {
		r.supervisor.StopAckWaiter(ackReq.Subscription, ackReq.AckIds, nil)
	}

	if modReq, ok := req.(*pubsubpb.ModifyAckDeadlineRequest); ok {
		r.supervisor.StopNackWaiter(modReq.Subscription, modReq.AckIds, nil)
	}

	return false, nil, nil
}
