package main

import (
	"os"

	"wtc/internal/audit"
	"wtc/internal/blade"
	"wtc/internal/brake"
	"wtc/internal/cable"
	"wtc/internal/console"
	"wtc/internal/conv"
	"wtc/internal/grid"
	"wtc/internal/ns"
	"wtc/internal/pitch"
	"wtc/internal/store"
	"wtc/internal/tower"
	"wtc/internal/wind"
	"wtc/internal/yaw"
)

func BuildServer(cfg Config) (*console.Server, error) {
	if err := os.MkdirAll(cfg.DataDir, 0o755); err != nil {
		return nil, err
	}
	st := store.New(cfg.DataDir)
	limits := ns.DefaultLimits()
	if err := limits.Validate(); err != nil {
		return nil, err
	}
	recorder := audit.New(st, cfg.AuditBufferSize)
	bladeSys := blade.New(limits.FeatherDeg)
	brakeSys := brake.NewSystem(st, limits, recorder)
	towerSys := tower.NewSystem(st, limits)
	cableSys := cable.NewSystem(st, limits)
	pitchSys := pitch.NewSystem(st, bladeSys, brakeSys, towerSys, recorder, limits)
	curve := conv.NewCurve(st, conv.DefaultCurve(limits))
	setpoint := pitch.NewSetpoint(limits.FeatherDeg)
	convSys := conv.NewSystem(st, pitchSys, recorder, limits, curve, setpoint)
	drive := yaw.NewDrive(0)
	yawSys := yaw.NewSystem(st, cableSys, recorder, limits)
	windSys := wind.NewSystem(drive, curve, limits, cfg.SampleWindow)
	gridSys := grid.NewSystem(convSys, recorder, 18)
	return console.NewServer(cfg.Addr, st, limits, recorder, bladeSys, brakeSys, pitchSys, yawSys, windSys, convSys, towerSys, cableSys, gridSys, drive), nil
}
