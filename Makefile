.PHONY: all
all: fmt doctoc test

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
	@go build -o ${GOPATH}/bin/rk cmd/rk/rk.go
	@echo "------------------------------------[Done]"

.PHONY: pkger
pkger:
	@echo "[pkger] Running pkger..."
	@pkger -o commands/pkg
	@echo "------------------------------------[Done]"

.PHONY: test
test:
	@echo "[test] Running go test..."
	@go test ./... -coverprofile coverage.txt 2>&1
	@go tool cover -html=coverage.txt
	@echo "------------------------------------[Done]"

