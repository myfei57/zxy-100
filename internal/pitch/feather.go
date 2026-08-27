package pitch

import (
	"wtc/internal/store"
)

func (s *System) Feather() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.angleDeg = s.limits.FeatherDeg
	s.blade.Feather()
	_ = s.st.Put(store.KeyPitchAngle, angleRecord{AngleDeg: s.angleDeg})
	s.audit.Append("pitch", "feather", "blades feathered")
	return nil
}
