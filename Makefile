include .env
export
IMAGE_NAME := davids_workout_app
CONTAINER_NAME := Davids_Workout_App_Container
DB_CONTAINER_NAME := davids_workout_db_container

.PHONY: build up up-db down shell clean logs migrate-up migrate-down sqlc run test-db-reset test test-cover test-race help

build:
	docker build -t $(IMAGE_NAME) .

up:
	docker compose up --build -d

up-db:
	docker compose up db -d

shell:
	docker exec -it $(CONTAINER_NAME) /bin/bash

down:
	docker compose down --remove-orphans

clean:
	docker compose down --volumes --remove-orphans

logs:
	docker compose logs -f

migrate-up:
	goose -dir internal/db/migrations postgres "$(DB_URL)" up

migrate-down:
	goose -dir internal/db/migrations postgres "$(DB_URL)" down

sqlc:
	sqlc generate

run:
	go run ./cmd/api/

test-db-reset:
	docker exec -i $(DB_CONTAINER_NAME) psql -U postgres -d postgres -c "DROP DATABASE IF EXISTS workout_tracker_test;"
	docker exec -i $(DB_CONTAINER_NAME) psql -U postgres -d postgres -c "CREATE DATABASE workout_tracker_test;"
	goose -dir internal/db/migrations postgres "$(TEST_DB_URL)" up

test: test-db-reset
	go test ./...

test-cover: test-db-reset
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

test-race: test-db-reset
	go test -race ./...

help:
	@echo "Available commands:"
	@echo "  make build         Build docker image"
	@echo "  make up            Start containers"
	@echo "  make up-db         Start only the db container"
	@echo "  make down          Stop containers"
	@echo "  make clean         Stop containers and wipe volumes"
	@echo "  make logs          Get docker logs"
	@echo "  make migrate-up    Run DB migrations"
	@echo "  make run           Run API locally"
	@echo "  make shell         Open shell in container"
	@echo "  make migrate-down  Roll back last migration"
	@echo "  make sqlc          Regenerate sqlc code"
	@echo "  make test          Run Go unit tests"
	@echo "  make test-cover    Run Go tests with coverage"
	@echo "  make test-race     Run Go tests with race"
