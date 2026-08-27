package pitch

import (
	"errors"
	"sync"

	"wtc/internal/audit"
	"wtc/internal/blade"
	"wtc/internal/brake"
	"wtc/internal/ns"
	"wtc/internal/store"
	"wtc/internal/tower"
)

var (
	ErrNoDurablePressure = errors.New("brake pressure is not durable")
	ErrTowerLatched      = errors.New("tower vibration latch is engaged")
	ErrBrakeLatched      = errors.New("brake latch is engaged")
)

type angleRecord struct {
	AngleDeg float64 `json:"angle_deg"`
}

type System struct {
	st       *store.Store
	blade    *blade.Blade
	brake    *brake.System
	tower    *tower.System
	audit    *audit.Recorder
	limits   ns.Limits
	mu       sync.Mutex
	angleDeg float64
	setpoint *Setpoint
}

func NewSystem(st *store.Store, bladeSys *blade.Blade, brakeSys *brake.System, towerSys *tower.System, recorder *audit.Recorder, limits ns.Limits) *System {
	system := &System{
		st:       st,
		blade:    bladeSys,
		brake:    brakeSys,
		tower:    towerSys,
		audit:    recorder,
		limits:   limits,
		setpoint: NewSetpoint(limits.FeatherDeg),
	}
	var record angleRecord
	if err := st.Get(store.KeyPitchAngle, &record); err == nil {
		system.angleDeg = record.AngleDeg
		system.blade.SetAngle(record.AngleDeg)
	}
	return system
}

func (s *System) Setpoint() *Setpoint {
	return s.setpoint
}

func clampDeg(deg, min, max float64) float64 {
	if deg < min {
		return min
	}
	if deg > max {
		return max
	}
	return deg
}
