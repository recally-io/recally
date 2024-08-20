# use the official Bun image
# see all versions at https://hub.docker.com/r/oven/bun/tags
FROM oven/bun:1 AS base
WORKDIR /usr/src/app

# install dependencies into temp directory
# this will cache them and speed up future builds
FROM base AS install
RUN mkdir -p /temp/dev
COPY web/package.json web/bun.lockb /temp/dev/
RUN cd /temp/dev && bun install --frozen-lockfile

# install with --production (exclude devDependencies)
RUN mkdir -p /temp/prod
COPY web/package.json web/bun.lockb /temp/prod/
RUN cd /temp/prod && bun install --frozen-lockfile --production

# copy node_modules from temp directory
# then copy all (non-ignored) project files into the image
FROM base AS prerelease
COPY --from=install /temp/dev/node_modules node_modules
COPY web .

# [optional] tests & build
ENV NODE_ENV=production
RUN bun run build

# Build go binary
FROM golang:1.22 AS build

WORKDIR /go/src/app

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
ENV CGO_ENABLED=0

COPY . .
COPY --from=prerelease /usr/src/app/dist web/dist
RUN ls web/dist && go build -ldflags="-s -w" -o /go/bin/app main.go

FROM gcr.io/distroless/base-debian12
WORKDIR /service/

COPY --from=build /go/bin/app .

CMD ["./app"]
