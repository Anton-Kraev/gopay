.PHONY:
.SILENT:

run:
	go run ./cmd/app/main.go

mock:
	go generate -run=mockgen ./...

test:
	go test .