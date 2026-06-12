SHELL := /bin/bash

BIN := prophet
DIST_DIR := dist
VERSION ?= 0.1.0-dev
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GOOS ?= $(shell go env GOOS 2>/dev/null || uname -s | tr A-Z a-z)
GOARCH ?= $(shell go env GOARCH 2>/dev/null || uname -m)
DIST_NAME := $(BIN)_$(VERSION)_$(GOOS)_$(GOARCH)
LDFLAGS := -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

.PHONY: help fmt build test vet tidy validate verify dist release-dry-run clean

help:
	@echo "Targets: fmt build test vet tidy validate verify dist release-dry-run clean"

fmt:
	gofmt -w ./cmd ./internal

build:
	mkdir -p bin
	go build -ldflags "$(LDFLAGS)" -o bin/$(BIN) ./cmd/prophet

test:
	go test ./...

vet:
	go vet ./...

tidy:
	go mod tidy

validate: build vet test
	bin/$(BIN) --help >/tmp/prophet-help.txt
	bin/$(BIN) version >/tmp/prophet-version.json
	bin/$(BIN) doctor >/tmp/prophet-doctor.json
	bin/$(BIN) self-test >/tmp/prophet-self-test.json
	bin/$(BIN) emit-evidence >/tmp/prophet-evidence.json
	bin/$(BIN) devtools profile list >/tmp/prophet-devtools-profiles.json
	bin/$(BIN) lab list >/tmp/prophet-labs.json
	bin/$(BIN) sourceos carry list >/tmp/prophet-sourceos-carry-list.json
	bin/$(BIN) holmes search "truth and evidence" >/tmp/prophet-holmes-search.json
	bin/$(BIN) model route --task summarize --privacy local-first >/tmp/prophet-model-route.json
	bin/$(BIN) guardrail test examples/policy.json examples/input.json >/tmp/prophet-guardrail-test.json
	bin/$(BIN) ledger validate >/tmp/prophet-ledger-validate.json
	bin/$(BIN) ledger records >/tmp/prophet-ledger-records.json
	bin/$(BIN) agent registry list >/tmp/prophet-agent-registry-list.json
	bin/$(BIN) spine list >/tmp/prophet-spine-list.json
	bin/$(BIN) spine gate --help >/tmp/prophet-spine-gate-help.txt
	bin/$(BIN) enrichment lifecycle >/tmp/prophet-enrichment-lifecycle.json
	bin/$(BIN) enrichment status >/tmp/prophet-enrichment-status.json
	bin/$(BIN) enrichment run --seed seed:enrichment/photo-v1 --dry-run >/tmp/prophet-enrichment-dry-run.json

verify: fmt validate

dist: validate
	mkdir -p $(DIST_DIR)
	cp bin/$(BIN) $(DIST_DIR)/$(DIST_NAME)
	(cd $(DIST_DIR) && sha256sum $(DIST_NAME) > $(DIST_NAME).sha256)

release-dry-run: dist
	@echo "release dry-run complete: $(DIST_DIR)/$(DIST_NAME)"

clean:
	rm -rf bin $(DIST_DIR)
