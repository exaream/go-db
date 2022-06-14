# Go Ops

## Overview
This is a tool made of Go for operating MySQL.

## Install
```shell
$ docker network create go_ops_network
$ docker-compose up -d
$ docker container exec -it go_ops bash
# cd /go/src/work/ops/
# go mod tidy
```

### Usage
Help
```shell
# cd /go/src/work/ops/cmd/sample
# go run main.go --help
```

Max arguments
```shell
# cd /go/src/work/ops/cmd/sample
# go run main.go --ini-path=credentials/foo.ini --section=sample --timeout=10 --user-id=2 --status=1
```
Minimum arguments
```shell
# cd /go/src/work/ops/cmd/sample
# go run main.go --user-id=2 --status=1
```

How to access phpMyAdmin
1. Check login info of sample user `opsuser` in `docker-compose.yml`
2. [http://localhost:13902/](http://localhost:13902/)
