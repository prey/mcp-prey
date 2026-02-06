package internal

import "testing"

func TestRequireID(t *testing.T) {
	if err := RequireID("", "deviceId"); err == nil {
		t.Fatalf("expected error for empty id")
	}
	if err := RequireID("abc", "deviceId"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
