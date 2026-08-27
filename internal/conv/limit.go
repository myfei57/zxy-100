package conv

func (s *System) ApplyLimit(powerKW float64) (uint64, float64) {
	angle := s.angleForPower(powerKW)
	seq, ack := s.setpoint.Apply(angle, "limiter")
	s.mu.Lock()
	s.lastLimit = powerKW
	s.limitHistory = append(s.limitHistory, powerKW)
	if len(s.limitHistory) > 32 {
		s.limitHistory = s.limitHistory[len(s.limitHistory)-32:]
	}
	s.mu.Unlock()
	return seq, ack
}

func (s *System) LastLimit() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.lastLimit
}

func (s *System) LimitHistory() []float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]float64(nil), s.limitHistory...)
}

func (s *System) angleForPower(powerKW float64) float64 {
	ratio := powerKW / s.limits.RatedPowerKW
	if ratio >= 1 {
		return s.limits.MinPitchDeg
	}
	if ratio <= 0 {
		return s.limits.MaxPitchDeg
	}
	return s.limits.MaxPitchDeg - ratio*(s.limits.MaxPitchDeg-s.limits.MinPitchDeg)
}
