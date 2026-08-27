package brake

import (
	"sync"

	"wtc/internal/audit"
	"wtc/internal/ns"
	"wtc/internal/store"
)

type pressureRecord struct {
	Bar     float64 `json:"bar"`
	Durable bool    `json:"durable"`
}

type latchRecord struct {
	Engaged bool `json:"engaged"`
}

type System struct {
	st      *store.Store
	limits  ns.Limits
	audit   *audit.Recorder
	mu      sync.Mutex
	bar     float64
	durable bool
	charged bool
	latch   bool
}

func NewSystem(st *store.Store, limits ns.Limits, recorder *audit.Recorder) *System {
	system := &System{st: st, limits: limits, audit: recorder}
	var pressure pressureRecord
	if err := st.Get(store.KeyBrakePressure, &pressure); err == nil {
		system.bar = pressure.Bar
		system.durable = pressure.Durable
	}
	var latch latchRecord
	if err := st.Get(store.KeyBrakeLatch, &latch); err == nil {
		system.latch = latch.Engaged
	}
	return system
}
