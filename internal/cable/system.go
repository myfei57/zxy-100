package cable

import (
	"sync"

	"wtc/internal/ns"
	"wtc/internal/store"
)

type wrapsRecord struct {
	Wraps int  `json:"wraps"`
	Alarm bool `json:"alarm"`
}

type System struct {
	st     *store.Store
	limits ns.Limits
	mu     sync.Mutex
	wraps  int
	alarm  bool
}

func NewSystem(st *store.Store, limits ns.Limits) *System {
	system := &System{st: st, limits: limits}
	var record wrapsRecord
	if err := st.Get(store.KeyCableWraps, &record); err == nil {
		system.wraps = record.Wraps
		// 报警由缠绕计数与限位推导而来，不信任落盘的 alarm 字段：
		// 旧版本在到达限位后才拉响、解缆一步就清零，重启后若计数仍处
		// 限位而 alarm 字段为 false，扭缆保护会继续被绕过。
		system.alarm = absInt(system.wraps) >= limits.TwistLimitTurns
	}
	return system
}

func (s *System) Wraps() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.wraps
}

func (s *System) Alarm() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.alarm
}
