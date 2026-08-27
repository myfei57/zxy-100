package grid

type GridReport struct {
	Online       bool    `json:"online"`
	VoltageKV    float64 `json:"voltage_kv"`
	UnderVoltage bool    `json:"under_voltage"`
}
