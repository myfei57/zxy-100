package verifycase

import (
	"path/filepath"
	"testing"

	"wtc/internal/audit"
	"wtc/internal/blade"
	"wtc/internal/brake"
	"wtc/internal/ns"
	"wtc/internal/pitch"
	"wtc/internal/store"
	"wtc/internal/tower"
)

func TestWtcPitchAfterBrakeDurable(t *testing.T) {
	limits := ns.DefaultLimits()
	st := store.New(filepath.Join(t.TempDir(), "state"))
	recorder := audit.New(st, 200)
	bladeSys := blade.New(limits.FeatherDeg)
	brakeSys := brake.NewSystem(st, limits, recorder)
	towerSys := tower.NewSystem(st, limits)
	pitchSys := pitch.NewSystem(st, bladeSys, brakeSys, towerSys, recorder, limits)
	if err := pitchSys.Move(30); err == nil {
		t.Fatal("pitch must refuse to move before the brake pressure is durable")
	}
	if err := brakeSys.SetPressure(limits.MinBrakeBar); err != nil {
		t.Fatalf("set pressure: %v", err)
	}
	if err := pitchSys.Move(30); err != nil {
		t.Fatalf("move after durable pressure: %v", err)
	}
	restarted := pitch.NewSystem(
		st,
		blade.New(limits.FeatherDeg),
		brake.NewSystem(st, limits, recorder),
		tower.NewSystem(st, limits),
		recorder,
		limits,
	)
	if got := restarted.Angle(); got != 30 {
		t.Fatalf("restart angle mismatch: %v", got)
	}
}
