FROM golang:1.25.8-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/trainer-service ./cmd/server

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/trainer-service .
COPY --from=builder /app/config ./config

EXPOSE 8080

CMD ["./trainer-service"]
