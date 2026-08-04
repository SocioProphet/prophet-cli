"""Prophet toolchain config — one platform-neutral source of truth, three native bindings.

`config.toml` is the single human-edited config (endpoints, contexts, tenancy, opt-ins, and secret
*references* — never secret values). From it this module:

  * `doctor(cfg)`  — a fail-closed check-set (the SAME checks on every platform): required keys present,
    the kube-context is NOT a production cluster, referenced secret files exist (and are 0600 on POSIX),
    the CLIs are on PATH. This is what makes "works as advertised" an enforced control, not a hope.
  * `emit(cfg, target=…)` — GENERATES the native rc so the three platforms never drift: a POSIX
    `.prophetrc` (`export`, XDG, `*_FILE`, 0600) or a Windows `prophet.psm1` (`$env:`, `%APPDATA%`,
    `cred:` via DPAPI/Credential Manager). The SourceOS Nix binder writes the POSIX form declaratively.

stdlib only (tomllib); runs identically on SourceOS / macOS / Windows.
"""
from __future__ import annotations

import os
import re
import stat
import sys
import tomllib
from pathlib import Path

# config.toml [section].key  ->  the real env var each tool reads (grounded in the estate).
ENV_MAP = {
    "profile.name": "PROPHET_PROFILE",
    "profile.endpoint": "SOURCEOS_ENDPOINT",
    "compute.kube_context": "SOURCEOS_KUBE_CONTEXT",
    "compute.kube_namespace": "SOURCEOS_KUBE_NAMESPACE",
    "compute.slurm_login": "SOURCEOS_SLURM_LOGIN",
    "hellgraph.http_port": "HELLGRAPH_HTTP_PORT",
    "identity.tenant": "SOURCEOS_TENANT",
    "identity.user": "SOURCEOS_USER",
    "identity.sensitivity": "SOURCEOS_SENSITIVITY",
}
REQUIRED = ["profile.name", "compute.kube_context", "compute.kube_namespace", "identity.tenant"]
# a kube-context containing any of these is almost certainly the shared prod cluster, not your sovereign one.
PROD_CONTEXT_MARKERS = ["prod", "production", "gke_", "_gke", "-prod-", "eks", "aks"]


def default_config_path() -> Path:
    if os.name == "nt":
        base = os.environ.get("APPDATA", str(Path.home() / "AppData" / "Roaming"))
        return Path(base) / "prophet" / "config.toml"
    base = os.environ.get("PROPHET_CONFIG") or os.environ.get("XDG_CONFIG_HOME") or str(Path.home() / ".config")
    # PROPHET_CONFIG may already be the prophet dir; XDG is the parent.
    p = Path(base)
    return (p / "config.toml") if p.name == "prophet" else (p / "prophet" / "config.toml")


def load(path: str | os.PathLike | None = None) -> dict:
    p = Path(path) if path else default_config_path()
    if not p.exists():
        raise FileNotFoundError(f"no prophet config at {p} — run `prophet config init` or ship it via the profile")
    with open(p, "rb") as f:
        return tomllib.load(f)


def _get(cfg: dict, dotted: str):
    cur = cfg
    for part in dotted.split("."):
        if not isinstance(cur, dict) or part not in cur:
            return None
        cur = cur[part]
    return cur


def _secret_refs(cfg: dict) -> dict:
    return cfg.get("secrets") or {}


def _expand_posix(value: str) -> str:
    # ${NOETICA_HOME}/x -> resolved for CHECKING (emit keeps the literal so the shell expands at load).
    return os.path.expandvars(value.replace("~", str(Path.home())))


