package console

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"wtc/internal/audit"
	"wtc/internal/conv"
	"wtc/internal/ns"
)

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func decodeJSON(r *http.Request, out any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(out)
}

type moveRequest struct {
	AngleDeg float64 `json:"angle_deg"`
}

type yawRequest struct {
	TargetDeg float64 `json:"target_deg"`
}

type pressureRequest struct {
	Bar float64 `json:"bar"`
}

type latchRequest struct {
	Engaged bool `json:"engaged"`
}

type windRequest struct {
	SpeedMS      float64 `json:"speed_ms"`
	DirectionDeg float64 `json:"direction_deg"`
}

type limitRequest struct {
	PowerKW float64 `json:"power_kw"`
}

type curveRequest struct {
	Points []conv.Point `json:"points"`
}

type vibrationRequest struct {
	LoadMM float64 `json:"load_mm"`
}

type voltageRequest struct {
	VoltageKV float64 `json:"voltage_kv"`
}

type pitchRampResponse struct {
	AngleDeg float64 `json:"angle_deg"`
	Steps    int     `json:"steps"`
}

type moveResponse struct {
	AngleDeg float64 `json:"angle_deg"`
	AckSeq   uint64  `json:"ack_seq"`
}

type yawTurnResponse struct {
	YawDeg float64 `json:"yaw_deg"`
	AckSeq uint64  `json:"ack_seq"`
}

type brakePressureResponse struct {
	Bar     float64 `json:"bar"`
	Durable bool    `json:"durable"`
}

type brakeChargeResponse struct {
	Bar     float64 `json:"bar"`
	Charged bool    `json:"charged"`
}

type brakeLatchResponse struct {
	Latched bool `json:"latched"`
}

type windSampleResponse struct {
	Sample    ns.WindSample `json:"sample"`
	SpeedKMH  float64       `json:"speed_kmh"`
	WindClass string        `json:"wind_class"`
	YawTarget float64       `json:"yaw_target"`
	AckSeq    uint64        `json:"ack_seq"`
}

type loadTargetResponse struct {
	PowerKW      float64 `json:"power_kw"`
	CurveVersion int     `json:"curve_version"`
	WindClass    string  `json:"wind_class"`
}

type limitResponse struct {
	AngleDeg    float64   `json:"angle_deg"`
	AckSeq      uint64    `json:"ack_seq"`
	LastLimitKW float64   `json:"last_limit_kw"`
	History     []float64 `json:"history"`
}

type tripResponse struct {
	Online bool `json:"online"`
	Closed bool `json:"closed"`
}

type recalibrateResponse struct {
	Version int `json:"version"`
}

type vibrationResponse struct {
	LoadMM     float64 `json:"load_mm"`
	Latched    bool    `json:"latched"`
	AlarmLevel string  `json:"alarm_level"`
}

type auditResponse struct {
	Events  []audit.Event `json:"events"`
	Count   int64         `json:"count"`
	Trimmed int           `json:"trimmed"`
}

type storeResponse struct {
	Keys  []string `json:"keys"`
	Count int      `json:"count"`
}

type curveResponse struct {
	Version int          `json:"version"`
	Points  []conv.Point `json:"points"`
}

func (s *Server) handlePitchMove(w http.ResponseWriter, r *http.Request) {
	var request moveRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.pitch.Move(request.AngleDeg); err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	seq, ack := s.pitch.Setpoint().Apply(request.AngleDeg, "manual")
	if !s.pitch.Setpoint().Ack(seq, ack) {
		writeError(w, http.StatusInternalServerError, "setpoint ack mismatch")
		return
	}
	writeJSON(w, http.StatusOK, moveResponse{AngleDeg: s.pitch.Angle(), AckSeq: seq})
}

func (s *Server) handlePitchRamp(w http.ResponseWriter, r *http.Request) {
	var request moveRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	final, err := s.pitch.Ramp(request.AngleDeg)
	if err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	steps := 0
	if s.limits.PitchRateDegPerSec > 0 {
		span := math.Abs(final - request.AngleDeg)
		steps = int(span / s.limits.PitchRateDegPerSec)
	}
	writeJSON(w, http.StatusOK, pitchRampResponse{AngleDeg: final, Steps: steps})
}

func (s *Server) handlePitchDemand(w http.ResponseWriter, r *http.Request) {
	var request moveRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.pitch.Demand(request.AngleDeg); err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]float64{"angle_deg": s.pitch.Angle()})
}

func (s *Server) handlePitchFeather(w http.ResponseWriter, r *http.Request) {
	if err := s.pitch.Feather(); err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"feathered": s.blade.Feathered()})
}

func (s *Server) handlePitchReset(w http.ResponseWriter, r *http.Request) {
	if err := s.pitch.Reset(); err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ready": s.pitch.Ready()})
}

func (s *Server) handleYawTurn(w http.ResponseWriter, r *http.Request) {
	var request yawRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	seq, ack := s.drive.Request(request.TargetDeg, "manual")
	if !s.drive.Ack(seq, ack) {
		writeError(w, http.StatusInternalServerError, "drive ack mismatch")
		return
	}
	if err := s.yaw.Turn(ack); err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, yawTurnResponse{YawDeg: s.yaw.Angle(), AckSeq: seq})
}

func (s *Server) handleYawApply(w http.ResponseWriter, r *http.Request) {
	if err := s.applyDrive(); err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, s.yaw.Report())
}

