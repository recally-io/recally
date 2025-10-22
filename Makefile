include .env

DATABASE_URL=postgresql://${DATABASE_USER}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=disable

lint: lint-ui lint-go

lint-go:
	@echo "Linting..."
	@go mod tidy
	@golangci-lint run --fix ./...
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

# Create a new Atlas migration
# Usage: make migrate-new name=your_migration_name
migrate-new:
	@echo "Creating Atlas migration..."
	@cd database && atlas migrate diff "$(name)" --env local
	@echo "New migration created in database/migrations/"

# Apply all pending migrations
# Usage: make migrate-up
migrate-up:
	@echo "Applying Atlas migrations..."
	@cd database && atlas migrate apply --env local --url "$(DATABASE_URL)"
	@echo "All pending migrations applied."

# Check migration status
# Usage: make migrate-status
migrate-status:
	@echo "Checking migration status..."
	@cd database && atlas migrate status --env local --url "$(DATABASE_URL)"

# Validate migration directory
# Usage: make migrate-validate
migrate-validate:
	@echo "Validating migrations..."
	@cd database && atlas migrate validate --env local

# Generate migration hash/checksum
# Usage: make migrate-hash
migrate-hash:
	@echo "Generating migration checksums..."
	@cd database && atlas migrate hash

psql:
	@echo "Connecting to database..."
	@docker compose exec -it postgres psql -U ${DATABASE_USER} -d ${DATABASE_NAME}

deploy:
	@echo "Deploying..."
	@dokploy app deploy

help:
	@echo "Available commands:"
	@echo "  lint: Lint the code (Go + UI)"
	@echo "  generate: Generate code (Go + SQL + Swagger)"
	@echo "  build: Build the code (Go + UI + Docs)"
	@echo "  test: Run tests"
	@echo "  run: Run the full application"
	@echo "  run-go: Run Go backend only"
	@echo "  run-ui: Run UI frontend only"
	@echo ""
	@echo "Database commands:"
	@echo "  db-up: Start PostgreSQL database"
	@echo "  psql: Connect to database CLI"
	@echo "  migrate-new: Create new migration (make migrate-new name=your_migration_name)"
	@echo "  migrate-up: Apply pending migrations"
	@echo "  migrate-status: Check migration status"
	@echo "  migrate-validate: Validate migration files"
	@echo "  migrate-hash: Generate migration checksums"
	@echo ""
	@echo "Docker commands:"
	@echo "  docker-build: Build with docker"
	@echo "  docker-up: Start with docker compose"
	@echo "  docker-down: Stop docker compose"
	@echo ""
	@echo "  help: Show this help message"
