include .env
export
IMAGE_NAME := Davids_Workout_App
CONTAINER_NAME := Davids_Workout_App_Container

.PHONY: build up down shell clean

build:
	docker build -t $(IMAGE_NAME) .

up:
	docker compose up -d

down:
	docker compose down

shell:
	docker exec -it $(CONTAINER_NAME) /bin/bash

clean:
	docker system prune -f

migrate-up:
	goose -dir db/migrations postgres "$(DB_URL)" up

migrate-down:
	goose -dir db/migrations postgres "$(DB_URL)" down

sqlc:
	sqlc generate
