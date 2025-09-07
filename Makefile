# GophKeeper Makefile

# Variables
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GO_VERSION := $(shell go version | awk '{print $$3}')

# Build flags
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT) -X main.goVersion=$(GO_VERSION)"

# Build directories
BIN_DIR := bin
CLIENT_BINARY := $(BIN_DIR)/client
SERVER_BINARY := $(BIN_DIR)/gophkeeper

# Default target
.PHONY: all
all: build

# Build both client and server
.PHONY: build
build: $(CLIENT_BINARY) $(SERVER_BINARY)

# Build client
$(CLIENT_BINARY):
	@echo "Building client..."
	@mkdir -p $(BIN_DIR)
	go build $(LDFLAGS) -o $(CLIENT_BINARY) ./cmd/client

# Build server
$(SERVER_BINARY):
	@echo "Building server..."
	@mkdir -p $(BIN_DIR)
	go build $(LDFLAGS) -o $(SERVER_BINARY) ./cmd/gophkeeper

# Build for multiple platforms
.PHONY: build-all
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(BIN_DIR)
	
	# Linux AMD64
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BIN_DIR)/client-linux-amd64 ./cmd/client
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BIN_DIR)/gophkeeper-linux-amd64 ./cmd/gophkeeper
	
	# Windows AMD64
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BIN_DIR)/client-windows-amd64.exe ./cmd/client
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BIN_DIR)/gophkeeper-windows-amd64.exe ./cmd/gophkeeper
	
	# macOS AMD64
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BIN_DIR)/client-darwin-amd64 ./cmd/client
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BIN_DIR)/gophkeeper-darwin-amd64 ./cmd/gophkeeper
	
	# macOS ARM64 (Apple Silicon)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BIN_DIR)/client-darwin-arm64 ./cmd/client
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BIN_DIR)/gophkeeper-darwin-arm64 ./cmd/gophkeeper

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -cover ./...

# Run tests with coverage report
.PHONY: test-coverage-html
test-coverage-html:
	@echo "Running tests with HTML coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BIN_DIR)
	rm -f coverage.out coverage.html

# Show version information
.PHONY: version
version:
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Go Version: $(GO_VERSION)"

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build              - Build both client and server"
	@echo "  build-all          - Build for multiple platforms (Linux, Windows, macOS)"
	@echo "  test               - Run tests"
	@echo "  test-coverage      - Run tests with coverage"
	@echo "  test-coverage-html - Run tests with HTML coverage report"
	@echo "  clean              - Clean build artifacts"
	@echo "  version            - Show version information"
	@echo "  help               - Show this help"
	@echo ""
	@echo "Variables:"
	@echo "  VERSION            - Version to build (default: git tag or 'dev')"
	@echo ""
	@echo "Examples:"
	@echo "  make build                    # Build with default version"
	@echo "  make VERSION=1.2.3 build     # Build with specific version"
	@echo "  make build-all               # Build for all platforms"
	@echo "  make test-coverage-html      # Generate coverage report"
