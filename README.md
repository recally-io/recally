# GoWorkers

Using Golang to build cloudflare workers.

## OpenAPI

Using AI to generate OpenAPI specs and then using [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) to generate the go server code.

## DB

Using AI to generate the DB schema and then using [sqlc](https://github.com/sqlc-dev/sqlc) to generate the go code.

- [sqlc](https://github.com/sqlc-dev/sqlc) generate go code from SQL
- [migrate](https://github.com/golang-migrate/migrate) to manage DB migrations
- [pgx](https://github.com/jackc/pgx) for postgres driver

## Setup

```
go build -o bin/go-workers cmd/httpserver/main.go

./bin/go-workers
```

## API

- /docs/ui - Swagger UI
- /docs/json - Swagger JSON
- /docs/redoc - Redoc UI
