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
- `--date string` - Date in `MM-DD-YYYY`, `MM/DD/YYYY`, or `M/D/YYYY` format (default: tomorrow)

**Examples:**

```shell
# Schedule a timed event
time-broker schedule event --title "Team Meeting" --timeRange "9:00AM-5:00PM"

# Schedule an all-day event
time-broker schedule event --title "Holiday" --date "12-25-2026"

# Schedule with all flags
time-broker schedule event --title "Focus Time" --description "Deep work session" --timeRange "2:00PM-4:00PM" --date "7/15/2026"
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
- `time-broker view event` - View a specific event by name (fuzzy match)
- `time-broker view day` - View a specific day's schedule
- `time-broker view availability` - View your availability

### view day

View all events for a specific day, ordered earliest to latest.

```shell
time-broker view day                    # shows today's events
time-broker view day --date 01-31-2027  # specific date
time-broker view day --date 1/31/2027   # M/D/YYYY also works
```

Date formats: `MM-DD-YYYY`, `MM/DD/YYYY`, or `M/D/YYYY`

**Output:**

```
* 9:00AM, Team Standup, Weekly sync
* 2:00PM, Dentist, Routine checkup

Timezone: America/New_York
```

### view availability

View your free time over a range of days.

```shell
time-broker view availability
time-broker view availability --range 3
time-broker view availability --startDay 01-15-2027 --startTime 9:00AM --endTime 5:00PM
time-broker view availability --startDay 7/15/2027 --range 3
```

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--range` | Number of days ahead (1-7) | 1 |
| `--startDay` | Start date in `MM-DD-YYYY`, `MM/DD/YYYY`, or `M/D/YYYY` format | today |
| `--startTime` | Start of time window (e.g. `9:00AM`) | 9:00AM |
| `--endTime` | End of time window (e.g. `5:00PM`) | 5:00PM |

**Output:**

```
* Monday 01-12-2026, [9:00AM - 10:30AM], [2:00PM - 3:00PM]
* Tuesday 01-13-2026, [10:00AM - 11:00AM]

Timezone: America/New_York
```

### view event

Search for an event by name using fuzzy matching (Levenshtein distance, 80% threshold).

```shell
time-broker view event --name "team meeting"
time-broker view event --name "standup" --day 1/31/2027
```

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--name` | Search term to match against event names (required) | - |
| `--day` | Day to search in `MM-DD-YYYY`, `MM/DD/YYYY`, or `M/D/YYYY` format | today |

**Output:**

```
* Team Standup
  Time: 9:00AM - 9:30AM
  Description: Weekly team sync
  Location: Conference Room A
```

## version

Print version information.

```shell
time-broker version
```
