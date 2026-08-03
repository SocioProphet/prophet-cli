# Running Archive

Index of `prophet-cli` repo-local working memory. This is the umbrella the
[`ARCHIVE_MANIFEST.json`](../ARCHIVE_MANIFEST.json) tracks: it points at the live
current-state and continuity documents so a next agent has one entry point.

## Tracks

- [IMPLEMENTATION_STATUS.md](IMPLEMENTATION_STATUS.md) — current-state snapshot of the
  command surface: real vs delegating vs scaffold, delegate engines, known gaps.
- [SESSION_WORKLOG.md](SESSION_WORKLOG.md) — append-only per-session change/verification log.
- RUNNING_ARCHIVE.md — this index.

## Discipline

- **Current-state, not aspiration.** If the manifest declares a track, the file exists.
- **Verify before you commit.** Record the exact green-gate used (build/vet/test invocation).
- **Façade boundary is load-bearing.** This repo owns command grammar and docs only; engine
  logic, transport, and receipt/enrollment semantics live in the engine repos. See
  [REPO_ROLE.md](REPO_ROLE.md).
- **Runtime vs working memory.** This archive is in-repo working memory. It is **not** a
  runtime ledger and not the estate ProofArtifact spine.
