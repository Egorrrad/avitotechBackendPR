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
	docker compose down

e2e-test:
	go clean -testcache && go test -v ./tests/...

load-test:
	@echo "Running load tests in isolated environment..."
	@if [ ! -f .env.test ]; then \
		echo ".env.test file not found"; \
		exit 1; \
	fi

	# Поднимаем сервисы с правильными переменными
	@docker compose --env-file .env.test -f docker-compose.load-test.yml up -d --build

	@source .env.test && \
			BASE_URL=http://localhost:$$HTTP_PORT \
			curl -f "${BASE_URL}/health" || { \
                echo "Service is not healthy" \
                exit 1 \
            }

	@chmod +x tests/load/run-load-test.sh
	@source .env.test && \
				BASE_URL=http://localhost:$$HTTP_PORT \
				OUTPUT_DIR=tests/load/results \
				./tests/load/run-load-test.sh || true

		@echo "Cleaning up test environment..."
		@docker compose -f docker-compose.load-test.yml down -v
		@echo "Load test completed"

.DEFAULT_GOAL := build