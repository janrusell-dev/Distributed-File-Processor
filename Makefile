ifneq (,$(wildcard ./.env))
	include	.env
	export
endif

DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@localhost:${DB_PORT}/${DB_NAME}?sslmode=disable"

.PHONY: migrate-up migrate-status migrate-down up down restart build logs ui test


migrate-up:
	@goose -dir internal/db/migrations postgres $(DB_URL) up

# Check migration status
migrate-status:
	@goose -dir internal/db/migrations postgres $(DB_URL) status

# Rollback one migration
migrate-down:
	@goose -dir internal/db/migrations postgres $(DB_URL) down

up:
	docker compose up -d

down:
ifeq ($(SERVICE),)
	docker compose down
else
	docker compose stop $(SERVICE)
endif

rebuild:
ifeq ($(SERVICE),)
	docker compose up -d --build
else
	docker compose up -d --build $(SERVICE)
endif

logs:
ifeq ($(SERVICE),)
	docker compose logs -f
else
	docker compose logs -f $(SERVICE)
endif

ui:
	xdg-open http://localhost:8080

test-file:
	go run cmd/client/main.go testdata/sample.txt



