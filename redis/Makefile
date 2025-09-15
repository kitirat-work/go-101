APP_NAME=go-redis-app
DOCKER_COMPOSE=docker-compose.yml

.PHONY: run build docker-up docker-down docker-logs

run:
	go run ./cmd/main.go

build:
	go build -o $(APP_NAME) ./cmd/main.go

docker-up:
	docker compose -f $(DOCKER_COMPOSE) up -d --build

docker-down:
	docker compose -f $(DOCKER_COMPOSE) down

docker-logs:
	docker compose -f $(DOCKER_COMPOSE) logs -f
