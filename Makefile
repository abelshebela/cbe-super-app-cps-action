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
	mockgen -source=./internal/port/outbound/event/init.go -destination=./mocks/port/outbound/event/init.go -package=mock_outbound
	mockgen -source=./internal/port/outbound/mini_app/outbound.go -destination=./mocks/port/outbound/miniapp/outbound.go -package=mock_outbound

mockgenInbound:
	mockgen -source=./internal/port/inbound/bulk_services/inbound.go -destination=./mocks/port/inbound/bulk_services/inbound.go -package=mock_inbound
	mockgen -source=./internal/port/inbound/event/init.go -destination=./mocks/port/inbound/event/inbound.go -package=mock_inbound
	mockgen -source=./internal/port/inbound/miniapp/inbound.go -destination=./mocks/port/inbound/miniapp/inbound.go -package=mock_inbound
	
mockgenService:
	mockgen -source=./internal/domain/action/service_impl.go -destination=./mocks/domain/action/service_impl.go -package=mock_domain
	mockgen -source=./internal/domain/event/service.go -destination=./mocks/domain/event/service.go -package=mock_domain
	mockgen -source=./internal/domain/miniapp/service_impl.go -destination=./mocks/domain/action/service_impl.go -package=mock_domain
	mockgen -source=./internal/domain/service/service.go -destination=./mocks/domain/service/service_impl.go -package=mock_domain

mockgenRepository:
	mockgen -source=./internal/domain/action/repository.go -destination=./mocks/domain/action/repository.go -package=mock_domain
	mockgen -source=./internal/domain/event/repository.go -destination=./mocks/domain/event/repository.go -package=mock_domain
	mockgen -source=./internal/domain/miniapp/repository.go -destination=./mocks/domain/miniapp/repository.go -package=mock_domain
	mockgen -source=./internal/domain/service/repository.go -destination=./mocks/domain/service/repository.go -package=mock_domain
	mockgen -source=./internal/domain/department/repository.go -destination=./internal/domain/department/mocks/repo_mock.go -package=mocks

mockgen_Mongo:
	mockgen -source=./internal/adapter/outbound/mongo/init.go -destination=./mocks/port/outbound/mongo/mongo.go -package=mock_mongo