# Phase 1 Sequence — Revision 2.1

1. Create repos:
   - sourceos-workspace
   - sourceos-control-plane
   - sourceos-profile-catalog
   - sourceos-nix-lifecycle
   - sourceos-policy-ledger
   - sourceos-sdk
   - sourceos-boot
   - homebrew-sourceos
2. Patch the **existing ProphetCLI/ProfitCLI repo** as a façade target; do not recreate it.
3. Seed repo skeletons and CI.
4. Copy frozen object set into owning repos.
5. Validate all JSON Schemas and examples in CI.
6. Implement pull-based Git lifecycle scaffolds using:
   - ConfigSource.v0
   - TokenDoor.v0
   - PolicyBundle.v0
7. Wire bootstrap CLI validation and verification commands in `sourceos-sdk/cmd/sourceos-bootstrap`.
8. Add façade commands and delegation docs in the existing ProphetCLI repo.
9. Add agent-hybrid docs, MCP boundary docs, workspace trust docs, and command grammar docs.
10. Import milestones and tasks.
11. Stop and request review before boot-provisioning or enrollment-token implementation.
