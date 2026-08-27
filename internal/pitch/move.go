package pitch

import (
	"fmt"

	"wtc/internal/store"
)

func (s *System) Move(deg float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// 启动联锁：制动压力未落盘完成前禁止变桨，避免机组重启时桨叶在压力建立前动作。
	if !s.brake.Durable() {
		return ErrNoDurablePressure
	}
	if s.tower.Latched() {
		return ErrTowerLatched
	}
	if s.brake.Latched() {
		return ErrBrakeLatched
	}
	deg = clampDeg(deg, s.limits.MinPitchDeg, s.limits.MaxPitchDeg)
	s.angleDeg = deg
	s.blade.SetAngle(deg)
	_ = s.st.Put(store.KeyPitchAngle, angleRecord{AngleDeg: deg})
	s.audit.Append("pitch", "move", fmt.Sprintf("angle=%.1f", deg))
	return nil
}
