package blade

import (
	"math"
	"sync"
)

type Blade struct {
	mu        sync.RWMutex
	angleDeg  float64
	feathered bool
}

func New(initialDeg float64) *Blade {
	return &Blade{angleDeg: clampDeg(initialDeg)}
}

func (b *Blade) SetAngle(deg float64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.angleDeg = clampDeg(deg)
	b.feathered = false
}

func (b *Blade) Feather() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.angleDeg = 90
	b.feathered = true
}

func (b *Blade) Feathered() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.feathered
}

func (b *Blade) Snapshot() Snapshot {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return Snapshot{AngleDeg: b.angleDeg, Feathered: b.feathered}
}

func clampDeg(deg float64) float64 {
	if deg < 0 {
		return 0
	}
	if deg > 90 {
		return 90
	}
	return math.Round(deg*10) / 10
}
