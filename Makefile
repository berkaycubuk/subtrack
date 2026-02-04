.PHONY: build install clean run-service deps help docker-cli docker-list docker-add docker-update docker-delete docker-check docker-health

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

docker-cli:
	@docker-compose run --rm cli ${ARGS}

docker-list:
	@make docker-cli ARGS=list

docker-add:
	@make docker-cli ARGS='add ${ARGS}'

docker-update:
	@make docker-cli ARGS='update ${ARGS}'

docker-delete:
	@make docker-cli ARGS='delete ${ARGS}'

docker-check:
	@make docker-cli ARGS=check

docker-health:
	@make docker-cli ARGS=health

help:
	@echo "SubTrack Makefile"
	@echo ""
	@echo "Local targets:"
	@echo "  build        - Build CLI and service binaries"
	@echo "  install      - Install CLI and service to GOPATH/bin"
	@echo "  run-service  - Run the service (requires building first)"
	@echo "  deps         - Install and tidy Go dependencies"
	@echo "  clean        - Remove build artifacts"
	@echo ""
	@echo "Docker CLI targets:"
	@echo "  docker-list            - List all subscriptions"
	@echo "  docker-add ARGS        - Add subscription (e.g., make docker-add ARGS=\"Netflix 15.99 USD monthly 15-02-2025\")"
	@echo "  docker-update ARGS     - Update subscription (e.g., make docker-update ARGS=\"1 Netflix 19.99 USD monthly 15-03-2025\")"
	@echo "  docker-delete ARGS     - Delete subscription (e.g., make docker-delete ARGS=1)"
	@echo "  docker-check           - Check upcoming payments"
	@echo "  docker-health          - Check Telegram bot health"
	@echo "  docker-cli ARGS        - Run custom CLI command"
	@echo ""
	@echo "  help         - Show this help message"
