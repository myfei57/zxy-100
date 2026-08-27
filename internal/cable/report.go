package cable

type CableReport struct {
	Wraps     int    `json:"wraps"`
	Alarm     bool   `json:"alarm"`
	Direction string `json:"direction"`
	AtLimit   bool   `json:"at_limit"`
}

func (s *System) Report() CableReport {
	return CableReport{
		Wraps:     s.Wraps(),
		Alarm:     s.Alarm(),
		Direction: s.Direction(),
		AtLimit:   s.AtLimit(),
	}
}
