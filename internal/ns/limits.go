package ns

import (
	"errors"
)

type Limits struct {
	MinPitchDeg        float64
	MaxPitchDeg        float64
	FeatherDeg         float64
	PitchRateDegPerSec float64
	MinBrakeBar        float64
	ChargeTargetBar    float64
	TwistLimitTurns    int
	MinSpeedRPM        float64
	MaxSpeedRPM        float64
	CutInWindMS        float64
	CutOutWindMS       float64
	RatedPowerKW       float64
	VibrationHighMM    float64
	VibrationRecoverMM float64
	YawDeadbandDeg     float64
}

func DefaultLimits() Limits {
	return Limits{
		MinPitchDeg:        2,
		MaxPitchDeg:        90,
		FeatherDeg:         90,
		PitchRateDegPerSec: 6,
		MinBrakeBar:        80,
		ChargeTargetBar:    120,
		TwistLimitTurns:    2,
		MinSpeedRPM:        4,
		MaxSpeedRPM:        14,
		CutInWindMS:        3,
		CutOutWindMS:       25,
		RatedPowerKW:       3000,
		VibrationHighMM:    30,
		VibrationRecoverMM: 12,
		YawDeadbandDeg:     3,
	}
}

func (l Limits) InWindWindow(speedMS float64) bool {
	return speedMS >= l.CutInWindMS && speedMS <= l.CutOutWindMS
}

func (l Limits) Validate() error {
	if l.MinPitchDeg < 0 || l.MaxPitchDeg > 90 || l.MinPitchDeg >= l.MaxPitchDeg {
		return errors.New("pitch limits are invalid")
	}
	if l.FeatherDeg > 90 || l.FeatherDeg < l.MaxPitchDeg {
		return errors.New("feather angle is outside the allowed range")
	}
	if l.PitchRateDegPerSec <= 0 {
		return errors.New("pitch rate must be positive")
	}
	if l.MinBrakeBar <= 0 || l.ChargeTargetBar < l.MinBrakeBar {
		return errors.New("brake pressure limits are invalid")
	}
	if l.TwistLimitTurns < 1 {
		return errors.New("cable twist limit must be at least one turn")
	}
	if l.MinSpeedRPM <= 0 || l.MaxSpeedRPM <= l.MinSpeedRPM {
		return errors.New("rotor speed limits are invalid")
	}
	if l.CutInWindMS <= 0 || l.CutOutWindMS <= l.CutInWindMS {
		return errors.New("wind window limits are invalid")
	}
	if l.RatedPowerKW <= 0 {
		return errors.New("rated power must be positive")
	}
	if l.VibrationHighMM <= l.VibrationRecoverMM {
		return errors.New("vibration thresholds are invalid")
	}
	if l.YawDeadbandDeg < 0 {
		return errors.New("yaw deadband cannot be negative")
	}
	return nil
}
