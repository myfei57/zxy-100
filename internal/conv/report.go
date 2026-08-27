package conv

type ConverterReport struct {
	Closed       bool    `json:"closed"`
	CurveVersion int     `json:"curve_version"`
	LastLimitKW  float64 `json:"last_limit_kw"`
}
