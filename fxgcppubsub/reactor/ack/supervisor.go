package ack

import (
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/reactor"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/subscription"
	"github.com/ankorstore/yokai/config"
)

var _ AckSupervisor = (*DefaultAckSupervisor)(nil)

type AckSupervisor interface {
	StartAckWaiter(subscriptionID string) *reactor.Waiter
	StopAckWaiter(subscriptionName string, ackIDs []string, err error)
}

type DefaultAckSupervisor struct {
	supervisor reactor.WaiterSupervisor
	config     *config.Config
}

func NewDefaultAckSupervisor(supervisor reactor.WaiterSupervisor, config *config.Config) *DefaultAckSupervisor {
	return &DefaultAckSupervisor{
		supervisor: supervisor,
		config:     config,
	}
}

func (s *DefaultAckSupervisor) StartAckWaiter(subscriptionID string) *reactor.Waiter {
	subscriptionName := subscription.NormalizeSubscriptionName(
		s.config.GetString("modules.gcppubsub.project.id"),
		subscriptionID,
	)

	return s.supervisor.StartWaiter(subscriptionName)
}

func (s *DefaultAckSupervisor) StopAckWaiter(subscriptionName string, ackIDs []string, err error) {
	s.supervisor.StopWaiter(subscriptionName, ackIDs, err)
}
