package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	buf.ReadFrom(r)
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
		if err := run([]string{"help"}); err != nil {
			t.Errorf("run() unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "Usage: time-broker <command>") {
		t.Errorf("expected help text, got: %s", got)
	}
}

func TestRunVersion(t *testing.T) {
	Version = "1.0.0-test"
	got := captureStdout(func() {
		if err := run([]string{"version"}); err != nil {
			t.Errorf("run() unexpected error: %v", err)
		}
	})
	want := "time-broker 1.0.0-test\n"
	if got != want {
		t.Errorf("version output = %q; want %q", got, want)
	}
}

func TestRunUnknownCommand(t *testing.T) {
	err := run([]string{"foobar"})
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

	if err := run([]string{"config"}); err != nil {
		t.Fatalf("run() returned error: %v", err)
	}
}

func TestRunConfigUnknownSubcommand(t *testing.T) {
	err := run([]string{"config", "foobar"})
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
	want := "time-broker 0.3.1-beta\n"
	if got != want {
		t.Errorf("runVersion() = %q; want %q", got, want)
	}
}

func TestInitReadsVersionFile(t *testing.T) {
	Version = "dev"
	dir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(oldWd)

	os.WriteFile("VERSION", []byte("2.0.0-rc1\n"), 0644)

	loadVersionFromFile()

	if Version != "2.0.0-rc1" {
		t.Errorf("Version = %q; want %q", Version, "2.0.0-rc1")
	}
}

func TestInitNoVersionFile(t *testing.T) {
	Version = "dev"
	dir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(oldWd)

	loadVersionFromFile()

	if Version != "dev" {
		t.Errorf("Version should stay 'dev' when no file, got %q", Version)
	}
}

func TestRunScheduleNotConfigured(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	err := run([]string{"schedule"})
	if err == nil {
		t.Fatal("expected error when not configured")
	}
}

func TestRunUpdateNotConfigured(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	err := run([]string{"update"})
	if err == nil {
		t.Fatal("expected error when not configured")
	}
}

func TestRunConfigHelp(t *testing.T) {
	got := captureStdout(func() {
		if err := run([]string{"config", "help"}); err != nil {
			t.Errorf("run() unexpected error: %v", err)
		}
	})
	if !strings.Contains(got, "Usage: time-broker config <subcommand>") {
		t.Errorf("expected config help text, got: %s", got)
	}
}
