package tower

type TowerReport struct {
	LoadMM     float64 `json:"load_mm"`
	Latched    bool    `json:"latched"`
	AlarmLevel string  `json:"alarm_level"`
}

func (s *System) Report() TowerReport {
	return TowerReport{
		LoadMM:     s.LoadMM(),
		Latched:    s.Latched(),
		AlarmLevel: s.AlarmLevel(),
	}
}
