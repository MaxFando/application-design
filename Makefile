PROJECT := github.com/MaxFando/application-design

appName = ms-application_design_test
compose = docker-compose -f docker-compose-local.yml -p $(appName)
compose-debug = docker-compose -f docker-compose-debug.yml -p $(appName)

up: down build
	@echo "Starting app..."
	$(compose) up -d
	@echo "Docker images built and started!"

build:
	@echo "Building images"
	$(compose) build
	@echo "Docker images built!"

down:
	@echo "Stopping docker compose..."
	$(compose) down
	@echo "Done!"

test:
	go test ./... -coverprofile=coverage.out -coverpkg ./internal/core/... && go tool cover -func coverage.out

lint:
	golangci-lint run -c .golangci.yaml

generate:
	@echo "Generating code..."
	go generate ./...
	@echo "Code generated!"
