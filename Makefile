# Makefile for macconv project

# Variables
GO := go
BINARY_NAME := macconv
PKG := .

# Build the project
build:
	$(GO) build -ldflags="-s -w" -o $(BINARY_NAME) $(PKG)

# Run the project
run:
	$(GO) run $(PKG)

# Clean up build artifacts
clean:
	rm -f $(BINARY_NAME)

# Test the project
test:
	$(GO) test ./...

# Default target
all: build run