#!/bin/bash

DEST_DIR="./"
DEST_FILE="$DEST_DIR/tailwindcss"

if command -v tailwindcss &>/dev/null; then
  echo "Tailwind CLI is available (command found)."
  exit 0
fi

if [ -f "$DEST_FILE" ]; then
  echo "Tailwind CLI is available (file found in $DEST_FILE)."
  exit 0
fi

# Tailwind not available. Download

# Check if the tailwindcss file exists in the destination directory
OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $OS in
linux)
  case $ARCH in
  x86_64) FILE="tailwindcss-linux-x64" ;;
  arm64) FILE="tailwindcss-linux-arm64" ;;
  armv7l) FILE="tailwindcss-linux-armv7" ;;
  *)
    echo "Unsupported architecture: $ARCH. Please download the latest release from GitHub."
    exit 1
    ;;
  esac
  ;;
darwin)
  case $ARCH in
  x86_64) FILE="tailwindcss-macos-x64" ;;
  arm64) FILE="tailwindcss-macos-arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH. Please download the latest release from GitHub."
    exit 1
    ;;
  esac
  ;;
mingw* | cygwin* | msys*)
  OS="windows"
  case $ARCH in
  x86_64) FILE="tailwindcss-windows-x64.exe" ;;
  arm64) FILE="tailwindcss-windows-arm64.exe" ;;
  *)
    echo "Unsupported architecture: $ARCH. Please download the latest release from GitHub."
    exit 1
    ;;
  esac
  ;;
*)
  echo "Unsupported operating system: $OS. Please download the latest release from GitHub."
  exit 1
  ;;
esac

# If the file doesn't exist, proceed with downloading from the static URL
DOWNLOAD_URL="https://github.com/tailwindlabs/tailwindcss/releases/latest/download/$FILE"
mkdir -p "$DEST_DIR"

echo "Downloading the latest Tailwind CSS CLI..."
curl -sL "$DOWNLOAD_URL" -o "$DEST_FILE"

# Check if the download was successful
if [ -f "$DEST_FILE" ]; then
  echo "Tailwind CLI has been successfully downloaded to $DEST_FILE."
else
  echo "Failed to download the Tailwind CLI. Please download it manually from GitHub."
  exit 1
fi
