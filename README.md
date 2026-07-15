<h1 align="center"><code>time-broker</code></h1>
<p align="center">
  <a href="https://github.com/emoral435/time-broker/actions/workflows/goreleaser.yml"
    ><img
      alt="Release"
      src="https://img.shields.io/github/actions/workflow/status/emoral435/time-broker/goreleaser.yml?style=flat-square&label=release"
  /></a>
  <a href="https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-blue?style=flat-square"
    ><img
      alt="Platform"
      src="https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-blue?style=flat-square"
  /></a>
  <a href="https://emoral435.github.io/time-broker/"
    ><img
      alt="Docs"
      src="https://img.shields.io/badge/docs-github.io-blue?style=flat-square"
  /></a>
</p>

<h3 align="center">Get back your time. Use a managed broker, hooked up to your calendar provider.</h3>

<p align="center">
  <img src="https://raw.githubusercontent.com/emoral435/time-broker/main/demo.gif" alt="time-broker demo" width="800" />
</p>

## Overview
time-broker is a TUI based command line tool that connects to popular online-calendar providers to help you manage when you are free, when you have upcoming events, and more.

time-broker was made so that you can easily chat with others and see when you are free to meet. No need to use when2meet when you can directly use a broker to manager when you are free, all at the reach of your terminal.


List of calendar providers supported:
* Google Calendar

Full documentation: https://emoral435.github.io/time-broker/

## Install

To install, you can use the following command:
```sh
curl -fsSL https://raw.githubusercontent.com/emoral435/time-broker/main/install.sh | sh
```

After installing, use `time-broker` or the shorthand `tb` to view the setup wizard.

## Uninstall

To uninstall, run:
```sh
time-broker uninstall
```

This will remove the binary, configuration, tokens, and any symlinks created by the installer.

## How it works
Upon installation, you can use `time-broker` to configure your calendar provider and additional personal settings.

From there, you can use `time-broker help` to see all the possible options you have, but for starters, you have...

```sh
time-broker schedule event # schedule a new event via your calendar provider
```

## Web Frontend

Alongside the TUI, time-broker includes a web frontend located in `frontend/` for a more graphical interface to view and manage your calendar.

```sh
# Start the dev server on localhost:3000
make frontend-dev

# Build for production
make frontend-build
```

The dev server proxies `/api` requests to the Go backend at `localhost:8080`. Configure the target in `frontend/vite.config.ts`.

## Development

### Prerequisites

- Go 1.25+
- golangci-lint
- Node.js (for frontend work)

### Makefile targets

| Target | Description |
|---|---|
| `build` | Build the Go binary to `bin/time-broker` |
| `lint` | Run golangci-lint on all Go files |
| `lint-fix` | Run golangci-lint with auto-fix |
| `vet` | Run `go vet` on all Go packages |
| `test` | Run all Go tests (bypasses cache) |
| `build-all` | Verify all Go packages compile |
| `frontend-dev` | Start the Vite dev server on port 3000 |
| `frontend-build` | Build the frontend for production |
| `frontend-lint` | Run oxlint on frontend TypeScript files |

### Precommit hooks

This project uses [lefthook](https://github.com/evilmartians/lefthook) to run
linting, vetting, building, and testing on every commit. All checks run in
parallel and only trigger when relevant files change.

To enable hooks after cloning:

```sh
make setup
```

This installs lefthook via Homebrew (if missing) and registers the git hooks.

## Contributing

This repository uses [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) for commit messages (e.g. `feat:`, `fix:`, `chore:`). All pull requests are squash-merged into `main`.

Releases are built with [GoReleaser](https://goreleaser.com). To create a release, push a semver tag:

```sh
git tag v1.2.0
git push origin v1.2.0
```

This triggers the goreleaser workflow, which builds binaries for all platforms and publishes a GitHub Release with checksums.

For full contributing guidelines, see the [contributing documentation](https://emoral435.github.io/time-broker/development/contributing/).

## Star History

<a href="https://www.star-history.com/?repos=emoral435%2Ftime-broker&type=date&legend=top-left">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/chart?repos=emoral435/time-broker&type=date&theme=dark&legend=top-left" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/chart?repos=emoral435/time-broker&type=date&legend=top-left" />
   <img alt="Star History Chart" src="https://api.star-history.com/chart?repos=emoral435/time-broker&type=date&legend=top-left" />
 </picture>
</a>