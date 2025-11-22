.PHONY: build

build:
	go build -v ./cmd/apipullreqs

.DEFAULT_GOAL := build

lint:
	@echo "Running golangci-lint"
	@golangci-lint run ./... --timeout=5m