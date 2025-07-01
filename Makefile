.PHONY: run build clean dev test

# Default target
all: run

# Run the application
run:
	go run cmd/server/main.go

# Build the application
build:
	go build -o htmx-app cmd/server/main.go

# Clean build artifacts
clean:
	rm -f htmx-app

# Development with hot reload (requires air)
dev:
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not found. Install with: go install github.com/cosmtrek/air@latest"; \
		echo "Running without hot reload..."; \
		go run cmd/server/main.go; \
	fi

# Run tests
test:
	go test ./...

# Install development dependencies
install-dev:
	go install github.com/cosmtrek/air@latest

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Help
help:
	@echo "Available commands:"
	@echo "  run        - Run the application"
	@echo "  build      - Build the application"
	@echo "  clean      - Clean build artifacts"
	@echo "  dev        - Run with hot reload (requires air)"
	@echo "  test       - Run tests"
	@echo "  install-dev- Install development dependencies"
	@echo "  fmt        - Format code"
	@echo "  lint       - Lint code"
	@echo "  help       - Show this help" 