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
	@go build -ldflags="-s -w"

.PHONY: clean
clean:  ## Clean workspace
	@rm -f ${APP}
	@rm -rf tmp

.PHONY: install
install:  ## Install project dependencies
	@go mod download

.PHONY: test
test: ## Run tests
	@go test -v -cover ./app/...

.PHONY: dev
dev:  ## Run the program in dev mode.
	@go run main.go run-server
