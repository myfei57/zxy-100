package pitch

func (s *System) Angle() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.angleDeg
}

func (s *System) MinReached() bool {
	return s.Angle() <= s.limits.MinPitchDeg+0.01
}
