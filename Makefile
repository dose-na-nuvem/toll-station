ENDPOINT:=http://localhost:18641/api/v1/traffic

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
	@go build -o ./_build/vehicles .

.PHONY: simule
simule:
	@echo  curl -X POST ${ENDPOINT} -H "Content-type: application/json;" -d '{"tag":"tag1"}' 
	@echo  curl -X POST ${ENDPOINT} -H "Content-type: application/json;" -d '{"tag":"tag2"}' 
	@echo  curl -X POST ${ENDPOINT} -H "Content-type: application/json;" -d '{"tag":"tag3"}' 

.PHONY: install-tools
install-tools:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.2

