# Go DB

## Overview
This is a tool made of Go for operating MySQL and PostgreSQL.

## Install
```shell
$ docker-compose up --build -d
$ docker exec -it go_db_app sh
```
Generate initial data.
```shell
$ cd /go/src/work/cmd/example
$ go run main.go --init --path=mysql.dsn
$ go run main.go --init --path=postgres.dsn
```

## Test
```shell
$ cd /go/src/work/
$ go test ./... -count=1
```

### Usage
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

Minimum arguments
```shell
$ cd /go/src/work/cmd/example
$ go run main.go --id=1 --before-sts=0 --after-sts=1
```

Max arguments
```shell
$ cd /go/src/work/cmd/example
$ go run main.go --type=ini --path=mysql.dsn --section=example_section --timeout=5s --id=1 --before-sts=0 --after-sts=1
```

How to access phpMyAdmin
1. Check login info of `exampleuser` in `docker-compose.yml`
2. Access [http://localhost:8880/](http://localhost:8880/)

How to access pgAdmin
1. Check login info of `pgadmin@example.com` in `docker-compose.yml`
2. Access [http://localhost:8888/](http://localhost:8888/)
