SHELL   := /bin/bash
GOOS    := $(shell go env GOOS)
GOARCH  := $(shell go env GOARCH)

.PHONY: all
all: vet build

.PHONY: build
build:
	go build ./cmd/qb

.PHONY: vet
vet:
	go vet

.PHONY: clean
clean:
	rm -f qb qb.exe

.PHONY: lint
lint:
	golangci-lint run
