include .env
export

CMD_DIR=./cmd/shortener

.PHONY: run test test-race migrate-up migrate-down

run:
	@go run $(CMD_DIR)

test:
	@go test ./...
	
test-race:
	@go test ./... -race


migrate-up:
	migrate -path migrations -database ${DATABASE_URL} up
migrate-down:
	migrate -path migrations -database ${DATABASE_URL} down