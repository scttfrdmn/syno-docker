# SynoDeploy Makefile
# Copyright 2025 Scott Friedman

.PHONY: build test clean install dev-install fmt vet lint staticcheck ineffassign misspell gocyclo deps quality-check pre-commit help

# Build variables
BINARY_NAME := synodeploy
BUILD_DIR := bin
MAIN_PACKAGE := ./main.go
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Go build flags
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)
BUILDFLAGS := -ldflags "$(LDFLAGS)"

# Tool versions
STATICCHECK_VERSION := latest
INEFFASSIGN_VERSION := latest
MISSPELL_VERSION := latest
GOCYCLO_VERSION := latest
GOLANGCI_LINT_VERSION := v1.55.2

# Default target
all: quality-check test build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(BUILDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)

# Run tests
test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)/
	rm -f coverage.out

# Install binary to system
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

# Development install (creates symlink)
dev-install: build
	@echo "Creating development symlink..."
	ln -sf $(PWD)/$(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Run golangci-lint
lint:
	@echo "Running golangci-lint..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(GOLANGCI_LINT_VERSION); \
	fi
	golangci-lint run

# Run staticcheck
staticcheck:
	@echo "Running staticcheck..."
	@if ! command -v staticcheck >/dev/null 2>&1; then \
		echo "Installing staticcheck..."; \
		go install honnef.co/go/tools/cmd/staticcheck@$(STATICCHECK_VERSION); \
	fi
	staticcheck ./...

# Run ineffassign
ineffassign:
	@echo "Running ineffassign..."
	@if ! command -v ineffassign >/dev/null 2>&1; then \
		echo "Installing ineffassign..."; \
		go install github.com/gordonklaus/ineffassign@$(INEFFASSIGN_VERSION); \
	fi
	ineffassign ./...

# Run misspell
misspell:
	@echo "Running misspell..."
	@if ! command -v misspell >/dev/null 2>&1; then \
		echo "Installing misspell..."; \
		go install github.com/client9/misspell/cmd/misspell@$(MISSPELL_VERSION); \
	fi
	misspell -error .

# Run gocyclo
gocyclo:
	@echo "Running gocyclo..."
	@if ! command -v gocyclo >/dev/null 2>&1; then \
		echo "Installing gocyclo..."; \
		go install github.com/fzipp/gocyclo/cmd/gocyclo@$(GOCYCLO_VERSION); \
	fi
	gocyclo -over 10 .

# Install development dependencies
deps:
	@echo "Installing development dependencies..."
	go mod download
	go install honnef.co/go/tools/cmd/staticcheck@$(STATICCHECK_VERSION)
	go install github.com/gordonklaus/ineffassign@$(INEFFASSIGN_VERSION)
	go install github.com/client9/misspell/cmd/misspell@$(MISSPELL_VERSION)
	go install github.com/fzipp/gocyclo/cmd/gocyclo@$(GOCYCLO_VERSION)
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(GOLANGCI_LINT_VERSION); \
	fi

# Run all quality checks (Go Report Card equivalent)
quality-check: fmt vet lint staticcheck ineffassign misspell gocyclo
	@echo "All quality checks passed! ✅"

# Pre-commit checks (used by git hooks)
pre-commit: quality-check test
	@echo "Pre-commit checks completed! ✅"

# Coverage report
coverage: test
	@echo "Generating coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build $(BUILDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PACKAGE)
	GOOS=darwin GOARCH=arm64 go build $(BUILDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PACKAGE)
	GOOS=linux GOARCH=amd64 go build $(BUILDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PACKAGE)
	GOOS=linux GOARCH=arm64 go build $(BUILDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PACKAGE)

# Release build
release: clean quality-check test build-all
	@echo "Release build completed! 📦"

# Show help
help:
	@echo "SynoDeploy Makefile Commands:"
	@echo ""
	@echo "Build Commands:"
	@echo "  build       Build the binary"
	@echo "  build-all   Build for multiple platforms"
	@echo "  release     Full release build with all checks"
	@echo "  clean       Remove build artifacts"
	@echo ""
	@echo "Development Commands:"
	@echo "  install     Install binary to /usr/local/bin"
	@echo "  dev-install Create development symlink"
	@echo "  deps        Install development dependencies"
	@echo ""
	@echo "Quality Commands:"
	@echo "  fmt         Format code"
	@echo "  vet         Run go vet"
	@echo "  lint        Run golangci-lint"
	@echo "  staticcheck Run staticcheck"
	@echo "  ineffassign Run ineffassign"
	@echo "  misspell    Run misspell"
	@echo "  gocyclo     Run gocyclo"
	@echo "  quality-check  Run all quality checks"
	@echo ""
	@echo "Testing Commands:"
	@echo "  test        Run tests"
	@echo "  coverage    Generate coverage report"
	@echo "  pre-commit  Run pre-commit checks"
	@echo ""
	@echo "Variables:"
	@echo "  VERSION     Version string (default: git describe)"
	@echo "  COMMIT      Git commit hash (default: git rev-parse)"
	@echo "  DATE        Build date (default: current UTC)"