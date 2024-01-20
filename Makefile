# Makefile for Go CLI Project

.DEFAULT_GOAL := help

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  install    - Install project dependencies"
	@echo "  build      - Build the Template CLI tool"
	@echo "  run        - Run the Template CLI tool"
	@echo "  clean      - Clean up build artifacts"

.PHONY: install
install:
	go get ./...

.PHONY: build
build: install
	go build -ldflags="-buildmode=exe -s -w -X main.version=dev" -mod=mod -o template-cli

.PHONY: run
run: build
	./template-cli

.PHONY: clean
clean:
	go clean -mod=mod
	rm -f ./template-cli
