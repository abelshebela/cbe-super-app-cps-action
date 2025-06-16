# Makefile for Go project
.PHONY: all format lint test mockgenInfra mockgenService mockgen_Mongo

# Default target: run all checks
all: format lint test

# Format code with gofmt and goimports
format:
	gofmt -s -w .
	goimports -w .

# Run linters with golangci-lint
lint:
	golangci-lint run ./...

# Run tests with race detection and coverage
test:
	go test ./... -race -cover

mockgenOutbound:
	mockgen -source=./internal/port/outbound/bulk_services/outbound.go -destination=./mocks/port/outbound/bulk_services/outbound.go -package=mock_outbound

mockgenInbound:
	mockgen -source=./internal/port/inbound/bulk_services/inbound.go -destination=./mocks/port/inbound/bulk_services/inbound.go -package=mock_inbound

mockgenService:
	mockgen -source=./internal/domain/action/service_impl.go -destination=./mocks/domain/action/service_impl.go -package=mock_services

mockgenRepository:
	mockgen -source=./internal/domain/action/repository.go -destination=./mocks/domain/repository/repository.go -package=mock_repository

mockgen_Mongo:
	mockgen -source=./internal/adapter/outbound/mongo/init.go -destination=./mocks/port/outbound/mongo/mongo.go -package=mock_mongo