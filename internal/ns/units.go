package ns

import (
	"math"
)

func SpeedKMH(speedMS float64) float64 {
	return speedMS * 3.6
}

func RadPerSec(rpm float64) float64 {
	return rpm * 0.10471975511965977
}

func WindClass(speedMS float64) string {
	switch {
	case speedMS < 3:
		return "calm"
	case speedMS < 12:
		return "below-rated"
	case speedMS <= 25:
		return "rated"
	default:
		return "storm"
	}
}

func DegreesToRadians(deg float64) float64 {
	return deg * math.Pi / 180
}