def doctor(cfg: dict, *, platform: str | None = None) -> list[dict]:
    """Fail-closed check-set. Returns [{level: ok|warn|error, name, message}]. Same logic everywhere;
    only the secret-permission check differs (POSIX 0600 vs Windows credential store)."""
    platform = platform or ("windows" if os.name == "nt" else "posix")
    out: list[dict] = []

    def add(level, name, msg):
        out.append({"level": level, "name": name, "message": msg})

    for key in REQUIRED:
        add("ok" if _get(cfg, key) not in (None, "") else "error", f"required:{key}",
            "set" if _get(cfg, key) not in (None, "") else "MISSING — a rollout cannot proceed without it")

    prof = _get(cfg, "profile.name")
    add("ok" if prof in ("box", "twin", "cloud") else "error", "profile.name",
        f"{prof!r}" if prof in ("box", "twin", "cloud") else f"{prof!r} not in box|twin|cloud")

    ctx = str(_get(cfg, "compute.kube_context") or "")
    if any(m in ctx.lower() for m in PROD_CONTEXT_MARKERS):
        add("error", "compute.kube_context",
            f"{ctx!r} looks like a PRODUCTION cluster — set your sovereign context, never the ambient default")
    elif ctx:
        add("ok", "compute.kube_context", ctx)

    for name, ref in _secret_refs(cfg).items():
        if platform == "windows":
            cred = (ref or {}).get("cred") if isinstance(ref, dict) else None
            add("ok" if cred else "warn", f"secret:{name}",
                f"cred:{cred} (DPAPI/Credential Manager)" if cred else "no windows credential ref")
            continue
        fileref = (ref or {}).get("file") if isinstance(ref, dict) else ref
        if not fileref:
            add("warn", f"secret:{name}", "no file reference on this platform")
            continue
        fp = Path(_expand_posix(str(fileref)))
        if not fp.exists():
            add("error", f"secret:{name}", f"referenced secret file missing: {fp}")
        else:
            mode = stat.S_IMODE(fp.stat().st_mode)
            add("ok" if mode == 0o600 else "error", f"secret:{name}",
                f"{fp} (0600)" if mode == 0o600 else f"{fp} is {oct(mode)} — must be 0600 (secrets are not world/group readable)")

    import shutil
    for cli in (cfg.get("require", {}).get("clis") or ["prophet", "sourceosctl"]):
        add("ok" if shutil.which(cli) else "warn", f"cli:{cli}",
            "on PATH" if shutil.which(cli) else "not on PATH (install the toolchain for this build)")

    return out


def emit(cfg: dict, *, target: str = "posix") -> str:
    """Generate the native rc from config.toml. `target` in {posix, powershell}."""
    env = {}
    for dotted, var in ENV_MAP.items():
        v = _get(cfg, dotted)
        if v is not None:
            env[var] = str(v)

    def secret_lines_posix():
        lines = []
        for name, ref in _secret_refs(cfg).items():
            fileref = (ref or {}).get("file") if isinstance(ref, dict) else ref
            if fileref:
                lines.append(f'export SOURCEOS_{name.upper()}_FILE="{fileref}"')
        return lines

    def secret_lines_ps():
        lines = []
        for name, ref in _secret_refs(cfg).items():
            cred = (ref or {}).get("cred") if isinstance(ref, dict) else None
            if cred:
                lines.append(f'$env:SOURCEOS_{name.upper()}_REF = "cred:{cred}"')
        return lines

    if target == "powershell":
        body = ["# prophet.psm1 — GENERATED from config.toml by `prophet doctor --emit powershell`. Do not hand-edit.",
                "# Imported from $PROFILE. CONFIG only; secrets are references (DPAPI/Credential Manager).",
                '$env:PROPHET_CONFIG = "$env:APPDATA\\prophet"',
                '$env:NOETICA_HOME  = if ($env:NOETICA_HOME) { $env:NOETICA_HOME } else { "$env:LOCALAPPDATA\\noetica" }']
        body += [f'$env:{k} = "{v}"' for k, v in env.items()]
        body += secret_lines_ps()
        body += ['if (Get-Command prophet -ErrorAction SilentlyContinue) { prophet doctor --quiet }']
        return "\n".join(body) + "\n"

    # posix
    body = ["# ~/.prophetrc — GENERATED from config.toml by `prophet doctor --emit posix`. Do not hand-edit.",
            "# Sourced by the workstation shell-spine. CONFIG only; secrets are *_FILE references (0600).",
            'export PROPHET_CONFIG="${XDG_CONFIG_HOME:-$HOME/.config}/prophet"',
            'export NOETICA_HOME="${NOETICA_HOME:-$HOME/.noetica}"']
    body += [f'export {k}="{v}"' for k, v in env.items()]
    body += secret_lines_posix()
    body += ['[ -f "$PROPHET_CONFIG/secrets.env" ] && { set -a; . "$PROPHET_CONFIG/secrets.env"; set +a; }',
             'command -v prophet >/dev/null 2>&1 && prophet doctor --quiet']
    return "\n".join(body) + "\n"


