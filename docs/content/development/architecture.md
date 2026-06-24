---
title: Architecture
weight: 1
---

## Overview

time-broker is structured as a standard Go CLI application:

```
cmd/
  time-broker/
    main.go              # Entry point
internal/
  providers/             # Calendar provider implementations
    google.go
  tui/                   # Terminal UI components
  commands/              # CLI command definitions
  config/                # Configuration management
```

## Design Principles

- **TUI-first**: The primary interface is the terminal; HTTP serving is secondary
- **Provider-agnostic**: Each calendar provider implements a common interface
- **Simple CLI**: Commands are flat (no nested subcommands where possible)
- **Minimal dependencies**: Prefer the standard library and well-maintained packages

## Provider Interface

All calendar providers implement the following contract:

```go
type Provider interface {
    Auth() error
    FreeSlots(day time.Time, minDuration time.Duration) ([]Slot, error)
    Book(slot Slot, title string) error
}
```
