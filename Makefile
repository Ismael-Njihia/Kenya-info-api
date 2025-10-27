.PHONY: help build run test clean docker-build docker-run swagger deps lint

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

deps: ## Download dependencies
	go mod download
	go mod tidy

build: ## Build the application
	go build -o bin/main cmd/api/main.go

run: ## Run the application
	go run cmd/api/main.go

test: ## Run tests
	go test -v -race -cover ./...

test-coverage: ## Run tests with coverage report
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint: ## Run linter
	golangci-lint run ./...

swagger: ## Generate swagger documentation
	swag init -g cmd/api/main.go -o docs

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html

docker-build: ## Build docker image
	docker build -t kenya-info-api:latest .

docker-run: ## Run docker container
	docker-compose up -d

docker-stop: ## Stop docker containers
	docker-compose down

docker-logs: ## View docker logs
	docker-compose logs -f

migrate-up: ## Run database migrations (placeholder)
	@echo "Migrations would run here"

migrate-down: ## Rollback database migrations (placeholder)
	@echo "Rollback would run here"

seed: ## Seed the database with sample data
	go run scripts/seed.go

dev: ## Run in development mode with hot reload (requires air)
	air

.DEFAULT_GOAL := help
