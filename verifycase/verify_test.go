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

func TestWtcGridTripFeatherOrder(t *testing.T) {
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
	if err := convSys.Trip(); err != nil {
		t.Fatalf("trip: %v", err)
	}
	if convSys.Closed() {
		t.Fatal("breaker must open on trip")
	}
	events := recorder.Recent(4)
	featherIndex := -1
	openIndex := -1
	for index, event := range events {
		if event.Action == "feather" {
			featherIndex = index
		}
		if event.Action == "breaker.open" {
			openIndex = index
		}
	}
	if featherIndex < 0 || openIndex < 0 {
		t.Fatalf("trip events missing: %v", events)
	}
	if featherIndex > openIndex {
		t.Fatal("breaker opened before the blades feathered")
	}
}
