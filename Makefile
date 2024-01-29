APP=go-frames-scores

.PHONY: build
build: install generate build_app

.PHONY: build_app
build_app: install
	cd cmd/${APP} && CGO_ENABLED=0 go build -o ../../bin/${APP} .

.PHONY: fmt
fmt:
	@gofmt -w .
	@goimports -w .

.PHONY: lint
lint:
	@command -v golangci-lint || go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.51.2 && golangci-lint run --timeout 5m

.PHONY: tidy
tidy:
	@go mod tidy
	@bash -c 'if [[ -n $$(git ls-files --other --exclude-standard --directory -- go.sum) ]]; then\
    	echo "go.sum was added by go mod tidy";\
    	exit 1;\
	fi'
	@git diff --exit-code -- go.sum go.mod

.PHONY: install
install:
	@go get ./...

.PHONY: tests
tests:
	@go test ./... -count=1

.PHONY: govulncheck
govulncheck:
	@if ! command -v ${HOME}/go/bin/govulncheck &> /dev/null; then \
		GO111MODULE=on go install golang.org/x/vuln/cmd/govulncheck@latest; \
	fi
	@${HOME}/go/bin/govulncheck ./...

.PHONY: generate
generate:
	@command -v stringer || go install golang.org/x/tools/cmd/stringer@latest
	@go generate $(shell go list ./...)
