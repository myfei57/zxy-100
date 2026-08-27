package yaw

import (
	"sync"
)

type driveState struct {
	target float64
	seq    uint64
	source string
}

type Drive struct {
	mu    sync.Mutex
	state driveState
}

func NewDrive(initialDeg float64) *Drive {
	return &Drive{state: driveState{target: initialDeg}}
}

func (d *Drive) Request(targetDeg float64, source string) (uint64, float64) {
	d.state.seq++
	d.state.target = targetDeg
	d.state.source = source
	return d.state.seq, d.state.target
}

func (d *Drive) Ack(seq uint64, targetDeg float64) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.state.seq == seq && d.state.target == targetDeg
}

func (d *Drive) Current() (float64, uint64, string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.state.target, d.state.seq, d.state.source
}
