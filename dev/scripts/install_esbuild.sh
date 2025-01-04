#!/bin/sh

DEST_DIR="$1"
TEMP_DIR="${DEST_DIR}/tmp_esbuild"
ESBUILD_BIN="${DEST_DIR}/esbuild"

# Check if esbuild is already installed
if command -v esbuild >/dev/null 2>&1 || [ -f "$ESBUILD_BIN" ]; then
  echo "✅ esbuild is already installed."
  exit 0
fi

mkdir -p "${TOOLS_DIR}"
echo "🚀 Cloning esbuild repository..."
git clone --depth 1 "https://github.com/evanw/esbuild.git" "${TEMP_DIR}"
cd "${TEMP_DIR}" || exit 1
echo "🔨 Building esbuild from source..."
go build ./cmd/esbuild
mv esbuild "${ESBUILD_BIN}"
chmod u+x "${ESBUILD_BIN}"

echo "🧹 Cleaning up temporary files..."
rm -rf "${TEMP_DIR}"

echo "✅ esbuild has been successfully built and installed at ${ESBUILD_BIN}."
