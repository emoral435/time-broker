package provider

import "time"

type Slot struct {
	Start time.Time
	End   time.Time
}

type Provider interface {
	Name() string
	Auth() error
	FreeSlots(day time.Time, minDuration time.Duration) ([]Slot, error)
	Book(title, description string, start, end time.Time, allDay bool) error
}
