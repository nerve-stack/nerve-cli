.DEFAULT_GOAL := test

.PHONY: fmt
fmt:
	@gofumpt -w -l .
	@goimports -w -l .
	@golangci-lint run --fix

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: test
test:
	@gotestsum -f testdox

.PHONY: run
run:
	go run ./cmd/nerve codegen server ./examples/example.yml -l go

.PHONY: build
build:
	go build -o ~/.cache/bin/nerve ./cmd/nerve

