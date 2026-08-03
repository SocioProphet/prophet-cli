#!/usr/bin/env python3
"""Validate the Prophet CLI 3x3 build-recipe contract.

Teeth (mirrors scripts/validate_lattice_studio_commands.py):
  * tier must be one of the enum declared in manifests/build-recipes.schema.json
    (the schema is the single source of truth -- validator and schema cannot drift);
  * every stage declares depends_on + a *real* health gate (required==true,
    non-empty check) -- a weaker gate is treated as missing;
  * the stage ordering is a DAG: no cycles, no forward edges (a stage may only
    depend on earlier stages) and no skip edges (a stage must depend on its
    immediate predecessor), so stage N+1 cannot start until stage N verifies.

Usage:
  validate_build_recipes.py [PATH]        validate one recipe file (default build/recipes.yaml)
  validate_build_recipes.py --self-test   run the accept fixture + every reject fixture
                                          and assert the teeth fire both ways
"""

from __future__ import annotations

import sys
from pathlib import Path

import yaml

REPO_ROOT = Path(__file__).resolve().parent.parent
SCHEMA_PATH = REPO_ROOT / "manifests" / "build-recipes.schema.json"
RECIPES_PATH = REPO_ROOT / "build" / "recipes.yaml"
FIXTURES_DIR = REPO_ROOT / "build" / "fixtures"


class RecipeError(ValueError):
    """Raised when a recipe document violates the contract."""


def require(condition: bool, message: str) -> None:
    if not condition:
        raise RecipeError(message)


def _load_schema() -> dict:
    import json

    schema = json.loads(SCHEMA_PATH.read_text(encoding="utf-8"))
    return schema


def _allowed_tiers(schema: dict) -> set[str]:
    tiers = schema["$defs"]["recipe"]["properties"]["tier"]["enum"]
    return set(tiers)


def _stage_required_keys(schema: dict) -> set[str]:
    return set(schema["$defs"]["stage"]["required"])


def _detect_cycle(nodes: list[str], edges: dict[str, list[str]]) -> list[str] | None:
    """Return a cycle path if the depends_on graph has one, else None."""
    WHITE, GRAY, BLACK = 0, 1, 2
    color = {n: WHITE for n in nodes}
    stack: list[str] = []

    def visit(n: str) -> list[str] | None:
        color[n] = GRAY
        stack.append(n)
        for m in edges.get(n, []):
            if m not in color:
                continue  # dangling ref handled elsewhere
            if color[m] == GRAY:
                return stack[stack.index(m):] + [m]
            if color[m] == WHITE:
                found = visit(m)
                if found:
                    return found
        stack.pop()
        color[n] = BLACK
        return None

    for n in nodes:
        if color[n] == WHITE:
            found = visit(n)
            if found:
                return found
    return None


