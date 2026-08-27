package brake

import (
	"wtc/internal/store"
)

func (s *System) SetPressure(bar float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	record := pressureRecord{Bar: bar, Durable: true}
	if err := s.st.Put(store.KeyBrakePressure, record); err != nil {
		return err
	}
	s.bar = bar
	s.durable = true
	s.charged = bar >= s.limits.ChargeTargetBar
	return nil
}

func (s *System) Durable() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.durable
}

func (s *System) Bar() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.bar
}
