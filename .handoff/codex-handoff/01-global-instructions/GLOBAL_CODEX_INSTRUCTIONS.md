# Global Codex Instructions — Revision 2.1

You are seeding and implementing **phase-1 SourceOS repositories** plus a **patch-only integration** to the existing ProphetCLI/ProfitCLI repository.

## Constitutional execution model
1. The deterministic CLI surface is canonical.
2. Agent surfaces are overlays that call deterministic tools; they must not invent parallel business logic.
3. The existing ProphetCLI repo is a **façade** only.
4. The bootstrap engine source home is **`sourceos-sdk/cmd/sourceos-bootstrap`**.
5. `homebrew-sourceos` is **distribution only**.

## Hard rules
1. Do not invent semantics beyond frozen specs.
2. Frozen objects in this revision are:
   - ExperienceProfile.v1
   - ReleaseSet.v1
   - BootReleaseSet.v1
   - Fingerprint.v1
   - ConfigSource.v0
   - TokenDoor.v0
   - PolicyBundle.v0
3. Remaining draft objects under `07-unfrozen-stubs/` are documentation placeholders only.
4. Phase-1 default Git update flow is **pull-based** with a **read-only token door**.
5. Homebrew is **Formula-first**; the Cask remains a scaffold only.
6. Never place secret material into Git, generated schemas, examples, or logs.
7. Do not put bootstrap implementation logic directly into the existing ProphetCLI repo.
8. Do not auto-load project MCP servers, hooks, or memory in untrusted workspaces.
9. Do not implement final boot enrollment, DeviceClaim, or PAL-Mac semantics in this revision.
10. All agent-executed mutations must remain explainable as deterministic command/tool invocations.

## CLI ergonomics policy
Use these canonical deterministic verbs:
- read side: `list`, `show`, `describe`, `explain`, `find`, `search`
- mutation: `create`, `update`, `delete`, `apply`, `wait`, `validate`, `verify`
- runtime: `status`, `logs`, `start`, `stop`, `restart`, `enable`, `disable`
- host bootstrap: `doctor`, `login`, `build`, `fetch`, `write`, `info`, `audit`

Use these global flags where applicable:
- `--profile`
- `--space`
- `--output` / `-o`
- `--query`
- `--quiet`
- `--debug`
- `--no-pager`

## Agent hybrid policy
1. Add agent affordances without replacing the deterministic surface.
2. Preserve three modes conceptually:
   - deterministic CLI
   - agent assist (`ask`, `plan`)
   - agent execute (`agent run`, approval-gated)
3. Treat MCP as an extension spine for tools/resources/prompts/roots, but do not auto-enable project MCP in untrusted workspaces.
4. Treat workspace memory as context only, not policy.
5. Bind roots, trust, approvals, and receipts together; planner agents must not directly mutate boot/system state.

## What Codex may implement now
- repo skeletons and READMEs
- frozen specs/schemas/examples/tests
- schema validation CI
- workspace bootstrap and lock file maintenance
- pull-based GitRefBuild scaffolds
- bootstrap CLI commands for doctor/login/build/validate/verify
- façade docs and delegation wrappers in the existing ProphetCLI repo
- agent runtime boundary docs, MCP wiring docs, and command-surface docs

## What Codex must not implement yet
- direct-push Git mutation flows
- final enrollment-token semantics
- final boot provisioning wire protocol
- PAL-Mac internals beyond docs/spec placeholders
- implicit plugin installation
- untrusted workspace auto-approval or auto-loaded hooks
