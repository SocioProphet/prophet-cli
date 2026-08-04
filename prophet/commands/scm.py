"""
prophet scm — GitHub <-> Gitea sovereignty: validate drift, dual-write, prepare the cutover.

The estate declared Gitea canonical and severed the pull-mirror, but the workflow never moved — so
GitHub holds the real state and Gitea is a frozen orphan. This surface makes the cutover CLEAN:
validate they are in sync before flipping, dual-write during the transition so they never diverge,
and never rip out `gh` (GitHub still runs CI meanwhile). It changes no remotes on its own.

Commands:
  prophet scm drift  [--repo OWNER/NAME]   compare GitHub vs Gitea HEAD — are we out of sync?
  prophet scm sync   [--branch B]          push the current branch to GitHub AND Gitea (dual-write)
  prophet scm status                       backend + auth + drift summary

Config (env): PROPHET_GITEA_HOST (default code.socioprophet.ai) · PROPHET_GITEA_TOKEN (read token,
mint in CI — never hardcode) · PROPHET_SCM=github|gitea|dual (default github).
"""
import json
import os
import subprocess
import urllib.request
from urllib.error import HTTPError, URLError

import click
from rich.console import Console
from rich.table import Table

console = Console()

GITEA_HOST = os.environ.get("PROPHET_GITEA_HOST", "code.socioprophet.ai")
BACKEND = os.environ.get("PROPHET_SCM", "github")


def _run(args) -> str:
    try:
        r = subprocess.run(args, capture_output=True, text=True, timeout=25)
        return r.stdout.strip() if r.returncode == 0 else ""
    except Exception:
        return ""


def _current_repo() -> str:
    """OWNER/NAME from gh, falling back to parsing the origin URL."""
    r = _run(["gh", "repo", "view", "--json", "nameWithOwner", "-q", ".nameWithOwner"])
    if r:
        return r
    url = _run(["git", "config", "--get", "remote.origin.url"])
    if url:
        tail = url.split("github.com")[-1].split("code.socioprophet.ai")[-1]
        return tail.lstrip(":/").removesuffix(".git")
    return ""


def _github_head(repo: str) -> str:
    return (_run(["gh", "api", f"repos/{repo}/commits/main", "-q", ".sha"]) or "")[:12]


def _gitea_head(repo: str) -> tuple:
    """(sha12, note) — note explains an empty sha (unauth / missing / unreachable)."""
    url = f"https://{GITEA_HOST}/api/v1/repos/{repo}/branches/main"
    req = urllib.request.Request(url)
    tok = os.environ.get("PROPHET_GITEA_TOKEN")
    if tok:
        req.add_header("Authorization", f"token {tok}")
    try:
        with urllib.request.urlopen(req, timeout=8) as resp:
            data = json.loads(resp.read())
            return ((data.get("commit") or {}).get("id", "")[:12], "")
    except HTTPError as e:
        return ("", f"HTTP {e.code}" + (" — set PROPHET_GITEA_TOKEN" if e.code in (401, 403) else ""))
    except (URLError, TimeoutError, OSError) as e:
        return ("", f"unreachable ({e})")
    except ValueError:
        return ("", "bad response")


@click.group()
def scm():
    """GitHub <-> Gitea sovereignty — validate drift, dual-write, prepare the cutover."""


@scm.command()
@click.option("--repo", default=None, help="OWNER/NAME (default: the current repo)")
def drift(repo):
    """Compare GitHub vs Gitea HEAD — is the 'canonical' Gitea actually in sync?"""
    repo = repo or _current_repo()
    if not repo:
        console.print("[red]not in a recognizable repo; pass --repo OWNER/NAME[/red]")
        raise SystemExit(2)
    gh = _github_head(repo)
    gt, note = _gitea_head(repo)
    t = Table(title=f"scm drift · {repo}", show_header=True, header_style="bold")
    t.add_column("remote"); t.add_column("main HEAD"); t.add_column("")
    t.add_row("github", gh or "?", "the real state (CI runs here)")
    t.add_row("gitea", gt or "—", note or "canonical-by-decree")
    console.print(t)
    if gt and gh == gt:
        console.print("[green]IN SYNC[/green] — safe to proceed with the cutover plan.")
    elif gt:
        console.print("[yellow]DRIFTED[/yellow] — GitHub is ahead of the 'canonical' Gitea. "
                      "Reconcile (push GitHub→Gitea via CI) before flipping.")
        raise SystemExit(1)
    else:
        console.print("[yellow]UNKNOWN[/yellow] — Gitea unreadable; treat as NOT clean until you can "
                      "confirm parity (set PROPHET_GITEA_TOKEN).")
        raise SystemExit(3)


@scm.command()
@click.option("--branch", default=None, help="branch to push (default: current)")
def sync(branch):
    """Dual-write: push the current branch to GitHub AND Gitea, so they never diverge."""
    branch = branch or _run(["git", "branch", "--show-current"])
    if not branch:
        console.print("[red]detached HEAD — checkout a branch first[/red]")
        raise SystemExit(2)
    console.print(f"pushing [bold]{branch}[/bold] to github (origin)…")
    if not _run(["git", "push", "origin", branch]) and subprocess.call(["git", "push", "origin", branch]):
        console.print("[red]github push failed[/red]"); raise SystemExit(1)
    remotes = _run(["git", "remote"]).split()
    if "gitea" in remotes:
        console.print(f"pushing [bold]{branch}[/bold] to gitea…")
        subprocess.call(["git", "push", "gitea", branch])
    else:
        repo = _current_repo()
        console.print("[yellow]no 'gitea' remote yet[/yellow] — add it (read the cutover plan first):\n"
                      f"  git remote add gitea https://{GITEA_HOST}/{repo}.git")


@scm.command()
def status():
    """Backend + GitHub auth + drift summary."""
    console.print(f"prophet scm — backend: [bold]{BACKEND}[/bold] · gitea host: {GITEA_HOST}")
    auth = _run(["gh", "auth", "status"]) or "github: not logged in"
    console.print(auth.splitlines()[0] if auth else "github: unknown")
    ctx = click.get_current_context()
    ctx.invoke(drift)
