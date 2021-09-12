.PHONY: all
all: fmt doctoc

.PHONY: fmt
fmt:
	@echo "[fmt] Formatting go project..."
	@gofmt -s -w . 2>&1
	@echo "------------------------------------[Done]"

.PHONY: doctoc
doctoc:
	@echo "[doctoc] Running doctoc..."
	@doctoc . 2>&1
	@echo "------------------------------------[Done]"

.PHONY: build
build:
	@echo "[build] Building to local..."
	@go build -o $GOPATH/bin/rk cmd/rk/rk.go
	@echo "------------------------------------[Done]"

.PHONY: pkger
pkger:
	@echo "[pkger] Running pkger..."
	@pkger -o commands/pkg
	@echo "------------------------------------[Done]"

