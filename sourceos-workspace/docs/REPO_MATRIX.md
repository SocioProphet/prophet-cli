# Repo Matrix (Phase 1)

| Repo | Role | Owned frozen objects |
|---|---|---|
| sourceos-workspace | Coordination and trust docs | None |
| sourceos-sdk | Bootstrap engine source home (`cmd/sourceos-bootstrap`) | None (references frozen objects) |
| sourceos-nix-lifecycle | Pull-based lifecycle scaffolding | ReleaseSet.v1, ConfigSource.v0, TokenDoor.v0 |
| sourceos-policy-ledger | Policy and ledger scaffolding | Fingerprint.v1, PolicyBundle.v0 |
| sourceos-boot | Boot channel/release scaffolding | BootReleaseSet.v1 |
| sourceos-profile-catalog | Profile catalog and bundles | ExperienceProfile.v1 |
| sourceos-control-plane | Control-plane object index/scaffold | Linked object model |
| homebrew-sourceos | Distribution/tap only | None |
| prophet-cli-existing | Façade-only patch target | None |
