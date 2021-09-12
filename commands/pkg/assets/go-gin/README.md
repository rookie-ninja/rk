# Demo
**[中文版](README_CN.md)**

Rk style golang server package with gin entry enabled.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Prerequisite](#prerequisite)
- [Install tools](#install-tools)
- [Build locally](#build-locally)
  - [With RK](#with-rk)
  - [With go build](#with-go-build)
- [Run locally](#run-locally)
  - [With RK](#with-rk-1)
  - [With go run](#with-go-run)
  - [With IDE](#with-ide)
- [Pack locally](#pack-locally)
  - [With RK](#with-rk-2)
- [Docker build](#docker-build)
  - [With RK](#with-rk-3)
  - [With docker build](#with-docker-build)
- [Docker run](#docker-run)
  - [With RK](#with-rk-4)
- [Clear locally](#clear-locally)
- [boot.yaml](#bootyaml)
- [build.yaml](#buildyaml)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Prerequisite
[rk](https://github.com/rookie-ninja/rk) command line tool is recommended in order to build project at enterprise level.

```shell script
go get github.com/rookie-ninja/rk/cmd/rk
```

## Install tools
In order to generate swagger for gin API, we are using [swag](https://github.com/swaggo/swag) command line tool.

- Install swag
```
$ rk install swag
```

[swag](https://github.com/swaggo/swag) need API function to be commented in order to generate swagger config file. swagger files
will be generated at ./docs folder by default.
- Example
``` go
// @Summary Greeter service
// @Id 1
// @version 1.0
// @produce application/json
// @Param name query string true "Input name"
// @Success 200 {object} GreeterResponse
// @Router /v1/greeter [get]
func Greeter(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, &GreeterResponse{
		Message: fmt.Sprintf("Hello %s!", ctx.Query("name")),
	})
}
```

## Build locally
### With RK
Edit build.yaml file with customization.
```shell script
$ rk build
```
- build.yaml
```yaml
---
build:
  type: go                           # Optional, default:go
#  main: ""                          # Optional, default: ./main.go
#  GOOS: ""                          # Optional, default: current OS
#  GOARCH: ""                        # Optional, default: current Arch
#  args: ""                          # Optional, default: "", arguments which will attached to [go build] command
  copy: ["docs"]                     # Optional, default: [], directories or files need to copy to [target] folder
  commands:
    before:
      - "swag init --parseInternal --propertyStrategy camelcase"
#   after: []                        # Optional, default: [], commands would be invoked after [go build] command locally
#  scripts:
#    before: []                      # Optional, default: [], scripts would be executed before [go build] command locally
#    after: []                       # Optional, default: [], scripts would be executed after [go build] command locally
docker:
#  build:
#    registry: ""                    # Optional, default: [package name]
#    tag: ""                         # Optional, default: [current git tag or branch-latestCommit]
#    args: [""]                      # Optional, default: "", docker args which will be attached to [docker build] command
  run:
    args: ["-p", "8080:8080"]        # Optional, default: "", docker args which will be attached to [docker run] command
```

./target folder will be generated with compiled binary file.
```shell script
└── target
    ├── bin
    │ └── rk-demo
    ├── boot.yaml
    └── docs
        ├── docs.go
        ├── swagger.json
        └── swagger.yaml
```

### With go build
```shell script
$ go build main.go
```

## Run locally
### With RK
```shell script
$ rk run
```
- RK tv: http://localhost:8080/rk/v1/tv/
- Swagger: http://localhost:8080/sw
- Prometheus metrics: http://localhost:8080/metrics
- Logs: stdout or target/logs folder

### With go run
```shell script
$ go run main.go
```

### With IDE
Write click main.go file and run it!

## Pack locally
### With RK
rk will build project and create a [packageName-version.tar.gz] file by compressing target folder generated.
```shell script
$ rk pack
```
```shell script
├── rk-demo-master-b233d43.tar.gz
```

## Docker build
### With RK
Docker file needs to be exist.
- Dockerfile
```dockerfile
FROM alpine

ENV WD=/usr/service/

RUN mkdir -p $WD
WORKDIR $WD
COPY target/ $WD

# By default, rk command will build main.go to suffix of module name in go.mod file
CMD ["bin/demo"]
```
Please add docker build args as needed in build.yaml file.
- build.yaml
```yaml
docker:
  build:
    registry: ""                    # Optional, default: [package name]
    tag: ""                         # Optional, default: [current git tag or branch-latestCommit]
    args: [""]                      # Optional, default: "", docker args which will be attached to [docker build] command
 
```
```shell script
$ rk docker build
```
```shell script
$ docker images
  rk-demo  master-b233d43   2a76b06fd150   16 seconds ago   45.7MB
```

### With docker build
[docker build](https://docs.docker.com/engine/reference/commandline/build/)

## Docker run
### With RK
Docker file needs to be exist.
- Dockerfile
```dockerfile
FROM alpine

ENV WD=/usr/service/

RUN mkdir -p $WD
WORKDIR $WD
COPY target/ $WD

# By default, rk command will build main.go to suffix of module name in go.mod file
CMD ["bin/demo"]
```

Please add docker run args as needed in build.yaml file.
- build.yaml
```yaml
docker:
  run:
    args: ["-p", "8080:8080"]        # Optional, default: "", docker args which will be attached to [docker run] command
```
```shell script
$ rk docker run
```

## Clear locally
Delete target/ folder.
```shell script
$ rk clear
```

## boot.yaml
Since we are using gin framework, please refer [rk-gin](https://github.com/rookie-ninja/rk-gin) for details.

## build.yaml
Since we are using rk cmd, please refer [rk](https://github.com/rookie-ninja/rk) for details.
