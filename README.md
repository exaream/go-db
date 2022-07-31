# Go DB

## Overview
A tool made of Go for operating MySQL and PostgreSQL.

## Install
```shell
$ git clone https://github.com/exaream/go-db.git
$ cd go-db
$ docker-compose up --build -d
$ docker exec -it go_db_app sh
```
Generate initial data in Docker container.
```shell
$ cd /go/src/work/cmd/example
$ go run main.go --init --path=mysql.dsn
$ go run main.go --init --path=pgsql.dsn
```

## Test
Run unit tests in Docker container.
```shell
$ cd /go/src/work/
$ go test ./... -count=1 -shuffle=on
```
Output coverage.
```shell
$ cd /go/src/work/
$ go test ./... -count=1 -coverprofile=cover.out
$ go tool cover -html=cover.out -o cover.html
```

## Usage
Help
```shell
$ cd /go/src/work/cmd/example
$ go run main.go --help
```

Version
```shell
$ cd /go/src/work/cmd/example
$ go run main.go --version
```

Command with minimum arguments
```shell
$ cd /go/src/work/cmd/example
$ go run main.go --id=1 --before-sts=0 --after-sts=1
```

Command with max arguments
```shell
$ cd /go/src/work/cmd/example
$ go run main.go --type=ini --path=mysql.dsn --section=example_section --timeout=5s --id=1 --before-sts=0 --after-sts=1
```

## DB

### How to access MySQL directly
```shell
$ docker container exec -it go_db_mysql sh
# mysql -h localhost -P 3306 -u exampleuser example_db -p
```

### How to access PostgreSQL directly
```shell
$ docker container exec -it go_db_pgsql sh
# psql -h localhost -p 5432 -U exampleuser example_db
```

### How to access phpMyAdmin
1. Check login info of `exampleuser` in `docker-compose.yml`
2. Access [http://localhost:8880/](http://localhost:8880/)

### How to access pgAdmin
1. Check login info of `pgadmin@example.com` in `docker-compose.yml`
2. Access [http://localhost:8888/](http://localhost:8888/)
