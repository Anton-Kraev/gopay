TARGET := ./bin

build-api:
	go build -o $(TARGET)/api ./cmd/api/main.go

run-api: build-api
	$(TARGET)/api \
		$(if $(ENV),--env=$(ENV)) \
		$(if $(GOPAY_HOST),--gopay-host=$(GOPAY_HOST)) \
		$(if $(GOPAY_PORT),--gopay-port=$(GOPAY_PORT)) \
		$(if $(DB_FILE_PATH),--db-file-path=$(DB_FILE_PATH)) \
		$(if $(DB_OPEN_TIMEOUT),--db-open-timeout=$(DB_OPEN_TIMEOUT)) \
		$(if $(YOOKASSA_CHECKOUT_URL),--yookassa-checkout-url=$(YOOKASSA_CHECKOUT_URL)) \
		$(if $(YOOKASSA_SHOP_ID),--yookassa-shop-id=$(YOOKASSA_SHOP_ID)) \
		$(if $(YOOKASSA_API_TOKEN),--yookassa-api-token=$(YOOKASSA_API_TOKEN))

build-bot:
	go build -o $(TARGET)/bot ./cmd/bot/main.go

run-bot: build-bot
	$(TARGET)/bot \
		$(if $(GOPAY_SERVER_URL),--gopay-server-url=$(GOPAY_SERVER_URL)) \
		$(if $(TG_BOT_TOKEN),--tg-bot-token=$(TG_BOT_TOKEN)) \
		$(if $(TG_ADMIN_IDS),--tg-admin-ids=$(TG_ADMIN_IDS))

mock:
	go generate -run=mockgen ./...

test:
	go test .
