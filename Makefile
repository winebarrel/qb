SHELL   := /bin/bash
VERSION := v0.1.0
GOOS    := $(shell go env GOOS)
GOARCH  := $(shell go env GOARCH)

.PHONY: all
all: vet build

.PHONY: build
build:
	go build -ldflags "-X main.version=$(VERSION)" ./cmd/qb

.PHONY: vet
vet:
	go vet

.PHONY: package
package: clean vet build
ifeq ($(GOOS),windows)
	zip qb_$(VERSION)_$(GOOS)_$(GOARCH).zip qb.exe
	sha1sum qb_$(VERSION)_$(GOOS)_$(GOARCH).zip > qb_$(VERSION)_$(GOOS)_$(GOARCH).zip.sha1sum
else
	gzip qb -c > qb_$(VERSION)_$(GOOS)_$(GOARCH).gz
	sha1sum qb_$(VERSION)_$(GOOS)_$(GOARCH).gz > qb_$(VERSION)_$(GOOS)_$(GOARCH).gz.sha1sum
endif

.PHONY: clean
clean:
	rm -f qb qb.exe
