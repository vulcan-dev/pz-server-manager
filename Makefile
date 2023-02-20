# Makefile for the pz-server-manager

GO_DIR ?= $(shell pwd)
GO_PKG ?= $(shell go list -e -f "{{ .ImportPath }}")
RM ?= rm -rf

GOOS ?= $(shell go env GOOS || echo linux)
GOARCH ?= $(shell go env GOARCH || echo amd64)
CGO_ENABLED ?= 0

GO_BUILD_FLAGS ?= -v -ldflags "-s -w"
ifeq ($(OS),Windows_NT)
GO_BUILD_ENV ?= $$Env:GOOS="$(GOOS)"; $$Env:GOARCH="$(GOARCH)"; $$Env:CGO_ENABLED=$(CGO_ENABLED);
else
GO_BUILD_ENV ?= GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED)
endif

.PHONY: all
all: build

.PHONE: init
init:
	$(info Initializing...)
	@go mod tidy

.PHONY: deps
deps: init
	$(info Downloading dependencies...)
	@go get -u github.com/sirupsen/logrus
	@go get -u github.com/gin-gonic/gin
	@go get -u github.com/gin-contrib/cors
	@go get -u gopkg.in/ini.v1

.PHONY: build
build:
	@go mod verify
	$(info Building pz-server-manager...)
	@go build $(GO_BUILD_FLAGS) -o bin/pz-server-manager cmd/pz-server-manager/main.go

run: build
	go run $(GO_BUILD_FLAGS) cmd/pz-server-manager/main.go

test:
	$(info Testing...)
	@go test -v ./...

bench:
	$(info Benchmarking...)
	@go test -bench=. ./...

.PHONY: clean
clean:
	$(RM) bin