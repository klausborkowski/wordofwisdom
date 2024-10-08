DOCKER_COMPOSE_FILE = docker-compose.yml
SERVER_DIR = cmd/server
CLIENT_DIR = cmd/client

.PHONY: help install test start-server start-client start up down restart build logs ps clean

help:
	@echo "Available targets:"
	@echo "  make install          - Download Go module dependencies"
	@echo "  make start-server     - Run the server application"
	@echo "  make start-client     - Run the client application"
	@echo "  make start            - Build and start Docker Compose services for server and client"
	@echo "  make up               - Start all Docker Compose services in detached mode"
	@echo "  make down             - Stop and remove Docker Compose services"
	@echo "  make restart          - Restart Docker Compose services"
	@echo "  make build            - Build Docker Compose services"
	@echo "  make logs             - Tail logs for Docker Compose services"
	@echo "  make ps               - List running Docker Compose containers"
	@echo "  make clean            - Remove stopped containers, unused volumes, and prune Docker system"

install:
	@echo "Downloading Go module dependencies..."
	go mod download


start-server:
	@echo "Starting server application..."
	go run $(SERVER_DIR)/main.go

start-client:
	@echo "Starting client application..."
	go run $(CLIENT_DIR)/main.go

start:
	@echo "Building and starting Docker Compose services for server and client..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) up --abort-on-container-exit --force-recreate --build server client


up:
	@echo "Starting all Docker Compose services..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

down:
	@echo "Stopping and removing all Docker Compose services..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

restart:
	@echo "Restarting Docker Compose services..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

build:
	@echo "Building Docker Compose services..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) build

logs:
	@echo "Tailing logs for Docker Compose services..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f

ps:
	@echo "Listing running Docker Compose containers..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) ps

clean:
	@echo "Removing Docker Compose services and pruning Docker system..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down --volumes --remove-orphans
	docker system prune -f

