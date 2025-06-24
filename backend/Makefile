.PHONY: help dev build test clean deps

# Variables
BINARY_NAME=stockle-api
MAIN_PACKAGE=./cmd/api

# Default target
help:
	@echo "Available commands:"
	@echo "  dev           - Run development server with hot reload"
	@echo "  build         - Build the application"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  deps          - Install dependencies"
	@echo "  clean         - Clean build artifacts"

# Development
dev:
	air

# Build
build:
	go build -o $(BINARY_NAME) $(MAIN_PACKAGE)

# Test
test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Dependencies
deps:
	go mod tidy
	go mod download

# Clean
clean:
	go clean
	rm -f $(BINARY_NAME) coverage.out coverage.html