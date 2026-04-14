# Development

This repository is a phase-1 facade CLI.

## Local workflow

Typical commands:
- `make fmt`
- `make test`
- `make vet`
- `make tidy`
- `make verify`

## Scope discipline

Do in this repo:
- command grammar
- wrapper commands
- facade docs
- command-surface tests

Do not do in this repo:
- bootstrap engine business logic
- boot or enrollment protocol implementation
- ownership of ReleaseSet, BootReleaseSet, ConfigSource, or TokenDoor semantics

## Branching

Use small topic branches and keep PRs reviewable.
