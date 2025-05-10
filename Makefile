# Tooling
MOCKGEN := go run go.uber.org/mock/mockgen@v0.5.2

# Build configuration
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
TARGET := ./bin

# Test configuration
GO_TEST_FLAGS ?=
PACKAGES ?= ./...

# API specific flags
API_FLAGS := $(if $(ENV),--env=$(ENV)) \
             $(if $(GOPAY_HOST),--gopay-host=$(GOPAY_HOST)) \
             $(if $(GOPAY_PORT),--gopay-port=$(GOPAY_PORT)) \
             $(if $(DB_FILE_PATH),--db-file-path=$(DB_FILE_PATH)) \
             $(if $(DB_OPEN_TIMEOUT),--db-open-timeout=$(DB_OPEN_TIMEOUT)) \
             $(if $(YOOKASSA_CHECKOUT_URL),--yookassa-checkout-url=$(YOOKASSA_CHECKOUT_URL)) \
             $(if $(YOOKASSA_SHOP_ID),--yookassa-shop-id=$(YOOKASSA_SHOP_ID)) \
             $(if $(YOOKASSA_API_TOKEN),--yookassa-api-token=$(YOOKASSA_API_TOKEN))

# Bot specific flags
BOT_FLAGS := $(if $(GOPAY_SERVER_URL),--gopay-server-url=$(GOPAY_SERVER_URL)) \
			 $(if $(TG_BOT_TOKEN),--tg-bot-token=$(TG_BOT_TOKEN)) \
			 $(if $(TG_ADMIN_IDS),--tg-admin-ids=$(TG_ADMIN_IDS))

.PHONY: build-api run-api build-bot run-bot mock test clean help

## build-api: Build the GoPay API
build-api:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(TARGET)/api ./cmd/api/main.go

## run-api: Build and run the GoPay API
run-api: build-api
	$(TARGET)/api $(API_FLAGS)

## build-bot: Build the Telegram bot
build-bot:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(TARGET)/bot ./cmd/bot/main.go

## run-bot: Build and run the Telegram bot
run-bot: build-bot
	$(TARGET)/bot $(BOT_FLAGS)

## mock: Generate mock files
mock:
	$(MOCKGEN) -version || go install go.uber.org/mock/mockgen@v0.5.2
	go generate -run=mockgen $(PACKAGES)

## test: Run unit tests
test: mock
	go test $(GO_TEST_FLAGS) $(PACKAGES)

## clean: Remove build and test artifacts
clean:
	rm -rf $(TARGET) mocks
	go clean -cache -testcache

## help: Display this help message
help:
	@echo "Available targets:"
	@awk '/^## / {sub(/^## /, "", $$0); print}' $(MAKEFILE_LIST) | column -t -s ':'
