.PHONY: run
run: build up

.PHONY: build
build:
	@docker build -t sms-gateway ./api

.PHONY: up
up:
	@echo "docker compose up"
	@docker compose -f docker-compose.yml up -d
	@sleep 30
	@echo "service is ready"

.PHONY: down
down:
	@echo "docker compose down"
	@docker compose -f docker-compose.yml down


