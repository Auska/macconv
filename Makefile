# Makefile for macconv project

# Variables
GO := go
BINARY_NAME := macconv
PKG := .
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "0.1.3")
BUILD_DATE := $(shell date -u +'%Y-%m-%d %H:%M:%S +0000')
LDFLAGS := -ldflags="-s -w -X 'main.version=$(VERSION)' -X 'main.buildDate=$(BUILD_DATE)'"

# Build the project
build:
	$(GO) build $(LDFLAGS) -o $(BINARY_NAME) $(PKG)

# Build for multiple platforms
build-all:
	$(GO) build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 $(PKG)
	GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe $(PKG)
	GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 $(PKG)
	GOOS=darwin GOARCH=arm64 $(GO) build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 $(PKG)

# Run the project
run:
	$(GO) run $(PKG)

# Clean up build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -rf dist/

# Test the project
test:
	$(GO) test ./pkg/...

# Test with coverage
test-coverage:
	$(GO) test -coverprofile=coverage.out ./pkg/...
	$(GO) tool cover -html=coverage.out -o coverage.html

# Lint the code
lint:
	golangci-lint run

# Format the code
fmt:
	$(GO) fmt ./...

# Vet the code
vet:
	$(GO) vet ./...

# Install dependencies
deps:
	$(GO) mod download
	$(GO) mod tidy

# Generate documentation
docs:
	godoc -http=:6060

# Default target
all: fmt vet test build

# Release target
release: clean fmt vet test build-all

# Development target
dev: fmt vet test run