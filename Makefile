.PHONY: init watch build clean check-env install-air install-pnpm

# Check if .env file exists and create from env.example if not
check-env:
	@if [ -f ".env" ]; then \
		echo ".env file exists, proceeding with the script..."; \
	else \
		if [ -f "env.example" ]; then \
			echo ".env file is missing, copying env.example to .env..."; \
			cp "env.example" ".env"; \
			echo ".env has been created from env.example"; \
		else \
			echo "Error: Neither .env nor env.example file exists. Please create an env.example file."; \
			exit 1; \
		fi \
	fi

# Install 'air' if not already installed
install-air:
	@if command -v air &>/dev/null; then \
		echo "'air' command is already installed."; \
		exit 0; \
	fi

	@if ! command -v go &>/dev/null; then \
		echo "Go is not installed. Please install Go first."; \
		exit 1; \
	fi

	echo "Installing 'air' using Go..."; \
	go install "github.com/air-verse/air@latest"; \

	@if command -v air &>/dev/null; then \
		echo "'air' has been successfully installed."; \
	else \
		echo "Failed to install 'air'. Please check for errors."; \
		exit 1; \
	fi

# Install 'pnpm' if not already installed
install-pnpm:
	@if ! command -v pnpm &>/dev/null; then \
		echo "pnpm is not installed. Please install curl first."; \
		exit 0; \
	fi

# Initialize environment, install air and pnpm, and tidy Go modules
init: check-env install-air install-pnpm
	@echo "Welcome to Tigerfly!"
	@echo "--------------------"
	@echo "Installing JS packages with pnpm..."
	@pnpm i
	@echo "Refreshing Go packages..."
	@go mod tidy

# Load environment variables from .env file (if exists)
-include .env

default: init

# Watch for changes with Tailwind and air
watch:
	@echo "Watching Tailwind..."
	@pnpm run "dev:tailwind" &
	@echo "Running Air..."
	@air

# Build the Tailwind CSS, Go app, and prepare the build directory
build:
	@echo "Building Tailwind CSS..."
	@pnpm run "build:tailwind"
	@echo "Copying directories to build..."
	@mkdir -p "${BUILD_DIR}"
	@cp -r "${VIEW_PATH}" "${BUILD_DIR}/"
	@cp .env "${BUILD_DIR}/.env"
	@echo "Copying database..."
	@cp "${DATABASE_FILE}" "${BUILD_DIR}/${DATABASE_FILE}"
	@echo "Building Go app..."
	@go build -o "${BUILD_DIR}/app" main.go

# Clean up build and temp directories
clean:
	@echo "Cleaning ${TEMP_DIR} directory..."
	@rm -rf "${TEMP_DIR}/"
	@echo "Cleaning build directory..."
	@rm -rf "${BUILD_DIR}/"
