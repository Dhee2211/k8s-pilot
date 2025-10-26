.PHONY: help build test lint clean install docker run-kind

# Variables
BINARY_NAME=kubectl-pilot
VERSION?=dev
COMMIT=$(shell git rev-parse --short HEAD)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${BUILD_TIME}"

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the binary
	go build ${LDFLAGS} -o ${BINARY_NAME} .

build-all: ## Build for all platforms
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-windows-amd64.exe .

install: build ## Install the binary
	cp ${BINARY_NAME} ${GOPATH}/bin/${BINARY_NAME}

test: ## Run tests
	go test -v -race -coverprofile=coverage.out ./...

test-integration: ## Run integration tests (requires Kind cluster)
	go test -v -tags=integration ./tests/...

coverage: test ## Generate test coverage report
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

lint: ## Run linters
	golangci-lint run ./...

fmt: ## Format code
	go fmt ./...
	gofmt -s -w .

vet: ## Run go vet
	go vet ./...

clean: ## Clean build artifacts
	rm -f ${BINARY_NAME}
	rm -rf bin/
	rm -f coverage.out coverage.html

deps: ## Download dependencies
	go mod download
	go mod tidy

update-deps: ## Update dependencies
	go get -u ./...
	go mod tidy

docker-build: ## Build Docker image
	docker build -t ${BINARY_NAME}:${VERSION} .
	docker tag ${BINARY_NAME}:${VERSION} ${BINARY_NAME}:latest

docker-run: ## Run Docker container
	docker run --rm -it ${BINARY_NAME}:latest

run-kind: ## Create a Kind cluster for testing
	kind create cluster --name pilot-test || true
	kubectl cluster-info --context kind-pilot-test

delete-kind: ## Delete the Kind test cluster
	kind delete cluster --name pilot-test

run: build ## Build and run the application
	./${BINARY_NAME} --help

run-example: build ## Run an example command
	./${BINARY_NAME} run "list pods" -v

diagnose-example: build ## Run diagnostics example
	./${BINARY_NAME} diagnose --all-namespaces

release: ## Create a new release (requires goreleaser)
	goreleaser release --clean

release-snapshot: ## Create a snapshot release
	goreleaser release --snapshot --clean

dev: ## Run in development mode with hot reload
	air -c .air.toml

benchmark: ## Run benchmarks
	go test -bench=. -benchmem ./...

# Development helpers
watch: ## Watch for changes and rebuild
	@while true; do \
		inotifywait -e modify,create,delete -r . && \
		make build && \
		echo "Rebuilt at $$(date)"; \
	done

.DEFAULT_GOAL := help
