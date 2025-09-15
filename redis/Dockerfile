FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o go-redis-app ./cmd/main.go

FROM golang:1.24-alpine
WORKDIR /app
COPY --from=builder /app/go-redis-app ./go-redis-app
COPY .env.example .env
EXPOSE 8080
CMD ["./go-redis-app"]
