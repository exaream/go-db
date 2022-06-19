# Go DB

## Overview
This is a tool made of Go for operating MySQL.

## Install
```shell
$ docker-compose up --build -d
$ docker container exec -it go_db_app bash
```

### Usage
Help
```shell
# cd /go/src/work/cmd/example
# go run main.go --help
```

Version
```shell
# cd /go/src/work/cmd/example
# go run main.go --version
```

Max arguments
```shell
# cd /go/src/work/cmd/example
# go run main.go --ini-path=example.ini --section=example --timeout=10 --user-id=2 --status=1
```
Minimum arguments
```shell
# cd /go/src/work/cmd/example
# go run main.go --user-id=2 --status=1
```

How to access phpMyAdmin
1. Check login info of example user `exampleuser` in `docker-compose.yml`
2. [http://localhost:13902/](http://localhost:13902/)
