package wind

import (
	"fmt"
	"math"
	"sync"

	"wtc/internal/conv"
	"wtc/internal/ns"
	"wtc/internal/yaw"
)

type System struct {
	drive   *yaw.Drive
	curve   *conv.Curve
	limits  ns.Limits
	mu      sync.Mutex
	last    ns.WindSample
	history []ns.WindSample
	window  int
}

func NewSystem(drive *yaw.Drive, curve *conv.Curve, limits ns.Limits, window int) *System {
	if window < 1 {
		window = 1
	}
	return &System{drive: drive, curve: curve, limits: limits, window: window}
}

func (s *System) Measure(speedMS, directionDeg float64) ns.WindSample {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.last.Sequence++
	s.last.SpeedMS = speedMS
	s.last.DirectionDeg = directionDeg
	s.history = append(s.history, s.last)
	if len(s.history) > s.window {
		s.history = s.history[len(s.history)-s.window:]
	}
	return s.last
}

func (s *System) LastSample() ns.WindSample {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.last
}

func (s *System) MeanSpeed() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.history) == 0 {
		return 0
	}
	total := 0.0
	for _, sample := range s.history {
		total += sample.SpeedMS
	}
	return total / float64(len(s.history))
}

func (s *System) SampleCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.history)
}

func (s *System) GustSpeedMS() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	maxSpeed := 0.0
	for _, sample := range s.history {
		if sample.SpeedMS > maxSpeed {
			maxSpeed = sample.SpeedMS
		}
	}
	return maxSpeed
}

func (s *System) MeanDirection() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.history) == 0 {
		return 0
	}
	sinSum := 0.0
	cosSum := 0.0
	for _, sample := range s.history {
		rad := ns.DegreesToRadians(sample.DirectionDeg)
		sinSum += math.Sin(rad)
		cosSum += math.Cos(rad)
	}
	deg := math.Atan2(sinSum, cosSum) * 180 / math.Pi
	if deg < 0 {
		deg += 360
	}
	return deg
}

func (s *System) DirectionHistogram() map[string]int {
	s.mu.Lock()
	defer s.mu.Unlock()
	bins := make(map[string]int)
	if len(s.history) == 0 {
		return bins
	}
	for _, sample := range s.history {
		bin := int(sample.DirectionDeg/22.5) % 16
		key := fmt.Sprintf("bin%02d", bin)
		bins[key]++
	}
	return bins
}

func (s *System) Report() WindReport {
	return WindReport{
		MeanSpeedMS:      s.MeanSpeed(),
		GustSpeedMS:      s.GustSpeedMS(),
		MeanDirectionDeg: s.MeanDirection(),
		SampleCount:      s.SampleCount(),
		AboveCutIn:       s.AboveCutIn(),
		SpeedVarianceMS:  s.SpeedVariance(),
		TrendSlope:       s.TrendSlope(),
		Histogram:        s.DirectionHistogram(),
	}
}

func (s *System) AboveCutIn() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	count := 0
	for _, sample := range s.history {
		if sample.SpeedMS >= s.limits.CutInWindMS {
			count++
		}
	}
	return count
}

func (s *System) SpeedVariance() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.history) < 2 {
		return 0
	}
	mean := 0.0
	for _, sample := range s.history {
		mean += sample.SpeedMS
	}
	mean /= float64(len(s.history))
	total := 0.0
	for _, sample := range s.history {
		delta := sample.SpeedMS - mean
		total += delta * delta
	}
	return total / float64(len(s.history))
}

func (s *System) TrendSlope() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	count := len(s.history)
	if count < 2 {
		return 0
	}
	sumX := 0.0
	sumY := 0.0
	sumXY := 0.0
	sumXX := 0.0
	for index, sample := range s.history {
		x := float64(index)
		sumX += x
		sumY += sample.SpeedMS
		sumXY += x * sample.SpeedMS
		sumXX += x * x
	}
	denominator := float64(count)*sumXX - sumX*sumX
	if denominator == 0 {
		return 0
	}
	return (float64(count)*sumXY - sumX*sumY) / denominator
}
