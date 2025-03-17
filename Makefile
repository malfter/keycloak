.DEFAULT_GOAL := help

.PHONY: help
help:  ## 💬 This help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: lint
lint:  ## 🔍 Run golangci-lint
	cd demo-apps/demo-backend && golangci-lint run ./...
	cd demo-apps/demo-webapp && golangci-lint run ./...

.PHONY: format
format:  ## 🪄  Run go fmt
	cd demo-apps/demo-backend && go fmt ./...
	cd demo-apps/demo-webapp && go fmt ./...

.PHONY: test
test:  ## 🧪 Run go test
	cd demo-apps/demo-backend && go test .
	cd demo-apps/demo-webapp && go test .
