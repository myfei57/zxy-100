package conv

import (
	"sync"

	"wtc/internal/ns"
	"wtc/internal/store"
)

type Point struct {
	SpeedMS float64 `json:"speed_ms"`
	PowerKW float64 `json:"power_kw"`
}

type curveRecord struct {
	Points  []Point `json:"points"`
	Version int     `json:"version"`
}

type Curve struct {
	st      *store.Store
	mu      sync.RWMutex
	points  []Point
	version int
}

func NewCurve(st *store.Store, fallback []Point) *Curve {
	curve := &Curve{st: st, points: append([]Point(nil), fallback...)}
	var record curveRecord
	if err := st.Get(store.KeyPowerCurve, &record); err == nil && len(record.Points) >= 2 {
		curve.points = record.Points
		curve.version = record.Version
	}
	return curve
}

func DefaultCurve(limits ns.Limits) []Point {
	return []Point{
		{SpeedMS: limits.CutInWindMS, PowerKW: 0},
		{SpeedMS: 12, PowerKW: limits.RatedPowerKW},
		{SpeedMS: limits.CutOutWindMS, PowerKW: limits.RatedPowerKW},
	}
}

func (c *Curve) Recalibrate(points []Point) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(points) < 2 {
		return ErrCurveTooShort
	}
	c.points = append([]Point(nil), points...)
	c.version++
	_ = c.st.Put(store.KeyPowerCurve, curveRecord{Points: c.points, Version: c.version})
	return nil
}

func (c *Curve) PowerAt(speedMS float64) float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return interpolate(c.points, speedMS)
}

func (c *Curve) Points() []Point {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return append([]Point(nil), c.points...)
}

func (c *Curve) Version() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.version
}

func interpolate(points []Point, speedMS float64) float64 {
	if len(points) == 0 {
		return 0
	}
	if speedMS <= points[0].SpeedMS {
		return points[0].PowerKW
	}
	last := points[len(points)-1]
	if speedMS >= last.SpeedMS {
		return last.PowerKW
	}
	for index := 1; index < len(points); index++ {
		prev := points[index-1]
		next := points[index]
		if speedMS <= next.SpeedMS {
			span := next.SpeedMS - prev.SpeedMS
			if span <= 0 {
				return next.PowerKW
			}
			ratio := (speedMS - prev.SpeedMS) / span
			return prev.PowerKW + ratio*(next.PowerKW-prev.PowerKW)
		}
	}
	return last.PowerKW
}
