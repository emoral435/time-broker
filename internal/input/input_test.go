package input

import (
	"testing"
	"time"
)

func TestParseDate(t *testing.T) {
	tests := []struct {
		input string
		want  time.Time
		err   bool
	}{
		{"01-31-2027", time.Date(2027, 1, 31, 0, 0, 0, 0, time.Local), false},
		{"12-25-2026", time.Date(2026, 12, 25, 0, 0, 0, 0, time.Local), false},
		{"01-01-2000", time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local), false},
		{"01/31/2027", time.Date(2027, 1, 31, 0, 0, 0, 0, time.Local), false},
		{"12/25/2026", time.Date(2026, 12, 25, 0, 0, 0, 0, time.Local), false},
		{"1/31/2027", time.Date(2027, 1, 31, 0, 0, 0, 0, time.Local), false},
		{"7/4/2026", time.Date(2026, 7, 4, 0, 0, 0, 0, time.Local), false},
		{"12/1/2026", time.Date(2026, 12, 1, 0, 0, 0, 0, time.Local), false},
		{"2027-01-31", time.Time{}, true},
		{"1-31-2027", time.Time{}, true},
		{"invalid", time.Time{}, true},
		{"", time.Time{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseDate(tt.input)
			if (err != nil) != tt.err {
				t.Errorf("ParseDate(%q) error = %v, wantErr %v", tt.input, err, tt.err)
				return
			}
			if !tt.err && !got.Equal(tt.want) {
				t.Errorf("ParseDate(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseDateWithSpaces(t *testing.T) {
	got, err := ParseDate("  01-31-2027  ")
	if err != nil {
		t.Fatalf("ParseDate with spaces unexpected error: %v", err)
	}
	want := time.Date(2027, 1, 31, 0, 0, 0, 0, time.Local)
	if !got.Equal(want) {
		t.Errorf("ParseDate(\"  01-31-2027  \") = %v, want %v", got, want)
	}
}

func TestParseDateFormats(t *testing.T) {
	got, err := ParseDate("1/15/2026")
	if err != nil {
		t.Fatalf("ParseDate M/D/YYYY unexpected error: %v", err)
	}
	want := time.Date(2026, 1, 15, 0, 0, 0, 0, time.Local)
	if !got.Equal(want) {
		t.Errorf("ParseDate(\"1/15/2026\") = %v, want %v", got, want)
	}

	got, err = ParseDate("01/15/2026")
	if err != nil {
		t.Fatalf("ParseDate MM/DD/YYYY unexpected error: %v", err)
	}
	if !got.Equal(want) {
		t.Errorf("ParseDate(\"01/15/2026\") = %v, want %v", got, want)
	}
}

func TestParseDateReturnsLocalTimezone(t *testing.T) {
	got, err := ParseDate("07-14-2026")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Location() != time.Local {
		t.Errorf("ParseDate location = %v, want %v", got.Location(), time.Local)
	}
}

func TestParseTime(t *testing.T) {
	tests := []struct {
		input string
		want  time.Duration
		err   bool
	}{
		{"3:00AM", 3 * time.Hour, false},
		{"3:00PM", 15 * time.Hour, false},
		{"12:00AM", 0, false},
		{"12:00PM", 12 * time.Hour, false},
		{"9:30AM", 9*time.Hour + 30*time.Minute, false},
		{"11:59PM", 23*time.Hour + 59*time.Minute, false},
		{"1:00am", 1 * time.Hour, false},
		{"1:00pm", 13 * time.Hour, false},
		{"0:00AM", 0, true},
		{"13:00AM", 0, true},
		{"3:60AM", 0, true},
		{"3:00", 0, true},
		{"3AM", 0, true},
		{"invalid", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseTime(tt.input)
			if (err != nil) != tt.err {
				t.Errorf("ParseTime(%q) error = %v, wantErr %v", tt.input, err, tt.err)
				return
			}
			if !tt.err && got != tt.want {
				t.Errorf("ParseTime(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestLevenshtein(t *testing.T) {
	t.Run("empty strings", func(t *testing.T) {
		if got := Levenshtein("", ""); got != 0 {
			t.Errorf("got %d, want 0", got)
		}
	})
	t.Run("first empty", func(t *testing.T) {
		if got := Levenshtein("foo", ""); got != 3 {
			t.Errorf("got %d, want 3", got)
		}
	})
	t.Run("second empty", func(t *testing.T) {
		if got := Levenshtein("", "bar"); got != 3 {
			t.Errorf("got %d, want 3", got)
		}
	})
	t.Run("identical", func(t *testing.T) {
		if got := Levenshtein("foo", "foo"); got != 0 {
			t.Errorf("got %d, want 0", got)
		}
	})
	t.Run("one deletion", func(t *testing.T) {
		if got := Levenshtein("foo", "fo"); got != 1 {
			t.Errorf("got %d, want 1", got)
		}
	})
	t.Run("two deletions", func(t *testing.T) {
		if got := Levenshtein("foo", "f"); got != 2 {
			t.Errorf("got %d, want 2", got)
		}
	})
	t.Run("classic", func(t *testing.T) {
		if got := Levenshtein("kitten", "sitting"); got != 3 {
			t.Errorf("got %d, want 3", got)
		}
	})
	t.Run("saturday sunday", func(t *testing.T) {
		if got := Levenshtein("saturday", "sunday"); got != 3 {
			t.Errorf("got %d, want 3", got)
		}
	})
	t.Run("long strings", func(t *testing.T) {
		if got := Levenshtein("rosettacode", "raisethysword"); got != 8 {
			t.Errorf("got %d, want 8", got)
		}
	})
}

func TestFuzzyMatch(t *testing.T) {
	t.Run("no match short query", func(t *testing.T) {
		if got := FuzzyMatch("standup", "Weekly Team Standup", 0.8); got {
			t.Error("expected false")
		}
	})
	t.Run("exact match", func(t *testing.T) {
		if got := FuzzyMatch("standup", "standup", 0.8); !got {
			t.Error("expected true")
		}
	})
	t.Run("case insensitive", func(t *testing.T) {
		if got := FuzzyMatch("Standup", "standup", 0.8); !got {
			t.Error("expected true")
		}
	})
	t.Run("unrelated query", func(t *testing.T) {
		if got := FuzzyMatch("sync", "standup", 0.8); got {
			t.Error("expected false")
		}
	})
	t.Run("empty query", func(t *testing.T) {
		if got := FuzzyMatch("", "standup", 0.8); got {
			t.Error("expected false")
		}
	})
	t.Run("near miss", func(t *testing.T) {
		if got := FuzzyMatch("qwer", "qwert", 0.8); !got {
			t.Error("expected true")
		}
	})
	t.Run("identical short", func(t *testing.T) {
		if got := FuzzyMatch("asdf", "asdf", 0.8); !got {
			t.Error("expected true")
		}
	})
	t.Run("one char added", func(t *testing.T) {
		if got := FuzzyMatch("mnop", "mnop", 0.8); !got {
			t.Error("expected true")
		}
	})
	t.Run("too different", func(t *testing.T) {
		if got := FuzzyMatch("greetings", "greetings earth", 0.8); got {
			t.Error("expected false")
		}
	})
	t.Run("partial with suffix", func(t *testing.T) {
		if got := FuzzyMatch("daily sync", "Daily Sync Notes", 0.8); got {
			t.Error("expected false")
		}
	})
	t.Run("exact multi word", func(t *testing.T) {
		if got := FuzzyMatch("daily sync", "Daily Sync", 0.8); !got {
			t.Error("expected true")
		}
	})
	t.Run("unrelated concept", func(t *testing.T) {
		if got := FuzzyMatch("planning", "Sprint Planning", 0.8); got {
			t.Error("expected false")
		}
	})
}
