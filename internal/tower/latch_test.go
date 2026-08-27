package tower

import (
	"path/filepath"
	"testing"

	"wtc/internal/ns"
	"wtc/internal/store"
)

// TestLatchHysteresis reproduces the field scenario: a tower vibration trip
// forces the turbine to a latched standstill, and once the vibration settles
// below the recovery threshold the latch must clear on its own so pitch can
// re-engage instead of sitting in standby all day.
func TestLatchHysteresis(t *testing.T) {
	limits := ns.DefaultLimits()
	st := store.New(filepath.Join(t.TempDir(), "store"))
	s := NewSystem(st, limits)

	// Below the recovery threshold: nothing latched, normal operation.
	s.Sample(limits.VibrationRecoverMM - 1)
	if s.Latched() {
		t.Fatal("latch must not engage below the high threshold")
	}

	// Vibration exceeds the high threshold: trip and latch.
	s.Sample(limits.VibrationHighMM + 1)
	if !s.Latched() {
		t.Fatal("latch must engage above the high threshold")
	}
	if level := s.AlarmLevel(); level != "trip" {
		t.Fatalf("alarm level while latched = %q, want %q", level, "trip")
	}

	// Wind drops, vibration falls back into the warn band but not yet below
	// the recovery threshold: the latch must stay engaged (hysteresis guard
	// so the turbine does not chatter around the trip point).
	s.Sample(limits.VibrationRecoverMM + 1)
	if !s.Latched() {
		t.Fatal("latch must hold while above the recovery threshold")
	}

	// Vibration settles below the recovery threshold: the latch auto-clears,
	// allowing the pitch interlock to re-engage without a manual reset.
	s.Sample(limits.VibrationRecoverMM - 1)
	if s.Latched() {
		t.Fatal("latch must clear once vibration drops below the recovery threshold")
	}
	if level := s.AlarmLevel(); level != "normal" {
		t.Fatalf("alarm level after recovery = %q, want %q", level, "normal")
	}
}

// TestLatchPersistsRecovery ensures a cleared latch is persisted, so a restart
// after recovery does not resurrect the trip.
func TestLatchPersistsRecovery(t *testing.T) {
	limits := ns.DefaultLimits()
	st := store.New(filepath.Join(t.TempDir(), "store"))

	s := NewSystem(st, limits)
	s.Sample(limits.VibrationHighMM + 1)
	if !s.Latched() {
		t.Fatal("latch must engage above the high threshold")
	}
	s.Sample(limits.VibrationRecoverMM - 1)
	if s.Latched() {
		t.Fatal("latch must clear below the recovery threshold")
	}

	// Reconstruct the system from the same store: a recovered latch must not
	// come back as engaged.
	reloaded := NewSystem(st, limits)
	if reloaded.Latched() {
		t.Fatal("recovered latch must not resurrect across a restart")
	}
}
