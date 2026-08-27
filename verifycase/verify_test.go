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

func TestWtcVibrationLatchReset(t *testing.T) {
	limits := ns.DefaultLimits()
	st := store.New(filepath.Join(t.TempDir(), "state"))
	recorder := audit.New(st, 200)
	bladeSys := blade.New(limits.FeatherDeg)
	brakeSys := brake.NewSystem(st, limits, recorder)
	towerSys := tower.NewSystem(st, limits)
	pitchSys := pitch.NewSystem(st, bladeSys, brakeSys, towerSys, recorder, limits)
	if err := brakeSys.SetPressure(limits.MinBrakeBar); err != nil {
		t.Fatalf("set pressure: %v", err)
	}
	towerSys.Sample(limits.VibrationHighMM + 10)
	if !towerSys.Latched() {
		t.Fatal("vibration latch must engage on high loads")
	}
	towerSys.Sample(limits.VibrationRecoverMM - 1)
	if towerSys.Latched() {
		t.Fatal("vibration latch must reset when loads recover")
	}
	if !pitchSys.Ready() {
		t.Fatal("pitch must be ready after vibration recovery")
	}
}
