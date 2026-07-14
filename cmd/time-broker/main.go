package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/charmbracelet/huh"

	"github.com/emoral435/time-broker/internal/config"
	"github.com/emoral435/time-broker/internal/provider/google"
)

const (
	helpAsStr = "help"
	devAsStr  = "dev"
	cfgAsStr  = "config"
	event     = "event"
	schedule  = "schedule"
	update    = "update"
	view      = "view"

	dateFormat = "01-02-2006"
	timeFormat = "3:04PM"
)

var Version = devAsStr

var runSetupWizardFn = runSetupWizard

var (
	githubAPIBase      = "https://api.github.com"
	githubDownloadBase = "https://github.com"
)

func init() {
	loadVersionFromFile()
}

func loadVersionFromFile() {
	if Version != devAsStr {
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
	case helpAsStr:
		runHelp()
	case "version":
		runVersion()
	case "auth":
		return runAuth()
	case cfgAsStr:
		return runConfig(args[1:])
	case schedule:
		return runSchedule(args[1:])
	case view:
		return runView(args[1:])
	case update:
		return runUpdate()
	case "uninstall":
		return runUninstall()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", args[0])
		runHelp()
		return fmt.Errorf("unknown command: %s", args[0])
	}

	// Here, we want to check if this is the first time the user has utilizer their broker. If they haven't, then they should run through their confioguration,
	// otherwise, they can be presented with their options
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if !config.IsConfigured(cfg) {
		if err := runInit(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	}

	return nil
}

func runHelp() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprint(w, "Usage: time-broker <command>\n\n")
	fmt.Fprint(w, "Commands:\n")
	fmt.Fprint(w, "  auth\tAuthenticate with your calendar provider\n")
	fmt.Fprint(w, "  config\tView or change configuration (run 'config help' for subcommands)\n")
	fmt.Fprint(w, "  help\tShow this help message\n")
	fmt.Fprint(w, "  schedule\tManage events on your calendar (run 'schedule help' for subcommands)\n")
	fmt.Fprint(w, "  uninstall\tRemove time-broker and its configuration\n")
	fmt.Fprint(w, "  update\tUpdate time-broker to the latest version\n")
	fmt.Fprint(w, "  version\tPrint version information\n")
	fmt.Fprint(w, "  view\tView events and availability on your calendar\n\n")
	fmt.Fprint(w, "Run 'time-broker help <command>' for more details.\n")
	w.Flush()
}

func runVersion() {
	fmt.Printf("time-broker version %s\n", Version)
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
	case helpAsStr:
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
		cfg, err = runSetupWizardFn()
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

func runSchedule(args []string) error {
	if len(args) == 0 {
		runScheduleHelp()
		return nil
	}

	switch args[0] {
	case helpAsStr:
		runScheduleHelp()
		return nil
	case event:
		return runScheduleEvent(args[1:])
	case "cancel":
		return runScheduleCancel()
	case "update":
		return runScheduleUpdate()
	default:
		fmt.Fprintf(os.Stderr, "Unknown schedule subcommand: %s\n\n", args[0])
		runScheduleHelp()
		return fmt.Errorf("unknown schedule subcommand: %s", args[0])
	}
}

func runScheduleHelp() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprint(w, "Usage: time-broker schedule <subcommand>\n\n")
	fmt.Fprint(w, "Manage events on your calendar.\n\n")
	fmt.Fprint(w, "Subcommands:\n")
	fmt.Fprint(w, "  help\tShow this help message\n")
	fmt.Fprint(w, "  event\tSchedule a new event\n")
	fmt.Fprint(w, "  cancel\tCancel an existing event\n")
	fmt.Fprint(w, "  update\tUpdate an existing event\n")
	w.Flush()
}

