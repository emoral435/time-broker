# Load Google Calendar credentials from .env if present.
# These are baked into the binary via ldflags at build time.
# The binary also checks GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET
# environment variables at runtime as a fallback.
ifneq (,$(wildcard .env))
include .env
endif

VERSION := $(shell cat version.txt)
XFLAGS := -X main.Version=$(VERSION)

ifdef GOOGLE_CLIENT_ID
XFLAGS += -X github.com/emoral435/time-broker/internal/provider/google.ClientID=$(GOOGLE_CLIENT_ID)
endif
ifdef GOOGLE_CLIENT_SECRET
XFLAGS += -X github.com/emoral435/time-broker/internal/provider/google.ClientSecret=$(GOOGLE_CLIENT_SECRET)
endif

.PHONY: build frontend-dev frontend-build frontend-lint run lint lint-fix vet test build-all setup install-hooks

build:
	go build -ldflags "$(XFLAGS)" -o bin/time-broker ./cmd/time-broker/

frontend-dev:
	cd frontend && npm run dev

frontend-build:
	cd frontend && npm ci && npm run build

run:
	./bin/time-broker

lint:
	golangci-lint run ./...

lint-fix:
	golangci-lint run --fix ./...

test:
	go test ./... -v -count=1

test-short:
	go test ./... -short -v -count=1

vet:
	go vet ./...

build-all:
	go build ./...

frontend-lint:
	cd frontend && npm run lint

setup: install-hooks

install-hooks:
	@which lefthook > /dev/null 2>&1 || brew install lefthook
	lefthook install