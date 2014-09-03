export GOPATH = $(shell pwd | sed -e 's/\/src\/.*$$//g')

test: build
	go test -v
.PHONY: test

benchmark: build
	go test -v -bench .
.PHONY: benchmark

integration-test: build
	mysql -u root -e 'SELECT VERSION();'
	mysql -u root -e 'CREATE DATABASE test;'
	mysql -u root -D test < integration_test.sql
	go test -v --enable_integration_test
.PHONY: integration-test

build: get
	go build
.PHONY: build

get: version
	go get github.com/go-sql-driver/mysql
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
