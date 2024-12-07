FIX_MIGRATION_QUERY="UPDATE schema_migrations SET dirty = false WHERE version = (SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1);"

POSTGRES_USER=song_manager
POSTGRES_PASSWORD=secretsong
POSTGRES_DATABASE=song_library
POSTGRES_PORT=5432
POSTGRES_HOST=localhost

DB_URL=postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=disable

create_db:
	@echo "Creating PostgreSQL database: $(POSTGRES_DATABASE)"
	createdb -h $(POSTGRES_HOST) -p $(POSTGRES_PORT) -U $(POSTGRES_USER) $(POSTGRES_DATABASE)

drop_db:
	@echo "Dropping PostgreSQL database: $(POSTGRES_DATABASE)"
	dropdb -h $(POSTGRES_HOST) -p $(POSTGRES_PORT) -U $(POSTGRES_USER) $(POSTGRES_DATABASE)

up_migrate:
	@echo "Running migrations UP for database: $(POSTGRES_DATABASE)"
	migrate -path db/migrations -database "$(DB_URL)" -verbose up

down_migrate:
	@echo "Rolling back migrations DOWN for database: $(POSTGRES_DATABASE)"
	migrate -path db/migrations -database "$(DB_URL)" -verbose down

up_migrate_last:
	@echo "Running the last migration UP for database: $(POSTGRES_DATABASE)"
	migrate -path db/migrations -database "$(DB_URL)" -verbose up 1

down_migrate_last:
	@echo "Rolling back last migration and resetting dirty state...\n"
	migrate -path db/migrations -database "$(DB_URL)" -verbose down 1 || true
	@echo "Executing custom query to reset migration dirty flag..."
	PGPASSWORD=$(POSTGRES_PASSWORD) psql -U $(POSTGRES_USER) -d $(POSTGRES_DATABASE) -c $(FIX_MIGRATION_QUERY)
	@echo "\nMigration reverted and dirty flag reset."

new_migration:
	@echo "Creating new migration with name: $(name)"
	migrate create -ext sql -dir db/migrations -seq $(name)

sqlc:
	@echo "Generating SQLC code"
	sqlc generate

mock:
	@echo "Generating mock store for the Store interface"
	mockgen -package mockdb -destination db/mock/store.go github.com/komron-dev/musicLibrary/db/sqlc Store

swag-gen:
	@echo "Generating Swagger documentation"
	swag init -g ./api/server.go -o ./api/docs

run:
	@echo "Running the Go application"
	go run main.go

import:
	@echo "Running 'go mod tidy' to clean up and synchronize dependencies"
	go mod tidy

.PHONY: create_db drop_db up_migrate up_migrate_last down_migrate down_migrate_last run import
