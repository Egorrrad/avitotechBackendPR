.PHONY: build

build:
	go build -v ./cmd/apipullreqs

.DEFAULT_GOAL := build