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
FROM debian:12-slim AS builder
RUN apt update && \
    apt install -y curl && \
    curl https://mise.run | MISE_INSTALL_PATH=/usr/local/bin/mise sh

WORKDIR /go/src/app
COPY mise.toml ./

RUN mise trust mise.toml && \
    mise install

COPY mise.toml go.mod go.sum ./
RUN mise x -- go mod download

COPY . .
COPY --from=ui-base /usr/src/app/dist web/dist
COPY --from=docs-base /usr/src/app/.vitepress/dist docs/.vitepress/dist
RUN mise build:go

# Final stage
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /service

COPY --from=builder /go/src/app/recally .

# Use non-root user for better security
USER nonroot:nonroot

# Expose the port the app runs on
# Expose the port specified by the PORT environment variable, defaulting to 1323
EXPOSE ${PORT:-1323}

CMD ["/service/recally"]
