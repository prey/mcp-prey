package tools

import (
	"context"
	"testing"

	"mcp-prey/internal"
)

func TestActionValidation(t *testing.T) {
	ctx := context.Background()

	if err := internal.RequireOneOf("start", "command", "start"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := internal.RequireOneOf("bad", "command", "start"); err == nil {
		t.Fatalf("expected error for invalid command")
	}
	if err := internal.RequireOneOf("alarm", "action_name", "alarm", "alert", "lock"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := internal.RequireOneOf("bad", "action_name", "alarm", "alert", "lock"); err == nil {
		t.Fatalf("expected error for invalid action")
	}

	_ = ctx
}
