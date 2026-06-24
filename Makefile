.PHONY: build frontend-dev frontend-build

build:
	go build -o bin/time-broker ./cmd/time-broker/

frontend-dev:
	cd frontend && npm run dev

frontend-build:
	cd frontend && npm run build
