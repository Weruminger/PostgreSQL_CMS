SHELL := /bin/bash
PKG := github.com/Weruminger/PostgreSQL_CMS/portal
BINARY := posgres-cms
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
BUILD   ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildDate=$(BUILD)

.PHONY: all test run build lint cover

all: test build

build:
	CGO_ENABLED=0 go build -ldflags '$(LDFLAGS)' -o bin/$(BINARY) ./cmd/go-ad-admin

test:
	go test ./... -cover -v

run: build
	./bin/$(BINARY)


# --- Go toolchain defaults ---
GO      ?= go
PKG     ?= ./...
BDDPKG  ?= ./internal/bdd

.PHONY: unit bdd bdd.v test cover tidy

tidy:
	$(GO) mod tidy

unit:
	$(GO) test $(PKG) -run '^(?!TestFeatures).*' -count=1

# nur die BDD-Suite
bdd:
	$(GO) test $(BDDPKG) -run TestFeatures -count=1

# BDD mit Verbose-Logs (t.Logf etc.)
bdd.v:
	$(GO) test -v $(BDDPKG) -run TestFeatures -count=1

# alles
test: tidy
	$(GO) test $(PKG) -count=1

# optional: Coverage aus allen Paketen
cover:
	$(GO) test $(PKG) -coverprofile=coverage.out -covermode=atomic
	$(GO) tool cover -func=coverage.out | tail -n 1



up:
	docker compose up -d --build
