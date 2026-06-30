---
title: Contributing
weight: 2
---

## Prerequisites

- Go 1.25 or later
- A calendar provider API key (for integration tests)

## Setup

```shell
git clone https://github.com/emoral435/time-broker
cd time-broker
go mod download
make setup     # installs lefthook and registers precommit hooks
make build     # verify the project compiles
```

## Precommit hooks

This project uses [lefthook](https://github.com/evilmartians/lefthook) to run
linting, vetting, building, and testing on every commit. All checks run in
parallel and only trigger when relevant files change:

| Hook | Trigger | What it runs |
|---|---|---|
| `go-vet` | `.go` files change | `go vet ./...` |
| `go-lint` | `.go` files change | `golangci-lint run ./...` |
| `go-build` | `.go` files change | `go build ./...` |
| `go-test` | `.go` files change | `go test ./... -count=1` |
| `frontend-lint` | `.ts/.tsx` files change | `oxlint` via `npm run lint` |
| `frontend-build` | `.ts/.tsx/.css/.json` files change | `vite build` via `npm run build` |

If any check fails, the commit is blocked. Run the failing check directly to
debug (e.g. `make lint`, `make test`).

To enable hooks after cloning: `make setup`. To bypass hooks for a quick
iteration (not recommended): `git commit --no-verify`.

## Makefile commands

Common development commands are available via the Makefile:

| Target | Description |
|---|---|
| `build` | Build the Go binary to `bin/time-broker` |
| `lint` | Run golangci-lint on all Go files |
| `lint-fix` | Run golangci-lint with auto-fix |
| `vet` | Run `go vet` on all Go packages |
| `test` | Run all Go tests (bypasses cache with `-count=1`) |
| `build-all` | Verify all Go packages compile |
| `frontend-dev` | Start the Vite dev server on port 3000 |
| `frontend-build` | Build the frontend for production |
| `frontend-lint` | Run oxlint on frontend TypeScript files |

## Submitting Changes

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run `make test` and `make lint` to verify locally
5. Submit a pull request

## Documentation

Documentation lives in the `docs/` directory and is built with Hugo + Hextra.
To preview documentation changes locally:

```shell
cd docs
hugo server -D
```

Then open http://localhost:1313.

## Releasing

1. Update the `VERSION` file at the project root to the new version (without `v` prefix, e.g. `0.2.0`).
2. Commit the change and open a pull request.
3. Merge the PR to `main`.

That's it. Merging a change to `VERSION` on `main` triggers the
[release](https://github.com/emoral435/time-broker/actions/workflows/release.yml)
workflow, which creates the tag, builds binaries for all platforms, generates
checksums, syncs the frontend version, and creates a GitHub Release with
auto-generated release notes.
