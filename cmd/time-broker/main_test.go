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

	"github.com/emoral435/time-broker/internal/config"
)

const (
	testUnknownCmd = "foobar"
	testVersion    = "1.0.0"
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

	if err := os.WriteFile("VERSION", []byte("2.0.0-rc1\n"), 0644); err != nil {
		t.Fatalf("write VERSION file: %v", err)
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

	err := run([]string{schedule, "event"})
	if err == nil {
		t.Fatal("expected error when not configured")
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
	if !strings.Contains(got, "view") {
		t.Errorf("expected 'view' in help text, got: %s", got)
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
