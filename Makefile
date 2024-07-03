PORT=8787


lint:
	@echo "Linting..."
	@golangci-lint run --fix ./...  --enable gofumpt


build:
	@echo "Building..."
	@go build -o bin/ ./cmd/httpserver

buildd:
	@echo "Build with docker"
	@docker build -t goworkers .

run:
	@echo "Running..."
	@go run ./cmd/httpserver

rund:
	@echo "Running with docker"
	@docker run -p $(PORT):$(PORT) --rm goworkers

