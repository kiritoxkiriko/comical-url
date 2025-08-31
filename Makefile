# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=shorturl
BINARY_UNIX=$(BINARY_NAME)_unix

# Build the application
.PHONY: build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...

# Install the application
.PHONY: install
install:
	$(GOCMD) install ./...

# Test the application
.PHONY: test
test:
	$(GOTEST) -v ./...

# Test with coverage
.PHONY: test-coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

# Clean build files
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f coverage.out

# Run the application
.PHONY: run
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME) serve

# Cross compilation
.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

# Docker commands
.PHONY: docker-build
docker-build:
	docker build -f deployments/docker/Dockerfile -t $(BINARY_NAME):latest .

.PHONY: docker-run
docker-run:
	docker-compose -f deployments/docker/docker-compose.yaml up

.PHONY: docker-stop
docker-stop:
	docker-compose -f deployments/docker/docker-compose.yaml down

.PHONY: docker-clean
docker-clean:
	docker-compose -f deployments/docker/docker-compose.yaml down -v
	docker rmi $(BINARY_NAME):latest

# Database migration
.PHONY: migrate
migrate:
	./$(BINARY_NAME) migrate

# Development helpers
.PHONY: deps
deps:
	$(GOMOD) download
	$(GOMOD) tidy

.PHONY: lint
lint:
	golangci-lint run

.PHONY: fmt
fmt:
	$(GOCMD) fmt ./...

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  install       - Install the application"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  clean         - Clean build files"
	@echo "  run           - Build and run the application"
	@echo "  build-linux   - Cross-compile for Linux"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run with Docker Compose"
	@echo "  docker-stop   - Stop Docker Compose"
	@echo "  docker-clean  - Clean Docker containers and images"
	@echo "  migrate       - Run database migration"
	@echo "  deps          - Download and tidy dependencies"
	@echo "  lint          - Run linter"
	@echo "  fmt           - Format code"
	@echo "  help          - Show this help"