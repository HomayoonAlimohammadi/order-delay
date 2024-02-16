# don't forget to change docker-compose.yaml accordingly
POSTGRES_URL ?= "postgres://user:pass@localhost:5432/order_delay?sslmode=disable"
MOCKERY_VERSION := v2.40.3

MOCK_PACKAGES := \
	internal/core \
	internal/database/querier

MOCK_FOLDERS = $(patsubst %,%/mocks,$(MOCK_PACKAGES))

.SECONDEXPANSION:
$(MOCK_FOLDERS): %/mocks: $$(wildcard $$*/*.go)
	rm -rf $@
	cd $(patsubst %/mocks,%,$@) && mockery --name ".*"

.PHONY: proto
proto:
	protoc -I=./proto --go_out=./proto --go_opt=paths=source_relative \
	 --go-grpc_out=./proto --go-grpc_opt=paths=source_relative order.proto

.PHONY: migrate-up
migrate-up: 
	migrate -database ${POSTGRES_URL} -path internal/database/migrations -verbose up

.PHONY: migrate-down
migrate-down: 
	migrate -database ${POSTGRES_URL} -path internal/database/migrations -verbose down -all

.PHONY: sqlc
sqlc:
	sqlc generate -f internal/database/sqlc.yaml

.PHONY: test
test:
	go test -race ./...

.PHONY: lint
lint:
	golangci-lint run . # TODO: add proper custom config

.PHONY: up
up:
	docker compose up -d
	$(MAKE) migrate-up

.PHONY: down
down:
	$(MAKE) migrate-down
	docker compose down

.PHONY: generate 
generate: proto sqlc mocks
	@echo Done!

.PHONY: mocks
mocks: mockery_check $(MOCK_FOLDERS)

cleanmocks:
	rm -rf $(shell find . -path "*/mocks")

mockery_check:
	@if [ "$$(mockery --version)" = $(MOCKERY_VERSION) ]; then \
		echo "Mockery version is correct"; \
	else \
		echo "Installing correct mockery version $(MOCKERY_VERSION)..."; \
		go install github.com/vektra/mockery/v2@$(MOCKERY_VERSION); \
	fi
