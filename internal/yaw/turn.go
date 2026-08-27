package yaw

import (
	"math"

	"wtc/internal/store"
)

func (s *System) Turn(targetDeg float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delta := targetDeg - s.angleDeg
	if math.Abs(delta) < s.limits.YawDeadbandDeg {
		return nil
	}
	if s.cable.Alarm() {
		return ErrTwistLimit
	}
	if err := s.cable.Account(delta); err != nil {
		return err
	}
	s.angleDeg = normalizeAngle(targetDeg)
	_ = s.st.Put(store.KeyYawAngle, angleRecord{AngleDeg: s.angleDeg})
	s.audit.Append("yaw", "turn", "nacelle rotated")
	return nil
}

func (s *System) Untwist() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.cable.Rewind(); err != nil {
		return err
	}
	s.audit.Append("yaw", "untwist", "cable unwound")
	return nil
}

func (s *System) ApplyTarget(targetDeg float64) error {
	delta := targetDeg - s.Angle()
	if math.Abs(delta) < s.limits.YawDeadbandDeg {
		return nil
	}
	return s.Turn(targetDeg)
}

func (s *System) Report() YawReport {
	return YawReport{
		AngleDeg: s.Angle(),
		Wraps:    s.cable.Wraps(),
		Alarm:    s.cable.Alarm(),
	}
}

func normalizeAngle(deg float64) float64 {
	for deg >= 360 {
		deg -= 360
	}
	for deg < 0 {
		deg += 360
	}
	return deg
}
