APP=dacabot

.PHONY: help
help:  ## This help
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[1m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: clean  ## Build the binary
	@echo "Building binary"
	@go build -ldflags="-s -w" -o dacabot main.go

.PHONY: clean
clean:  ## Remove cached files and dirs from workspace
	@echo "Cleaning workspace"
	@rm -f ${APP}
	@rm -rf tmp

.PHONY: dev
dev:  ## Run the program in dev mode.
	@go run main.go run-server

.PHONY: install
install:  ## Install project dependencies
	@echo "Downloading modules"
	@go mod download

.PHONY: test
test: ## Run tests
	@echo "Running tests"
	@go test ./app
