package conv

import (
	"math"
	"path/filepath"
	"testing"

	"wtc/internal/ns"
	"wtc/internal/store"
)

func TestCurveInterpolation(t *testing.T) {
	st := store.New(filepath.Join(t.TempDir(), "curve"))
	curve := NewCurve(st, []Point{
		{SpeedMS: 3, PowerKW: 0},
		{SpeedMS: 12, PowerKW: 3000},
		{SpeedMS: 25, PowerKW: 3000},
	})
	if got := curve.PowerAt(12); got != 3000 {
		t.Fatalf("rated point mismatch: %v", got)
	}
	if got := curve.PowerAt(25); got != 3000 {
		t.Fatalf("cut-out point mismatch: %v", got)
	}
	mid := curve.PowerAt(7.5)
	if math.Abs(mid-1500) > 0.01 {
		t.Fatalf("mid interpolation mismatch: %v", mid)
	}
	if curve.Version() != 0 {
		t.Fatalf("fresh curve must start at version 0")
	}
}

func TestCurveRecalibratePersists(t *testing.T) {
	root := filepath.Join(t.TempDir(), "curve")
	st := store.New(root)
	limits := ns.DefaultLimits()
	curve := NewCurve(st, DefaultCurve(limits))
	points := []Point{
		{SpeedMS: 3, PowerKW: 0},
		{SpeedMS: 10, PowerKW: 800},
		{SpeedMS: 25, PowerKW: 1500},
	}
	if err := curve.Recalibrate(points); err != nil {
		t.Fatalf("recalibrate: %v", err)
	}
	reloaded := NewCurve(st, DefaultCurve(limits))
	if got := reloaded.Version(); got != 1 {
		t.Fatalf("reloaded version mismatch: %d", got)
	}
	if got := reloaded.PowerAt(25); got != 1500 {
		t.Fatalf("reloaded curve mismatch: %v", got)
	}
}
