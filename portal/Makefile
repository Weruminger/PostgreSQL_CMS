VERSION ?= v0.1.0
BUILD_DIR := ../work/build/$(VERSION)
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

APP := portal
TEST_PKGS := ./...
COVER_DIR := ../work/coverage
COVER_PROFILE := $(COVER_DIR)/coverage.out
COVER_HTML := $(COVER_DIR)/coverage.html
SQL_DIR := ../db
TEST_DB_URL := postgres://cms_owner:changeme@localhost:55432/cms?sslmode=disable

.PHONY: build clean release test cover coverhtml test-db-up test-db-wait test-db-migrate test-db-down itest

build:
	@echo "==> Building $(APP) for $(GOOS)/$(GOARCH)"
	mkdir -p $(BUILD_DIR)/$(GOOS)-$(GOARCH)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BUILD_DIR)/$(GOOS)-$(GOARCH)/$(APP) ./cmd/portal

clean:
	rm -rf ../work/build/* ../work/coverage/*

release: clean
	$(MAKE) build GOOS=darwin GOARCH=arm64
	$(MAKE) build GOOS=linux GOARCH=amd64
	$(MAKE) build GOOS=windows GOARCH=amd64

test:
	@echo "==> Unit tests"
	go test -count=1 $(TEST_PKGS)

cover:
	@mkdir -p $(COVER_DIR)
	@echo "==> Tests with coverage"
	go test -count=1 -covermode=atomic -coverprofile=$(COVER_PROFILE) $(TEST_PKGS)
	@echo "==> Coverage summary"
	go tool cover -func=$(COVER_PROFILE)

coverhtml: cover
	@echo "==> Generating HTML report"
	go tool cover -html=$(COVER_PROFILE) -o $(COVER_HTML)
	@echo "Open: $(COVER_HTML)"

# ---------- Integration tests (PostgreSQL container) ----------
test-db-up:
	docker compose -f ../docker-compose.test.yml up -d pgtest

test-db-wait:
	@echo "waiting for postgres on 55432..."
	@for i in $$(seq 1 60); do \
		PGPASSWORD=changeme psql -h localhost -p 55432 -U cms_owner -d cms -c 'select 1' >/dev/null 2>&1 && break; \
		sleep 1; \
	done

test-db-migrate:
	@ls -1 $(SQL_DIR)/*.sql | sort | while read f; do \
		echo ">> $$f"; \
		PGPASSWORD=changeme psql -h localhost -p 55432 -U cms_owner -d cms -v ON_ERROR_STOP=1 -f $$f; \
	done

test-db-down:
	docker compose -f ../docker-compose.test.yml down -v

itest: test-db-up test-db-wait test-db-migrate
	@echo "==> integration tests"
	PGHOST=localhost PGPORT=55432 PGDATABASE=cms PGUSER=cms_owner PGPASSWORD=changeme \
		go test -count=1 -tags=integration $(TEST_PKGS)
