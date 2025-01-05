#!/bin/sh

DEST_DIR="$1"
DEST_FILE="$DEST_DIR/tailwindcss"

# Check if Tailwind CLI is available via command or if the binary exists in the destination directory
if command -v tailwindcss >/dev/null 2>&1 || [ -f "$DEST_FILE" ]; then
  echo "✅ tailwind cli is available."
  exit 0
fi

# Determine OS and architecture for downloading the appropriate binary
OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case $OS in
linux)
  case $ARCH in
  x86_64) FILE="tailwindcss-linux-x64" ;;
  arm64) FILE="tailwindcss-linux-arm64" ;;
  armv7l) FILE="tailwindcss-linux-armv7" ;;
  *)
    echo "❌ Unsupported architecture: $ARCH. Please download the latest release from GitHub."
    exit 1
    ;;
  esac
  ;;
darwin)
  case $ARCH in
  x86_64) FILE="tailwindcss-macos-x64" ;;
  arm64) FILE="tailwindcss-macos-arm64" ;;
  *)
    echo "❌ Unsupported architecture: $ARCH. Please download the latest release from GitHub."
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
    echo "❌ Unsupported architecture: $ARCH. Please download the latest release from GitHub."
    exit 1
    ;;
  esac
  ;;
*)
  echo "❌ Unsupported operating system: $OS. Please download the latest release from GitHub."
  exit 1
  ;;
esac

# Download the Tailwind CLI binary
DOWNLOAD_URL="https://github.com/tailwindlabs/tailwindcss/releases/latest/download/$FILE"
mkdir -p "$DEST_DIR"
echo "🚀 Downloading the latest Tailwind CSS CLI..."
curl -sL "$DOWNLOAD_URL" -o "$DEST_FILE"
chmod u+x "$DEST_FILE"

# Verify the download
if [ -f "$DEST_FILE" ]; then
  echo "✅ tailwind cli has been successfully downloaded to $DEST_FILE."
else
  echo "❌ Failed to download the tailwind cli. Please download it manually from GitHub."
  exit 1
fi
