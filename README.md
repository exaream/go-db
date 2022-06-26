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

Minimum arguments
```shell
# cd /go/src/work/cmd/example
# go run main.go --id=1 --beforeSts=0 --AfterSts=1
```

Max arguments
```shell
# cd /go/src/work/cmd/example
# go run main.go --type=ini --path=example.dsn --section=example_section --timeout=10s --id=1 --beforeSts=0 --AfterSts=1
```

How to access phpMyAdmin
1. Check login info of example user `exampleuser` in `docker-compose.yml`
2. [http://localhost:13902/](http://localhost:13902/)
