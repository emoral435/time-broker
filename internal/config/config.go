package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Provider     string `toml:"provider"`
	WeekStartDay string `toml:"week_start_day"`
}

func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home dir: %w", err)
	}
	return filepath.Join(home, ".time-broker", "config"), nil
}

func EnsureDir() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("get home dir: %w", err)
	}
	dir := filepath.Join(home, ".time-broker")
	return os.MkdirAll(dir, 0755)
}

func Load() (*Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	cfg := &Config{}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("read config: %w", err)
	}

	if err := toml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return cfg, nil
}

func Save(c *Config) error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}

	data := fmt.Sprintf(`# Provider - Calendar service to use.
# Currently supported: "google"
# Leave empty to be prompted on next run.
provider = %q

# Week start day - First day of your calendar week.
# Options: "monday", "sunday"
# Leave empty to be prompted on next run.
week_start_day = %q
`, c.Provider, c.WeekStartDay)

	return os.WriteFile(path, []byte(data), 0644)
}

func IsConfigured(c *Config) bool {
	return c.Provider != "" && c.WeekStartDay != ""
}
