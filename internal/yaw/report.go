package yaw

type YawReport struct {
	AngleDeg float64 `json:"angle_deg"`
	Wraps    int     `json:"wraps"`
	Alarm    bool    `json:"alarm"`
}
