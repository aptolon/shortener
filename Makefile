include .env
export

CMD_DIR=./cmd/shortener

.PHONY: run test test-race migrate-up migrate-down docker-build docker-up docker-down clean

build:
	@go build -o $(BIN_PATH) $(CMD_PATH)
run:
	@go run $(CMD_DIR)
clean:
	rm -rf ./bin

test:
	@go test ./...
	
test-race:
	@go test ./... -race

migrate-up:
	migrate -path migrations -database ${DATABASE_URL} up
migrate-down:
	migrate -path migrations -database ${DATABASE_URL} down

docker-build:
	docker compose build
docker-up:
	docker compose up
docker-down:
	docker compose down

