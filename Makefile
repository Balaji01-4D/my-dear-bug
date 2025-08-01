# My Dear Bug - Build and Run Commands

.PHONY: help build run clean install-deps test fmt lint

# Default target
help: ## Show all available commands
	@echo "My Dear Bug API - Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Dependencies
install-deps: ## Install all Go dependencies
	@echo "Installing Go dependencies..."
	go mod tidy
	@echo "Dependencies installed successfully"

# Build
build: ## Build the application binary
	@echo "Building application..."
	@mkdir -p build
	go build -o build/app main.go
	@echo "Build completed! Binary: ./build/app"

# Run
run: ## Start the development server
	@echo "Starting My Dear Bug API server..."
	@echo "API Server: http://localhost:8080"
	go run main.go

# Testing
test: ## Run the test suite
	@echo "Running tests..."
	go test ./...

# Code Quality
fmt: ## Format Go code
	@echo "Formatting Go code..."
	go fmt ./...
