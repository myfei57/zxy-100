package verifycase

import (
	"path/filepath"
	"testing"

	"wtc/internal/audit"
	"wtc/internal/blade"
	"wtc/internal/brake"
	"wtc/internal/conv"
	"wtc/internal/ns"
	"wtc/internal/pitch"
	"wtc/internal/store"
	"wtc/internal/tower"
)

func TestWtcCutInAfterPitchMin(t *testing.T) {
	limits := ns.DefaultLimits()
	st := store.New(filepath.Join(t.TempDir(), "state"))
	recorder := audit.New(st, 200)
	bladeSys := blade.New(limits.FeatherDeg)
	brakeSys := brake.NewSystem(st, limits, recorder)
	towerSys := tower.NewSystem(st, limits)
	pitchSys := pitch.NewSystem(st, bladeSys, brakeSys, towerSys, recorder, limits)
	curve := conv.NewCurve(st, conv.DefaultCurve(limits))
	setpoint := pitch.NewSetpoint(limits.FeatherDeg)
	convSys := conv.NewSystem(st, pitchSys, recorder, limits, curve, setpoint)
	if err := brakeSys.SetPressure(limits.MinBrakeBar); err != nil {
		t.Fatalf("set pressure: %v", err)
	}
	if err := pitchSys.Move(45); err != nil {
		t.Fatalf("move to 45: %v", err)
	}
	if err := convSys.Close(); err == nil {
		t.Fatal("breaker must not close before the pitch reaches the minimum")
	}
	if convSys.Closed() {
		t.Fatal("breaker must stay open when the pitch is above the minimum")
	}
	if err := pitchSys.Move(limits.MinPitchDeg); err != nil {
		t.Fatalf("move to minimum: %v", err)
	}
	if err := convSys.Close(); err != nil {
		t.Fatalf("close after pitch minimum: %v", err)
	}
}
