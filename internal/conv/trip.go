package conv

import (
	"wtc/internal/store"
)

func (s *System) Trip() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = false
	_ = s.st.Put(store.KeyBreakerState, breakerRecord{Closed: false})
	s.audit.Append("conv", "breaker.open", "grid decoupled")
	return s.pitch.Feather()
}
