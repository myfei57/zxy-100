package conv

import (
	"wtc/internal/store"
)

func (s *System) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.pitch.MinReached() {
		return ErrPitchNotMin
	}
	if !s.pitch.Ready() {
		return ErrNotReady
	}
	s.closed = true
	_ = s.st.Put(store.KeyBreakerState, breakerRecord{Closed: true})
	s.audit.Append("conv", "breaker.close", "grid coupled")
	return nil
}

func (s *System) Open() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = false
	_ = s.st.Put(store.KeyBreakerState, breakerRecord{Closed: false})
	s.audit.Append("conv", "breaker.open", "grid decoupled")
	return nil
}

func (s *System) Closed() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.closed
}

func (s *System) Report() ConverterReport {
	return ConverterReport{
		Closed:       s.Closed(),
		CurveVersion: s.CurveVersion(),
		LastLimitKW:  s.LastLimit(),
	}
}
