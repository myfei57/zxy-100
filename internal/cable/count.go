package cable

import (
	"math"

	"wtc/internal/store"
)

func (s *System) AtLimit() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return absInt(s.wraps) >= s.limits.TwistLimitTurns
}

func (s *System) Account(deltaDeg float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	turns := int(math.Round(math.Abs(deltaDeg) / 360))
	if deltaDeg < 0 {
		turns = -turns
	}
	s.wraps += turns
	if s.wraps > s.limits.TwistLimitTurns {
		s.alarm = true
	}
	if s.wraps < -s.limits.TwistLimitTurns {
		s.alarm = true
	}
	return s.st.Put(store.KeyCableWraps, wrapsRecord{Wraps: s.wraps, Alarm: s.alarm})
}

func (s *System) Rewind() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.wraps > 0 {
		s.wraps--
	}
	if s.wraps < 0 {
		s.wraps++
	}
	s.alarm = false
	return s.st.Put(store.KeyCableWraps, wrapsRecord{Wraps: s.wraps, Alarm: s.alarm})
}

func (s *System) Direction() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.wraps > 0 {
		return "cw"
	}
	if s.wraps < 0 {
		return "ccw"
	}
	return "none"
}

func absInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
