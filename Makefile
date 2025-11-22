.PHONY: build

build:
	go build -v ./cmd/prservice


lint:
	@echo "Running golangci-lint"
	@golangci-lint run ./... --timeout=5m

run:
	docker compose up --build -d

stop:
	docker compose stop

down:
	docker compose down

.DEFAULT_GOAL := build