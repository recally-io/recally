# Use the official Bun image
# See all versions at https://hub.docker.com/r/oven/bun/tags
FROM oven/bun:1 AS base
WORKDIR /usr/src/app

# Install dependencies into temp directory
# This will cache them and speed up future builds
COPY web/package.json web/bun.lockb /temp/
RUN cd /temp && bun install --frozen-lockfile && \
    mkdir -p /usr/src/app/node_modules && \
    cp -r /temp/node_modules /usr/src/app/ && \
    rm -rf /temp

# Then copy all (non-ignored) project files into the image
COPY web .

# Run tests and build
ENV NODE_ENV=production
RUN bun test && bun run build

# Build Go binary
FROM golang:1.23-alpine AS build
WORKDIR /go/src/app

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest && \
    go install github.com/swaggo/swag/cmd/swag@latest && \
    go install github.com/kevinburke/go-bindata/v4/...@latest

COPY go.mod go.sum ./
RUN go mod download
ENV CGO_ENABLED=0 GOOS=linux

COPY . .
COPY --from=base /usr/src/app/dist web/dist
RUN go generate ./... && \
    go-bindata -prefix "database/migrations/" -pkg migrations -o database/bindata.go database/migrations/ && \
    sqlc generate && \
    swag init -g internal/port/httpserver/router.go && \
    go build -ldflags="-s -w" -o /go/bin/app main.go

# Final stage
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /service

COPY --from=build /go/bin/app .

# Use non-root user for better security
USER nonroot:nonroot

# Expose the port the app runs on
# Expose the port specified by the PORT environment variable, defaulting to 1323
EXPOSE ${PORT:-1323}

CMD ["./app"]
