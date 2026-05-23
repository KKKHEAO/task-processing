.PHONY: up down logs clean migrate kafka-create-topics infra-up infra-down infra-logs build build-api build-outboxer build-worker up-services down-services

# === Сборка ===
build-api:
	go build -ldflags="-s -w" -o bin/api ./apps/api/cmd/main.go

build-outboxer:
	go build -ldflags="-s -w" -o bin/outboxer ./apps/outboxer/cmd/main.go

build-worker:
	go build -ldflags="-s -w" -o bin/worker ./apps/worker/cmd/main.go

build: build-api build-outboxer build-worker

up-services:
	docker-compose up -d --build --force-recreate api outboxer worker

down-services:
	docker-compose down api outboxer worker

# === Docker ===
up: infra-up up-services

down: infra-down down-services

logs:
	docker-compose logs -f

logs-%:
	docker-compose logs -f $*

clean:
	docker-compose down -v

# === Инфраструктура ===
migrate:
	docker-compose run --rm migrate

kafka-create-topics:
	bash kafka.sh

infra-up:
	docker-compose up -d postgres migrate kafka kafka-ui
	@sleep 3
	bash kafka.sh

infra-down:
	docker-compose down postgres migrate kafka kafka-ui -v

infra-logs:
	docker-compose logs -f postgres kafka kafka-ui
