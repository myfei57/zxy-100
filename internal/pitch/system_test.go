package pitch

import (
	"errors"
	"path/filepath"
	"testing"

	"wtc/internal/audit"
	"wtc/internal/blade"
	"wtc/internal/brake"
	"wtc/internal/ns"
	"wtc/internal/store"
	"wtc/internal/tower"
)

func newTestSystem(t *testing.T) (*System, *brake.System) {
	t.Helper()
	st := store.New(filepath.Join(t.TempDir(), "pitch"))
	limits := ns.DefaultLimits()
	recorder := audit.New(st, 200)
	bladeSys := blade.New(limits.FeatherDeg)
	brakeSys := brake.NewSystem(st, limits, recorder)
	towerSys := tower.NewSystem(st, limits)
	return NewSystem(st, bladeSys, brakeSys, towerSys, recorder, limits), brakeSys
}

// TestMoveRejectedBeforeDurablePressure locks the startup interlock: the
// blades must not move until the brake pressure has been durably saved.
func TestMoveRejectedBeforeDurablePressure(t *testing.T) {
	sys, brakeSys := newTestSystem(t)
	if brakeSys.Durable() {
		t.Fatal("brake must not be durable before pressure is saved")
	}
	if err := sys.Move(15); !errors.Is(err, ErrNoDurablePressure) {
		t.Fatalf("move before durable pressure must be rejected with ErrNoDurablePressure, got %v", err)
	}
	if got := sys.Angle(); got != 0 {
		t.Fatalf("angle must not change after a rejected move, got %v", got)
	}
	// Once the pressure is durably saved, the interlock lifts.
	if err := brakeSys.SetPressure(120); err != nil {
		t.Fatalf("set pressure: %v", err)
	}
	if err := sys.Move(15); err != nil {
		t.Fatalf("move after durable pressure must succeed, got %v", err)
	}
	if got := sys.Angle(); got != 15 {
		t.Fatalf("angle must reflect the accepted move, got %v", got)
	}
}

// TestDemandRejectedBeforeDurablePressure ensures Demand enforces the same
// interlock even though it separately charges the accumulator.
func TestDemandRejectedBeforeDurablePressure(t *testing.T) {
	sys, brakeSys := newTestSystem(t)
	if err := brakeSys.Charge(); err != nil {
		t.Fatalf("charge: %v", err)
	}
	if err := sys.Demand(15); !errors.Is(err, ErrNoDurablePressure) {
		t.Fatalf("demand after charge but before durable pressure must be rejected with ErrNoDurablePressure, got %v", err)
	}
	if got := sys.Angle(); got != 0 {
		t.Fatalf("angle must not change after a rejected demand, got %v", got)
	}
	if err := brakeSys.SetPressure(120); err != nil {
		t.Fatalf("set pressure: %v", err)
	}
	if err := sys.Demand(15); err != nil {
		t.Fatalf("demand after durable pressure must succeed, got %v", err)
	}
	if got := sys.Angle(); got != 15 {
		t.Fatalf("angle must reflect the accepted demand, got %v", got)
	}
}

// TestDurablePressureSurvivesRestart reproduces the reported reboot sequence:
// pressure is saved, the controller restarts, and the interlock must stay
// cleared so a subsequent move is allowed — but only because save finished.
func TestDurablePressureSurvivesRestart(t *testing.T) {
	root := filepath.Join(t.TempDir(), "pitch")
	st := store.New(root)
	limits := ns.DefaultLimits()
	recorder := audit.New(st, 200)
	bladeSys := blade.New(limits.FeatherDeg)
	brakeSys := brake.NewSystem(st, limits, recorder)
	towerSys := tower.NewSystem(st, limits)
	sys := NewSystem(st, bladeSys, brakeSys, towerSys, recorder, limits)
	if err := brakeSys.SetPressure(120); err != nil {
		t.Fatalf("set pressure: %v", err)
	}
	_ = sys

	// Rebuild every subsystem from the same store, simulating a restart.
	reloadedBrake := brake.NewSystem(st, limits, recorder)
	reloaded := NewSystem(st, bladeSys, reloadedBrake, towerSys, recorder, limits)
	if !reloadedBrake.Durable() {
		t.Fatal("durable pressure must survive restart when the save completed")
	}
	if err := reloaded.Move(20); err != nil {
		t.Fatalf("move after restart with durable pressure must succeed, got %v", err)
	}
}

// TestMoveRejectedAfterRestartWithoutSave reproduces the failure: the
// controller restarted before the pressure save finished, so the interlock
// must stay engaged and reject the move.
func TestMoveRejectedAfterRestartWithoutSave(t *testing.T) {
	root := filepath.Join(t.TempDir(), "pitch")
	st := store.New(root)
	limits := ns.DefaultLimits()
	recorder := audit.New(st, 200)
	bladeSys := blade.New(limits.FeatherDeg)
	brakeSys := brake.NewSystem(st, limits, recorder)
	towerSys := tower.NewSystem(st, limits)
	sys := NewSystem(st, bladeSys, brakeSys, towerSys, recorder, limits)
	// Charge the accumulator but never durably save the pressure.
	_ = brakeSys.Charge()
	_ = sys // silence unused

	// Restart: nothing was persisted for pressure, so durable stays false.
	reloadedBrake := brake.NewSystem(st, limits, recorder)
	reloaded := NewSystem(st, bladeSys, reloadedBrake, towerSys, recorder, limits)
	if reloadedBrake.Durable() {
		t.Fatal("durable must be false when the pressure save never completed")
	}
	if err := reloaded.Move(20); !errors.Is(err, ErrNoDurablePressure) {
		t.Fatalf("move after restart without durable pressure must be rejected, got %v", err)
	}
}
