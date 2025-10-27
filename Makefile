.PHONY: help build run test clean docker-build docker-up docker-down swagger install

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

install: ## Install dependencies
	go mod download
	go mod tidy

build: ## Build the application
	go build -o bin/api ./cmd/api

run: ## Run the application
	go run ./cmd/api/main.go

test: ## Run tests
	go test -v -cover ./...

test-coverage: ## Run tests with coverage report
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html

swagger: ## Generate Swagger documentation
	swag init -g cmd/api/main.go -o docs

docker-build: ## Build Docker image
	docker-compose build

docker-up: ## Start Docker containers
	docker-compose up -d

docker-down: ## Stop Docker containers
	docker-compose down

docker-logs: ## View Docker logs
	docker-compose logs -f

lint: ## Run linter
	golangci-lint run ./...

seed: ## Seed database with sample data
	go run ./scripts/seed.go
