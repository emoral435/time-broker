package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsConfigured(t *testing.T) {
	tests := []struct {
		name   string
		cfg    *Config
		wanted bool
	}{
		{name: "empty", cfg: &Config{}, wanted: false},
		{name: "provider only", cfg: &Config{Provider: "google"}, wanted: false},
		{name: "week_start_day only", cfg: &Config{WeekStartDay: "monday"}, wanted: false},
		{name: "full", cfg: &Config{Provider: "google", WeekStartDay: "monday"}, wanted: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsConfigured(tc.cfg); got != tc.wanted {
				t.Errorf("IsConfigured(%+v) = %v; want %v", tc.cfg, got, tc.wanted)
			}
		})
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	cfg := &Config{Provider: "google", WeekStartDay: "sunday"}
	if err := Save(cfg); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	got, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if got.Provider != cfg.Provider {
		t.Errorf("Provider = %q; want %q", got.Provider, cfg.Provider)
	}
	if got.WeekStartDay != cfg.WeekStartDay {
		t.Errorf("WeekStartDay = %q; want %q", got.WeekStartDay, cfg.WeekStartDay)
	}
}

func TestSaveAndLoadEmpty(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	cfg := &Config{}
	if err := Save(cfg); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	got, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if got.Provider != "" || got.WeekStartDay != "" {
		t.Errorf("expected empty config, got %+v", got)
	}
}

func TestLoadFileNotFound(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg == nil {
		t.Fatal("Load() returned nil config")
	}
}

func TestEnsureDir(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	if err := EnsureDir(); err != nil {
		t.Fatalf("EnsureDir() error: %v", err)
	}

	expected := filepath.Join(dir, ".time-broker")
	if _, err := os.Stat(expected); os.IsNotExist(err) {
		t.Errorf("directory %s was not created", expected)
	}
}

func TestEnsureDirIdempotent(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	if err := EnsureDir(); err != nil {
		t.Fatalf("first EnsureDir() error: %v", err)
	}
	if err := EnsureDir(); err != nil {
		t.Fatalf("second EnsureDir() error: %v", err)
	}
}

func TestConfigPath(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	path, err := ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath() error: %v", err)
	}
	want := filepath.Join(dir, ".time-broker", "config")
	if path != want {
		t.Errorf("ConfigPath() = %q; want %q", path, want)
	}
}
