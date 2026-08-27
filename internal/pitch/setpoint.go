package pitch

import (
	"sync"
)

type setpointState struct {
	angle  float64
	seq    uint64
	source string
}

type Setpoint struct {
	mu    sync.Mutex
	state setpointState
}

func NewSetpoint(initialDeg float64) *Setpoint {
	return &Setpoint{state: setpointState{angle: initialDeg}}
}

func (sp *Setpoint) Apply(deg float64, source string) (uint64, float64) {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	sp.state.seq++
	sp.state.angle = deg
	sp.state.source = source
	return sp.state.seq, sp.state.angle
}

func (sp *Setpoint) Ack(seq uint64, deg float64) bool {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	return sp.state.seq == seq && sp.state.angle == deg
}

func (sp *Setpoint) Current() (float64, uint64, string) {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	return sp.state.angle, sp.state.seq, sp.state.source
}
