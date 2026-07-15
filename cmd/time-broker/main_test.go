package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/emoral435/time-broker/internal/config"
	"github.com/emoral435/time-broker/internal/input"
)

const (
	testUnknownCmd = "foobar"
	testVersion    = "1.0.0"
	testEvent      = "event"
	testInvalid    = "invalid"
	testDateFlag   = "--date"
)

func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	return buf.String()
}

func TestRunNoArgs(t *testing.T) {
	got := captureStdout(func() {
		if err := run([]string{}); err != nil {
			t.Errorf("run() unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "Usage: time-broker <command>") {
		t.Errorf("expected help text, got: %s", got)
	}
}

func TestRunHelp(t *testing.T) {
	got := captureStdout(func() {
		if err := run([]string{helpAsStr}); err != nil {
			t.Errorf("run() unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "Usage: time-broker <command>") {
		t.Errorf("expected help text, got: %s", got)
	}
}

func TestRunVersion(t *testing.T) {
	got := captureStdout(func() {
		if err := run([]string{"version"}); err != nil {
			t.Errorf("run() unexpected error: %v", err)
		}
	})
	want := "time-broker version"
	if !strings.HasPrefix(got, want) {
		t.Errorf("version did not contain the correct prefix, string produces: %q; want prefix %q", got, want)
	}
}

func TestRunUnknownCommand(t *testing.T) {
	err := run([]string{testUnknownCmd})
	if err == nil {
		t.Fatal("expected error for unknown command")
	}
	if !strings.Contains(err.Error(), "unknown command") {
		t.Errorf("error = %q; want 'unknown command'", err)
	}
}

func TestRunAuthNoConfig(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	err := run([]string{"auth"})
	if err == nil {
		t.Fatal("expected error when not configured")
	}
	if !strings.Contains(err.Error(), "no configuration found") {
		t.Errorf("error = %q; want 'no configuration found'", err)
	}
}

func TestRunConfigNoArgsNoConfig(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	if err := run([]string{cfgAsStr}); err != nil {
		t.Fatalf("run() returned error: %v", err)
	}
}

func TestRunConfigUnknownSubcommand(t *testing.T) {
	err := run([]string{cfgAsStr, testUnknownCmd})
	if err == nil {
		t.Fatal("expected error for unknown config subcommand")
	}
	if !strings.Contains(err.Error(), "unknown config subcommand") {
		t.Errorf("error = %q; want 'unknown config subcommand'", err)
	}
}

func TestRunVersionOutput(t *testing.T) {
	Version = "0.3.1-beta"
	got := captureStdout(runVersion)
	want := "time-broker version 0.3.1-beta\n"
	if got != want {
		t.Errorf("runVersion() = %q; want %q", got, want)
	}
}

func TestInitReadsVersionFile(t *testing.T) {
	Version = devAsStr
	dir := t.TempDir()
	oldWd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() {
		if err := os.Chdir(oldWd); err != nil {
			t.Logf("failed to restore wd: %v", err)
		}
	}()

	if err := os.WriteFile("version.txt", []byte("2.0.0-rc1\n"), 0644); err != nil {
		t.Fatalf("write version.txt file: %v", err)
	}

	loadVersionFromFile()

	if Version != "2.0.0-rc1" {
		t.Errorf("Version = %q; want %q", Version, "2.0.0-rc1")
	}
}

func TestInitNoVersionFile(t *testing.T) {
	Version = devAsStr
	dir := t.TempDir()
	oldWd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() {
		if err := os.Chdir(oldWd); err != nil {
			t.Logf("failed to restore wd: %v", err)
		}
	}()

	loadVersionFromFile()

	if Version != devAsStr {
		t.Errorf("Version should stay 'dev' when no file, got %q", Version)
	}
}

func TestRunScheduleNotConfigured(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	oldWizard := runSetupWizardFn
	runSetupWizardFn = func() (*config.Config, error) {
		return nil, errors.New("wizard disabled in tests")
	}
	t.Cleanup(func() { runSetupWizardFn = oldWizard })

	got := captureStdout(func() {
		if err := run([]string{schedule, testEvent}); err != nil {
			t.Errorf("run() unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "Event Details") {
		t.Errorf("expected event details output, got: %s", got)
	}
	if !strings.Contains(got, "Cancelled") {
		t.Errorf("expected 'Cancelled' in output, got: %s", got)
	}
}

func TestRunScheduleNoArgsShowsHelp(t *testing.T) {
	got := captureStdout(func() {
		if err := run([]string{schedule}); err != nil {
			t.Errorf("run() unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "Usage: time-broker schedule <subcommand>") {
		t.Errorf("expected schedule help text, got: %s", got)
	}
}

func TestRunScheduleHelp(t *testing.T) {
	got := captureStdout(func() {
		if err := run([]string{schedule, helpAsStr}); err != nil {
			t.Errorf("run() unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "Usage: time-broker schedule <subcommand>") {
		t.Errorf("expected schedule help text, got: %s", got)
	}
	if !strings.Contains(got, "event") {
		t.Errorf("expected 'event' in help text, got: %s", got)
	}
	if !strings.Contains(got, "cancel") {
		t.Errorf("expected 'cancel' in help text, got: %s", got)
	}
	if !strings.Contains(got, "update") {
		t.Errorf("expected 'update' in help text, got: %s", got)
	}
}

func TestRunScheduleUnknownSubcommand(t *testing.T) {
	err := run([]string{schedule, testUnknownCmd})
	if err == nil {
		t.Fatal("expected error for unknown schedule subcommand")
	}
	if !strings.Contains(err.Error(), "unknown schedule subcommand") {
		t.Errorf("error = %q; want 'unknown schedule subcommand'", err)
	}
}

func TestRunViewNoArgsShowsHelp(t *testing.T) {
	got := captureStdout(func() {
		if err := run([]string{view}); err != nil {
			t.Errorf("run() unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "Usage: time-broker view <subcommand>") {
		t.Errorf("expected view help text, got: %s", got)
	}
}

func TestRunViewHelp(t *testing.T) {
	got := captureStdout(func() {
		if err := run([]string{view, helpAsStr}); err != nil {
			t.Errorf("run() unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "Usage: time-broker view <subcommand>") {
		t.Errorf("expected view help text, got: %s", got)
	}
	if !strings.Contains(got, "event") {
		t.Errorf("expected 'event' in help text, got: %s", got)
	}
	if !strings.Contains(got, "day") {
		t.Errorf("expected 'day' in help text, got: %s", got)
	}
	if !strings.Contains(got, "availability") {
		t.Errorf("expected 'availability' in help text, got: %s", got)
	}
}

func TestRunViewUnknownSubcommand(t *testing.T) {
	err := run([]string{view, testUnknownCmd})
	if err == nil {
		t.Fatal("expected error for unknown view subcommand")
	}
	if !strings.Contains(err.Error(), "unknown view subcommand") {
		t.Errorf("error = %q; want 'unknown view subcommand'", err)
	}
}

func TestRunViewNotConfigured(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	oldWizard := runSetupWizardFn
	runSetupWizardFn = func() (*config.Config, error) {
		return nil, errors.New("wizard disabled in tests")
	}
	t.Cleanup(func() { runSetupWizardFn = oldWizard })

	err := run([]string{view, "day"})
	if err == nil {
		t.Fatal("expected error when not configured")
	}
}

func TestRunUpdateNotConfigured(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	oldAPI := githubAPIBase
	githubAPIBase = server.URL
	t.Cleanup(func() { githubAPIBase = oldAPI })

	err := run([]string{update})
	if err == nil {
		t.Fatal("expected error when GitHub API is unreachable")
	}
}

func TestLatestVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/emoral435/time-broker/releases/latest" {
			t.Errorf("unexpected path: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_ = json.NewEncoder(w).Encode(githubRelease{
			TagName: "v1.2.3",
			Body:    "Release notes here",
		})
	}))
	defer server.Close()

	oldAPI := githubAPIBase
	githubAPIBase = server.URL
	t.Cleanup(func() { githubAPIBase = oldAPI })

	got, err := latestVersion()
	if err != nil {
		t.Fatalf("latestVersion() error: %v", err)
	}
	if got != "v1.2.3" {
		t.Errorf("latestVersion() = %q; want %q", got, "v1.2.3")
	}
}

func TestRunUpdateAlreadyLatest(t *testing.T) {
	Version = testVersion

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(githubRelease{
			TagName: "1.0.0",
			Body:    "",
		})
	}))
	defer server.Close()

	oldAPI := githubAPIBase
	githubAPIBase = server.URL
	t.Cleanup(func() { githubAPIBase = oldAPI })

	got := captureStdout(func() {
		if err := run([]string{update}); err != nil {
			t.Errorf("run() unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "Already up to date") {
		t.Errorf("expected 'Already up to date', got: %s", got)
	}
}

func TestRunUpdateNewVersionAvailable(t *testing.T) {
	Version = testVersion

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(githubRelease{
			TagName: "2.0.0",
			Body:    "New features",
		})
	}))
	defer server.Close()

	oldAPI := githubAPIBase
	githubAPIBase = server.URL
	t.Cleanup(func() { githubAPIBase = oldAPI })

	got := captureStdout(func() {
		err := run([]string{update})
		if err == nil {
			t.Error("expected error during download (no real binary to replace)")
		}
	})
	if !strings.Contains(got, "Updating from 1.0.0 to 2.0.0") {
		t.Errorf("expected update message, got: %s", got)
	}
}

func TestRunConfigHelp(t *testing.T) {
	got := captureStdout(func() {
		if err := run([]string{cfgAsStr, helpAsStr}); err != nil {
			t.Errorf("run() unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "Usage: time-broker config <subcommand>") {
		t.Errorf("expected config help text, got: %s", got)
	}
}

func TestRunScheduleEventHelp(t *testing.T) {
	got := captureStdout(func() {
		if err := run([]string{schedule, testEvent, helpAsStr}); err != nil {
			t.Errorf("run() unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "Usage: time-broker schedule event [flags]") {
		t.Errorf("expected schedule event help text, got: %s", got)
	}
	if !strings.Contains(got, "--title") {
		t.Errorf("expected '--title' in help text, got: %s", got)
	}
	if !strings.Contains(got, "--description") {
		t.Errorf("expected '--description' in help text, got: %s", got)
	}
	if !strings.Contains(got, "--timeRange") {
		t.Errorf("expected '--timeRange' in help text, got: %s", got)
	}
	if !strings.Contains(got, "--date") {
		t.Errorf("expected '--date' in help text, got: %s", got)
	}
}

func TestParseTimeRange(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		startH  int
		startM  int
		endH    int
		endM    int
	}{
		{
			name:    "valid range",
			input:   "9:00AM-5:00PM",
			wantErr: false,
			startH:  9, startM: 0,
			endH: 17, endM: 0,
		},
		{
			name:    "valid range same period",
			input:   "2:00PM-4:00PM",
			wantErr: false,
			startH:  14, startM: 0,
			endH: 16, endM: 0,
		},
		{
			name:    "valid range with minutes",
			input:   "9:30AM-11:15AM",
			wantErr: false,
			startH:  9, startM: 30,
			endH: 11, endM: 15,
		},
		{
			name:    "no dash",
			input:   "9:00AM",
			wantErr: true,
		},
		{
			name:    "missing AM/PM",
			input:   "9:00-5:00PM",
			wantErr: true,
		},
		{
			name:    "invalid format",
			input:   testInvalid,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end, err := parseTimeRange(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("parseTimeRange(%q) expected error, got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseTimeRange(%q) unexpected error: %v", tt.input, err)
			}
			startH := int(start.Hours())
			startM := int(start.Minutes()) % 60
			if startH != tt.startH || startM != tt.startM {
				t.Errorf("start = %d:%d, want %d:%d", startH, startM, tt.startH, tt.startM)
			}
			endH := int(end.Hours())
			endM := int(end.Minutes()) % 60
			if endH != tt.endH || endM != tt.endM {
				t.Errorf("end = %d:%d, want %d:%d", endH, endM, tt.endH, tt.endM)
			}
		})
	}
}

func TestParseDateFlag(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		month   time.Month
		day     int
		year    int
	}{
		{
			name:    "valid date",
			input:   "07-14-2026",
			wantErr: false,
			month:   time.July,
			day:     14,
			year:    2026,
		},
		{
			name:    "valid date january",
			input:   "01-01-2026",
			wantErr: false,
			month:   time.January,
			day:     1,
			year:    2026,
		},
		{
			name:    "valid date M/D/YYYY",
			input:   "7/14/2026",
			wantErr: false,
			month:   time.July,
			day:     14,
			year:    2026,
		},
		{
			name:    "valid date MM/DD/YYYY",
			input:   "01/01/2026",
			wantErr: false,
			month:   time.January,
			day:     1,
			year:    2026,
		},
		{
			name:    "invalid format",
			input:   "2026-07-14",
			wantErr: true,
		},
		{
			name:    "invalid date",
			input:   "13-01-2026",
			wantErr: true,
		},
		{
			name:    "invalid day",
			input:   "07-32-2026",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := input.ParseDate(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseDate(%q) expected error, got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseDate(%q) unexpected error: %v", tt.input, err)
			}
			if d.Month() != tt.month {
				t.Errorf("month = %v, want %v", d.Month(), tt.month)
			}
			if d.Day() != tt.day {
				t.Errorf("day = %d, want %d", d.Day(), tt.day)
			}
			if d.Year() != tt.year {
				t.Errorf("year = %d, want %d", d.Year(), tt.year)
			}
		})
	}
}

func TestDefaultDate(t *testing.T) {
	got := defaultDate()
	expected := time.Now().AddDate(0, 0, 1).Format(dateFormat)
	if got != expected {
		t.Errorf("defaultDate() = %q, want %q", got, expected)
	}
}

func TestRunScheduleEventDisplaysDetails(t *testing.T) {
	got := captureStdout(func() {
		if err := run([]string{schedule, testEvent,
			"--title", "Test Meeting",
			"--description", "Test Description",
			testDateFlag, "07-14-2026",
			"--timeRange", "9:00AM-5:00PM",
		}); err != nil {
			t.Errorf("run() unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "Event Details") {
		t.Errorf("expected 'Event Details' in output, got: %s", got)
	}
	if !strings.Contains(got, "Test Meeting") {
		t.Errorf("expected 'Test Meeting' in output, got: %s", got)
	}
	if !strings.Contains(got, "Test Description") {
		t.Errorf("expected 'Test Description' in output, got: %s", got)
	}
	if !strings.Contains(got, "9:00 AM - 5:00 PM") {
		t.Errorf("expected '9:00 AM - 5:00 PM' in output, got: %s", got)
	}
	if !strings.Contains(got, "Tuesday, July 14, 2026") {
		t.Errorf("expected 'Tuesday, July 14, 2026' in output, got: %s", got)
	}
}

func TestRunScheduleEventAllDay(t *testing.T) {
	got := captureStdout(func() {
		if err := run([]string{schedule, testEvent,
			"--title", "Holiday",
			testDateFlag, "12-25-2026",
		}); err != nil {
			t.Errorf("run() unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "All day") {
		t.Errorf("expected 'All day' in output, got: %s", got)
	}
	if !strings.Contains(got, "Holiday") {
		t.Errorf("expected 'Holiday' in output, got: %s", got)
	}
}

func TestRunScheduleEventInvalidTimeRange(t *testing.T) {
	err := run([]string{schedule, testEvent, "--timeRange", testInvalid})
	if err == nil {
		t.Fatal("expected error for invalid time range")
	}
	if !strings.Contains(err.Error(), "invalid time range") {
		t.Errorf("error = %q; want 'invalid time range'", err)
	}
}

func TestRunScheduleEventInvalidDate(t *testing.T) {
	err := run([]string{schedule, testEvent, testDateFlag, testInvalid})
	if err == nil {
		t.Fatal("expected error for invalid date")
	}
	if !strings.Contains(err.Error(), "invalid date") {
		t.Errorf("error = %q; want 'invalid date'", err)
	}
}
