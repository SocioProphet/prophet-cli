"""Tests for the Prophet facade CLI."""

from __future__ import annotations

import subprocess

from prophet_cli.cli import main


class Completed:
    def __init__(self, returncode: int = 0):
        self.returncode = returncode


def test_sourceos_local_model_delegates_to_sourceosctl(monkeypatch):
    calls: list[list[str]] = []

    monkeypatch.setattr("shutil.which", lambda binary: f"/usr/bin/{binary}")

    def fake_run(cmd, check=False):
        calls.append(cmd)
        assert check is False
        return Completed(0)

    monkeypatch.setattr(subprocess, "run", fake_run)

    rc = main(["sourceos", "local-model", "doctor"])

    assert rc == 0
    assert calls == [["/usr/bin/sourceosctl", "local-model", "doctor"]]


def test_sourceos_office_delegates_to_sourceosctl(monkeypatch):
    calls: list[list[str]] = []

    monkeypatch.setattr("shutil.which", lambda binary: f"/usr/bin/{binary}")

    def fake_run(cmd, check=False):
        calls.append(cmd)
        assert check is False
        return Completed(0)

    monkeypatch.setattr(subprocess, "run", fake_run)

    rc = main(["sourceos", "office", "doctor"])

    assert rc == 0
    assert calls == [["/usr/bin/sourceosctl", "office", "doctor"]]


def test_sourceos_agent_machine_delegates_to_sourceosctl(monkeypatch):
    calls: list[list[str]] = []

    monkeypatch.setattr("shutil.which", lambda binary: f"/usr/bin/{binary}")

    def fake_run(cmd, check=False):
        calls.append(cmd)
        return Completed(0)

    monkeypatch.setattr(subprocess, "run", fake_run)

    rc = main(["sourceos", "agent-machine", "mounts", "plan"])

    assert rc == 0
    assert calls == [["/usr/bin/sourceosctl", "agent-machine", "mounts", "plan"]]


def test_sourceos_agent_term_delegates_to_agent_term(monkeypatch):
    calls: list[list[str]] = []

    monkeypatch.setattr("shutil.which", lambda binary: f"/usr/bin/{binary}")

    def fake_run(cmd, check=False):
        calls.append(cmd)
        return Completed(0)

    monkeypatch.setattr(subprocess, "run", fake_run)

    rc = main(["sourceos", "agent-term", "office", "inspect", "!prophet-workspace", "/tmp/demo.pptx"])

    assert rc == 0
    assert calls == [["/usr/bin/agent-term", "office", "inspect", "!prophet-workspace", "/tmp/demo.pptx"]]


def test_missing_delegate_returns_127(monkeypatch, capsys):
    monkeypatch.setattr("shutil.which", lambda binary: None)

    rc = main(["sourceos", "office", "doctor"])

    captured = capsys.readouterr()
    assert rc == 127
    assert "required delegate not found: sourceosctl" in captured.err
    assert "sourceos-devtools" in captured.err
