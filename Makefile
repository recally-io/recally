lint:
	@echo "Linting..."
	@golangci-lint run --fix ./...  --enable gofumpt


build:
	@echo "Building..."
	@go build -o bin/ ./cmd/httpserver

run: build
	@echo "Running..."
	@go run ./cmd/httpserver
