# Contributing to prophet-cli

Thank you for your interest in contributing! This document explains the conventions and boundaries for this repository.

## What this repository is

`prophet-cli` is a **façade-only CLI** in phase 1. It owns the command surface and global flag grammar for the Prophet / SourceOS platform. It does **not** own bootstrap business logic, enrollment semantics, or SourceOS object definitions — those live in their respective upstream repositories.

## How to contribute

1. **Fork and branch** — create a feature branch from `main`.
2. **Keep changes small and reviewable** — prefer focused, single-purpose commits.
3. **Run the test suite** before opening a pull request:
   ```bash
   go mod tidy
   go test ./...
   go vet ./...
   ```
4. **Open a pull request** against `main` with a clear description of what the change does and why.

## What we welcome

- Improved command grammar and `--help` text
- Wrapper docs and examples in `docs/`
- Compileable scaffolds that demonstrate delegation patterns
- Additional tests and usage examples
- Clearer error messages and output formatting

## Boundaries — please do not

- Embed bootstrap business logic in this repo
- Invent boot or enrollment semantics here
- Add hidden side effects to wrapper commands
- Auto-enable project MCP servers in untrusted workspaces
- Bypass or extend the frozen object set from the SourceOS handoff materials
- Implement competing command-language semantics that shadow the deterministic verbs

## Code style

- Standard Go formatting (`gofmt`).
- Keep cobra `Use`, `Short`, `Long`, and `Example` fields consistent with the verb taxonomy in [docs/CLI_SURFACE_POLICY.md](docs/CLI_SURFACE_POLICY.md).
- New commands must emit structured output via the shared `emit` helper so that `--output json` always works.

## Documentation

- Update [docs/COMMANDS.md](docs/COMMANDS.md) when adding or changing commands.
- Update delegation docs if bootstrap wiring changes.
- Keep phase-1 scaffold status visible in command output until a command is fully wired.
