# Путь до docker compose файлов
COMPOSE_DIR=./deployments

# env файлы
ENV_FILE=.env
ENV_TEST_FILE=.env-test
EXAMPLE_ENV_FILE=exampleenv
EXAMPLE_ENV_TEST_FILE=exampleenvtest

# Docker Compose файлы
COMPOSE_FILE=$(COMPOSE_DIR)/docker-compose.yaml
COMPOSE_TEST_FILE=$(COMPOSE_DIR)/docker-compose-test.yaml

# Проверка и копирование .env, если не существует
init-env:
	@if [ ! -f $(ENV_FILE) ]; then \
		echo "Creating $(ENV_FILE) from $(EXAMPLE_ENV_FILE)..."; \
		cp $(EXAMPLE_ENV_FILE) $(ENV_FILE); \
	fi
	@if [ ! -f $(ENV_TEST_FILE) ]; then \
		echo "Creating $(ENV_TEST_FILE) from $(EXAMPLE_ENV_TEST_FILE)..."; \
		cp $(EXAMPLE_ENV_TEST_FILE) $(ENV_TEST_FILE); \
	fi

up: init-env
	@echo "Starting main project with $(ENV_FILE)..."
	@docker compose -f $(COMPOSE_FILE) up -d --build

down:
	@echo "Stopping main project..."
	@docker compose -f $(COMPOSE_FILE) down

up-test: init-env up-test-start check-tests-logs

up-test-start: init-env
	@echo "Starting test environment with $(ENV_TEST_FILE)..."
	@docker compose -f $(COMPOSE_TEST_FILE) up -d --build
	
down-test:
	@echo "Stopping test environment..."
	@docker compose -f $(COMPOSE_TEST_FILE) down -v

rebuild: init-env
	@echo "Rebuilding all services..."
	@docker compose -f $(COMPOSE_FILE) up -d --build

rebuild-test: init-env
	@echo "Rebuilding test services..."
	@docker compose -f $(COMPOSE_TEST_FILE) up -d --build

check-tests-logs:
	@echo "Checking logs for test container..."
	@if docker ps -a --format '{{.Names}}' | grep -q "^deployments-test-1$$"; then \
		if [ "`docker inspect -f '{{.State.Running}}' deployments-test-1`" = "true" ]; then \
			echo "Container is running. Waiting and showing logs..."; \
			sleep 50; \
			docker logs --tail=50 deployments-test-1; \
		else \
			echo "Container exists but is not running."; \
			docker logs --tail=50 deployments-test-1; \
		fi \
	else \
		echo "Container 'deployments-test-1' does not exist."; \
	fi

.PHONY: up down up-test up-test-start down-test rebuild rebuild-test check-tests-logs init-env