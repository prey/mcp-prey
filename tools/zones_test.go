package tools

import "testing"

func TestValidateZoneTrigger(t *testing.T) {
	if err := validateZoneTrigger(ZoneTriggerParams{}); err == nil {
		t.Fatalf("expected error for empty trigger")
	}
	if err := validateZoneTrigger(ZoneTriggerParams{Context: "when_in", ActionName: "alarm"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := validateZoneTrigger(ZoneTriggerParams{Context: "bad", ActionName: "alarm"}); err == nil {
		t.Fatalf("expected error for invalid context")
	}
	if err := validateZoneTrigger(ZoneTriggerParams{Context: "when_in", ActionName: "bad"}); err == nil {
		t.Fatalf("expected error for invalid action")
	}
}

func TestValidateZoneNotifications(t *testing.T) {
	if err := validateZoneNotifications(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := validateZoneNotifications(&ZoneNotificationParams{WhenIn: "on", WhenOut: "off"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := validateZoneNotifications(&ZoneNotificationParams{WhenIn: "bad"}); err == nil {
		t.Fatalf("expected error for invalid when_in")
	}
	if err := validateZoneNotifications(&ZoneNotificationParams{WhenOut: "bad"}); err == nil {
		t.Fatalf("expected error for invalid when_out")
	}
}
