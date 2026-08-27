package yaw

import (
	"errors"
	"sync"

	"wtc/internal/audit"
	"wtc/internal/cable"
	"wtc/internal/ns"
	"wtc/internal/store"
)

var ErrTwistLimit = errors.New("yaw paused at the cable twist limit")

type angleRecord struct {
	AngleDeg float64 `json:"angle_deg"`
}

type System struct {
	st       *store.Store
	cable    *cable.System
	audit    *audit.Recorder
	limits   ns.Limits
	mu       sync.Mutex
	angleDeg float64
}

func NewSystem(st *store.Store, cableSys *cable.System, recorder *audit.Recorder, limits ns.Limits) *System {
	system := &System{st: st, cable: cableSys, audit: recorder, limits: limits}
	var record angleRecord
	if err := st.Get(store.KeyYawAngle, &record); err == nil {
		system.angleDeg = record.AngleDeg
	}
	return system
}

func (s *System) Angle() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.angleDeg
}
