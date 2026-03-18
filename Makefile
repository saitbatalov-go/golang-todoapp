include .env
export

export PROJECT_ROOT= $(shell pwd)

env-up:
	@docker-compose up -d todoapp-postgres

env-down:
	@docker-compose down todoapp-postgres

env-cleanup:
	@read -p "Очистить все volume окружения? Опасно потеря данных? [y/N] " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose stop todoapp-postgres && \
		docker compose rm -f todoapp-postgres && \
		rm -rf out/pgdata && \
		echo "Файлы окружения очищены"; \
	else \
		echo "Очистка отменена"; \
	fi

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Переменная seq не определена, например: make migrate-create seq=init"; \
		exit 1; \
	fi; \
	docker compose run --rm todoapp-postgres-migrate \
	create \
	-ext sql \
	-dir /migrations \
	-seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Переменная action не определена, например: make migrate-action action=up"; \
		exit 1; \
	fi; \
	docker compose run --rm todoapp-postgres-migrate \
	-path /migrations \
	-database postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@todoapp-postgres:5432/$(POSTGRES_DB)?sslmode=disable \
	"$(action)"