func runScheduleEvent(args []string) error {
	if len(args) > 0 && args[0] == helpAsStr {
		runScheduleEventHelp()
		return nil
	}

	fs := flag.NewFlagSet("schedule event", flag.ContinueOnError)
	title := fs.String("title", "Event Title", "event title")
	description := fs.String("description", "Event Description", "event description")
	timeRange := fs.String("timeRange", "", "time range (e.g., 9:00AM-5:00PM)")
	date := fs.String("date", defaultDate(), "date in MM-DD-YYYY format (default: tomorrow)")

	fs.SetOutput(os.Stderr)

	if err := fs.Parse(args); err != nil {
		return err
	}

	parsedDate, err := parseDateFlag(*date)
	if err != nil {
		return fmt.Errorf("invalid date %q: %w", *date, err)
	}

	allDay := *timeRange == ""
	var start, end time.Time

	if !allDay {
		start, end, err = parseTimeRange(*timeRange)
		if err != nil {
			return fmt.Errorf("invalid time range %q: %w", *timeRange, err)
		}
		start = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(),
			start.Hour(), start.Minute(), 0, 0, time.Local)
		end = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(),
			end.Hour(), end.Minute(), 0, 0, time.Local)
	} else {
		start = parsedDate
		end = parsedDate.AddDate(0, 0, 1)
	}

	fmt.Println("\nEvent Details:")
	fmt.Printf("  Title:       %s\n", *title)
	fmt.Printf("  Description: %s\n", *description)
	fmt.Printf("  Date:        %s\n", parsedDate.Format("Monday, January 2, 2006"))
	if allDay {
		fmt.Println("  Time:        All day")
	} else {
		fmt.Printf("  Time:        %s - %s\n", start.Format("3:04 PM"), end.Format("3:04 PM"))
	}
	fmt.Println()

	if !confirmAction("Proceed? [y/N] ") {
		fmt.Println("Cancelled.")
		return nil
	}

	cfg, err := ensureConfigured()
	if err != nil {
		return err
	}

	switch cfg.Provider {
	case "google":
		g := google.New()
		if err := g.Book(*title, *description, start, end, allDay); err != nil {
			return fmt.Errorf("failed to book event: %w", err)
		}
	default:
		return fmt.Errorf("unknown provider: %s", cfg.Provider)
	}

	fmt.Printf("Successfully booked %q\n", *title)
	return nil
}

func runScheduleEventHelp() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprint(w, "Usage: time-broker schedule event [flags]\n\n")
	fmt.Fprint(w, "Schedule a new event on your calendar.\n\n")
	fmt.Fprint(w, "Flags:\n")
	fmt.Fprint(w, "  --title string\tEvent title (default \"Event Title\")\n")
	fmt.Fprint(w, "  --description string\tEvent description (default \"Event Description\")\n")
	fmt.Fprint(w, "  --timeRange string\tTime range in H:MMAM-H:MMPM format (default: all day)\n")
	fmt.Fprint(w, "  --date string\t\tDate in MM-DD-YYYY format (default: tomorrow)\n\n")
	fmt.Fprint(w, "Examples:\n")
	fmt.Fprint(w, "  time-broker schedule event --title \"Team Meeting\" --timeRange \"9:00AM-5:00PM\"\n")
	fmt.Fprint(w, "  time-broker schedule event --title \"Holiday\" --date \"12-25-2026\"\n")
	fmt.Fprint(w, "  time-broker schedule event --title \"Focus Time\" --timeRange \"2:00PM-4:00PM\" --date \"07-15-2026\"\n")
	w.Flush()
}

func parseTimeRange(s string) (start, end time.Time, err error) {
	parts := strings.SplitN(s, "-", 2)
	if len(parts) != 2 {
		return time.Time{}, time.Time{}, fmt.Errorf("expected format H:MMAM-H:MMPM")
	}

	start, err = time.Parse(timeFormat, strings.TrimSpace(parts[0]))
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid start time %q: %w", parts[0], err)
	}

	end, err = time.Parse(timeFormat, strings.TrimSpace(parts[1]))
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid end time %q: %w", parts[1], err)
	}

	return start, end, nil
}

func parseDateFlag(s string) (time.Time, error) {
	return time.Parse(dateFormat, s)
}

func defaultDate() string {
	return time.Now().AddDate(0, 0, 1).Format(dateFormat)
}

func confirmAction(prompt string) bool {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

func runScheduleCancel() error {
	_, err := ensureConfigured()
	if err != nil {
		return err
	}
	fmt.Println("schedule cancel: not yet implemented")
	return nil
}

func runScheduleUpdate() error {
	_, err := ensureConfigured()
	if err != nil {
		return err
	}
	fmt.Println("schedule update: not yet implemented")
	return nil
}

func runView(args []string) error {
	if len(args) == 0 {
		runViewHelp()
		return nil
	}

	switch args[0] {
	case helpAsStr:
		runViewHelp()
		return nil
	case event:
		return runViewEvent()
	case "day":
		return runViewDay(args[1:])
	case "availability":
		return runViewAvailability()
	default:
		fmt.Fprintf(os.Stderr, "Unknown view subcommand: %s\n\n", args[0])
		runViewHelp()
		return fmt.Errorf("unknown view subcommand: %s", args[0])
	}
}

func runViewHelp() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprint(w, "Usage: time-broker view <subcommand>\n\n")
	fmt.Fprint(w, "View events and availability on your calendar.\n\n")
	fmt.Fprint(w, "Subcommands:\n")
	fmt.Fprint(w, "  help\tShow this help message\n")
	fmt.Fprint(w, "  event\tView a specific event\n")
	fmt.Fprint(w, "  day\tView a specific day's schedule\n")
	fmt.Fprint(w, "  availability\tView your availability\n")
	w.Flush()
}

