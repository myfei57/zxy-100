package tower

import (
	"sync"

	"wtc/internal/ns"
	"wtc/internal/store"
)

type latchRecord struct {
	Engaged bool `json:"engaged"`
}

type System struct {
	st     *store.Store
	limits ns.Limits
	mu     sync.Mutex
	latch  bool
	loadMM float64
}

func NewSystem(st *store.Store, limits ns.Limits) *System {
	system := &System{st: st, limits: limits}
	var latch latchRecord
	if err := st.Get(store.KeyTowerLatch, &latch); err == nil {
		system.latch = latch.Engaged
	}
	return system
}

func (s *System) Sample(mm float64) {
	s.mu.Lock()
	s.loadMM = mm
	s.mu.Unlock()
	s.checkTrip()
	s.checkRecovery()
}
