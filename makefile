.PHONY: up down logs clean migrate kafka-create-topics infra-up infra-down infra-logs

up:
	docker-compose up -d --build

down:
	docker-compose down

logs:
	docker-compose logs -f

logs-%s:
	docker-compose logs -f $*

clean:
	docker-compose down -v

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