def _print_doctor(results: list[dict], *, as_json: bool, quiet: bool) -> int:
    errors = [r for r in results if r["level"] == "error"]
    if as_json:
        import json
        print(json.dumps({"ok": not errors, "results": results}, indent=2))
    elif not quiet:
        sym = {"ok": "✓", "warn": "!", "error": "✗"}
        for r in results:
            print(f"  {sym.get(r['level'],'?')} {r['name']}: {r['message']}", file=sys.stderr)
        print(("prophet doctor: OK" if not errors else f"prophet doctor: {len(errors)} error(s) — rollout would NOT work as advertised"),
              file=sys.stderr)
    elif errors:
        print(f"prophet doctor: {len(errors)} error(s)", file=sys.stderr)
    return 0 if not errors else 1


TEMPLATE = '''# config.toml — the SourceOS / Prophet toolchain config. ONE source of truth; the native rc files
# (.prophetrc on POSIX, prophet.psm1 on Windows) are GENERATED from this by `prophet doctor --emit`.
# CONFIG only — secret VALUES never live here; secrets are references.
schema = "prophet.config/v0"

[profile]
name = "box"        # box | twin | cloud
endpoint = "box"    # box = local/LAN, twin = always-on cloud rendezvous

[compute]
kube_context = "sourceos-box"     # MUST be your sovereign context — never the ambient prod cluster
kube_namespace = "you-default"
# slurm_login = "slurm.local"     # opt-in HPC

[hellgraph]
http_port = 8787
superpeer = false                 # p2p federation opt-in

[identity]
tenant = "you"
user = "dev"
sensitivity = "normal"            # normal | sensitive

[secrets]
# references, never values. POSIX resolves `file` (checked for existence + 0600); Windows uses `cred`
# (DPAPI / Credential Manager). Most already live in the .noetica state home.
signing_key   = { file = "${NOETICA_HOME}/sovereign-root.key", cred = "sourceos-signing" }
anthropic_key = { file = "${NOETICA_HOME}/anthropic.key",      cred = "anthropic" }

[require]
clis = ["prophet", "sourceosctl"]
'''


def run_config(argv: list[str]) -> int:
    """`prophet config [path|show|init]`."""
    cmd = argv[0] if argv else "path"
    p = default_config_path()
    if cmd == "path":
        print(p)
        return 0
    if cmd == "init":
        if p.exists():
            print(f"prophet config: {p} already exists — not overwriting", file=sys.stderr)
            return 1
        p.parent.mkdir(parents=True, exist_ok=True)
        p.write_text(TEMPLATE)
        print(f"wrote {p} — edit it, then `prophet doctor`")
        return 0
    if cmd == "show":
        try:
            print(p.read_text())
            return 0
        except OSError as e:
            print(f"prophet config: {e}", file=sys.stderr)
            return 1
    print(f"prophet config: unknown subcommand {cmd!r} (path|show|init)", file=sys.stderr)
    return 2


def run_doctor(argv: list[str]) -> int:
    """`prophet doctor [--json] [--quiet] [--emit posix|powershell] [--config PATH]`."""
    import argparse
    ap = argparse.ArgumentParser(prog="prophet doctor")
    ap.add_argument("--json", action="store_true")
    ap.add_argument("--quiet", action="store_true")
    ap.add_argument("--emit", choices=["posix", "powershell"])
    ap.add_argument("--config")
    a = ap.parse_args(argv)
    try:
        cfg = load(a.config)
    except FileNotFoundError as e:
        print(f"prophet doctor: {e}", file=sys.stderr)
        return 1
    if a.emit:
        sys.stdout.write(emit(cfg, target=a.emit))
        return 0
    return _print_doctor(doctor(cfg), as_json=a.json, quiet=a.quiet)
