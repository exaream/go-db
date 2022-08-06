# Go DB

## Overview
`dbutil` package is a tool made of Go for operating MySQL and PostgreSQL.

## Setup
The following command builds some Docker container and generates `users` table in MySQL and PostgreSQL.
```shell
$ git clone https://github.com/exaream/go-db.git
$ cd go-db
$ docker-compose up --build -d
```
Login `go_db_app` container for using Go.
```shell
$ docker exec -it go_db_app sh
```
Generate initial data in Docker container as you like.
```shell
$ cd /go/src/work/_examples/example
$ go run main.go --init-data --path=mysql.dsn
$ go run main.go --init-data --path=pgsql.dsn
```

## Test
Run unit tests in Docker container.
(Can NOT use `-race` option due to DB conflict)
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
Show help
```shell
$ cd /go/src/work/cmd/example
$ go run main.go --help
```

Show version
```shell
$ cd /go/src/work/cmd/example
$ go run main.go --version
```

Use minimum arguments
```shell
$ cd /go/src/work/cmd/example
$ go run main.go --id=1 --before-sts=0 --after-sts=1
```

Use all arguments
```shell
$ cd /go/src/work/cmd/example
$ go run main.go --type=ini --path=mysql.dsn --section=example_section --timeout=5s --id=1 --before-sts=0 --after-sts=1
```

## DB

Access MySQL directly
```shell
$ docker container exec -it go_db_mysql sh
# mysql -h localhost -P 3306 -u exampleuser example_db -p
```

Access PostgreSQL directly
```shell
$ docker container exec -it go_db_pgsql sh
# psql -h localhost -p 5432 -U exampleuser example_db
```

Access phpMyAdmin
1. Check login info of `exampleuser` in `docker-compose.yml`
2. Access [http://localhost:8880/](http://localhost:8880/)

Aaccess pgAdmin
1. Check login info of `pgadmin@example.com` in `docker-compose.yml`
2. Access [http://localhost:8888/](http://localhost:8888/)

## TODO
* Create a mechanism to avoid DB conflicts during testing by referring to [spool](https://github.com/cloudspannerecosystem/spool).
