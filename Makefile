.DEFAULT_GOAL := up
SHELL := /bin/sh

EXAMPLES_DIR := $(CURDIR)/examples
BIN_DIR := $(CURDIR)/examples/example/cmd/example
MYSQL_DIR := $(CURDIR)/_development/mysql/storage
POSTGRES_DIR := $(CURDIR)/_development/pgsql/storage
PGADMIN_DIR := $(CURDIR)/_development/pgadmin

############################################
# Run before logging in go_db_app container.
############################################
.PHONY: up
up: ## start server
	@docker compose up --build -d

.PHONY: login
login: ## login server
	@docker exec -it go_db_app sh

.PHONY: down
down: ## stop server
	@docker compose down

.PHONY: clean
clean: ## remove binaries, tools and generated files
	@docker compose down
	@rm -rf $(MYSQL_DIR)/* \
	       $(POSTGRES_DIR)/* \
	       $(PGADMIN_DIR)/* \
	       $(BIN_DIR)/example \
	       $(CURDIR)/cover.out \
	       $(CURDIR)/cover.html
	@touch $(MYSQL_DIR)/.gitkeep
	@touch $(PGADMIN_DIR)/.gitkeep

############################################
# Run after logging in go_db_app container.
############################################
.PHONY: setup
setup: ## set up DB
	@go run $(BIN_DIR)/main.go --setup --path=$(BIN_DIR)/mysql.dsn
	@go run $(BIN_DIR)/main.go --setup --path=$(BIN_DIR)/pgsql.dsn

.PHONY: vuln
vuln: ## check vulnerability
	@govulncheck ./...

.PHONY: lint
lint: ## un linter
	@golangci-lint run ./...


.PHONY: build
build: ## generate binary
	@go build -buildvcs=false -o $(BIN_DIR)/example $(BIN_DIR)/main.go

.PHONY: test
test: ## run all tests
	@go test ./... -count=1 -shuffle=on

.PHONY: cover
cover: ## generate coverage report
	@go test ./... -count=1 -coverprofile=cover.out
	@go tool cover -html=cover.out -o cover.html

.PHONY: help
help: ## print help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
