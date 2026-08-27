package wind

type WindReport struct {
	MeanSpeedMS      float64        `json:"mean_speed_ms"`
	GustSpeedMS      float64        `json:"gust_speed_ms"`
	MeanDirectionDeg float64        `json:"mean_direction_deg"`
	SampleCount      int            `json:"sample_count"`
	AboveCutIn       int            `json:"above_cut_in"`
	SpeedVarianceMS  float64        `json:"speed_variance_ms"`
	TrendSlope       float64        `json:"trend_slope"`
	Histogram        map[string]int `json:"histogram"`
}
