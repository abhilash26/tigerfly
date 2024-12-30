.PHONY: init watch build clean

init: check-env install-air install-tailwind
	@echo "Welcome to Tigerfly!"
	@echo "--------------------"
	@echo "Refreshing Go packages..."
	@go mod tidy

check-env:
	@./scripts/copy-env.sh

install-air:
	@./scripts/install-air.sh

install-tailwind:
	@./scripts/install-tailwind.sh

# Load environment variables from .env file (if exists)
-include .env

default: init

watch:
	@echo "Starting Tailwind watch..."
	@./tailwindcss -i "./assets/css/main.css" -o "./static/css/main.css" --watch &
	@echo "Running Air..."
	@air

build:
	@echo "Building Tailwind CSS..."
	@./tailwindcss -i "./assets/css/main.css" -o "./static/css/main.css"
	@echo "Copying directories to build..."
	@mkdir -p "${BUILD_DIR}"
	@cp -r "${VIEW_PATH}" "${BUILD_DIR}/"
	@cp .env "${BUILD_DIR}/.env"
	@echo "Copying database..."
	@cp "${DATABASE_FILE}" "${BUILD_DIR}/${DATABASE_FILE}"
	@echo "Building Go app..."
	@go build -o "${BUILD_DIR}/app" main.go

clean:
	@echo "Cleaning ${TEMP_DIR} directory..."
	@rm -rf "${TEMP_DIR}/"
	@echo "Cleaning build directory..."
	@rm -rf "${BUILD_DIR}/"
