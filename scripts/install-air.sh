#!/bin/bash

if command -v air &>/dev/null; then
  echo "'air' command is already installed."
  exit 0
fi

if ! command -v go &>/dev/null; then
  echo "Go is not installed. Please install Go first."
  exit 1
fi

echo "Installing 'air' using Go..."
go install github.com/cosmtrek/air@latest

if command -v air &>/dev/null; then
  echo "'air' has been successfully installed."
else
  echo "Failed to install 'air'. Please check for errors."
  exit 1
fi
