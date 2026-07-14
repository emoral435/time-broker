---
title: Contributing
weight: 2
---

## Prerequisites

- Go 1.25 or later
- Node.js and npm (for frontend work)
- A Google Calendar OAuth client ID and secret (for local testing)

## Setup

```shell
git clone https://github.com/emoral435/time-broker
cd time-broker
go mod download
make setup     # installs lefthook and registers precommit hooks
make build     # verify the project compiles
```

## Local Development

### Google Calendar credentials

time-broker needs Google Calendar API credentials to run. There are two ways to
provide them during local development:

**Option 1: `.env` file (recommended for local builds)**

Copy the example file and fill in your credentials:

```shell
cp .env.example .env
# Edit .env with your GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET
```

`make build` reads `.env` and bakes the values into the binary via linker
flags. You only need to re-run `make build` if you change the `.env` file.

**Option 2: Shell environment variables (no rebuild needed)**

Export the variables in your current shell, then run the binary directly:

```shell
export GOOGLE_CLIENT_ID="your-client-id"
export GOOGLE_CLIENT_SECRET="your-client-secret"
make run
```

At runtime, the binary checks for `GOOGLE_CLIENT_ID` and
`GOOGLE_CLIENT_SECRET` environment variables as a fallback when no build-time
values are present.

### Obtaining credentials

1. Go to the [Google Cloud Console](https://console.cloud.google.com/apis/credentials).
2. Create an OAuth 2.0 Client ID (type: Web Application).
3. Add `http://localhost:8085/callback` as an authorized redirect URI.
4. Copy the Client ID and Client Secret into your `.env` file or export them.

### Running

```shell
make build && make run
```

Or, if you exported the environment variables directly:

```shell
make run
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
