"""Tests for the Prophet governed-runner smoke facade."""

from __future__ import annotations

import subprocess

from prophet_cli.cli import main


class Completed:
    def __init__(self, returncode: int = 0):
        self.returncode = returncode


def test_governed_runner_smoke_delegates_to_sp_run(monkeypatch):
    calls: list[list[str]] = []

    monkeypatch.setattr("shutil.which", lambda binary: f"/usr/bin/{binary}")

    def fake_run(cmd, check=False):
        calls.append(cmd)
        assert check is False
        return Completed(0)

    monkeypatch.setattr(subprocess, "run", fake_run)

    rc = main([
        "governed-runner",
        "smoke",
        "--output-dir",
        ".socioprophet/smoke/governed-runner",
    ])

    assert rc == 0
    assert calls == [[
        "/usr/bin/sp-run",
        "smoke",
        "--output-dir",
        ".socioprophet/smoke/governed-runner",
    ]]
