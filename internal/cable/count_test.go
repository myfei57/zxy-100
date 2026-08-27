package cable

import (
	"path/filepath"
	"testing"

	"wtc/internal/ns"
	"wtc/internal/store"
)

func newTestSystem(t *testing.T) *System {
	t.Helper()
	st := store.New(filepath.Join(t.TempDir(), "cable"))
	return NewSystem(st, ns.DefaultLimits())
}

// 缠绕计数到达限位（等于上限）必须立即拉响报警，而不是超过上限之后才拉响，
// 否则偏航会在到达限位后再多转一圈才被扭缆保护拦下。
func TestAlarmRaisedAtLimitNotBeyond(t *testing.T) {
	system := newTestSystem(t)
	limit := ns.DefaultLimits().TwistLimitTurns

	// 每圈 360°，正向缠绕到上限-1 圈：未到限位，不应报警。
	for i := 0; i < limit-1; i++ {
		if err := system.Account(360); err != nil {
			t.Fatalf("account below limit: %v", err)
		}
	}
	if system.Alarm() {
		t.Fatalf("alarm must be clear below the limit, wraps=%d", system.Wraps())
	}
	if system.AtLimit() {
		t.Fatalf("must not be at limit below the limit, wraps=%d", system.Wraps())
	}

	// 再缠绕一圈到达上限：Alarm 与 AtLimit 必须同时为真。
	if err := system.Account(360); err != nil {
		t.Fatalf("account to limit: %v", err)
	}
	if system.Wraps() != limit {
		t.Fatalf("wraps must reach the limit, got %d want %d", system.Wraps(), limit)
	}
	if !system.AtLimit() {
		t.Fatalf("must be at limit when wraps==%d", limit)
	}
	if !system.Alarm() {
		t.Fatalf("alarm must be raised at the limit (wraps==%d), not only beyond it", limit)
	}
}

// 对称：负向缠绕到达下限同样应在限位即报警。
func TestAlarmRaisedAtNegativeLimit(t *testing.T) {
	system := newTestSystem(t)
	limit := ns.DefaultLimits().TwistLimitTurns

	if err := system.Account(-360 * float64(limit)); err != nil {
		t.Fatalf("account to negative limit: %v", err)
	}
	if system.Wraps() != -limit {
		t.Fatalf("wraps must reach the negative limit, got %d want %d", system.Wraps(), -limit)
	}
	if !system.AtLimit() {
		t.Fatalf("must be at limit when wraps==-%d", limit)
	}
	if !system.Alarm() {
		t.Fatalf("alarm must be raised at the negative limit (wraps==-%d)", limit)
	}
}

// 解缆回绕后，只要计数仍处限位，报警就不得解除——同方向偏航不得继续执行。
// 只有回到安全范围（绝对值低于上限）后才解除报警。
func TestRewindKeepsAlarmUntilSafe(t *testing.T) {
	system := newTestSystem(t)
	limit := ns.DefaultLimits().TwistLimitTurns

	// 缠绕到 limit+1（历史上旧逻辑需超过上限才拉响，此处模拟越限场景）。
	if err := system.Account(360 * float64(limit+1)); err != nil {
		t.Fatalf("account beyond limit: %v", err)
	}
	if !system.Alarm() {
		t.Fatal("alarm must be raised beyond the limit")
	}

	// 回绕一步：计数从 limit+1 回到 limit，仍处限位，报警必须保持。
	if err := system.Rewind(); err != nil {
		t.Fatalf("rewind: %v", err)
	}
	if system.Wraps() != limit {
		t.Fatalf("wraps must be back at the limit, got %d want %d", system.Wraps(), limit)
	}
	if !system.AtLimit() {
		t.Fatalf("must still be at limit after one rewind, wraps=%d", system.Wraps())
	}
	if !system.Alarm() {
		t.Fatalf("alarm must persist while wraps==%d (still at limit)", limit)
	}

	// 再回绕一步：计数回到 limit-1，进入安全范围，报警解除。
	if err := system.Rewind(); err != nil {
		t.Fatalf("rewind to safe: %v", err)
	}
	if system.AtLimit() {
		t.Fatalf("must not be at limit in the safe range, wraps=%d", system.Wraps())
	}
	if system.Alarm() {
		t.Fatalf("alarm must clear once wraps back in safe range, wraps=%d", system.Wraps())
	}
}

// 持久化记录里的 alarm 字段不可信：重启后必须按 wraps 与限位重新推导。
// 模拟落盘一条 wraps==limit 但 alarm=false 的旧记录，重新加载后必须报警。
func TestNewSystemDerivesAlarmFromWraps(t *testing.T) {
	dir := t.TempDir()
	st := store.New(filepath.Join(dir, "cable"))
	limit := ns.DefaultLimits().TwistLimitTurns

	// 写入一份旧版本的、不一致的持久化记录：计数已到限位，但 alarm=false。
	if err := st.Put(store.KeyCableWraps, wrapsRecord{Wraps: limit, Alarm: false}); err != nil {
		t.Fatalf("seed stale record: %v", err)
	}

	system := NewSystem(st, ns.DefaultLimits())
	if system.Wraps() != limit {
		t.Fatalf("wraps must reload from store, got %d want %d", system.Wraps(), limit)
	}
	if !system.AtLimit() {
		t.Fatalf("reloaded system must be at limit, wraps=%d", limit)
	}
	if !system.Alarm() {
		t.Fatalf("alarm must be re-derived from wraps after reload (wraps==%d), not trusted from store", limit)
	}
}

// AtLimit 与 Alarm 在所有缠绕计数下必须保持一致，二者不得再出现阈值错位。
func TestAtLimitAndAlarmInvariant(t *testing.T) {
	system := newTestSystem(t)
	for _, wraps := range []int{-5, -3, -2, -1, 0, 1, 2, 3, 5} {
		// 直接写底层字段以覆盖各计数取值。
		system.mu.Lock()
		system.wraps = wraps
		system.alarm = absInt(system.wraps) >= ns.DefaultLimits().TwistLimitTurns
		system.mu.Unlock()

		if system.AtLimit() != system.Alarm() {
			t.Fatalf("AtLimit/Alarm mismatch at wraps=%d: atlimit=%v alarm=%v",
				wraps, system.AtLimit(), system.Alarm())
		}
	}
}
