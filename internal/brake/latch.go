package brake

import (
	"wtc/internal/store"
)

func (s *System) Engage() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.latch = true
	_ = s.st.Put(store.KeyBrakeLatch, latchRecord{Engaged: true})
}

func (s *System) Release() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.latch = false
	_ = s.st.Put(store.KeyBrakeLatch, latchRecord{Engaged: false})
}

func (s *System) Latched() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.latch
}

func (s *System) ReleaseIfRecovered() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return nil
}
