package yaw

import (
	"errors"
	"path/filepath"
	"testing"

	"wtc/internal/audit"
	"wtc/internal/cable"
	"wtc/internal/ns"
	"wtc/internal/store"
)

func newTestYaw(t *testing.T) (*System, *cable.System) {
	t.Helper()
	st := store.New(filepath.Join(t.TempDir(), "yaw"))
	limits := ns.DefaultLimits()
	cableSys := cable.NewSystem(st, limits)
	recorder := audit.New(st, 200)
	return NewSystem(st, cableSys, recorder, limits), cableSys
}

// 缠绕计数到达上限后，同方向偏航必须被扭缆保护拦截，不得继续转动。
// 这条不变量对应现场故障：大风天连续同向偏航，计数到限位后偏航仍在转，
// 直到最后扭缆保护跳闸、机组停机一整天。
func TestTurnBlockedAtCableLimit(t *testing.T) {
	system, cableSys := newTestYaw(t)
	limit := ns.DefaultLimits().TwistLimitTurns

	// 连续正向偏航，每次一整圈，把缠绕计数累到上限。
	for i := 0; i < limit; i++ {
		if err := system.Turn(system.Angle() + 360); err != nil {
			t.Fatalf("turn %d below limit must succeed: %v", i+1, err)
		}
	}
	if !cableSys.AtLimit() {
		t.Fatalf("cable must be at limit after %d wraps, got wraps=%d", limit, cableSys.Wraps())
	}
	if !cableSys.Alarm() {
		t.Fatalf("twist alarm must be raised at the limit (wraps==%d)", limit)
	}

	// 再向同方向偏航：必须被保护拦下，且角度、计数不再变化。
	beforeAngle := system.Angle()
	beforeWraps := cableSys.Wraps()
	err := system.Turn(system.Angle() + 360)
	if !errors.Is(err, ErrTwistLimit) {
		t.Fatalf("same-direction turn at the limit must return ErrTwistLimit, got %v", err)
	}
	if system.Angle() != beforeAngle {
		t.Fatalf("angle must not advance past the limit: got %v want %v", system.Angle(), beforeAngle)
	}
	if cableSys.Wraps() != beforeWraps {
		t.Fatalf("wraps must not increase past the limit: got %d want %d", cableSys.Wraps(), beforeWraps)
	}
}

// 解缆回绕到安全范围后，偏航才允许恢复执行。
func TestTurnAllowedAfterUntwistToSafe(t *testing.T) {
	system, cableSys := newTestYaw(t)
	limit := ns.DefaultLimits().TwistLimitTurns

	// 缠绕到限位并触发保护。
	for i := 0; i < limit; i++ {
		if err := system.Turn(system.Angle() + 360); err != nil {
			t.Fatalf("turn below limit: %v", err)
		}
	}
	if err := system.Turn(system.Angle() + 360); !errors.Is(err, ErrTwistLimit) {
		t.Fatalf("turn at limit must be blocked: %v", err)
	}

	// 回绕一步后仍处限位时，保护不得提前解除。
	if err := system.Untwist(); err != nil {
		t.Fatalf("untwist: %v", err)
	}
	if cableSys.AtLimit() {
		t.Fatalf("after one untwist wraps must drop below limit, got wraps=%d", cableSys.Wraps())
	}
	// 回到安全范围后，同方向偏航重新可用。
	if err := system.Turn(system.Angle() + 360); err != nil {
		t.Fatalf("turn must be allowed once back in safe range: %v", err)
	}
}
