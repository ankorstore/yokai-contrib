package ack

import (
	"fmt"

	"github.com/ankorstore/yokai-contrib/fxgcppubsub/reactor"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/subscription"
	"github.com/ankorstore/yokai/config"
)

const (
	Ack  = "ack"
	Nack = "nack"
)

var _ AckSupervisor = (*DefaultAckSupervisor)(nil)

// AckSupervisor is a reactor supervisor that reacts to acks ans nacks.
type AckSupervisor interface {
	StartAckWaiter(subscriptionID string) *reactor.Waiter
	StopAckWaiter(subscriptionName string, ackIDs []string, err error)
	StartNackWaiter(subscriptionID string) *reactor.Waiter
	StopNackWaiter(subscriptionName string, ackIDs []string, err error)
}

// DefaultAckSupervisor is the default AckSupervisor implementation.
type DefaultAckSupervisor struct {
	supervisor reactor.WaiterSupervisor
	config     *config.Config
}

// NewDefaultAckSupervisor returns a new DefaultAckSupervisor instance.
func NewDefaultAckSupervisor(supervisor reactor.WaiterSupervisor, config *config.Config) *DefaultAckSupervisor {
	return &DefaultAckSupervisor{
		supervisor: supervisor,
		config:     config,
	}
}

// StartAckWaiter starts an ack waiter on a provided subscriptionID.
func (s *DefaultAckSupervisor) StartAckWaiter(subscriptionID string) *reactor.Waiter {
	return s.startWaiter(subscriptionID, Ack)
}

// StopAckWaiter stop an ack waiter for a provided subscriptionName.
func (s *DefaultAckSupervisor) StopAckWaiter(subscriptionName string, ackIDs []string, err error) {
	s.stopWaiter(subscriptionName, Ack, ackIDs, err)
}

// StartNackWaiter starts a nack waiter on a provided subscriptionID.
func (s *DefaultAckSupervisor) StartNackWaiter(subscriptionID string) *reactor.Waiter {
	return s.startWaiter(subscriptionID, Nack)
}

// StopNackWaiter stop a nack waiter for a provided subscriptionName.
func (s *DefaultAckSupervisor) StopNackWaiter(subscriptionName string, ackIDs []string, err error) {
	s.stopWaiter(subscriptionName, Nack, ackIDs, err)
}

func (s *DefaultAckSupervisor) startWaiter(subscriptionID string, kind string) *reactor.Waiter {
	subscriptionName := subscription.NormalizeSubscriptionName(
		s.config.GetString("modules.gcppubsub.project.id"),
		subscriptionID,
	)

	return s.supervisor.StartWaiter(fmt.Sprintf("%s::%s", kind, subscriptionName))
}

func (s *DefaultAckSupervisor) stopWaiter(subscriptionName string, kind string, ackIDs []string, err error) {
	s.supervisor.StopWaiter(fmt.Sprintf("%s::%s", kind, subscriptionName), ackIDs, err)
}
