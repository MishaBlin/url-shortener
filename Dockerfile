FROM golang:1.24 AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY config ./config
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -o url-shortener ./cmd/url-shortener

FROM alpine:3.21
RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /src/url-shortener .
COPY config/local.yaml ./config/local.yaml

ENV CONFIG_PATH=/app/config/local.yaml
EXPOSE 8080

CMD ["./url-shortener"]
