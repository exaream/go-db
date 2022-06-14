# Go Ops

## Overview
This is a tool made of Go for operating MySQL.

## Install
```shell
$ docker network create go_ops_network
$ docker-compose up -d
$ docker container exec -it go_ops bash
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
