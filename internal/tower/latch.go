package tower

import (
	"wtc/internal/store"
)

func (s *System) Latched() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.latch
}

func (s *System) LoadMM() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.loadMM
}

func (s *System) checkTrip() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.latch && s.loadMM >= s.limits.VibrationHighMM {
		s.latch = true
		_ = s.st.Put(store.KeyTowerLatch, latchRecord{Engaged: true})
	}
}

func (s *System) checkRecovery() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.latch && s.loadMM < s.limits.VibrationRecoverMM {
		s.latch = false
		_ = s.st.Put(store.KeyTowerLatch, latchRecord{Engaged: false})
	}
}

func (s *System) ClearIfRecovered() error {
	s.checkRecovery()
	return nil
}

func (s *System) AlarmLevel() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.latch {
		return "trip"
	}
	if s.loadMM >= s.limits.VibrationHighMM {
		return "high"
	}
	if s.loadMM >= s.limits.VibrationRecoverMM {
		return "warn"
	}
	return "normal"
}
