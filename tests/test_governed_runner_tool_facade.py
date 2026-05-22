"""Tests for Prophet governed-runner tool facade commands."""

from __future__ import annotations

import subprocess

from prophet_cli.cli import main


class Completed:
    def __init__(self, returncode: int = 0):
        self.returncode = returncode


def test_governed_runner_tool_list_tools_delegates_to_sp_run(monkeypatch):
    calls: list[list[str]] = []

    monkeypatch.setattr("shutil.which", lambda binary: f"/usr/bin/{binary}")

    def fake_run(cmd, check=False):
        calls.append(cmd)
        assert check is False
        return Completed(0)

    monkeypatch.setattr(subprocess, "run", fake_run)

    rc = main(["governed-runner", "tool", "list-tools"])

    assert rc == 0
    assert calls == [["/usr/bin/sp-run", "tool", "list-tools"]]


def test_governed_runner_tool_call_delegates_to_sp_run(monkeypatch):
    calls: list[list[str]] = []

    monkeypatch.setattr("shutil.which", lambda binary: f"/usr/bin/{binary}")

    def fake_run(cmd, check=False):
        calls.append(cmd)
        assert check is False
        return Completed(0)

    monkeypatch.setattr(subprocess, "run", fake_run)

    rc = main([
        "governed-runner",
        "tool",
        "call",
        "governed_runner.doctor",
        "--args-json",
        "{}",
    ])

    assert rc == 0
    assert calls == [[
        "/usr/bin/sp-run",
        "tool",
        "call",
        "governed_runner.doctor",
        "--args-json",
        "{}",
    ]]
