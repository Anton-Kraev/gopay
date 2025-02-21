.PHONY:
.SILENT:

run:
	go run ./cmd/api/main.go

mock:
	go generate -run=mockgen ./...

test:
	go test .