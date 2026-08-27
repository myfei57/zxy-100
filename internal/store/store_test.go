package store

import (
	"path/filepath"
	"testing"
)

func TestStoreRoundTrip(t *testing.T) {
	st := New(filepath.Join(t.TempDir(), "store"))
	type payload struct {
		Value int `json:"value"`
	}
	if err := st.Put("demo/key", payload{Value: 42}); err != nil {
		t.Fatalf("put: %v", err)
	}
	if !st.Exists("demo/key") {
		t.Fatal("key must exist after put")
	}
	var got payload
	if err := st.Get("demo/key", &got); err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Value != 42 {
		t.Fatalf("value mismatch: %d", got.Value)
	}
	if err := st.Delete("demo/key"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if st.Exists("demo/key") {
		t.Fatal("key must not exist after delete")
	}
	var missing payload
	if err := st.Get("demo/missing", &missing); err == nil {
		t.Fatal("get must fail for a missing key")
	}
}
