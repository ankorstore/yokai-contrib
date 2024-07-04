package reactor

import (
	"sync"
)

// WaiterSupervisor is the supervisor for Waiter instances, affected to targets.
type WaiterSupervisor struct {
	waiters map[string]*Waiter
	mutex   sync.RWMutex
}

// NewWaiterSupervisor returns a new WaiterSupervisor instance.
func NewWaiterSupervisor() *WaiterSupervisor {
	return &WaiterSupervisor{
		waiters: make(map[string]*Waiter),
	}
}

// StartWaiter starts a Waiter for a target.
func (s *WaiterSupervisor) StartWaiter(target string) *Waiter {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	waiter := NewWaiter()

	s.waiters[target] = waiter

	return waiter
}

// StopWaiter stops a Waiter for a target with result.
func (s *WaiterSupervisor) StopWaiter(target string, data any, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if waiter, found := s.waiters[target]; found {
		waiter.Stop(data, err)
	}
}
