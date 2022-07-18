# Go DB

## Overview
This is a tool made of Go for operating MySQL.

## Install
```shell
$ docker-compose up --build -d
$ docker container exec -it go_db_app bash
```

## Test
Set 1 (process) to `-p` option to avoid conflicts when updating DB by multiple packages.
```shell
# cd /go/src/work/
# go test ./... -count=1 -p=1
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
# go run main.go --id=1 --before-sts=0 --after-sts=1
```

Max arguments
```shell
# cd /go/src/work/cmd/example
# go run main.go --type=ini --path=example.dsn --section=example_section --timeout=10s --id=1 --before-sts=0 --after-sts=1
```

How to access phpMyAdmin
1. Check login info of example user `exampleuser` in `docker-compose.yml`
2. [http://localhost:13902/](http://localhost:13902/)
