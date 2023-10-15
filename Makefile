.DEFAULT_GOAL := all

CMD_DISCOVER_HTTP_ENDPOINT = yq ".server.http.endpoint" local/config.yaml
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

.PHONY: simule
simule: k6-http

.PHONY: k6-http
k6-http:
	@docker run -it --rm --network host -v `pwd`:/app -w /app grafana/k6 run -e HOSTNAME=$(HTTP_HOSTNAME) -i 50 -u 50 tests/k6-http/script.js

.PHONY: gen-k6-http
gen-k6-http:
# CUIDADO: irá sobrescrever quaisquer scripts na pasta de destino
	@docker run -u 1000 --rm -v .:/workdir -w /workdir openapitools/openapi-generator-cli generate -g k6 -i openapiv2/traffic/v1/traffic.swagger.json -o local/k6-http

.PHONY: install-tools
install-tools:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.2

.PHONY: trafego-http
trafego-http:
	@curl -X POST -H "Content-type: application/json" -d '{"tag":"abcdefg1234567890"}' http://$(HTTP_HOSTNAME)
