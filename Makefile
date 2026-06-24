.PHONY: build frontend-dev frontend-build
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

.PHONY: build run

build:
	go build -o bin/time-broker ./cmd/time-broker/
	go build -ldflags="-X main.Version=$(VERSION)" -o bin/time-broker ./cmd/time-broker/

frontend-dev:
	cd frontend && npm run dev

frontend-build:
	cd frontend && npm run build

run:
	./bin/time-broker
