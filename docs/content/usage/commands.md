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

(not yet implemented)

```shell
time-broker schedule
```

## update

(not yet implemented)

```shell
time-broker update
```

## get

(not yet implemented)

```shell
time-broker get
```

## version

Print version information.

```shell
time-broker version
```
