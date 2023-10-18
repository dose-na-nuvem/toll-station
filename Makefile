.DEFAULT_GOAL := all

# Depende de algumas ferramentas
EXECUTABLES = ls yq docker
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "Falta $(exec) no AMBIENTE, verifique sua instalação ")))

CMD_DISCOVER_HTTP_ENDPOINT = yq ".server.http.endpoint" tests/config.yaml
HTTP_HOSTNAME := $$( $(CMD_DISCOVER_HTTP_ENDPOINT) )

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

.PHONY: k6-http
k6-http:
	@echo "Lembre-se de iniciar o serviço antes ;-)"
	@docker run -it --rm --network host -v `pwd`:/app -w /app grafana/k6 run -e HOSTNAME=$(HTTP_HOSTNAME) -i 10 -u 10 tests/k6-http/script.js

.PHONY: generate-k6-http
generate-k6-http:
# CUIDADO: irá sobrescrever quaisquer scripts na pasta de destino
	@docker run -u 1000 --rm -v .:/workdir -w /workdir openapitools/openapi-generator-cli generate -g k6 -i openapiv2/traffic/v1/traffic.swagger.json -o local/k6-http

.PHONY: install-tools
install-tools:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.2

