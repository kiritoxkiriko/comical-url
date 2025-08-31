#!/bin/bash

# Build script for Short URL Service

set -e

echo "Building Short URL Service..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Go is not installed or not in PATH"
    exit 1
fi

# Get the project root directory
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_ROOT"

# Clean previous builds
echo "Cleaning previous builds..."
rm -f shorturl shorturl_*

# Download dependencies
echo "Downloading dependencies..."
go mod download
go mod tidy

# Run tests
echo "Running tests..."
go test -v ./...

# Build for current platform
echo "Building for current platform..."
go build -o shorturl -v ./...

# Build for Linux (useful for Docker)
echo "Building for Linux..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o shorturl_linux -v

echo "Build completed successfully!"
echo "Binaries created:"
echo "  - shorturl (current platform)"
echo "  - shorturl_linux (Linux amd64)"