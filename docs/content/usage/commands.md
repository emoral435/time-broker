---
title: Commands
weight: 1
---

## auth

Authenticate with your calendar provider.

```shell
time-broker auth
```

Opens a browser to authorize time-broker to access your calendar. Tokens
are saved to `~/.time-broker/tokens.json` and auto-refreshed.

## config

View or change configuration.

**Subcommands:**

- `time-broker config help` - Show config subcommands
- `time-broker config init` - Run the setup wizard
- `time-broker config list` - Show all configuration options and current values

## help

Show the help message listing all available commands.

```shell
time-broker help
```

## schedule

Manage events on your calendar.

**Subcommands:**

- `time-broker schedule help` - Show schedule subcommands
- `time-broker schedule event` - Schedule a new event
- `time-broker schedule cancel` - Cancel an existing event
- `time-broker schedule update` - Update an existing event

### schedule event

Schedule a new event on your calendar.

```shell
time-broker schedule event --title "Team Meeting" --timeRange "9:00AM-5:00PM"
```

**Flags:**

- `--title string` - Event title (default: "Event Title")
- `--description string` - Event description (default: "Event Description")
- `--timeRange string` - Time range in `H:MMAM-H:MMPM` format (default: all day)
- `--date string` - Date in `MM-DD-YYYY` format (default: tomorrow)

**Examples:**

```shell
# Schedule a timed event
time-broker schedule event --title "Team Meeting" --timeRange "9:00AM-5:00PM"

# Schedule an all-day event
time-broker schedule event --title "Holiday" --date "12-25-2026"

# Schedule with all flags
time-broker schedule event --title "Focus Time" --description "Deep work session" --timeRange "2:00PM-4:00PM" --date "07-15-2026"
```

## update

Update time-broker to the latest version.

```shell
time-broker update
```

## view

View events and availability on your calendar.

**Subcommands:**

- `time-broker view help` - Show view subcommands
- `time-broker view event` - View a specific event
- `time-broker view day` - View a specific day's schedule
- `time-broker view availability` - View your availability

## version

Print version information.

```shell
time-broker version
```
