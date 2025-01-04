#!/bin/sh

# Check if 'air' is already installed
if command -v air >/dev/null 2>&1; then
  echo "✅ 'air' command is already installed."
  exit 0
fi

# Check if 'go' is installed
if ! command -v go >/dev/null 2>&1; then
  echo "❌ Go is not installed. Please install Go first."
  exit 1
fi

# Install 'air' using Go
echo "🚀 Installing 'air' using Go..."
go install "github.com/cosmtrek/air@latest"

# Verify installation
if command -v air >/dev/null 2>&1; then
  echo "✅ 'air' has been successfully installed."
else
  echo "❌ Failed to install 'air'. Please check for errors."
  exit 1
fi
