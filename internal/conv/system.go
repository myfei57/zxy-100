package conv

import (
	"errors"
	"sync"

	"wtc/internal/audit"
	"wtc/internal/ns"
	"wtc/internal/pitch"
	"wtc/internal/store"
)

var (
	ErrPitchNotMin   = errors.New("pitch is not at the minimum angle")
	ErrNotReady      = errors.New("turbine is not ready")
	ErrCurveTooShort = errors.New("power curve needs at least two points")
)

type breakerRecord struct {
	Closed bool `json:"closed"`
}

type System struct {
	st           *store.Store
	pitch        *pitch.System
	audit        *audit.Recorder
	limits       ns.Limits
	mu           sync.Mutex
	closed       bool
	setpoint     *pitch.Setpoint
	curve        *Curve
	lastLimit    float64
	limitHistory []float64
}

func NewSystem(st *store.Store, pitchSys *pitch.System, recorder *audit.Recorder, limits ns.Limits, curve *Curve, setpoint *pitch.Setpoint) *System {
	system := &System{
		st:       st,
		pitch:    pitchSys,
		audit:    recorder,
		limits:   limits,
		curve:    curve,
		setpoint: setpoint,
	}
	var record breakerRecord
	if err := st.Get(store.KeyBreakerState, &record); err == nil {
		system.closed = record.Closed
	}
	return system
}

func (s *System) CurveVersion() int {
	return s.curve.Version()
}

func (s *System) Recalibrate(points []Point) error {
	return s.curve.Recalibrate(points)
}

func (s *System) Points() []Point {
	return s.curve.Points()
}
