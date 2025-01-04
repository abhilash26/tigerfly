# Define variables

ROOT_DIR="$(shell pwd)"
TEMP_DIR="${ROOT_DIR}/tmp"
BUILD_DIR="${ROOT_DIR}/cmd"
SCRIPTS_DIR="${ROOT_DIR}/dev/scripts"
TOOLS_DIR="${ROOT_DIR}/dev/tools"

INPUT_CSS="${ROOT_DIR}/assets/css/main.css"
OUTPUT_CSS="${ROOT_DIR}/static/css/main.css"
INPUT_JS="${ROOT_DIR}/assets/js/main.js"
OUTPUT_JS="${ROOT_DIR}/static/js/main.js"

.PHONY: init watch build clean check-env install-air install-tailwind install-esbuild refresh run-css run-js watch-css watch-js watch-go build-css build-js build-go

# Check if .env file exists and create from env.example if not
check-env:
	@sh "${SCRIPTS_DIR}/check_env.sh" "${ROOT_DIR}"

# Install 'air' if not already installed
install-air:
	@sh "${SCRIPTS_DIR}/install_air.sh"

# Install 'tailwindcss' if not already installed
install-tailwind:
	@sh "${SCRIPTS_DIR}/install_tailwind_cli.sh" "${TOOLS_DIR}"

# Install 'esbuild' if not already installed
install-esbuild:
	@sh "${SCRIPTS_DIR}/install_esbuild.sh" "${TOOLS_DIR}"

# Refresh Go modules
refresh:
	@echo "🔄 Refreshing Go modules..."
	@go mod tidy

# Initialize environment, install necessary tools, and set up project
init: check-env install-air install-tailwind install-esbuild refresh
	@echo "-----------------------"
	@echo "🎉 Welcome to Tigerfly!"
	@echo "-----------------------"

# Load environment variables from .env file (if exists)
-include .env

# Default target
default: init

# Run for changes with Tailwind and Esbuild
run-css:
	@echo "🚀 Running CSS with Tailwind..."
	@${TOOLS_DIR}/tailwindcss -i "${INPUT_CSS}" -o "${OUTPUT_CSS}"

run-js:
	@echo "🚀 Running JS with Esbuild..."
	@${TOOLS_DIR}/esbuild "${INPUT_JS}" --outfile="${OUTPUT_JS}"

run: run-css run-js watch-go

# Watch for changes with Tailwind, Esbuild, and Go app with air
watch-css:
	@echo "👀 Watching CSS with Tailwind..."
	@${TOOLS_DIR}/tailwindcss -i "${INPUT_CSS}" -o "${OUTPUT_CSS}" --watch

watch-js:
	@echo "👀 Watching JS with Esbuild..."
	@${TOOLS_DIR}/esbuild "${INPUT_JS}" --outfile="${OUTPUT_JS}" --watch=forever

watch-go:
	@echo "🚀 Running Air..."
	@air

# Build the Tailwind CSS, Esbuild, and Go app, and prepare the build directory
build-css:
	@echo "🔨 Building CSS with Tailwind..."
	@${TOOLS_DIR}/tailwindcss -i "${INPUT_CSS}" -o "${OUTPUT_CSS}" --minify

build-js:
	@echo "🔨 Building JS with Esbuild..."
	@${TOOLS_DIR}/esbuild "${INPUT_JS}" --minify --bundle --outfile="${OUTPUT_JS}"

build-go:
	@echo "📂 Copying directories to build..."
	@mkdir -p "${BUILD_DIR}"
	@cp -r "${VIEW_PATH}" "${BUILD_DIR}/"
	@cp .env "${BUILD_DIR}/.env"
	@echo "📦 Copying database..."
	@cp "${DATABASE_FILE}" "${BUILD_DIR}/${DATABASE_FILE}"
	@echo "⚙️  Building Go app..."
	@go build -o "${BUILD_DIR}/app" main.go

build: build-css build-js build-go

# Clean up build and temp directories
clean:
	@echo "🧹 Cleaning ${TEMP_DIR} directory..."
	@rm -rf "${TEMP_DIR}/"
	@echo "🧹 Cleaning build directory..."
	@rm -rf "${BUILD_DIR}/"
