VERSION := $(shell cat VERSION)
LDFLAGS := -ldflags "-X main.Version=$(VERSION)"

.PHONY: build frontend-dev frontend-build run

build:
	go build $(LDFLAGS) -o bin/time-broker ./cmd/time-broker/

frontend-dev:
	cd frontend && npm run dev

frontend-build:
	cd frontend && npm run build

run:
	./bin/time-broker
