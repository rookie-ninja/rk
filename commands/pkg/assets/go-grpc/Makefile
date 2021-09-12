.PHONY: all
all: fmt buf test doctoc lint

.PHONY: lint
lint:
	@echo "[golangci-lint] Running golangci-lint..."
	@golangci-lint run 2>&1
	@echo "------------------------------------[Done]"

.PHONY: fmt
fmt:
	@echo "[fmt] Formatting go project..."
	@gofmt -s -w . 2>&1
	@echo "------------------------------------[Done]"

.PHONY: test
test:
	@echo "[test] Running go test..."
	@go test $(go list ./... | grep -v vendor) -coverprofile coverage.txt -covermode=atomic 2>&1
	@echo "------------------------------------[Done]"

.PHONY: doctoc
doctoc:
	@echo "[doctoc] Running doctoc..."
	@doctoc . 2>&1
	@echo "------------------------------------[Done]"

.PHONY: buf
buf:
	@echo "[buf] Running buf..."
	@buf generate --path api/v1
	@echo "------------------------------------[Done]"

