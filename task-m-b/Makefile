APP_NAME := task-manager
CMD_APP  := ./cmd/app
BIN_DIR  := bin
BIN      := $(BIN_DIR)/app
GO       := go

# Load env vars from .env if present (for compose/migrations)
ifneq (,$(wildcard .env))
include .env
export
endif

.PHONY: help fmt vet tidy build build-local run test mocks vendor docker-build up down ps logs migrate db-shell health clean

help:
	@echo "Common targets:"
	@echo "  make build         - build Linux static binary to $(BIN)"
	@echo "  make build-local   - build for current OS/arch"
	@echo "  make run           - run app locally with config/local.yaml"
	@echo "  make test          - run tests"
	@echo "  make mocks         - generate mocks (go:generate)"
	@echo "  make docker-build  - build Docker image $(APP_NAME)"
	@echo "  make up            - docker compose up (app+db+migrate)"
	@echo "  make down          - docker compose down"
	@echo "  make logs          - tail app logs"
	@echo "  make migrate       - run migrations service once"
	@echo "  make db-shell      - psql into db container"
	@echo "  make health        - curl app health endpoint"
	@echo "  make tidy vendor   - tidy modules / vendor deps"
	@echo "  make clean         - remove build artifacts"

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

tidy:
	$(GO) mod tidy

$(BIN):
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -ldflags "-s -w" -o $(BIN) $(CMD_APP)

build: $(BIN)

build-local:
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN) $(CMD_APP)

run:
	CONFIG_PATH=config/local.yaml $(GO) run $(CMD_APP)

test:
	$(GO) test ./...

mocks:
	$(GO) generate ./...

vendor:
	$(GO) mod vendor

docker-build:
	docker build -t $(APP_NAME) .

up:
	docker compose up -d --build

down:
	docker compose down

ps:
	docker compose ps

logs:
	docker compose logs -f app

run-migrate:

	docker compose run --rm migrate


db-shell:
	# Open psql shell inside the db container
	docker compose exec db psql -U $(DB_USER) -d $(DB_NAME)

health:
	# Check app health endpoint via published port
	curl -sf http://localhost:8080/healthz || true

clean:
	rm -rf $(BIN_DIR)

