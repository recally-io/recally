include .env

DATABASE_URL=postgresql://${DATABASE_USER}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=disable

lint: lint-ui lint-go

lint-go:
	@echo "Linting..."
	@go mod tidy
	@golangci-lint run --fix ./...  --enable gofumpt
	@swag fmt

lint-ui:
	@echo "Linting web..."
	@cd web && bun run lint

generate: generate-go generate-sql generate-spec

generate-go:
	@echo "Generating go..."
	@go generate ./...

generate-sql:
	@echo "Generating sql..."
	@go-bindata -prefix "database/migrations/" -pkg migrations -o database/bindata.go database/migrations/
	@sqlc generate

generate-spec:
	@echo "Generate API spec ..."
	@swag init -g internal/port/httpserver/router.go -o docs/swagger
	# @echo "Generating SDK ..."
	# @mkdir -p web/src/sdk && rm -rf web/src/sdk
	# @openapi-generator generate --skip-validate-spec -i docs/swagger/swagger.yaml -g typescript-fetch -o web/src/sdk

build: build-ui build-docs build-go

build-go: generate
	@echo "Building go..."
	@go build -o bin/app main.go

build-docs: 
	@echo "Building docs..."
	@cd docs && bun run docs:build

build-ui:
	@echo "Building UI..."
	@cd web && bun run build

test:
	@echo "Testing..."
	@go test ./...

run: build
	@echo "Running go..."
	@./bin/app

run-ui:
	@echo "Running web..."
	@cd web && bun run dev

run-go: build-go
	@echo "Running dev..."
	@DEBUG_UI=true ./bin/app

ngrok:
	@echo "Running ngrok..."
	@ngrok http 1323

db-up:
	@echo "Starting database..."
	@docker compose up -d postgres

docker-build:
	@echo "Building with docker"
	@docker compose build

docker-run:
	@echo "Running with docker"
	@docker compose up --build

docker-up:
	@echo "Running with docker"
	@docker compose up --build -d

docker-down:
	@echo "Stopping docker"
	@docker compose down

# Create a new migration file
# Usage: make migrate-new name=your_migration_name
migrate-new:
	@echo "Creating migration..."
	@migrate create -ext sql -dir database/migrations -seq "$(name)"
	@echo "New migration created: database/migrations/*_$(name).sql"

# Run all pending migrations or up to a specific version
# Usage: make migrate-up [version=X]
migrate-up:
	@echo "Migrating up..."
	@if [ -z "$(version)" ]; then \
		migrate -path database/migrations -database "$(DATABASE_URL)" up; \
		echo "All pending migrations applied."; \
	else \
		migrate -path database/migrations -database "$(DATABASE_URL)" up $(version); \
		echo "Migrated up to version $(version)."; \
	fi

# Revert all migrations or down to a specific version
# Usage: make migrate-down [version=X]
migrate-down:
	@echo "Migrating down..."
	@if [ -z "$(version)" ]; then \
		migrate -path database/migrations -database "$(DATABASE_URL)" down; \
		echo "All migrations reverted."; \
	else \
		migrate -path database/migrations -database "$(DATABASE_URL)" down $(version); \
		echo "Migrated down to version $(version)."; \
	fi

# Drop all tables in the database
# Usage: make migrate-drop
migrate-drop:
	@echo "Dropping all migrations..."
	@migrate -path database/migrations -database "$(DATABASE_URL)" drop
	@echo "All migrations dropped."

# Force set the database version
# Usage: make migrate-force version=X
migrate-force:
	@echo "Forcing migration version to $(version)..."
	@migrate -path database/migrations -database "$(DATABASE_URL)" force "$(version)"
	@echo "Database version forcibly set to $(version)."

psql:
	@echo "Connecting to database..."
	@docker compose exec -it postgres psql -U ${DATABASE_USER} -d ${DATABASE_NAME}

deploy:
	@echo "Deploying..."
	@dokploy app deploy

help:
	@echo "Available commands:"
	@echo "  lint: Lint the code"
	@echo "  build: Build the code"
	@echo "  buildd: Build the code with docker"
	@echo "  run: Run the code"
	@echo "  rund: Run the code with docker"
	@echo "  migrate-new: Create a new migration, run with 'make migrate-new name=your_migration_name'"
	@echo "  migrate-up: Migrate up"
	@echo "  migrate-down: Migrate down"
	@echo "  help: Show this help message"
