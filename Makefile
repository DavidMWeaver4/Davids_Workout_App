include .env
export
IMAGE_NAME := Davids_Workout_App
CONTAINER_NAME := Davids_Workout_App_Container

.PHONY: build up down shell clean

build:
	docker build -t $(IMAGE_NAME) .

up:
	docker compose up -d

shell:
	docker exec -it $(CONTAINER_NAME) /bin/bash

down:
	docker compose down --remove-orphans

clean:
	docker compose down --volumes --remove-orphans

migrate-up:
	goose -dir internal/db/migrations postgres "$(DB_URL)" up

migrate-down:
	goose -dir internal/db/migrations postgres "$(DB_URL)" down

sqlc:
	sqlc generate

run:
	go run ./cmd/api/

help:
	@echo "Available commands:"
	@echo "  make build         Build docker image"
	@echo "  make up            Start containers"
	@echo "  make down          Stop containers"
	@echo "  make clean         Stop containers and wipe volumes"
	@echo "  make migrate-up    Run DB migrations"
	@echo "  make run           Run API locally"
	@echo "  make shell         Open shell in container"
	@echo "  make migrate-down  Roll back last migration"
	@echo "  make sqlc          Regenerate sqlc code"
