package pitch

func (s *System) Reset() error {
	if err := s.brake.ReleaseIfRecovered(); err != nil {
		return err
	}
	if err := s.tower.ClearIfRecovered(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.audit.Append("pitch", "reset", "startup interlocks cleared")
	return nil
}
