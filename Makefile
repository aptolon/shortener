CMD_DIR=./cmd/shortener

.PHONY: run test test-race

run:
	@go run $(CMD_DIR)

test:
	@go test ./...

test-race:
	@go test ./... -race