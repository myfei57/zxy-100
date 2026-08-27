package ns

import (
	"testing"
)

func TestDefaultLimitsWindow(t *testing.T) {
	limits := DefaultLimits()
	if limits.CutInWindMS >= limits.CutOutWindMS {
		t.Fatalf("cut-in must stay below cut-out: %v >= %v", limits.CutInWindMS, limits.CutOutWindMS)
	}
	if limits.MinPitchDeg >= limits.FeatherDeg {
		t.Fatal("minimum pitch must stay below the feather angle")
	}
	if !limits.InWindWindow(10) {
		t.Fatal("10 m/s must be inside the wind window")
	}
	if limits.InWindWindow(2) || limits.InWindWindow(30) {
		t.Fatal("out-of-window speeds must be rejected")
	}
}
