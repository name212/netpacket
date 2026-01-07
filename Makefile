# Copyright 2026
# license that can be found in the LICENSE file.

GOLANGCI_VERSION = 2.7.2
GOFUMPT_VERSION=0.9.2

PLATFORM_NAME := $(shell uname -m)

OS_NAME := $(shell uname)
ifndef OS
	ifeq ($(UNAME), Linux)
		OS = linux
	else ifeq ($(UNAME), Darwin)
		OS = darwin
	endif
endif

ifeq ($(OS_NAME), Linux)
	GOFUMPT_PLATFORM = linux
else ifeq ($(OS_NAME), Darwin)
	GOFUMPT_PLATFORM = darwin
endif

ifeq ($(PLATFORM_NAME), x86_64)
	GOFUMPT_ARCH = amd64
else ifeq ($(PLATFORM_NAME), arm64)
	GOFUMPT_ARCH = arm64
endif

.PHONY: bin/gofumpt bin/golangci-lint clean curl-installed go-installed

bin:
	mkdir -p bin

curl-installed:
	command -v curl > /dev/null

go-installed:
	command -v go > /dev/null
	go version

bin/golangci-lint: curl-installed bin
	if ! ./hack/check_binary.sh "golangci-lint" "--version" "$(GOLANGCI_VERSION)"; then \
  	  echo "Install golangci-lint"; \
	  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | BINARY=golangci-lint bash -s -- v${GOLANGCI_VERSION}; \
	  chmod +x "./bin/golangci-lint"; \
	fi

bin/gofumpt: curl-installed bin
	if ! ./hack/check_binary.sh "gofumpt" "-version" "$(GOFUMPT_VERSION)"; then \
  	  echo "Install gofumpt"; \
	  curl -sSfLo "bin/gofumpt" https://github.com/mvdan/gofumpt/releases/download/v$(GOFUMPT_VERSION)/gofumpt_v$(GOFUMPT_VERSION)_$(GOFUMPT_PLATFORM)_$(GOFUMPT_ARCH); \
	  chmod +x "./bin/gofumpt"; \
	fi

deps: bin bin/golangci-lint bin/gofumpt

lint: bin/golangci-lint
	./bin/golangci-lint run ./... -c .golangci.yaml

lint/fix: bin/golangci-lint
	./bin/golangci-lint run ./... -c .golangci.yaml --fix

fmt: bin/gofumpt
	 ./bin/gofumpt .


test: go-installed
	cd tests; go test -v -p 1 ./...

all: bin deps fmt lint test

clean:
	rm -rf ./bin