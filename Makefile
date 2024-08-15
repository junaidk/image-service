# Variables
DOCKER_COMPOSE_FILE = docker-compose.yml
APP_SERVICE = app
GO_CMD = go
DOCKER_COMPOSE_CMD = docker-compose
DOCKER_CMD = docker

# Build the Go application
.PHONY: build
build:
	$(GO_CMD) build -o bin/app ./cmd/server

# Build the Go application image
.PHONY: image
image:
	$(DOCKER_CMD) build -t image-service .

# Run the application with Docker Compose
.PHONY: up
up:
	$(DOCKER_COMPOSE_CMD) up -d

# Stop the application
.PHONY: down
down:
	$(DOCKER_COMPOSE_CMD) down

# Restart the application
.PHONY: restart
restart: down up

# Run tests
.PHONY: test
test:
	$(GO_CMD) test ./...

# Clean up generated files
.PHONY: clean
clean:
	rm -rf bin/

# View application logs
.PHONY: logs
logs:
	$(DOCKER_COMPOSE_CMD) logs -f $(APP_CONTAINER)

# Remove Docker containers, networks, volumes, and images
.PHONY: down-clear
down-clear:
	$(DOCKER_COMPOSE_CMD) down --volumes --rmi all --remove-orphans
