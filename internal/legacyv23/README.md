# internal/legacyv23

This package namespace is reserved for **selective porting** from the unpublished Prophet CLI v2.3 attempt.

It exists so we can preserve lineage without compiling legacy code directly into the new façade by accident.

Rules:
- do not copy legacy files here verbatim and wire them into production builds without review
- first record the legacy file under `docs/legacy/prophet_cli_v2_3/`
- then port only the useful parts into canonical packages

Primary porting targets:
- `internal/cmd/`
- `internal/a2a/`
- `internal/config/`
- `internal/runtime/`
- `internal/transport/`
- `internal/receipts/`
