# Makefile for Project Scaffolder CLI

.DEFAULT_GOAL := help

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  install    - Install project dependencies"
	@echo "  build      - Build the Project Scaffolder CLI tool"
	@echo "  run        - Run the Project Scaffolder CLI tool"
	@echo "  clean      - Clean up build artifacts"

.PHONY: install
install:
	go get ./...

.PHONY: build
build: install
	go build -ldflags="-buildmode=exe -s -w -X main.version=dev" -mod=mod -o scaffolder

.PHONY: run
run: build
	./scaffolder

.PHONY: clean
clean:
	go clean -mod=mod
	rm -f ./scaffolder