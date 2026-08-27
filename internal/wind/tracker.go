package wind

import (
	"wtc/internal/ns"
)

func (s *System) Track(sample ns.WindSample) (uint64, float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.last = sample
	return s.drive.Request(sample.DirectionDeg, "tracker")
}
