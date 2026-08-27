package pitch

import (
	"fmt"
	"math"

	"wtc/internal/store"
)

func (s *System) Ramp(targetDeg float64) (float64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.brake.Durable() {
		return s.angleDeg, ErrNoDurablePressure
	}
	if s.tower.Latched() {
		return s.angleDeg, ErrTowerLatched
	}
	if s.brake.Latched() {
		return s.angleDeg, ErrBrakeLatched
	}
	targetDeg = clampDeg(targetDeg, s.limits.MinPitchDeg, s.limits.MaxPitchDeg)
	steps := MoveSteps(s.angleDeg, targetDeg, s.limits.PitchRateDegPerSec)
	current := s.angleDeg
	for _, next := range steps {
		current = next
		s.angleDeg = current
		s.blade.SetAngle(current)
		_ = s.st.Put(store.KeyPitchAngle, angleRecord{AngleDeg: current})
	}
	s.audit.Append("pitch", "ramp", fmt.Sprintf("target=%.1f reached=%.1f", targetDeg, current))
	return current, nil
}

func MoveSteps(from, to, rate float64) []float64 {
	if rate <= 0 {
		return []float64{to}
	}
	steps := make([]float64, 0, 64)
	current := from
	step := rate
	if to < from {
		step = -rate
	}
	for current != to {
		next := current + step
		if step > 0 && next > to {
			next = to
		}
		if step < 0 && next < to {
			next = to
		}
		current = math.Round(next*10) / 10
		steps = append(steps, current)
	}
	return steps
}
