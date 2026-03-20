# Загружаем .env файл
include .env
export

# Динамический PROJECT_ROOT
PROJECT_ROOT := $(shell pwd)
export PROJECT_ROOT

# Общая переменная для подключения к БД
DB_URL := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@todoapp-postgres:5432/$(POSTGRES_DB)?sslmode=disable

env-up:
	@docker-compose --env-file .env up -d todoapp-postgres
	@sleep 3
	@make migrate-up

env-down:
	@docker-compose --env-file .env down todoapp-postgres

env-cleanup:
	@read -p "Очистить все volume окружения? Опасно потеря данных? [y/N] " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose stop todoapp-postgres && \
		docker compose rm -f todoapp-postgres && \
		sudo rm -rf out/pgdata && \
		echo "Файлы окружения очищены"; \
	else \
		echo "Очистка отменена"; \
	fi

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Переменная seq не определена, например: make migrate-create seq=init"; \
		exit 1; \
	fi; \
	docker-compose --env-file .env run --rm todoapp-postgres-migrate \
	create \
	-ext sql \
	-dir /migrations \
	-seq "$(seq)"

migrate-status:
	@docker-compose --env-file .env run --rm todoapp-postgres-migrate \
	-path /migrations \
	-database "$(DB_URL)" \
	version

migrate-up:
	@docker-compose --env-file .env run --rm todoapp-postgres-migrate \
	-path /migrations \
	-database "$(DB_URL)" \
	up

migrate-down:
	@docker-compose --env-file .env run --rm todoapp-postgres-migrate \
	-path /migrations \
	-database "$(DB_URL)" \
	down

migrate-reset:
	@docker-compose --env-file .env run --rm todoapp-postgres-migrate \
	-path /migrations \
	-database "$(DB_URL)" \
	drop -f || true
	@docker-compose --env-file .env run --rm todoapp-postgres-migrate \
	-path /migrations \
	-database "$(DB_URL)" \
	up

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder 