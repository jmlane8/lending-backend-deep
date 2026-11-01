# ---------- local settings ----------

# env.local holds the localhost settings. Every value is exported to the recipes below.
-include env.local
export

COMPOSE = docker compose --env-file env.local

# ---------- env helpers ----------

.PHONY: env
env: ## Copy .env.example → .env (run once; adjust values as needed)
ifeq ($(OS),Windows_NT)
	copy .env.example .env
else
	cp .env.example .env
endif

# ---------- run ----------

# Settings come from env.local (see DEVELOPMENT.md).
.PHONY: run
run: ## Run the HTTP server (migrations run automatically at startup via golang-migrate)
	go run ./cmd/http

# ---------- local database ----------

.PHONY: db-up
db-up: ## Start local PostgreSQL in Docker and wait until it is ready
	$(COMPOSE) up -d --wait

.PHONY: db-seed
db-seed: ## Load the sandbox partners (needs the tables, so run the server first)
	$(COMPOSE) exec -T postgres psql -U $(POSTGRES_USERNAME) -d $(POSTGRES_DATABASE) -v ON_ERROR_STOP=1 -f /seed/sandbox_partner.sql

.PHONY: db-down
db-down: ## Stop local PostgreSQL and keep its data
	$(COMPOSE) down

.PHONY: db-reset
db-reset: ## Stop local PostgreSQL and delete its data
	$(COMPOSE) down -v

# ---------- OpenAPI spec ----------

REDOCLY_VERSION = 2.25.1

ifeq ($(OS),Windows_NT)
  REDOCLY_CHECK := $(shell redocly --version 2>NUL)
else
  REDOCLY_CHECK := $(shell redocly --version 2>/dev/null)
endif

.PHONY: openapi
openapi: ## Bundle OpenAPI YAML → resources/gen/openapi/openapi.json (embedded in binary)
ifndef REDOCLY_CHECK
	npm install -g @redocly/cli@$(REDOCLY_VERSION)
endif
	redocly bundle ./resources/openapi/openapi.yml \
	  --output ./resources/gen/openapi/openapi.json --ext json

# ---------- test ----------

COVERAGE_OUTPUT      = coverage.out
COVERAGE_OUTPUT_HTML = coverage.html

.PHONY: test
test: ## Run tests with race detector and generate HTML coverage report
	go test -v -race -failfast \
	  -coverprofile=$(COVERAGE_OUTPUT) -coverpkg=./... \
	  ./...
	go tool cover -html=$(COVERAGE_OUTPUT) -o $(COVERAGE_OUTPUT_HTML)

# ---------- deps ----------

.PHONY: vendor
vendor: ## Tidy go.mod and vendor all dependencies
	go mod tidy && go mod vendor

# ---------- help ----------

.PHONY: help
help: ## Show this help message
ifeq ($(OS),Windows_NT)
	@powershell -NoProfile -Command "Get-Content '$(MAKEFILE_LIST)' | Where-Object { $$_ -match '^[a-zA-Z_/-]+:.*?## ' } | ForEach-Object { $$parts = $$_ -split ':.*?## '; '{0,-22} {1}' -f $$parts[0], $$parts[1] }"
else
	@grep -E '^[a-zA-Z_/-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	  | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-22s\033[0m %s\n", $$1, $$2}'
endif

.DEFAULT_GOAL := help
