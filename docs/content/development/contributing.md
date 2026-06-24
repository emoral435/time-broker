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
go build -o bin/time-broker ./cmd/time-broker/
```

## Running Tests

```shell
go test ./...
```

## Submitting Changes

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run `go mod tidy` and `go test ./...`
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
[tag-on-version-change](https://github.com/emoral435/time-broker/actions/workflows/tag-on-version-change.yml)
workflow, which creates a `v$(VERSION)` tag. The tag then triggers the
[release](https://github.com/emoral435/time-broker/actions/workflows/release.yml)
workflow, which builds binaries for all platforms, generates checksums, syncs
the frontend version, and creates a GitHub Release with auto-generated release
notes.
