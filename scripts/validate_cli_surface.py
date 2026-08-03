#!/usr/bin/env python3
"""Reconcile the Prophet CLI spec command tree against CLI_SURFACE_POLICY.md.

Teeth:
  * `verbs` in manifests/cli-surface.json must equal the canonical verb model
    parsed from docs/CLI_SURFACE_POLICY.md -- drift fires *both ways*
    (a verb declared but not in policy, or in policy but not declared);
  * every spec noun resolves only to canonical verbs (no competing command
    language);
  * per the phase-1 policy, an empty subtree must not be scaffolded: an
    `adopted` noun needs a real delegate and >=1 verb; a `deferred` noun must
    claim no verbs and name a follow-up;
  * the whole spec command tree must be accounted for (no silently dropped noun).

Usage:
  validate_cli_surface.py               reconcile the manifest against policy
  validate_cli_surface.py --self-test   also prove the drift teeth fire both ways
"""

from __future__ import annotations

import json
import re
import sys
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parent.parent
MANIFEST_PATH = REPO_ROOT / "manifests" / "cli-surface.json"
POLICY_PATH = REPO_ROOT / "docs" / "CLI_SURFACE_POLICY.md"

# The spec command tree (Prophet CLI & System Architecture spec, intake 2026-07-31).
# The reconciliation must account for every one of these nouns.
SPEC_NOUNS = {
    "context", "config", "build", "fs", "graph", "docker", "conda", "minio",
    "rsync", "kafka", "work", "vocab", "gib", "search", "opencog", "genesis",
    "inception", "twin", "bridge", "agents",
}

# Policy section header -> manifest verb-category key.
SECTION_TO_CATEGORY = {
    "Read verbs": "read",
    "Mutation verbs": "mutation",
    "Runtime verbs": "runtime",
    "Bootstrap verbs": "bootstrap",
}


class SurfaceError(ValueError):
    """Raised when the CLI surface drifts from policy."""


def require(condition: bool, message: str) -> None:
    if not condition:
        raise SurfaceError(message)


def parse_policy(text: str) -> dict[str, list[str]]:
    """Extract {category: [verb, ...]} from CLI_SURFACE_POLICY.md."""
    verbs: dict[str, list[str]] = {}
    current: str | None = None
    for line in text.splitlines():
        header = re.match(r"^##\s+(.*?)\s*$", line)
        if header:
            current = SECTION_TO_CATEGORY.get(header.group(1).strip())
            if current is not None:
                verbs.setdefault(current, [])
            continue
        if current is None:
            continue
        bullet = re.match(r"^-\s+`([a-z][a-z0-9-]*)`\s*$", line)
        if bullet:
            verbs[current].append(bullet.group(1))
    return verbs


def reconcile(manifest: dict, policy_verbs: dict[str, list[str]]) -> None:
    require(manifest.get("apiVersion") == "cli.socioprophet.dev/v1", "apiVersion mismatch")
    require(manifest.get("kind") == "ProphetCliSurfaceReconciliation", "kind mismatch")

    declared = manifest.get("verbs")
    require(isinstance(declared, dict), "manifest.verbs must be an object")

    # Verb-model drift, both ways, per category.
    require(
        set(declared) == set(policy_verbs),
        f"verb categories drift: manifest={sorted(declared)} policy={sorted(policy_verbs)}",
    )
    for category, pol in policy_verbs.items():
        man = set(declared.get(category, []))
        pol_set = set(pol)
        not_in_policy = man - pol_set
        not_declared = pol_set - man
        require(not not_in_policy, f"{category}: verbs declared but not in policy: {sorted(not_in_policy)}")
        require(not not_declared, f"{category}: verbs in policy but not declared: {sorted(not_declared)}")

    all_policy_verbs = {v for vs in policy_verbs.values() for v in vs}

    # Noun teeth.
    nouns = manifest.get("nouns")
    require(isinstance(nouns, list) and nouns, "manifest.nouns must be a non-empty list")
    seen: set[str] = set()
    for entry in nouns:
        noun = entry.get("noun")
        require(isinstance(noun, str) and noun, f"noun entry missing name: {entry}")
        seen.add(noun)
        status = entry.get("status")
        require(status in {"adopted", "deferred"}, f"noun {noun}: status must be adopted|deferred, got {status!r}")
        nverbs = entry.get("verbs", [])
        require(isinstance(nverbs, list), f"noun {noun}: verbs must be a list")

        # No competing command language: every verb must be canonical.
        stray = set(nverbs) - all_policy_verbs
        require(not stray, f"noun {noun}: non-canonical verbs (competing command language): {sorted(stray)}")

        if status == "adopted":
            delegate = entry.get("delegate")
            require(isinstance(delegate, str) and delegate.strip(), f"noun {noun}: adopted noun needs a delegate")
            require(nverbs, f"noun {noun}: adopted noun must resolve to >=1 canonical verb")
        else:  # deferred
            require(not nverbs, f"noun {noun}: deferred noun must not scaffold verbs (phase-1 policy)")
            require(entry.get("followup"), f"noun {noun}: deferred noun must name a follow-up")

    # Whole spec command tree accounted for, both ways.
    missing = SPEC_NOUNS - seen
    extra = seen - SPEC_NOUNS
    require(not missing, f"spec nouns not reconciled: {sorted(missing)}")
    require(not extra, f"manifest declares nouns outside the spec command tree: {sorted(extra)}")


