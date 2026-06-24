package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/charmbracelet/huh"

	"github.com/emoral435/time-broker/internal/config"
)

var Version = "dev"

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "help":
			runHelp()
		case "version":
			runVersion()
		case "init":
			runInit()
		case "config":
			runConfig(os.Args[1:])
		case "schedule", "update", "get":
			runWithConfig(os.Args[1])
		default:
			fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", os.Args[1])
			runHelp()
			os.Exit(1)
		}
		return
	}
	runHelp()
}

func runHelp() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprint(w, "Usage: time-broker <command>\n\n")
	fmt.Fprint(w, "Commands:\n")
	fmt.Fprint(w, "  config\tView or change configuration (run 'config help' for subcommands)\n")
	fmt.Fprint(w, "  get\t(not yet implemented)\n")
	fmt.Fprint(w, "  help\tShow this help message\n")
	fmt.Fprint(w, "  init\tSet up time-broker for first use\n")
	fmt.Fprint(w, "  schedule\tSchedule a meeting or view availability\n")
	fmt.Fprint(w, "  update\tCheck for updates\n")
	fmt.Fprint(w, "  version\tPrint version information\n\n")
	fmt.Fprint(w, "Run 'time-broker help <command>' for more details.\n")
	w.Flush()
}

func runVersion() {
	fmt.Printf("time-broker %s\n", Version)
}

func runInit() {
	cfg, err := runSetupWizard()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if err := config.Save(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Configuration saved to ~/.time-broker/config")
}

func runConfig(args []string) {
	if len(args) == 1 {
		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if !config.IsConfigured(cfg) {
			fmt.Println("No configuration found. Run 'time-broker config init' to set up.")
			return
		}
		fmt.Printf("provider: %s\nweek_start_day: %s\n", cfg.Provider, cfg.WeekStartDay)
		return
	}

	switch args[1] {
	case "help":
		runConfigHelp()
	case "init":
		runInit()
	case "list":
		runConfigList()
	default:
		fmt.Fprintf(os.Stderr, "Unknown config subcommand: %s\n\n", args[1])
		runConfigHelp()
		os.Exit(1)
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

func runConfigList() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("provider\tCalendar service to use")
	if cfg.Provider != "" {
		fmt.Printf(" (currently: %s)", cfg.Provider)
	}
	fmt.Println()
	fmt.Printf("week_start_day\tFirst day of your calendar week")
	if cfg.WeekStartDay != "" {
		fmt.Printf(" (currently: %s)", cfg.WeekStartDay)
	}
	fmt.Println()
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

func runWithConfig(cmd string) {
	cfg, err := ensureConfigured()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	switch cmd {
	case "schedule":
		fmt.Printf("schedule: not yet implemented (configured for %s)\n", cfg.Provider)
	case "update":
		fmt.Println("update: not yet implemented")
	case "get":
		fmt.Println("get: not yet implemented")
	}
}

func runSetupWizard() (*config.Config, error) {
	var provider string
	var weekStartDay string
	var setupAlias bool

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
		huh.NewGroup(
			huh.NewConfirm().
				Title("Create 'tb' alias?").
				Description("Add 'alias tb=time-broker' to your shell config so you can use 'tb' as a shorthand.").
				Value(&setupAlias),
		),
	).WithTheme(huh.ThemeCharm())

	if err := form.Run(); err != nil {
		return nil, fmt.Errorf("setup wizard: %w", err)
	}

	if setupAlias {
		if err := addShellAlias(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not create alias: %v\n", err)
		}
	}

	return &config.Config{
		Provider:     provider,
		WeekStartDay: weekStartDay,
	}, nil
}

func addShellAlias() error {
	shell := os.Getenv("SHELL")
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	var rcFile string
	var aliasLine string

	switch {
	case strings.HasSuffix(shell, "/zsh"):
		rcFile = filepath.Join(home, ".zshrc")
		aliasLine = "alias tb='time-broker'"
	case strings.HasSuffix(shell, "/bash"):
		rcFile = filepath.Join(home, ".bashrc")
		aliasLine = "alias tb='time-broker'"
	case strings.HasSuffix(shell, "/fish"):
		rcFile = filepath.Join(home, ".config", "fish", "config.fish")
		aliasLine = "alias tb=time-broker"
	default:
		fmt.Println("To create the 'tb' alias manually, add to your shell config:")
		fmt.Println("  alias tb='time-broker'")
		return nil
	}

	if data, err := os.ReadFile(rcFile); err == nil {
		if strings.Contains(string(data), "alias tb=") {
			fmt.Println("Alias 'tb' already configured in", rcFile)
			return nil
		}
	}

	f, err := os.OpenFile(rcFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("open %s: %w", rcFile, err)
	}
	defer f.Close()

	if _, err := f.WriteString("\n# time-broker alias\n" + aliasLine + "\n"); err != nil {
		return err
	}

	fmt.Printf("Alias 'tb' added to %s\n", rcFile)
	fmt.Println("Run 'source " + rcFile + "' or restart your terminal to use 'tb'.")
	return nil
}
