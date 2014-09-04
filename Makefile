export GOPATH = $(shell pwd | sed -e 's/\/src\/.*$$//g')

test: build
	go test -v
.PHONY: test

benchmark: build
	go test -v -bench .
.PHONY: benchmark

coverage: build
	go get code.google.com/p/go.tools/cmd/cover
	go test -coverprofile=/tmp/coverage
	go tool cover -html=/tmp/coverage
.PHONY: coverage

build: get
	go build
.PHONY: build

get: version
	go get
.PHONY: get

version:
	@go version
.PHONY: version

format:
	gofmt -w ./
.PHONY: format

document:
	godoc -http=:6060

info:
	@echo "GOPATH=$${GOPATH}"
.PHONY: info
