package pitch

type PitchReport struct {
	AngleDeg    float64 `json:"angle_deg"`
	MinReached  bool    `json:"min_reached"`
	Feathered   bool    `json:"feathered"`
	Ready       bool    `json:"ready"`
	SetpointDeg float64 `json:"setpoint_deg"`
}

func (s *System) Report() PitchReport {
	angle, _, _ := s.Setpoint().Current()
	return PitchReport{
		AngleDeg:    s.Angle(),
		MinReached:  s.MinReached(),
		Feathered:   s.blade.Feathered(),
		Ready:       s.Ready(),
		SetpointDeg: angle,
	}
}
