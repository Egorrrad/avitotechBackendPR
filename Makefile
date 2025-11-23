.PHONY: build

build:
	go build -v ./cmd/prservice


lint:
	@echo "Running golangci-lint"
	@golangci-lint run ./... --config .golangci.yml --timeout=5m

lint-fix:
	@echo "Running golangci-lint fix"
	golangci-lint run ./... --config .golangci.yml --fix

run:
	docker compose up --build -d

stop:
	docker compose stop

down:
	docker compose down

.DEFAULT_GOAL := build