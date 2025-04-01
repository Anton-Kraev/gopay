.PHONY:
.SILENT:

api:
	go run ./cmd/api/main.go

mock:
	go generate -run=mockgen ./...

test:
	go test .
