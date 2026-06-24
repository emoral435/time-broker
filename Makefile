VERSION := $(shell cat VERSION)
COMMIT  := $(shell git describe --always --dirty 2>/dev/null || echo "unknown")
DATE    := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

.PHONY: build frontend-dev frontend-build

build:
	go build -o bin/time-broker ./cmd/time-broker/
	go build $(LDFLAGS) -o bin/time-broker ./cmd/time-broker/

frontend-dev:
	cd frontend && npm run dev

frontend-build:
	cd frontend && npm run build

run:
	./bin/time-broker
