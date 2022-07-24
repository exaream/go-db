.PHONY: install up login down clean build test cover

MAKEFILE_DIR := $(shell pwd)
CMD_DIR := $(MAKEFILE_DIR)/cmd/example
MYSQL_DIR := $(MAKEFILE_DIR)/_local/mysql/storage
POSTGRES_DIR := $(MAKEFILE_DIR)/_local/postgres/storage
PGADMIN_DIR := $(MAKEFILE_DIR)/_local/pgadmin/

############################################
# Run outside of Docker container.
############################################

install:
	docker compose up -d --build

up:
	docker compose up -d

login:
	docker exec -it go_db_app sh

down:
	docker compose down

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

build:
	go build -o $(CMD_DIR)/example $(CMD_DIR)/main.go

test:
	go test ./... -count=1

cover:
	go test ./... -count=1 -p=1 -coverprofile=cover.out
	go tool cover -html=cover.out -o cover.html
