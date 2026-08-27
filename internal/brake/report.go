package brake

type BrakeReport struct {
	Bar     float64 `json:"bar"`
	Durable bool    `json:"durable"`
	Charged bool    `json:"charged"`
	Latched bool    `json:"latched"`
}

func (s *System) Report() BrakeReport {
	return BrakeReport{
		Bar:     s.Bar(),
		Durable: s.Durable(),
		Charged: s.Charged(),
		Latched: s.Latched(),
	}
}
