SHELL := /bin/bash

.PHONY: help fmt test vet tidy verify

help:
	@echo "Targets: fmt test vet tidy verify"

fmt:
	gofmt -w ./cmd ./internal

test:
	go test ./...

vet:
	go vet ./...

tidy:
	go mod tidy

verify: fmt vet test
