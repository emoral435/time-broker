package provider

import (
	"testing"
	"time"
)

func TestSlotZeroValue(t *testing.T) {
	var s Slot
	if !s.Start.IsZero() {
		t.Error("Slot.Start should be zero value")
	}
	if !s.End.IsZero() {
		t.Error("Slot.End should be zero value")
	}
}

func TestSlotDuration(t *testing.T) {
	start := time.Date(2026, 6, 26, 9, 0, 0, 0, time.UTC)
	end := time.Date(2026, 6, 26, 10, 0, 0, 0, time.UTC)
	s := Slot{Start: start, End: end}
	if s.End.Sub(s.Start) != time.Hour {
		t.Errorf("expected duration 1h, got %v", s.End.Sub(s.Start))
	}
}
