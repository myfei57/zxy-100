package brake

func (s *System) Charge() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.bar >= s.limits.ChargeTargetBar {
		s.charged = true
		return nil
	}
	s.bar = s.limits.ChargeTargetBar
	s.charged = true
	s.audit.Append("brake", "charge", "accumulator charged")
	return nil
}

func (s *System) Charged() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.charged
}
