package verifycase

import (
	"path/filepath"
	"testing"

	"wtc/internal/conv"
	"wtc/internal/ns"
	"wtc/internal/store"
	"wtc/internal/wind"
	"wtc/internal/yaw"
)

func TestWtcPowerCurveFresh(t *testing.T) {
	limits := ns.DefaultLimits()
	st := store.New(filepath.Join(t.TempDir(), "state"))
	curve := conv.NewCurve(st, conv.DefaultCurve(limits))
	drive := yaw.NewDrive(0)
	windSys := wind.NewSystem(drive, curve, limits, 10)
	old := windSys.LoadSetpoint(12).PowerKW
	points := []conv.Point{
		{SpeedMS: limits.CutInWindMS, PowerKW: 0},
		{SpeedMS: 10, PowerKW: 800},
		{SpeedMS: limits.CutOutWindMS, PowerKW: 1500},
	}
	if err := curve.Recalibrate(points); err != nil {
		t.Fatalf("recalibrate: %v", err)
	}
	fresh := windSys.LoadSetpoint(12).PowerKW
	if fresh >= old {
		t.Fatal("load setpoint must use the recalibrated curve")
	}
	if curve.Version() != 1 {
		t.Fatalf("curve version must advance, got %d", curve.Version())
	}
}