func runViewEvent() error {
	_, err := ensureConfigured()
	if err != nil {
		return err
	}
	fmt.Println("view event: not yet implemented")
	return nil
}

func runViewDay(args []string) error {
	_, err := ensureConfigured()
	if err != nil {
		return err
	}
	if len(args) == 0 {
		fmt.Println("view day: not yet implemented (defaults to today)")
	} else {
		fmt.Printf("view day: not yet implemented (date: %s)\n", args[0])
	}
	return nil
}

func runViewAvailability() error {
	_, err := ensureConfigured()
	if err != nil {
		return err
	}
	fmt.Println("view availability: not yet implemented")
	return nil
}

func runUpdate() error {
	fmt.Printf("time-broker %s\n", Version)
	fmt.Println("Checking for updates...")

	latest, err := latestVersion()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	if Version == latest {
		fmt.Println("Already up to date.")
		return nil
	}

	fmt.Printf("Updating from %s to %s...\n", Version, latest)
	if err := downloadAndUpdate(latest); err != nil {
		return fmt.Errorf("update failed: %w", err)
	}
	fmt.Printf("Successfully updated to %s\n", latest)
	return nil
}

const githubRepo = "emoral435/time-broker"

type githubRelease struct {
	TagName string `json:"tag_name"`
	Body    string `json:"body"`
}

func latestVersion() (string, error) {
	url := fmt.Sprintf("%s/repos/%s/releases/latest", githubAPIBase, githubRepo)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	return release.TagName, nil
}

func downloadAndUpdate(version string) error {
	osName := runtime.GOOS
	arch := runtime.GOARCH

	tarball := fmt.Sprintf("time-broker-%s-%s-%s.tar.gz", version, osName, arch)
	url := fmt.Sprintf("%s/%s/releases/download/%s/%s", githubDownloadBase, githubRepo, version, tarball)

	tmpDir, err := os.MkdirTemp("", "time-broker-update-*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	tarballPath := filepath.Join(tmpDir, tarball)
	if err := downloadFile(url, tarballPath); err != nil {
		return fmt.Errorf("failed to download %s: %w", tarball, err)
	}

	if err := exec.Command("tar", "xzf", tarballPath, "-C", tmpDir).Run(); err != nil {
		return fmt.Errorf("failed to extract tarball: %w", err)
	}

	newBinary := filepath.Join(tmpDir, "time-broker")
	if _, err := os.Stat(newBinary); err != nil {
		return fmt.Errorf("extracted binary not found: %w", err)
	}

	currentBinary, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to find current binary: %w", err)
	}
	currentBinary, err = filepath.EvalSymlinks(currentBinary)
	if err != nil {
		return fmt.Errorf("failed to resolve current binary path: %w", err)
	}

	if err := os.Rename(newBinary, currentBinary); err != nil {
		return fmt.Errorf("failed to replace binary at %s: %w", currentBinary, err)
	}

	return nil
}

func runUninstall() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to find home directory: %w", err)
	}

	dataDir := filepath.Join(home, ".time-broker")

	// Determine symlink directories to check, matching install.sh logic.
	linkDirs := []string{}
	localBin := filepath.Join(home, ".local", "bin")
	if _, err := os.Stat(localBin); err == nil {
		linkDirs = append(linkDirs, localBin)
	}
	linkDirs = append(linkDirs, "/usr/local/bin")

	binNames := []string{"time-broker", "tb"}
	var symlinksToRemove []string

	for _, dir := range linkDirs {
		for _, name := range binNames {
			path := filepath.Join(dir, name)
			link, err := os.Readlink(path)
			if err != nil {
				continue
			}
			// Only remove symlinks that point into our data directory.
			if strings.Contains(link, dataDir) || strings.HasPrefix(link, dataDir) {
				symlinksToRemove = append(symlinksToRemove, path)
			}
		}
	}

	fmt.Println("The following will be removed:")
	fmt.Printf("  Directory: %s\n", dataDir)
	if len(symlinksToRemove) > 0 {
		for _, s := range symlinksToRemove {
			fmt.Printf("  Symlink:   %s\n", s)
		}
	}
	fmt.Println()

	if !confirmAction("Proceed with uninstall? [y/N] ") {
		fmt.Println("Uninstall cancelled.")
		return nil
	}

	// Remove symlinks.
	for _, s := range symlinksToRemove {
		if err := os.Remove(s); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not remove %s: %v\n", s, err)
		} else {
			fmt.Printf("Removed %s\n", s)
		}
	}

	// Remove data directory.
	if err := os.RemoveAll(dataDir); err != nil {
		return fmt.Errorf("failed to remove %s: %w", dataDir, err)
	}
	fmt.Printf("Removed %s\n", dataDir)

	fmt.Println("time-broker has been uninstalled.")
	return nil
}

func downloadFile(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download returned status %d", resp.StatusCode)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
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
