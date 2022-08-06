# Go DB

## Overview

`dbutil` package is Go's CLI tool for operating MySQL and PostgreSQL.

## Install

Build some Docker containers and generate `users` table in MySQL and PostgreSQL.
```shell
$ git clone https://github.com/exaream/go-db.git
$ cd go-db
$ docker-compose up --build -d
```
Login `go_db_app` container for using Go.
```shell
$ docker exec -it go_db_app sh
```

## Setup initial data
Generate initial data in Docker containers as you like.
```shell
$ cd /go/src/work/cmd/example
$ go run main.go --setup --path=mysql.dsn
$ go run main.go --setup --path=pgsql.dsn
```

## Usage

### Sample

`example` package is a sample for updating `status` column of `users` table.  
Move to the following directory in `go_db_app` container.
```shell
$ cd /go/src/work/cmd/example
```
Show help
```shell
$ go run main.go --help
usage: example [<flags>]

An example command made of Go to operate MySQL and PostgreSQL.

Flags:
  --help                       Show context-sensitive help (also try --help-long and --help-man).
  --type="ini"                 Set a config type.
  --path="mysql.dsn"           Set a config file path.
  --section="example_section"  Set a config section name.
  --timeout=10s                Set a timeout value. e.g. 5s
  --id=0                       Set an ID.
  --before-sts=0               Set a before status.
  --after-sts=0                Set a after status.
  --setup                      Set true if you want to initialize data.
  --version                    Show application version.

```

Show version
```shell
$ go run main.go --version
```

Use minimum arguments
```shell
$ go run main.go --id=1 --before-sts=0 --after-sts=1
```

Use all arguments
```shell
$ go run main.go --type=ini --path=mysql.dsn --section=example_section --timeout=5s --id=1 --before-sts=0 --after-sts=1
```

### DB

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

## Test
Move to the working directory in `go_db_app` container.
```shell
$ cd /go/src/work/
```
Run unit tests in Docker container.
(Can NOT use `-race` option due to DB conflict)
```shell
$ go test ./... -count=1 -shuffle=on
```
Output coverage.
```shell
$ go test ./... -count=1 -coverprofile=cover.out
$ go tool cover -html=cover.out -o cover.html
```

## TODO
* Create a mechanism to avoid DB conflicts during testing by referring to [spool](https://github.com/cloudspannerecosystem/spool).
