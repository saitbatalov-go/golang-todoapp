
include .env
export

export PROJECT_ROOT=$(shell pwd)


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
		docker compose stop todoapp-postgres port-forwarder && \
		sudo rm -rf ${PROJECT_ROOT}/out/pgdata && \
		echo "Файлы окружения очищены"; \
	else \
		echo "Очистка отменена"; \
	fi

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Переменная seq не опреде-runлена, например: make migrate-create seq=init"; \
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
log-cleanup:
	@read -p "Очистить все log файлы? Опасно утеря логов? [y/N] " ans; \
	if [ "$$ans" = "y" ]; then \
		sudo rm -rf ${PROJECT_ROOT}/out/logs && \
		echo "Файлы логов очищены"; \
	else \
		echo "Очистка логов отменена"; \
	fi
todoapp-run:
	@export LOGGER_FOLDER=$(PROJECT_ROOT)/out/logs && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/todoapp/main.go

todoapp-deploy:
	@docker compose up -d --build todoapp

todoapp-undeploy:
	@docker compose stop todoapp

swagger-gen:
	@docker compose run --rm swagger \
		init \
		-g cmd/todoapp/main.go \
		-o docs \
		--parseInternal \
		--parseDependency

ps:
	@docker compose ps
	

	# ========== НОВЫЕ КОМАНДЫ ДЛЯ KAFKA/REDPANDA ==========

# Запуск всех сервисов (PostgreSQL, Redis, Redpanda, приложение, Telegram)
kafka-up:
	@echo "🚀 Запуск всех сервисов..."
	@docker-compose --env-file .env up -d todoapp-postgres todoapp-redis redpanda redpanda-console todoapp telegram
	@sleep 5
	@echo "✅ Все сервисы запущены"
	@echo "📊 Redpanda Console: http://localhost:8080"

# Остановка всех сервисов
kafka-down:
	@docker-compose --env-file .env down

# Перезапуск Telegram consumer
telegram-restart:
	@docker-compose --env-file .env restart telegram

# Логи Telegram consumer
telegram-logs:
	@docker-compose --env-file .env logs -f telegram

# Логи Redpanda
redpanda-logs:
	@docker-compose --env-file .env logs -f redpanda

# Создание топика вручную (если нужно)
redpanda-create-topic:
	@docker exec -it todoapp-redpanda rpk topic create todoapp --partitions 3 --replicas 1 || echo "Топик уже существует"

# Просмотр сообщений в топике
redpanda-consume:
	@docker exec -it todoapp-redpanda rpk topic consume todoapp --offset 0 --partitions 0

# Просмотр всех топиков
redpanda-topics:
	@docker exec -it todoapp-redpanda rpk topic list

# Полный запуск для разработки (с пересборкой)
dev-up:
	@echo "🔨 Пересборка и запуск..."
	@docker-compose --env-file .env down
	@docker-compose --env-file .env build --no-cache todoapp telegram
	@docker-compose --env-file .env up -d
	@echo "✅ Готово!"
	@echo "📊 API: http://localhost:5050"
	@echo "📊 Redpanda Console: http://localhost:8080"

# Проверка статуса всех сервисов
kafka-ps:
	@docker-compose --env-file .env ps