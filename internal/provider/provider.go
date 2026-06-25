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
	Book(title string, start, end time.Time) error
}
