package internal

import "testing"

func TestRequireOneOf(t *testing.T) {
	if err := RequireOneOf("", "command", "start"); err == nil {
		t.Fatalf("expected error for empty value")
	}
	if err := RequireOneOf("stop", "command", "start"); err == nil {
		t.Fatalf("expected error for invalid value")
	}
	if err := RequireOneOf("start", "command", "start"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
