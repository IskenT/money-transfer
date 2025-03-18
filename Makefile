.PHONY: build run docs test clean

# Project name
PROJECT_NAME := money-transfer
BUILD_DIR := build

build:
	@echo "Building application..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(PROJECT_NAME) cmd/server/main.go

run:
	@go run cmd/server/main.go

docs:
	@echo "Generating Swagger documentation..."
	@$(shell go env GOPATH)/bin/swag init -g cmd/server/main.go --parseDependency --output docs

test:
	@echo "Running tests..."
	@go test -v ./...

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@go clean
