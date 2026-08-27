package ns

type WindSample struct {
	SpeedMS      float64 `json:"speed_ms"`
	DirectionDeg float64 `json:"direction_deg"`
	Sequence     uint64  `json:"sequence"`
}

type LoadTarget struct {
	PowerKW      float64 `json:"power_kw"`
	CurveVersion int     `json:"curve_version"`
}

type Status struct {
	State           string  `json:"state"`
	PitchDeg        float64 `json:"pitch_deg"`
	PitchSetpoint   float64 `json:"pitch_setpoint"`
	Feathered       bool    `json:"feathered"`
	YawDeg          float64 `json:"yaw_deg"`
	YawTarget       float64 `json:"yaw_target"`
	BrakeBar        float64 `json:"brake_bar"`
	BrakeDurable    bool    `json:"brake_durable"`
	BrakeCharged    bool    `json:"brake_charged"`
	BreakerClosed   bool    `json:"breaker_closed"`
	SpeedRPM        float64 `json:"speed_rpm"`
	RotorRadPerSec  float64 `json:"rotor_rad_per_sec"`
	CableWraps      int     `json:"cable_wraps"`
	CableDirection  string  `json:"cable_direction"`
	TwistAlarm      bool    `json:"twist_alarm"`
	VibrationMM     float64 `json:"vibration_mm"`
	VibrationLatch  bool    `json:"vibration_latch"`
	TowerAlarmLevel string  `json:"tower_alarm_level"`
	GridOnline      bool    `json:"grid_online"`
	VoltageKV       float64 `json:"voltage_kv"`
	UnderVoltage    bool    `json:"under_voltage"`
	WindMeanMS      float64 `json:"wind_mean_ms"`
	WindGustMS      float64 `json:"wind_gust_ms"`
	WindDirDeg      float64 `json:"wind_dir_deg"`
	CurveVersion    int     `json:"curve_version"`
	LastLimitKW     float64 `json:"last_limit_kw"`
	AuditCount      int64   `json:"audit_count"`
	Ready           bool    `json:"ready"`
}
