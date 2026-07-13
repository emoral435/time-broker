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
- `time-broker schedule view` - View upcoming events

## update

(not yet implemented)

```shell
time-broker update
```

## version

Print version information.

```shell
time-broker version
```
