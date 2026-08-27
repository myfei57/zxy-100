package pitch

func (s *System) Ready() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.brake.Latched() {
		return false
	}
	if s.tower.Latched() {
		return false
	}
	return s.brake.Durable()
}
