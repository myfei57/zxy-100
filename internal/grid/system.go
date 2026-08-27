package grid

import (
	"fmt"
	"sync"

	"wtc/internal/audit"
	"wtc/internal/conv"
)

type System struct {
	conv         *conv.System
	audit        *audit.Recorder
	mu           sync.Mutex
	online       bool
	minVoltageKV float64
	voltageKV    float64
	underVoltage bool
}

func NewSystem(convSys *conv.System, recorder *audit.Recorder, minVoltageKV float64) *System {
	return &System{conv: convSys, audit: recorder, online: true, minVoltageKV: minVoltageKV}
}

func (s *System) Loss() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.online = false
	s.audit.Append("grid", "loss", "grid supply lost")
	return s.conv.Trip()
}

func (s *System) Online() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.online
}

func (s *System) SampleVoltage(kv float64) error {
	s.mu.Lock()
	wasOnline := s.online
	s.voltageKV = kv
	s.underVoltage = kv < s.minVoltageKV
	s.mu.Unlock()
	if wasOnline && s.underVoltage {
		return s.Loss()
	}
	s.audit.Append("grid", "voltage", fmt.Sprintf("kv=%.1f", kv))
	return nil
}

func (s *System) VoltageKV() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.voltageKV
}

func (s *System) UnderVoltage() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.underVoltage
}

func (s *System) Report() GridReport {
	return GridReport{
		Online:       s.Online(),
		VoltageKV:    s.VoltageKV(),
		UnderVoltage: s.UnderVoltage(),
	}
}