func (s *Server) handleYawUntwist(w http.ResponseWriter, r *http.Request) {
	if err := s.yaw.Untwist(); err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, s.cable.Report())
}

func (s *Server) handleBrakePressure(w http.ResponseWriter, r *http.Request) {
	var request pressureRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.brake.SetPressure(request.Bar); err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, brakePressureResponse{Bar: s.brake.Bar(), Durable: s.brake.Durable()})
}

func (s *Server) handleBrakeCharge(w http.ResponseWriter, r *http.Request) {
	if err := s.brake.Charge(); err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, brakeChargeResponse{Bar: s.brake.Bar(), Charged: s.brake.Charged()})
}

func (s *Server) handleBrakeLatch(w http.ResponseWriter, r *http.Request) {
	var request latchRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if request.Engaged {
		s.brake.Engage()
	} else {
		s.brake.Release()
	}
	writeJSON(w, http.StatusOK, brakeLatchResponse{Latched: s.brake.Latched()})
}

func (s *Server) handleWindSample(w http.ResponseWriter, r *http.Request) {
	var request windRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	sample := s.wind.Measure(request.SpeedMS, request.DirectionDeg)
	seq, target := s.wind.Track(sample)
	writeJSON(w, http.StatusOK, windSampleResponse{
		Sample:    sample,
		SpeedKMH:  ns.SpeedKMH(sample.SpeedMS),
		WindClass: ns.WindClass(sample.SpeedMS),
		YawTarget: target,
		AckSeq:    seq,
	})
}

func (s *Server) handleWindTrend(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.wind.Report())
}

func (s *Server) handleLoadSetpoint(w http.ResponseWriter, r *http.Request) {
	var request windRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	target := s.wind.LoadSetpoint(request.SpeedMS)
	writeJSON(w, http.StatusOK, loadTargetResponse{
		PowerKW:      target.PowerKW,
		CurveVersion: target.CurveVersion,
		WindClass:    ns.WindClass(request.SpeedMS),
	})
}

func (s *Server) handleCutIn(w http.ResponseWriter, r *http.Request) {
	sample := s.wind.LastSample()
	if !s.limits.InWindWindow(sample.SpeedMS) {
		writeError(w, http.StatusConflict, "wind speed is outside the cut-in window")
		return
	}
	if err := s.conv.Close(); err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"closed": s.conv.Closed()})
}

func (s *Server) handleBreakerOpen(w http.ResponseWriter, r *http.Request) {
	if err := s.conv.Open(); err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"closed": s.conv.Closed()})
}

func (s *Server) handleLimit(w http.ResponseWriter, r *http.Request) {
	var request limitRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	seq, angle := s.conv.ApplyLimit(request.PowerKW)
	writeJSON(w, http.StatusOK, limitResponse{
		AngleDeg:    angle,
		AckSeq:      seq,
		LastLimitKW: s.conv.LastLimit(),
		History:     s.conv.LimitHistory(),
	})
}

func (s *Server) handleTrip(w http.ResponseWriter, r *http.Request) {
	if err := s.grid.Loss(); err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, tripResponse{Online: s.grid.Online(), Closed: s.conv.Closed()})
}

func (s *Server) handleRecalibrate(w http.ResponseWriter, r *http.Request) {
	var request curveRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.conv.Recalibrate(request.Points); err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, recalibrateResponse{Version: s.conv.CurveVersion()})
}

func (s *Server) handleVibration(w http.ResponseWriter, r *http.Request) {
	var request vibrationRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	s.tower.Sample(request.LoadMM)
	writeJSON(w, http.StatusOK, vibrationResponse{
		LoadMM:     s.tower.LoadMM(),
		Latched:    s.tower.Latched(),
		AlarmLevel: s.tower.AlarmLevel(),
	})
}

func (s *Server) handleAudit(w http.ResponseWriter, r *http.Request) {
	limit := 50
	if raw := r.URL.Query().Get("limit"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	source := r.URL.Query().Get("source")
	action := r.URL.Query().Get("action")
	var events []audit.Event
	switch {
	case source != "":
		events = s.audit.BySource(source, limit)
	case action != "":
		events = s.audit.ByAction(action, limit)
	default:
		events = s.audit.Recent(limit)
	}
	trimmed := s.audit.Trim(int64(s.audit.BufferSize()))
	writeJSON(w, http.StatusOK, auditResponse{Events: events, Count: s.audit.Count(), Trimmed: trimmed})
}

func (s *Server) handleAuditSummary(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"summary": s.audit.Summary(), "count": s.audit.Count()})
}

func (s *Server) handleStore(w http.ResponseWriter, r *http.Request) {
	keys, err := s.st.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, storeResponse{Keys: keys, Count: len(keys)})
}

func (s *Server) handlePowerCurve(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, curveResponse{Version: s.conv.CurveVersion(), Points: s.conv.Points()})
}

func (s *Server) handleGridVoltage(w http.ResponseWriter, r *http.Request) {
	var request voltageRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.grid.SampleVoltage(request.VoltageKV); err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, s.grid.Report())
}

func (s *Server) handlePitchStatus(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.pitch.Report())
}

func (s *Server) handleBrakeStatus(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.brake.Report())
}

func (s *Server) handleCableStatus(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.cable.Report())
}

func (s *Server) handleTowerStatus(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.tower.Report())
}
