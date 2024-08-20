include .env

DATABASE_URL=postgresql://${DATABASE_USER}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=disable

lint: lint-ui lint-go

lint-go:
	@echo "Linting..."
	@go mod tidy
	@golangci-lint run --fix ./...  --enable gofumpt

lint-ui:
	@echo "Linting web..."
	@prettier ./web --write

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
	@swag init -g internal/port/httpserver/router.go
	# @echo "Generating SDK ..."
	# @mkdir -p web/src/sdk && rm -rf web/src/sdk
	# @openapi-generator generate --skip-validate-spec -i docs/swagger.yaml -g typescript-fetch -o web/src/sdk

build: build-ui build-go

build-go:
	@echo "Building go..."
	@go build -o bin/app main.go

build-ui:
	@echo "Building UI..."
	@cd web && bun run build

test:
	@echo "Testing..."
	@go test ./...

run: run-go

run-go:
	@echo "Running go..."
	@go run main.go

run-ui:
	@echo "Running web..."
	@cd web && bun run dev

run-dev:
	@echo "Running dev..."
	@ENV=dev go run main.go

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

migrate-new:
	@echo "Creating migration..."
	@migrate create -ext sql -dir database/migrations -seq "$(name)"

migrate-up:
	@echo "Migrating up..."
	@migrate -path database/migrations -database "$(DATABASE_URL)" up

migrate-down:
	@echo "Migrating down..."
	@migrate -path database/migrations -database "$(DATABASE_URL)" down

psql:
	@echo "Connecting to database..."
	@docker compose exec -it postgres psql -U ${DATABASE_USER} -d ${DATABASE_NAME}

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
