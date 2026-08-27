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

func TestWtcBrakeLatchClear(t *testing.T) {
	limits := ns.DefaultLimits()
	st := store.New(filepath.Join(t.TempDir(), "state"))
	recorder := audit.New(st, 200)
	bladeSys := blade.New(limits.FeatherDeg)
	brakeSys := brake.NewSystem(st, limits, recorder)
	towerSys := tower.NewSystem(st, limits)
	pitchSys := pitch.NewSystem(st, bladeSys, brakeSys, towerSys, recorder, limits)
	brakeSys.Engage()
	if !brakeSys.Latched() {
		t.Fatal("brake latch must engage")
	}
	if err := brakeSys.SetPressure(limits.MinBrakeBar + 10); err != nil {
		t.Fatalf("set pressure: %v", err)
	}
	if err := pitchSys.Reset(); err != nil {
		t.Fatalf("reset: %v", err)
	}
	if brakeSys.Latched() {
		t.Fatal("brake latch must clear when the pressure recovers")
	}
	if !pitchSys.Ready() {
		t.Fatal("pitch must be ready after the brake latch clears")
	}
}
