package provider

import "time"

type Slot struct {
	Start time.Time
	End   time.Time
}

type Event struct {
	Title       string
	Description string
	Location    string
	Start       time.Time
	End         time.Time
	AllDay      bool
}

type Provider interface {
	Name() string
	Auth() error
	EnsureAuthenticated() error
	Book(title, description string, start, end time.Time, allDay bool) error
	FreeSlots(start, end time.Time, minDuration time.Duration) ([]Slot, error)
	EventsForDay(day time.Time) ([]Event, error)
	Timezone() *time.Location
}
