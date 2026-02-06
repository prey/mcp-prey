package internal

import "testing"

func TestMaskSensitive(t *testing.T) {
	input := map[string]any{
		"token": "abc",
		"nested": map[string]any{
			"password": "secret",
		},
		"list": []any{
			map[string]any{"api_key": "123"},
		},
	}
	masked := MaskSensitive(input).(map[string]any)
	if masked["token"] != "***" {
		t.Fatalf("expected token masked")
	}
	nested := masked["nested"].(map[string]any)
	if nested["password"] != "***" {
		t.Fatalf("expected password masked")
	}
	list := masked["list"].([]any)
	item := list[0].(map[string]any)
	if item["api_key"] != "***" {
		t.Fatalf("expected api_key masked")
	}
}
