MAKEFILE_DIR := $(shell pwd)
CMD_DIR := $(MAKEFILE_DIR)/cmd/example
MYSQL_DIR := $(MAKEFILE_DIR)/_local/mysql/storage
POSTGRES_DIR := $(MAKEFILE_DIR)/_local/pgsql/storage
PGADMIN_DIR := $(MAKEFILE_DIR)/_local/pgadmin/

############################################
# Run outside of Docker container.
############################################
.PHONY: setup
setup:
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
# Run inside of Docker container.
############################################
.PHONY: init-data
init-data:
	go run $(CMD_DIR)/main.go --init-data --path=$(CMD_DIR)/mysql.dsn
	go run $(CMD_DIR)/main.go --init-data --path=$(CMD_DIR)/pgsql.dsn

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
