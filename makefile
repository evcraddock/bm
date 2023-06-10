.PHONY: help
help: ## Show this help
	@egrep -h '\s##\s' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: build-linux
build-linux: ## Build the application
	GOOS=linux GOARCH=amd64 go build -o .build/linux/bm cmd/bm/bm.go

.PHONY: build-mac
build-mac: ## Build the app for macos
	GOOS=darwin GOARCH=amd64 go build -o .build/mac/bm cmd/bm/bm.go

# .PHONY: install
# install: build ## Install the application
# 	cp .build/bm ${GOPATH}/bin/bm

.PHONY: fmt
fmt: ## Run all formatings
	go mod vendor
	go mod tidy
	go fmt ./...

.PHONY: test-all
test-all: ## Run all tests
	go test ./...

.PHONY: coverage
coverage: ## Show test coverage
	@go test -coverprofile=coverage.out ./... > /dev/null
	go tool cover -func=coverage.out
	@rm coverage.out

.PHONY: test
test: test-all coverage ## Show Test Results
