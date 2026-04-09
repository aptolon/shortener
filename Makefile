include .env
-include .env.test
export

CMD_DIR=./cmd/shortener

.PHONY: run test test-race migrate-up migrate-down test-db-up test-db-down test-integration docker-build docker-up docker-down

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


test-db-up:
	docker compose -f docker-compose.test.yml up -d --remove-orphans

test-db-down:
	docker compose -f docker-compose.test.yml down -v --remove-orphans

test-integration: test-db-up
	sleep 3
	migrate -path migrations -database "$(TEST_DATABASE_URL)" up
	TEST_DATABASE_URL="$(TEST_DATABASE_URL)" go test ./... -v -race
	$(MAKE) test-db-down

docker-build:
	docker compose build
docker-up:
	docker compose up
docker-down:
	docker compose down

