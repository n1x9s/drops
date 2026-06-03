APP=secondbrain
GOOSE_DRIVER?=postgres
GOOSE_DBSTRING?=postgres://secondbrain:secondbrain@localhost:5432/secondbrain?sslmode=disable
GOOSE_MIGRATION_DIR?=./migrations

.PHONY: run test tidy migrate-up migrate-down migrate-status docker-build

run:
	go run ./cmd/api

test:
	go test ./...

tidy:
	go mod tidy

migrate-up:
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING="$(GOOSE_DBSTRING)" GOOSE_MIGRATION_DIR=$(GOOSE_MIGRATION_DIR) go run github.com/pressly/goose/v3/cmd/goose up

migrate-down:
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING="$(GOOSE_DBSTRING)" GOOSE_MIGRATION_DIR=$(GOOSE_MIGRATION_DIR) go run github.com/pressly/goose/v3/cmd/goose down

migrate-status:
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING="$(GOOSE_DBSTRING)" GOOSE_MIGRATION_DIR=$(GOOSE_MIGRATION_DIR) go run github.com/pressly/goose/v3/cmd/goose status

docker-build:
	docker build -t $(APP):local .
