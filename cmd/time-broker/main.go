package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/charmbracelet/huh"

	"github.com/emoral435/time-broker/internal/config"
	"github.com/emoral435/time-broker/internal/provider/google"
)

var Version = "dev"

func init() {
	loadVersionFromFile()
}

func loadVersionFromFile() {
	if Version != "dev" {
		return
	}
	data, err := os.ReadFile("VERSION")
	if err != nil {
		return
	}
	if v := strings.TrimSpace(string(data)); v != "" {
		Version = v
	}
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		runHelp()
		return nil
	}

	switch args[0] {
	case "help":
		runHelp()
	case "version":
		runVersion()
	case "auth":
		return runAuth()
	case "config":
		return runConfig(args[1:])
	case "schedule", "update", "get":
		return runWithConfig(args[0])
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", args[0])
		runHelp()
		return fmt.Errorf("unknown command: %s", args[0])
	}
	return nil
}

func runHelp() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprint(w, "Usage: time-broker <command>\n\n")
	fmt.Fprint(w, "Commands:\n")
	fmt.Fprint(w, "  auth\tAuthenticate with your calendar provider\n")
	fmt.Fprint(w, "  config\tView or change configuration (run 'config help' for subcommands)\n")
	fmt.Fprint(w, "  get\t(not yet implemented)\n")
	fmt.Fprint(w, "  help\tShow this help message\n")
	fmt.Fprint(w, "  schedule\tSchedule a meeting or view availability\n")
	fmt.Fprint(w, "  update\tCheck for updates\n")
	fmt.Fprint(w, "  version\tPrint version information\n\n")
	fmt.Fprint(w, "Run 'time-broker help <command>' for more details.\n")
	w.Flush()
}

func runVersion() {
	fmt.Printf("time-broker %s\n", Version)
}

func runAuth() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if !config.IsConfigured(cfg) {
		return fmt.Errorf("no configuration found. Run 'time-broker config init' first")
	}

	switch cfg.Provider {
	case "google":
		g := google.New()
		if err := g.Auth(); err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}
	default:
		return fmt.Errorf("unknown provider: %s", cfg.Provider)
	}
	fmt.Println("Authenticated successfully.")
	return nil
}

func runInit() error {
	cfg, err := runSetupWizard()
	if err != nil {
		return err
	}
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("error saving config: %w", err)
	}
	fmt.Println("Configuration saved to ~/.time-broker/config")
	return nil
}

func runConfig(args []string) error {
	if len(args) == 0 {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if !config.IsConfigured(cfg) {
			fmt.Println("No configuration found. Run 'time-broker config init' to set up.")
			return nil
		}
		fmt.Printf("provider: %s\nweek_start_day: %s\n", cfg.Provider, cfg.WeekStartDay)
		return nil
	}

	switch args[0] {
	case "help":
		runConfigHelp()
		return nil
	case "init":
		return runInit()
	case "list":
		return runConfigList()
	default:
		fmt.Fprintf(os.Stderr, "Unknown config subcommand: %s\n\n", args[0])
		runConfigHelp()
		return fmt.Errorf("unknown config subcommand: %s", args[0])
	}
}

func runConfigHelp() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprint(w, "Usage: time-broker config <subcommand>\n\n")
	fmt.Fprint(w, "Subcommands:\n")
	fmt.Fprint(w, "  help\tShow this help message\n")
	fmt.Fprint(w, "  init\tSet up time-broker for first use\n")
	fmt.Fprint(w, "  list\tShow all configuration options\n")
	w.Flush()
}

func runConfigList() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprint(w, "Key\tDescription\tCurrent Value\n")
	fmt.Fprint(w, "---\t-----------\t-------------\n")
	fmt.Fprintf(w, "provider\tCalendar service to use")
	if cfg.Provider != "" {
		fmt.Fprintf(w, "\t%s", cfg.Provider)
	} else {
		fmt.Fprint(w, "\t(not set)")
	}
	fmt.Fprint(w, "\n")
	fmt.Fprintf(w, "week_start_day\tFirst day of your calendar week")
	if cfg.WeekStartDay != "" {
		fmt.Fprintf(w, "\t%s", cfg.WeekStartDay)
	} else {
		fmt.Fprint(w, "\t(not set)")
	}
	fmt.Fprint(w, "\n")
	w.Flush()

	if !config.IsConfigured(cfg) {
		fmt.Println("\nRun 'time-broker config init' to set up configuration.")
	}
	return nil
}

func ensureConfigured() (*config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	if !config.IsConfigured(cfg) {
		cfg, err = runSetupWizard()
		if err != nil {
			return nil, err
		}
		if err := config.Save(cfg); err != nil {
			return nil, err
		}
		fmt.Println("Configuration saved to ~/.time-broker/config")
	}
	return cfg, nil
}

func runWithConfig(cmd string) error {
	cfg, err := ensureConfigured()
	if err != nil {
		return err
	}
	switch cmd {
	case "schedule":
		fmt.Printf("schedule: not yet implemented (configured for %s)\n", cfg.Provider)
	case "update":
		fmt.Println("update: not yet implemented")
	case "get":
		fmt.Println("get: not yet implemented")
	}
	return nil
}

func runSetupWizard() (*config.Config, error) {
	var provider string
	var weekStartDay string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Calendar Provider").
				Description("Choose your calendar service").
				Options(huh.NewOption("Google Calendar", "google")).
				Value(&provider),
			huh.NewSelect[string]().
				Title("Week Start Day").
				Description("First day of your calendar week").
				Options(
					huh.NewOption("Monday", "monday"),
					huh.NewOption("Sunday", "sunday"),
				).
				Value(&weekStartDay),
		),
	).WithTheme(huh.ThemeCharm())

	if err := form.Run(); err != nil {
		return nil, fmt.Errorf("setup wizard: %w", err)
	}

	return &config.Config{
		Provider:     provider,
		WeekStartDay: weekStartDay,
	}, nil
}
