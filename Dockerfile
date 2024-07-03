FROM golang:1.22 as build

WORKDIR /go/src/app

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
ENV CGO_ENABLED=0

COPY . .
RUN go build -ldflags="-s -w" -o /go/bin/app ./cmd/httpserver


FROM gcr.io/distroless/base-debian12
WORKDIR /service/

COPY --from=build /go/bin/app .

CMD ["./app"]
