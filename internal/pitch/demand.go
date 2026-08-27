package pitch

import (
	"fmt"

	"wtc/internal/store"
)

func (s *System) Demand(deg float64) error {
	if err := s.brake.Charge(); err != nil {
		return err
	}
	s.mu.Lock()
	deg = clampDeg(deg, s.limits.MinPitchDeg, s.limits.MaxPitchDeg)
	s.angleDeg = deg
	s.blade.SetAngle(deg)
	_ = s.st.Put(store.KeyPitchAngle, angleRecord{AngleDeg: deg})
	s.mu.Unlock()
	s.audit.Append("pitch", "demand", fmt.Sprintf("angle=%.1f", deg))
	return nil
}
