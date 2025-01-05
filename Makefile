# Define variables
ROOT_DIR := $(shell pwd)
TEMP_DIR := ${ROOT_DIR}/tmp
BUILD_DIR := ${ROOT_DIR}/cmd
SCRIPTS_DIR := ${ROOT_DIR}/dev/scripts
TOOLS_DIR := ${ROOT_DIR}/dev/tools

INPUT_CSS := ${ROOT_DIR}/assets/css/main.css
OUTPUT_CSS := ${ROOT_DIR}/static/css/main.css
INPUT_JS := ${ROOT_DIR}/assets/js/main.js
OUTPUT_JS := ${ROOT_DIR}/static/js/main.js

.PHONY: init watch build clean check-env install-requirements refresh run watch-css watch-js watch-go build

# Check if .env file exists and create from env.example if not
check-env:
	@sh "${SCRIPTS_DIR}/check_env.sh" "${ROOT_DIR}"

# Install 'requirements' if not already installed
install-requirements:
	@sh "${SCRIPTS_DIR}/go_install.sh" "air" "github.com/air-verse/air@latest"
	@sh "${SCRIPTS_DIR}/go_install.sh" "templ" "github.com/a-h/templ/cmd/templ@latest"
	@sh "${SCRIPTS_DIR}/install_tailwind_cli.sh" "${TOOLS_DIR}"
	@sh "${SCRIPTS_DIR}/go_install.sh" "sqlc" "github.com/sqlc-dev/sqlc/cmd/sqlc@latest"
	@sh "${SCRIPTS_DIR}/install_esbuild.sh" "${TOOLS_DIR}"

# Refresh Go modules
refresh:
	@echo "🔄 Refreshing go modules..."
	@go mod tidy

# Initialize environment, install necessary tools, and set up project
init: check-env install-requirements refresh
	@echo "-----------------------"
	@echo "🎉 Welcome to Tigerfly!"
	@echo "-----------------------"

# Load environment variables from .env file (if exists)
-include .env

# Default target
default: init

# Run for changes with Tailwind and Esbuild
run:
	@echo "🚀 Running css with tailwind..."
	@${TOOLS_DIR}/tailwindcss -i "${INPUT_CSS}" -o "${OUTPUT_CSS}"
	@echo "🚀 Running js with esbuild..."
	@${TOOLS_DIR}/esbuild "${INPUT_JS}" --outfile="${OUTPUT_JS}"
	@echo "🚀 Running go with air..."
	@air

# Watch for changes with Tailwind, Esbuild, and Go app with air
watch-css:
	@echo "👀 Watching css with tailwind..."
	@${TOOLS_DIR}/tailwindcss -i "${INPUT_CSS}" -o "${OUTPUT_CSS}" --watch

watch-js:
	@echo "👀 Watching js with esbuild..."
	@${TOOLS_DIR}/esbuild "${INPUT_JS}" --outfile="${OUTPUT_JS}" --watch

watch-go:
	@echo "🚀 Running air..."
	@air

# Build the Tailwind CSS, Esbuild, and Go app, and prepare the build directory
build:
	@echo "🔨 Building css with tailwind..."
	@${TOOLS_DIR}/tailwindcss -i "${INPUT_CSS}" -o "${OUTPUT_CSS}" --minify
	@echo "🔨 Building js with esbuild..."
	@${TOOLS_DIR}/esbuild "${INPUT_JS}" --minify --bundle --outfile="${OUTPUT_JS}"
	@echo "📂 Preparing build directory..."
	@mkdir -p "${BUILD_DIR}"
	@cp -r "${VIEW_PATH}" "${BUILD_DIR}/"
	@cp .env "${BUILD_DIR}/.env"
	@echo "📦 Copying database..."
	@cp "${DATABASE_FILE}" "${BUILD_DIR}/${DATABASE_FILE}"
	@echo "⚙️  Building go app..."
	@go build -o "${BUILD_DIR}/app" main.go

# Clean up build and temp directories
clean:
	@echo "🧹 Cleaning ${TEMP_DIR} directory..."
	@rm -rf "${TEMP_DIR}/"
	@echo "🧹 Cleaning build directory..."
	@rm -rf "${BUILD_DIR}/"
