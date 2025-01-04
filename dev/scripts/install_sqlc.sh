#!/bin/sh

# Check if 'sqlc' is already installed
if command -v sqlc >/dev/null 2>&1; then
  echo "✅ sqlc is already installed."
  exit 0
fi

# Check if 'go' is installed
if ! command -v go >/dev/null 2>&1; then
  echo "❌ go is not installed. Please install go first."
  exit 1
fi

# Install 'sqlc' using Go
echo "🚀 Installing sqlc using go..."
go install "github.com/sqlc-dev/sqlc/cmd/sqlc@latest"

# Verify installation
if command -v sqlc >/dev/null 2>&1; then
  echo "✅ sqlc has been successfully installed."
else
  echo "❌ Failed to install sqlc. Please check for errors."
  exit 1
fi
