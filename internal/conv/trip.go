package conv

import (
	"wtc/internal/store"
)

func (s *System) Trip() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// 脱网保护顺序：桨叶完全顺桨之后方可断开电网，
	// 否则空载风轮在脱网瞬间会加速冲高转速。
	if err := s.pitch.Feather(); err != nil {
		return err
	}
	s.closed = false
	_ = s.st.Put(store.KeyBreakerState, breakerRecord{Closed: false})
	s.audit.Append("conv", "breaker.open", "grid decoupled")
	return nil
}