def _validate_recipe(recipe: dict, allowed_tiers: set[str], stage_keys: set[str]) -> None:
    tier = recipe.get("tier")
    require(tier in allowed_tiers, f"unknown tier {tier!r}; allowed={sorted(allowed_tiers)}")

    stages = recipe.get("stages")
    require(isinstance(stages, list) and stages, f"tier {tier}: stages must be a non-empty list")

    ids = [s.get("id") for s in stages]
    require(all(isinstance(i, str) and i for i in ids), f"tier {tier}: every stage needs a string id")
    require(len(set(ids)) == len(ids), f"tier {tier}: duplicate stage ids in {ids}")
    index = {sid: i for i, sid in enumerate(ids)}

    # Stage shape + health-gate teeth.
    for stage in stages:
        sid = stage["id"]
        missing = stage_keys - set(stage)
        require(not missing, f"tier {tier} stage {sid}: missing keys {sorted(missing)}")
        require(isinstance(stage["depends_on"], list), f"tier {tier} stage {sid}: depends_on must be a list")
        health = stage["health"]
        require(isinstance(health, dict), f"tier {tier} stage {sid}: health must be an object")
        # A gate is real only when required==true AND check is non-empty.
        require(
            health.get("required") is True and isinstance(health.get("check"), str) and health["check"].strip(),
            f"tier {tier} stage {sid}: missing health gate (need required:true and a non-empty check)",
        )

    # Referential integrity.
    edges: dict[str, list[str]] = {}
    for stage in stages:
        sid = stage["id"]
        deps = stage["depends_on"]
        for d in deps:
            require(d in index, f"tier {tier} stage {sid}: depends_on unknown stage {d!r}")
        edges[sid] = list(deps)

    # Cycle teeth (runs before ordering so a true cycle reports as a cycle).
    cycle = _detect_cycle(ids, edges)
    require(cycle is None, f"tier {tier}: cycle in depends_on: {' -> '.join(cycle) if cycle else ''}")

    # Ordering teeth: no forward edges, no skip edges.
    for pos, stage in enumerate(stages):
        sid = stage["id"]
        deps = stage["depends_on"]
        for d in deps:
            require(index[d] < pos, f"tier {tier} stage {sid}: forward edge -- depends on later stage {d!r}")
        if pos == 0:
            require(not deps, f"tier {tier} stage {sid}: first stage must not depend on anything")
        else:
            predecessor = ids[pos - 1]
            require(
                predecessor in deps,
                f"tier {tier} stage {sid}: skip edge -- must depend on immediate predecessor {predecessor!r}",
            )


def validate(path: Path) -> None:
    schema = _load_schema()
    allowed_tiers = _allowed_tiers(schema)
    stage_keys = _stage_required_keys(schema)

    doc = yaml.safe_load(path.read_text(encoding="utf-8"))
    require(isinstance(doc, dict), "document root must be a mapping")
    require(doc.get("apiVersion") == "cli.socioprophet.dev/v1", "apiVersion mismatch")
    require(doc.get("kind") == "ProphetCliBuildRecipes", "kind mismatch")
    recipes = doc.get("recipes")
    require(isinstance(recipes, list) and recipes, "recipes must be a non-empty list")

    for recipe in recipes:
        require(isinstance(recipe, dict), "each recipe must be a mapping")
        _validate_recipe(recipe, allowed_tiers, stage_keys)


# --- self-test: prove the teeth fire both ways --------------------------------

REJECT_FIXTURES = {
    "reject-forward-dep.yaml": "forward edge",
    "reject-missing-gate.yaml": "missing health gate",
    "reject-unknown-tier.yaml": "unknown tier",
    "reject-cycle.yaml": "cycle in depends_on",
}


def self_test() -> int:
    failures: list[str] = []

    # Accept: the canonical 3x3 recipe must validate.
    try:
        validate(RECIPES_PATH)
        print(f"PASS accept {RECIPES_PATH.relative_to(REPO_ROOT)}")
    except RecipeError as exc:
        failures.append(f"accept fixture {RECIPES_PATH.name} should VERIFY but was rejected: {exc}")

    # Reject: each bad fixture must fire, with the expected reason.
    for name, expected in REJECT_FIXTURES.items():
        path = FIXTURES_DIR / name
        try:
            validate(path)
            failures.append(f"reject fixture {name} should be REJECTED but VERIFIED (negative control breached)")
        except RecipeError as exc:
            if expected in str(exc):
                print(f"PASS reject {name}: {exc}")
            else:
                failures.append(f"reject fixture {name} fired wrong reason: expected ~{expected!r}, got: {exc}")

    if failures:
        for f in failures:
            print(f"FAIL {f}", file=sys.stderr)
        return 1
    print("PASS all build-recipe teeth verified")
    return 0


def main(argv: list[str]) -> int:
    if "--self-test" in argv:
        return self_test()
    target = Path(argv[1]) if len(argv) > 1 else RECIPES_PATH
    try:
        validate(target)
        print(f"PASS {target}")
        return 0
    except Exception as exc:  # noqa: BLE001
        print(f"FAIL {target}: {exc}", file=sys.stderr)
        return 1


if __name__ == "__main__":
    raise SystemExit(main(sys.argv))
