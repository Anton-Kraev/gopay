.PHONY:
.SILENT:

api:
	go run ./cmd/api/main.go

bot:
	go run ./cmd/bot/main.go

mock:
	go generate -run=mockgen ./...

test:
	go test .
