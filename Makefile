include .env
export

export PROJECT_ROOT= $(shell pwd)

env-up:
	docker-compose up -d todoapp-postgres

env-down:
	docker-compose down todoapp-postgres

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