def validate() -> None:
    manifest = json.loads(MANIFEST_PATH.read_text(encoding="utf-8"))
    policy_verbs = parse_policy(POLICY_PATH.read_text(encoding="utf-8"))
    require(policy_verbs, "no verbs parsed from policy -- policy file changed shape?")
    reconcile(manifest, policy_verbs)


# --- self-test: prove the drift teeth fire both ways --------------------------

def self_test() -> int:
    import copy

    manifest = json.loads(MANIFEST_PATH.read_text(encoding="utf-8"))
    policy_verbs = parse_policy(POLICY_PATH.read_text(encoding="utf-8"))
    failures: list[str] = []

    def expect_reject(mut: dict, label: str, needle: str) -> None:
        try:
            reconcile(mut, policy_verbs)
            failures.append(f"{label}: should be REJECTED but reconciled (negative control breached)")
        except SurfaceError as exc:
            if needle in str(exc):
                print(f"PASS reject {label}: {exc}")
            else:
                failures.append(f"{label}: fired wrong reason: expected ~{needle!r}, got: {exc}")

    # Accept: the real manifest must reconcile against the real policy.
    try:
        reconcile(manifest, policy_verbs)
        print("PASS accept manifests/cli-surface.json")
    except SurfaceError as exc:
        failures.append(f"accept: real manifest should reconcile but was rejected: {exc}")

    # (a) command present but not in policy.
    m = copy.deepcopy(manifest)
    m["verbs"]["read"] = m["verbs"]["read"] + ["teleport"]
    expect_reject(m, "extra-verb-not-in-policy", "declared but not in policy")

    # (b) in policy but not declared.
    m = copy.deepcopy(manifest)
    m["verbs"]["read"] = [v for v in m["verbs"]["read"] if v != "list"]
    expect_reject(m, "policy-verb-not-declared", "in policy but not declared")

    # (c) noun uses a competing (non-canonical) verb.
    m = copy.deepcopy(manifest)
    m["nouns"][0] = dict(m["nouns"][0], verbs=["frobnicate"])
    expect_reject(m, "competing-command-language", "competing command language")

    # (d) deferred noun scaffolds an empty subtree.
    m = copy.deepcopy(manifest)
    for i, n in enumerate(m["nouns"]):
        if n["status"] == "deferred":
            m["nouns"][i] = dict(n, verbs=["list"])
            break
    expect_reject(m, "scaffolded-deferred-noun", "must not scaffold")

    if failures:
        for f in failures:
            print(f"FAIL {f}", file=sys.stderr)
        return 1
    print("PASS all CLI-surface reconciliation teeth verified")
    return 0


def main(argv: list[str]) -> int:
    if "--self-test" in argv:
        return self_test()
    try:
        validate()
        print("PASS manifests/cli-surface.json reconciled with docs/CLI_SURFACE_POLICY.md")
        return 0
    except Exception as exc:  # noqa: BLE001
        print(f"FAIL cli-surface reconciliation: {exc}", file=sys.stderr)
        return 1


if __name__ == "__main__":
    raise SystemExit(main(sys.argv))
