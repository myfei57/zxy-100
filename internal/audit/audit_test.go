package audit

import (
	"path/filepath"
	"testing"

	"wtc/internal/store"
)

func TestRecorderRecentOrder(t *testing.T) {
	st := store.New(filepath.Join(t.TempDir(), "audit"))
	recorder := New(st, 100)
	if got := recorder.Count(); got != 0 {
		t.Fatalf("fresh recorder must be empty, got %d", got)
	}
	first := recorder.Append("pitch", "move", "angle=30.0")
	second := recorder.Append("brake", "charge", "accumulator charged")
	if first.Seq >= second.Seq {
		t.Fatalf("sequence must increase: %d >= %d", first.Seq, second.Seq)
	}
	events := recorder.Recent(10)
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[0].Action != "move" || events[1].Action != "charge" {
		t.Fatalf("event order mismatch: %s %s", events[0].Action, events[1].Action)
	}
	if events[0].ID == "" || events[1].ID == "" {
		t.Fatal("events must carry ids")
	}
}
