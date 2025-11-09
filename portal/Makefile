
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
SHELL := /bin/bash
PKG := github.com/Weruminger/PostgreSQL_CMS
BINARY := portal
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
BUILD   ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildDate=$(BUILD)
BUILD_DIR := ../work/build/$(VERSION)

.PHONY: all test run build lint cover

APP := portal

TEST_PKGS := ./...
COVER_PROFILE := ../work/coverage/coverage.out
COVER_HTML := ../work/coverage/coverage.html

test:
	@mkdir -p ../work/coverage
	@echo "==> Running unit tests"
	go test -count=1 $(TEST_PKGS)

cover:
	@mkdir -p ../work/coverage
	@echo "==> Running tests with coverage"
	go test -count=1 -covermode=atomic -coverprofile=$(COVER_PROFILE) $(TEST_PKGS)
	@echo "==> Coverage summary"
	go tool cover -func=$(COVER_PROFILE)

coverhtml: cover
	@echo "==> Generating HTML report"
	go tool cover -html=$(COVER_PROFILE) -o $(COVER_HTML)
	@echo "Open: $(COVER_HTML)"

build:
	@echo "==> Building $(APP) for $(GOOS)/$(GOARCH)"
	mkdir -p $(BUILD_DIR)/$(GOOS)-$(GOARCH)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags '$(LDFLAGS)' -o $(BUILD_DIR)/$(GOOS)-$(GOARCH)/$(APP) ./cmd/portal

clean:
	rm -rf ../work/build/*

release: clean
	$(MAKE) build GOOS=darwin GOARCH=arm64
	$(MAKE) build GOOS=linux GOARCH=amd64
	$(MAKE) build GOOS=windows GOARCH=amd64
