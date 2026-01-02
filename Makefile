migrate-up:
	goose -dir internal/db/migrations postgres $(DB_URL) up

# Check migration status
migrate-status:
	goose -dir internal/db/migrations postgres $(DB_URL) status

# Rollback one migration
migrate-down:
	goose -dir internal/db/migrations postgres $(DB_URL) down