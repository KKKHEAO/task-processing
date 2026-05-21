.PHONY: up down logs clean migrate kafka-create-topics infra-up infra-down infra-logs build build-api build-outboxer build-worker

# === Сборка ===
build-api:
	go build -ldflags="-s -w" -o bin/api ./apps/api/cmd/main.go

build-outboxer:
	go build -ldflags="-s -w" -o bin/outboxer ./apps/outboxer/cmd/main.go

build-worker:
	go build -ldflags="-s -w" -o bin/worker ./apps/worker/cmd/main.go

build: build-api build-outboxer build-worker

# === Docker ===
up:
	docker-compose up -d --build

down:
	docker-compose down

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

infra-down:
	docker-compose down postgres migrate kafka kafka-ui -v

infra-logs:
	docker-compose logs -f postgres kafka kafka-ui
