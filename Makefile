.PHONY: build

build:
	go build -v ./cmd/app


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

down-volume:
	docker compose down -v

e2e-test:
	go clean -testcache && go test -v ./tests/e2e/...

load-test:
	@echo "Running load tests in isolated environment..."
	@if [ ! -f .env.test ]; then \
		echo ".env.test file not found"; \
		exit 1; \
	fi

	@docker compose --env-file .env.test -f docker-compose.load-test.yml up -d --build

	@echo "Waiting 20 seconds for services to start..."
	@sleep 20

	@source .env.test && \
	BASE_URL="http://localhost:$$HTTP_PORT" && \
	if ! curl -f "$$BASE_URL/health"; then \
	  echo "Service is not healthy"; \
	  echo "$$BASE_URL/health"; \
	  docker compose -f docker-compose.load-test.yml down -v; \
	  exit 1; \
    fi

	@chmod +x tests/load/run-load-test.sh
	@source .env.test && \
	BASE_URL=http://localhost:$$HTTP_PORT \
	OUTPUT_DIR=results \
	./tests/load/run-load-test.sh || true

	@echo "Cleaning up test environment..."
	@docker compose -f docker-compose.load-test.yml down -v
	@echo "Load test completed"

.DEFAULT_GOAL := build