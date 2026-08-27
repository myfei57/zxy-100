package cable

import (
	"math"

	"wtc/internal/store"
)

func (s *System) AtLimit() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return absInt(s.wraps) >= s.limits.TwistLimitTurns
}

func (s *System) Account(deltaDeg float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	turns := int(math.Round(math.Abs(deltaDeg) / 360))
	if deltaDeg < 0 {
		turns = -turns
	}
	s.wraps += turns
	// 报警阈值必须与 AtLimit 一致：缠绕计数到达上限即拉响，而不是
	// 超过上限之后才拉响，否则偏航会在到达限位后继续多转一圈。
	s.alarm = absInt(s.wraps) >= s.limits.TwistLimitTurns
	return s.st.Put(store.KeyCableWraps, wrapsRecord{Wraps: s.wraps, Alarm: s.alarm})
}

func (s *System) Rewind() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.wraps > 0 {
		s.wraps--
	}
	if s.wraps < 0 {
		s.wraps++
	}
	// 解缆后报警是否解除取决于缠绕是否已回到安全范围，而不是一刀清零，
	// 否则只要回绕一步、计数仍处限位，同方向偏航就会绕过保护继续执行。
	s.alarm = absInt(s.wraps) >= s.limits.TwistLimitTurns
	return s.st.Put(store.KeyCableWraps, wrapsRecord{Wraps: s.wraps, Alarm: s.alarm})
}

func (s *System) Direction() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.wraps > 0 {
		return "cw"
	}
	if s.wraps < 0 {
		return "ccw"
	}
	return "none"
}

func absInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
