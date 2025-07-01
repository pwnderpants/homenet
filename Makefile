# Default target
all: build

# Build the application
build:
	go build -o homenet cmd/server/main.go

# Clean build artifacts
clean:
	rm -f homenet
dev:
	@if command -v air > /dev/null; then \
		air cmd/server/main.go; \
	else \
		echo "Air not found. Install with: go install github.com/air-verse/air@latest"; \
		echo "Running without hot reload..."; \
		go run cmd/server/main.go; \
	fi

# Help
help:
	@echo "Available commands:"
	@echo "  build      - Build the application"
	@echo "  clean      - Clean build artifacts"
	@echo "  dev        - Run with hot reload (requires air)"
	@echo "  help       - Show this help" 
