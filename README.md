# rk-cmd
Rookie command line tools. It contains lots of useful utility tools including downloading dependencies, generate Go files
from proto files etc.

rk lets you:
- Install third party command lint tools
- Generate Go files, gRpc go files and swagger json files

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

  - [Installation](#installation)
  - [Quick Start](#quick-start)
  - [Command Overview](#command-overview)
    - [rk help](#rk-help)
    - [rk install](#rk-install)
    - [rk gen](#rk-gen)
  - [Development Status: Stable](#development-status-stable)
  - [Contributing](#contributing)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Installation
```shell script
go install github.com/rookie-ninja/rk
```

## Quick Start
We'll start with a general overview of the commands. 
There are more commands, and we will get into] usage below, but this shows the basic functionality.

```shell script
rk help
rk install # Install third party command line tools
rk gen # Generate files from proto
```

## Command Overview

### rk help
```shell script
rk help
```
Print help message

### rk install
```shell script
rk install 
```

Subcommands
```shell script
COMMANDS:
   protobuf                 install protobuf on local machine
   protoc-gen-go            install protoc-gen-go on local machine
   protoc-gen-grpc-gateway  install protoc-gen-grpc-gateway on local machine
   protoc-gen-swagger       install protoc-gen-swagger on local machine
   golint                   install golint on local machine
   golangci-lint            install golangci-lint on local machine
   gocov                    install gocov on local machine
   mockgen                  install mockgen on local machine
   rk-std                   install rk standard environment on local machine
   help, h                  Shows a list of commands or help for one command
```

### rk gen
```shell script
rk gen # Generate files from proto
```

Subcommands
```shell script
COMMANDS:
   pb-go    generate go file from proto
   help, h  Shows a list of commands or help for one command
```

## Development Status: Stable

## Contributing
We encourage and support an active, healthy community of contributors &mdash;
including you! Details are in the [contribution guide](CONTRIBUTING.md) and
the [code of conduct](CODE_OF_CONDUCT.md). The rk maintainers keep an eye on
issues and pull requests, but you can also report any negative conduct to
dongxuny@gmail.com. That email list is a private, safe space; even the zap
maintainers don't have access, so don't hesitate to hold us to a high
standard.

<hr>

Released under the [MIT License](LICENSE).
