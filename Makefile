.PHONY: build

build:
	go build -v ./cmd/apipullreqs


lint:
	@echo "Running golangci-lint"
	@golangci-lint run ./... --timeout=5m

.DEFAULT_GOAL := build