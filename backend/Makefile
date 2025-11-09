.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make run                - Run application locally"
	@echo "  make build              - Build application"
	@echo "  make test               - Run tests"
	@echo "  make download           - Download dependencies"
	@echo "  make docker-up          - Start all services in Docker"
	@echo "  make docker-down        - Stop all Docker services"
	@echo "  make docker-build       - Build Docker images"
	@echo "  make docker-logs        - View Docker logs"
	@echo "  make docker-ps          - List Docker containers"
	@echo "  make docker-exec-api    - Shell into API container"
	@echo "  make docker-exec-db     - Shell into DB container"
	@echo "  make migrate-up         - Run migrations up"
	@echo "  make migrate-down       - Run migrations down"
	@echo "  make migrate-create     - Create new migration"
	@echo "  make db-shell           - PostgreSQL shell"
	@echo "  make redis-cli          - Redis CLI"
	@echo "  make clean              - Clean build artifacts"

.PHONY: download
download:
	go mod download
	go mod tidy

.PHONY: run
run:
	go run cmd/api/main.go

.PHONY: build
build:
	go build -o autosave-backend cmd/api/main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: test-coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: lint
lint:
	golangci-lint run

.PHONY: docker-up
docker-up:
	cd .. \&\& docker-compose up -d

.PHONY: docker-down
docker-down:
	cd .. \&\& docker-compose down

.PHONY: docker-build
docker-build:
	cd .. \&\& docker-compose build

.PHONY: docker-logs
docker-logs:
	cd .. \&\& docker-compose logs -f

.PHONY: docker-ps
docker-ps:
	cd .. \&\& docker-compose ps

.PHONY: docker-exec-api
docker-exec-api:
	cd .. \&\& docker-compose exec api sh

.PHONY: docker-exec-db
docker-exec-db:
	cd .. \&\& docker-compose exec postgres bash

.PHONY: migrate-up
migrate-up:
	cd .. \&\& docker-compose --profile tools run --rm migrate

.PHONY: migrate-down
migrate-down:
	cd .. \&\& docker-compose run --rm migrate -path /migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@postgres:5432/${DB_NAME}?sslmode=disable" down 1

.PHONY: migrate-create
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

.PHONY: db-shell
db-shell:
	cd .. \&\& docker-compose exec postgres psql -U ${DB_USER} -d ${DB_NAME}

.PHONY: redis-cli
redis-cli:
	cd .. \&\& docker-compose exec redis redis-cli

.PHONY: clean
clean:
	rm -rf autosave-backend tmp/ *.out coverage.html

.PHONY: docker-clean
docker-clean: docker-down
	cd .. \&\& docker-compose down -v
	docker system prune -f

.PHONY: install-tools
install-tools:
	go install github.com/cosmtrek/air@latest
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: prod-build
prod-build:
	docker build -f Dockerfile -t autosave-backend:latest .

.PHONY: prod-up
prod-up:
	cd .. \&\& docker-compose -f cd .. \&\& docker-compose.yml -f cd .. \&\& docker-compose.prod.yml up -d

.PHONY: prod-down
prod-down:
	cd .. \&\& docker-compose -f cd .. \&\& docker-compose.yml -f cd .. \&\& docker-compose.prod.yml down


