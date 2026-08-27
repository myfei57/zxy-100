package wind

import (
	"wtc/internal/ns"
)

func (s *System) LoadSetpoint(speedMS float64) ns.LoadTarget {
	power := s.curve.PowerAt(speedMS)
	if power > s.limits.RatedPowerKW {
		power = s.limits.RatedPowerKW
	}
	if power < 0 {
		power = 0
	}
	return ns.LoadTarget{PowerKW: power, CurveVersion: s.curve.Version()}
}
