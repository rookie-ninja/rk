# Demo
**[English](README.md)**

Golang 微服务标准包，使用 rk-boot 来启动 gin 服务。

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [前提条件](#%E5%89%8D%E6%8F%90%E6%9D%A1%E4%BB%B6)
- [安装工具](#%E5%AE%89%E8%A3%85%E5%B7%A5%E5%85%B7)
- [本地编译](#%E6%9C%AC%E5%9C%B0%E7%BC%96%E8%AF%91)
  - [通过 RK](#%E9%80%9A%E8%BF%87-rk)
  - [通过 go build 命令](#%E9%80%9A%E8%BF%87-go-build-%E5%91%BD%E4%BB%A4)
- [本地运行](#%E6%9C%AC%E5%9C%B0%E8%BF%90%E8%A1%8C)
  - [通过 RK](#%E9%80%9A%E8%BF%87-rk-1)
  - [通过 go run](#%E9%80%9A%E8%BF%87-go-run)
  - [通过 IDE](#%E9%80%9A%E8%BF%87-ide)
- [本地打包](#%E6%9C%AC%E5%9C%B0%E6%89%93%E5%8C%85)
  - [通过 RK](#%E9%80%9A%E8%BF%87-rk-2)
- [Docker 编译](#docker-%E7%BC%96%E8%AF%91)
  - [通过 RK](#%E9%80%9A%E8%BF%87-rk-3)
  - [通过 docker build](#%E9%80%9A%E8%BF%87-docker-build)
- [本地运行 Docker](#%E6%9C%AC%E5%9C%B0%E8%BF%90%E8%A1%8C-docker)
  - [通过 RK](#%E9%80%9A%E8%BF%87-rk-4)
- [清理本地缓存](#%E6%B8%85%E7%90%86%E6%9C%AC%E5%9C%B0%E7%BC%93%E5%AD%98)
- [boot.yaml](#bootyaml)
- [build.yaml](#buildyaml)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 前提条件
推荐下载 [rk](https://github.com/rookie-ninja/rk) 命令行工具来配置 gin 所需要的环境。

```shell script
go get -u github.com/rookie-ninja/rk/cmd/rk
```

## 安装工具
为了生成 swagger，我们需要安装 swag 工具。

- 安装
```
$ rk install swag
```

[swag](https://github.com/swaggo/swag) 通过代码中的注视来生成 Swagger 参数文件，并且把文件生成到 docs/ 文件夹中。
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

## 本地编译
### 通过 RK
请编辑 build.yaml 文件来修改默认编译过程。
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
  copy: ["api"]                      # Optional, default: [], directories or files need to copy to [target] folder
  commands:
    before:
      - "buf generate --path api/v1"
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

./target 文件夹将会被创建，里面包含了编译好的二进制可运行文件。
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

### 通过 go build 命令
```shell script
$ go build main.go
```

## 本地运行
### 通过 RK
```shell script
$ rk run
```
- RK tv: http://localhost:8080/rk/v1/tv/
- Swagger: http://localhost:8080/sw
- Prometheus metrics: http://localhost:8080/metrics
- 日志: stdout 或者 target/logs 文件夹

### 通过 go run
```shell script
$ go run main.go
```

### 通过 IDE
右键点击 main.go 文件来运行。

## 本地打包
### 通过 RK
rk 命令行会编译并且创建 tag.gz 格式的压缩包，命名格式：[packageName-version.tar.gz]
```shell script
$ rk pack
```
```shell script
├── rk-demo-master-b233d43.tar.gz
```

## Docker 编译
### 通过 RK
Dockerfile 是必须的前提条件。
- Dockerfile
```dockerfile
FROM alpine

ENV WD=/usr/service/

RUN mkdir -p $WD
WORKDIR $WD
COPY target/ $WD

# 默认情况下，rk 命令行会在 go.mod 文件中寻找 module，并且以 module 的名字作为 main.go 的编译输出。
CMD ["bin/demo"]
```
请在 args 中添加 docker 命令行所需要的参数。
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

### 通过 docker build
[docker build](https://docs.docker.com/engine/reference/commandline/build/)

## 本地运行 Docker
### 通过 RK
Dockerfile 是必须的前提条件。
- Dockerfile
```dockerfile
FROM alpine

ENV WD=/usr/service/

RUN mkdir -p $WD
WORKDIR $WD
COPY target/ $WD

# 默认情况下，rk 命令行会在 go.mod 文件中寻找 module，并且以 module 的名字作为 main.go 的编译输出。
CMD ["bin/demo"]
```

请在 args 中添加 docker 命令行所需要的参数。
- build.yaml
```yaml
docker:
  run:
    args: ["-p", "8080:8080"]        # Optional, default: "", docker args which will be attached to [docker run] command
```
```shell script
$ rk docker run
```

## 清理本地缓存
rk 命令行将会删除 target/ 文件夹。
```shell script
$ rk clear
```

## boot.yaml
我们采用 gin 作为底层框架，请参考：[rk-gin](https://github.com/rookie-ninja/rk-gin)

## build.yaml
我们采用 rk 命令行进行编译，请参考：[rk](https://github.com/rookie-ninja/rk)
