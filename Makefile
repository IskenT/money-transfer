.PHONY: build run docs test clean db-setup db-start db-stop migrate-check migrate-up migrate-down

# Project name
PROJECT_NAME := money-transfer
BUILD_DIR := build

build:
	@echo "Building application..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(PROJECT_NAME) cmd/server/main.go
	@go build -o $(BUILD_DIR)/migrate cmd/migrate/migrate.go

run: db-setup
	@echo "Starting the application..."
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

# Combined database and migration setup
db-setup: db-start migrate-check
	@echo "Database setup complete"

# Start 
db-start:
	@echo "Starting database using Podman..."
	@if ! command -v podman &> /dev/null; then \
		echo "Error: Podman is not installed. Please install Podman first."; \
		exit 1; \
	fi
	
	@# Check if the container image is available
	@if ! podman image exists postgres:16-alpine; then \
		echo "Pulling PostgreSQL image..."; \
		podman pull postgres:16-alpine || { echo "Error: Failed to pull PostgreSQL image. Check your internet connection."; exit 1; }; \
	fi
	
	@# Check if volume exists, create if not
	@if ! podman volume exists money_transfer_data; then \
		echo "Creating volume..."; \
		podman volume create money_transfer_data; \
	fi
	
	@# Check if container exists
	@if ! podman container exists money-transfer-db; then \
		echo "Creating and starting database container..."; \
		podman run -d --name money-transfer-db \
			-p 5432:5432 \
			-e POSTGRES_USER=money_transfer \
			-e POSTGRES_PASSWORD=password \
			-e POSTGRES_DB=money_transfer \
			-v money_transfer_data:/var/lib/postgresql/data \
			postgres:16-alpine || { echo "Error: Failed to start database container."; exit 1; }; \
	elif [ "$$(podman container inspect -f '{{.State.Running}}' money-transfer-db 2>/dev/null)" != "true" ]; then \
		echo "Starting existing database container..."; \
		podman start money-transfer-db || { echo "Error: Failed to start existing database container."; exit 1; }; \
	else \
		echo "Database container is already running."; \
	fi
	
	@echo "Waiting for database to be ready..."
	@for i in $$(seq 1 10); do \
		if podman exec money-transfer-db pg_isready -U money_transfer &> /dev/null; then \
			echo "Database is ready!"; \
			break; \
		fi; \
		if [ $$i -eq 10 ]; then \
			echo "Error: Database failed to become ready within the timeout period."; \
			exit 1; \
		fi; \
		echo "Waiting for database to start (attempt $$i/10)..."; \
		sleep 3; \
	done

# Check 
migrate-check:
	@echo "Checking and applying migrations if needed..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/migrate cmd/migrate/migrate.go
	
	@# Wait for database to be fully available before migrations
	@for i in $$(seq 1 10); do \
		if PGPASSWORD=password psql -h localhost -U money_transfer -d money_transfer -c "SELECT 1" &> /dev/null; then \
			echo "Database connection successful, running migrations..."; \
			$(BUILD_DIR)/migrate up; \
			break; \
		fi; \
		if [ $$i -eq 10 ]; then \
			echo "Error: Could not connect to database for migrations."; \
			exit 1; \
		fi; \
		echo "Waiting for database connection (attempt $$i/10)..."; \
		sleep 2; \
	done

db-stop:
	@echo "Stopping database..."
	@if podman container exists money-transfer-db && [ "$$(podman container inspect -f '{{.State.Running}}' money-transfer-db 2>/dev/null)" = "true" ]; then \
		podman stop money-transfer-db; \
		echo "Database stopped."; \
	else \
		echo "Database is not running."; \
	fi

#  migration commands
migrate-up:
	@echo "Running migrations up..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/migrate cmd/migrate/migrate.go
	@$(BUILD_DIR)/migrate up

migrate-down:
	@echo "Running migrations down..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/migrate cmd/migrate/migrate.go
	@$(BUILD_DIR)/migrate down

# Reset 
reset: db-stop clean
	@echo "Removing database container and volumes..."
	@if podman container exists money-transfer-db; then \
		podman rm -f money-transfer-db 2>/dev/null || true; \
		echo "Container removed."; \
	fi
	@podman volume rm -f money_transfer_data 2>/dev/null || true
	@echo "Reset completed"