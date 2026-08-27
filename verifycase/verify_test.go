package verifycase

import (
	"path/filepath"
	"testing"

	"wtc/internal/audit"
	"wtc/internal/cable"
	"wtc/internal/ns"
	"wtc/internal/store"
	"wtc/internal/yaw"
)

func TestWtcYawUntwistAtLimit(t *testing.T) {
	limits := ns.DefaultLimits()
	st := store.New(filepath.Join(t.TempDir(), "state"))
	recorder := audit.New(st, 200)
	cableSys := cable.NewSystem(st, limits)
	yawSys := yaw.NewSystem(st, cableSys, recorder, limits)
	if err := yawSys.Turn(360); err != nil {
		t.Fatalf("first turn: %v", err)
	}
	if err := yawSys.Turn(360); err != nil {
		t.Fatalf("second turn: %v", err)
	}
	if got := cableSys.Wraps(); got != limits.TwistLimitTurns {
		t.Fatalf("wraps at the twist limit: %d", got)
	}
	if err := yawSys.Turn(360); err == nil {
		t.Fatal("yaw must pause at the cable twist limit")
	}
	if got := cableSys.Wraps(); got != limits.TwistLimitTurns {
		t.Fatalf("wraps changed past the limit: %d", got)
	}
	if err := yawSys.Untwist(); err != nil {
		t.Fatalf("untwist: %v", err)
	}
	if got := cableSys.Wraps(); got != limits.TwistLimitTurns-1 {
		t.Fatalf("wraps after untwist: %d", got)
	}
}
