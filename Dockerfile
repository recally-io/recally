# Use the official Bun image
# See all versions at https://hub.docker.com/r/oven/bun/tags
# Build UI
FROM oven/bun:1 AS ui-base
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

# Build the application
ENV NODE_ENV=production
RUN bun run build

# Buld docs using Vitepress
FROM oven/bun:1 AS docs-base
WORKDIR /usr/src/app

RUN apt update && apt install -y git

# Install dependencies into temp directory
# This will cache them and speed up future builds
COPY docs/package.json docs/bun.lockb /temp/
RUN cd /temp && bun install --frozen-lockfile && \
    mkdir -p /usr/src/app/node_modules && \
    cp -r /temp/node_modules /usr/src/app/ && \
    rm -rf /temp

# Then copy all (non-ignored) project files into the image
COPY docs .

# Run tests and build
ENV NODE_ENV=production
RUN bun run docs:build


# Build Go binary
FROM golang:1.24-alpine AS build
WORKDIR /go/src/app

# Install Atlas CLI and other tools
RUN apk add --no-cache curl && \
    curl -sSf https://atlasgo.sh | sh -s -- --yes && \
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest && \
    go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download
ENV CGO_ENABLED=0 GOOS=linux

COPY . .
COPY --from=ui-base /usr/src/app/dist web/dist
COPY --from=docs-base /usr/src/app/.vitepress/dist docs/.vitepress/dist
RUN go generate ./... && \
    sqlc generate && \
    swag init -g internal/port/httpserver/router.go -o docs/swagger && \
    go build -ldflags="-s -w" -o /go/bin/app main.go

# Final stage
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /service

# Copy Atlas CLI binary and migrations
COPY --from=build /usr/local/bin/atlas /usr/local/bin/atlas
COPY --from=build /go/src/app/database/migrations /service/database/migrations

COPY --from=build /go/bin/app .

# Use non-root user for better security
USER nonroot:nonroot

# Expose the port the app runs on
# Expose the port specified by the PORT environment variable, defaulting to 1323
EXPOSE ${PORT:-1323}

CMD ["./app"]
