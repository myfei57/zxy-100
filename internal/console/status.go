package console

import (
	"net/http"

	"wtc/internal/ns"
)

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	angle, seq, source := s.pitch.Setpoint().Current()
	_ = seq
	_ = source
	yawTarget, _, _ := s.drive.Current()
	snapshot := s.blade.Snapshot()
	rpm := s.speed()
	status := ns.Status{
		State:           s.state(),
		PitchDeg:        s.pitch.Angle(),
		PitchSetpoint:   angle,
		Feathered:       snapshot.Feathered,
		YawDeg:          s.yaw.Angle(),
		YawTarget:       yawTarget,
		BrakeBar:        s.brake.Bar(),
		BrakeDurable:    s.brake.Durable(),
		BrakeCharged:    s.brake.Charged(),
		BreakerClosed:   s.conv.Closed(),
		SpeedRPM:        rpm,
		RotorRadPerSec:  ns.RadPerSec(rpm),
		CableWraps:      s.cable.Wraps(),
		CableDirection:  s.cable.Direction(),
		TwistAlarm:      s.cable.Alarm(),
		VibrationMM:     s.tower.LoadMM(),
		VibrationLatch:  s.tower.Latched(),
		TowerAlarmLevel: s.tower.AlarmLevel(),
		GridOnline:      s.grid.Online(),
		VoltageKV:       s.grid.VoltageKV(),
		UnderVoltage:    s.grid.UnderVoltage(),
		WindMeanMS:      s.wind.MeanSpeed(),
		WindGustMS:      s.wind.GustSpeedMS(),
		WindDirDeg:      s.wind.MeanDirection(),
		CurveVersion:    s.conv.CurveVersion(),
		LastLimitKW:     s.conv.LastLimit(),
		AuditCount:      s.audit.Count(),
		Ready:           s.pitch.Ready(),
	}
	writeJSON(w, http.StatusOK, status)
}

func (s *Server) state() string {
	switch {
	case !s.grid.Online():
		return "grid-loss"
	case s.conv.Closed():
		return "generating"
	case s.cable.Alarm():
		return "twist-alarm"
	case s.tower.Latched():
		return "vibration-latch"
	case s.pitch.MinReached():
		return "cut-in-ready"
	default:
		return "standby"
	}
}

func (s *Server) speed() float64 {
	mean := s.wind.MeanSpeed()
	if mean <= 0 {
		return 0
	}
	ratio := mean / 1.6
	rpm := ratio * s.limits.MaxSpeedRPM
	if rpm < s.limits.MinSpeedRPM {
		return 0
	}
	if rpm > s.limits.MaxSpeedRPM {
		return s.limits.MaxSpeedRPM
	}
	return rpm
}
