package console

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"wtc/internal/audit"
	"wtc/internal/blade"
	"wtc/internal/brake"
	"wtc/internal/cable"
	"wtc/internal/conv"
	"wtc/internal/grid"
	"wtc/internal/ns"
	"wtc/internal/pitch"
	"wtc/internal/store"
	"wtc/internal/tower"
	"wtc/internal/wind"
	"wtc/internal/yaw"
)

func buildTestServer(t *testing.T) *Server {
	t.Helper()
	st := store.New(filepath.Join(t.TempDir(), "server"))
	limits := ns.DefaultLimits()
	recorder := audit.New(st, 200)
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
	windSys := wind.NewSystem(drive, curve, limits, 10)
	gridSys := grid.NewSystem(convSys, recorder, 18)
	return NewServer(":0", st, limits, recorder, bladeSys, brakeSys, pitchSys, yawSys, windSys, convSys, towerSys, cableSys, gridSys, drive)
}

func TestStatusEndpoint(t *testing.T) {
	server := buildTestServer(t)
	request := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status code: %d", response.Code)
	}
	var status ns.Status
	if err := json.Unmarshal(response.Body.Bytes(), &status); err != nil {
		t.Fatalf("decode status: %v", err)
	}
	if !status.GridOnline {
		t.Fatal("grid must be online initially")
	}
	if status.Ready {
		t.Fatal("turbine must not be ready without durable brake pressure")
	}
}

func TestTripEndpointInvariant(t *testing.T) {
	server := buildTestServer(t)
	request := httptest.NewRequest(http.MethodPost, "/api/trip", nil)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("trip must succeed, got %d", response.Code)
	}
	var body struct {
		Online bool `json:"online"`
		Closed bool `json:"closed"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode trip: %v", err)
	}
	if body.Online || body.Closed {
		t.Fatalf("trip must leave the grid offline and the breaker open: %+v", body)
	}
}
