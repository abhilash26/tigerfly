#!/bin/sh

PKG="$1"
PKG_URL="$2"

# Check if PKG is already installed
if command -v "${PKG}" >/dev/null 2>&1; then
  echo "✅ ${PKG} is already installed."
  exit 0
fi

# Install PKG using Go
echo "🚀 Installing ${PKG} using go..."
go install "${PKG_URL}"

# Verify installation
if command -v "${PKG}" >/dev/null 2>&1; then
  echo "✅ ${PKG} has been successfully installed."
else
  echo "❌ Failed to install ${PKG}. Please check for errors."
  exit 1
fi
