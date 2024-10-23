```makefile
.PHONY: all build test clean docker-build docker-up docker-down migrate-up migrate-down

# Go related variables
BINARY_NAME=contractor-api
MAIN_PACKAGE=./cmd/api

# Docker related variables
DOCKER_COMPOSE=docker-compose

all: test build

build:
	go build -o ${BINARY_NAME} ${MAIN_PACKAGE}

test:
	go test -v ./...

clean:
	go clean
	rm -f ${BINARY_NAME}

docker-build:
	${DOCKER_COMPOSE} build

docker-up:
	${DOCKER_COMPOSE} up -d

docker-down:
	${DOCKER_COMPOSE} down

migrate-up:
	migrate -path migrations/postgres -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}" up

migrate-down:
	migrate -path migrations/postgres -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}" down
```