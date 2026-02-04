.PHONY: build install clean run-service deps help

build:
	@echo "Building CLI..."
	@go build -o bin/subtrack-cli ./cmd/cli
	@echo "Building service..."
	@go build -o bin/subtrack-service ./cmd/service
	@echo "Build complete!"

install:
	@echo "Installing CLI..."
	@go install ./cmd/cli
	@echo "Installing service..."
	@go install ./cmd/service
	@echo "Install complete!"

run-service:
	@./bin/subtrack-service

deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download
	@echo "Dependencies installed!"

clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@echo "Clean complete!"

help:
	@echo "SubTrack Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  build        - Build CLI and service binaries"
	@echo "  install      - Install CLI and service to GOPATH/bin"
	@echo "  run-service  - Run the service (requires building first)"
	@echo "  deps         - Install and tidy Go dependencies"
	@echo "  clean        - Remove build artifacts"
	@echo "  help         - Show this help message"
