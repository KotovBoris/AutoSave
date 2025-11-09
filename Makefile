.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make run       - Run the application"
	@echo "  make build     - Build the application"
	@echo "  make test      - Run tests"
	@echo "  make download  - Download dependencies"

.PHONY: download
download:
	go mod download
	go mod tidy

.PHONY: run
run:
	go run cmd/api/main.go

.PHONY: build
build:
	go build -o bin/autosave cmd/api/main.go

.PHONY: test
test:
	go test -v ./...
