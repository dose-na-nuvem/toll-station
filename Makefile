.DEFAULT_GOAL := all

.PHONY: all
all: lint test build

.PHONY: lint
lint:
	@golangci-lint --timeout 120s run ./...

.PHONY: test
test:
	@go test -v ./...

.PHONY: build
build:
	@go build -o ./_build/pedagio .

.PHONY: install-tools
install-tools:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.2

.PHONY: trafego-http
trafego-http:
	@curl -X POST -H "Content-type: application/json" -d '{"tag":"abcdefg1234567890"}' http://localhost:61000
