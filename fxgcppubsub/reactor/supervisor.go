package reactor

import (
	"sync"
)

var _ WaiterSupervisor = (*DefaultWaiterSupervisor)(nil)

// WaiterSupervisor is the interface for waiters supervisors.
type WaiterSupervisor interface {
	StartWaiter(target string) *Waiter
	StopWaiter(target string, data any, err error)
}

// DefaultWaiterSupervisor is the default WaiterSupervisor implementation.
type DefaultWaiterSupervisor struct {
	waiters map[string]*Waiter
	mutex   sync.RWMutex
}

// NewDefaultWaiterSupervisor returns a new DefaultWaiterSupervisor instance.
func NewDefaultWaiterSupervisor() *DefaultWaiterSupervisor {
	return &DefaultWaiterSupervisor{
		waiters: make(map[string]*Waiter),
	}
}

// StartWaiter starts a Waiter for a target.
func (s *DefaultWaiterSupervisor) StartWaiter(target string) *Waiter {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	waiter := NewWaiter()

	s.waiters[target] = waiter

	return waiter
}

// StopWaiter stops a Waiter for a target with result.
func (s *DefaultWaiterSupervisor) StopWaiter(target string, data any, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if waiter, found := s.waiters[target]; found {
		waiter.Stop(data, err)
	}
}
