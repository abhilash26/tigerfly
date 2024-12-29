.PHONY: init run build clean

OS="$(uname)"

default: init

init:
	@echo "Installing go air tool"
	@go install "github.com/air-verse/air@latest"
	@echo "Installing npm packages using pnpm"
	@pnpm i