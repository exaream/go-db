.PHONY: install login clean build test cover

CMD_DIR := $(shell pwd)/cmd/example
MYSQL_DIR := $(shell pwd)/_local/mysql/storage
POSTGRES_DIR := $(shell pwd)/_local/postgres/storage

############################################
# Run outside of Docker container.
############################################

install:
	docker compose up -d --build

login:
	docker exec -it go_db_app sh

clean:
	docker compose down
	rm -rf $(MYSQL_DIR)/*
	rm -rf $(POSTGRES_DIR)/*
	rm -rf $(CMD_DIR)/example
	rm -rf $(shell pwd)/cover.out
	rm -rf $(shell pwd)/cover.html
	touch $(MYSQL_DIR)/.gitkeep
	touch $(POSTGRES_DIR)/.gitkeep

############################################
# Run inside of Docker container.
############################################

build:
	go build -o $(CMD_DIR)/example $(CMD_DIR)/main.go

test:
	go test ./... -count=1 -p=1

cover:
	go test ./... -count=1 -p=1 -coverprofile=cover.out
	go tool cover -html=cover.out -o cover.html
