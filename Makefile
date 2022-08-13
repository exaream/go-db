MAKEFILE_DIR := $(shell pwd)
EXAMPLES_DIR := $(MAKEFILE_DIR)/cmd
CMD_DIR := $(EXAMPLES_DIR)/example
MYSQL_DIR := $(MAKEFILE_DIR)/_development/mysql/storage
POSTGRES_DIR := $(MAKEFILE_DIR)/_development/pgsql/storage
PGADMIN_DIR := $(MAKEFILE_DIR)/_development/pgadmin

############################################
# Run outside of go_db_app container.
############################################
.PHONY: install
install:
	docker compose up -d --build

.PHONY: up
up:
	docker compose up -d

.PHONY: login
login:
	docker exec -it go_db_app sh

.PHONY: down
down:
	docker compose down

.PHONY: clean
clean:
	docker compose down
	rm -rf $(MYSQL_DIR)/* \
	       $(POSTGRES_DIR)/* \
	       $(PGADMIN_DIR)/* \
	       $(CMD_DIR)/example \
	       $(MAKEFILE_DIR)/cover.out \
	       $(MAKEFILE_DIR)/cover.html
	touch $(MYSQL_DIR)/.gitkeep
	touch $(PGADMIN_DIR)/.gitkeep

############################################
# Run inside of go_db_app container.
############################################
.PHONY: setup
setup:
	go run $(CMD_DIR)/main.go --setup --path=$(CMD_DIR)/mysql.dsn
	go run $(CMD_DIR)/main.go --setup --path=$(CMD_DIR)/pgsql.dsn

.PHONY: check
check:
	govulncheck ./...
	golangci-lint run

.PHONY: build
build:
	go build -o $(CMD_DIR)/example $(CMD_DIR)/main.go

.PHONY: test
test:
	go test ./... -count=1

.PHONY: shuffle
shuffle:
	go test ./... -count=1 -shuffle=on

.PHONY: cover
cover:
	go test ./... -count=1 -coverprofile=cover.out
	go tool cover -html=cover.out -o cover.html
