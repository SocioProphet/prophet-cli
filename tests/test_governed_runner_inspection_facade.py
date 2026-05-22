"""Tests for Prophet governed-runner inspection facade commands."""

from __future__ import annotations

import subprocess

from prophet_cli.cli import main


class Completed:
    def __init__(self, returncode: int = 0):
        self.returncode = returncode


def test_governed_runner_list_delegates_to_sp_run(monkeypatch):
    calls: list[list[str]] = []

    monkeypatch.setattr("shutil.which", lambda binary: f"/usr/bin/{binary}")

    def fake_run(cmd, check=False):
        calls.append(cmd)
        assert check is False
        return Completed(0)

    monkeypatch.setattr(subprocess, "run", fake_run)

    rc = main([
        "governed-runner",
        "list",
        "--runs-root",
        ".socioprophet/smoke/governed-runner",
    ])

    assert rc == 0
    assert calls == [[
        "/usr/bin/sp-run",
        "list",
        "--runs-root",
        ".socioprophet/smoke/governed-runner",
    ]]


def test_governed_runner_status_delegates_to_sp_run(monkeypatch):
    calls: list[list[str]] = []

    monkeypatch.setattr("shutil.which", lambda binary: f"/usr/bin/{binary}")

    def fake_run(cmd, check=False):
        calls.append(cmd)
        assert check is False
        return Completed(0)

    monkeypatch.setattr(subprocess, "run", fake_run)

    rc = main([
        "governed-runner",
        "status",
        ".socioprophet/smoke/governed-runner/run",
    ])

    assert rc == 0
    assert calls == [[
        "/usr/bin/sp-run",
        "status",
        ".socioprophet/smoke/governed-runner/run",
    ]]


def test_governed_runner_inspect_delegates_to_sp_run(monkeypatch):
    calls: list[list[str]] = []

    monkeypatch.setattr("shutil.which", lambda binary: f"/usr/bin/{binary}")

    def fake_run(cmd, check=False):
        calls.append(cmd)
        assert check is False
        return Completed(0)

    monkeypatch.setattr(subprocess, "run", fake_run)

    rc = main([
        "governed-runner",
        "inspect",
        ".socioprophet/smoke/governed-runner/run",
    ])

    assert rc == 0
    assert calls == [[
        "/usr/bin/sp-run",
        "inspect",
        ".socioprophet/smoke/governed-runner/run",
    ]]
