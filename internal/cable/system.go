package cable

import (
	"sync"

	"wtc/internal/ns"
	"wtc/internal/store"
)

type wrapsRecord struct {
	Wraps int  `json:"wraps"`
	Alarm bool `json:"alarm"`
}

type System struct {
	st     *store.Store
	limits ns.Limits
	mu     sync.Mutex
	wraps  int
	alarm  bool
}

func NewSystem(st *store.Store, limits ns.Limits) *System {
	system := &System{st: st, limits: limits}
	var record wrapsRecord
	if err := st.Get(store.KeyCableWraps, &record); err == nil {
		system.wraps = record.Wraps
		system.alarm = record.Alarm
	}
	return system
}

func (s *System) Wraps() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.wraps
}

func (s *System) Alarm() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.alarm
}
