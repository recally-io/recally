include .env

lint:
	@echo "Linting..."
	@golangci-lint run --fix ./...  --enable gofumpt

generate:
	@echo "Generating..."
	@go generate ./...
	@go-bindata -prefix "database/migrations/" -pkg migrations -o database/bindata.go database/migrations/

build: generate lint
	@echo "Building..."
	@go build -o bin/app main.go

test: lint
	@echo "Testing..."
	@go test ./...

run: build
	@echo "Running..."
	@./bin/app

docker-build:
	@echo "Building with docker"
	@docker compose --build

docker-up:
	@echo "Running with docker"
	@docker compose up --build -d

docker-down:
	@echo "Stopping docker"
	@docker compose down

migrate-new:
	@echo "Creating migration..."
	@migrate create -ext sql -dir db/migrations -seq "$(name)"

# migrate-up:
# 	@echo "Migrating up..."
# 	@migrate -path db/migrations -database "$(DATABASE_URL)" up

# migrate-down:
# 	@echo "Migrating down..."
# 	@migrate -path db/migrations -database "$(DATABASE_URL)" down

sqlc:
	@echo "Generating sqlc..."
	@sqlc generate

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
