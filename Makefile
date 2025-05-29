# Tooling
MOCKGEN := go.uber.org/mock/mockgen@v0.5.2
SWAG := github.com/swaggo/swag/cmd/swag@v1.16.4
GOLANGCI_LINT := github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6

# Build configuration
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
TARGET := ./bin

# Test configuration
GO_TEST_FLAGS ?=
PACKAGES ?= ./...

# Lint configuration
LINT_FLAGS ?=

.PHONY: build-api build-bot docs mock test lint clean help

## build-api: Build the GoPay API
build-api: docs
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(TARGET)/api ./cmd/api/main.go

## build-bot: Build the Telegram bot
build-bot:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(TARGET)/bot ./cmd/bot/main.go

## docs: Generate Swagger documentation for the GoPay API
docs:
	go run $(SWAG) -v || go install $(SWAG)
	go run $(SWAG) init -g internal/http/server/server.go -o ./docs

## mock: Generate mock files
mock:
	go run $(MOCKGEN) -version || go install $(MOCKGEN)
	go run $(MOCKGEN) -package=mocks -source=./gopay.go -destination=./mocks/gopay_mocks.go

## test: Run unit tests
test: docs mock
	go test $(GO_TEST_FLAGS) $(PACKAGES)

## lint: Run linters
lint: docs mock
	go run $(GOLANGCI_LINT) version || go install $(GOLANGCI_LINT)
	go run $(GOLANGCI_LINT) run $(LINT_FLAGS)

## clean: Remove build and test artifacts
clean:
	rm -rf $(TARGET) mocks docs
	go clean -cache -testcache

## help: Display this help message
help:
	@echo "Available targets:"
	@awk '/^## / {sub(/^## /, "", $$0); print}' $(MAKEFILE_LIST) | column -t -s ':'
