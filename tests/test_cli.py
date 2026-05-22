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


def test_sourceos_network_delegates_to_sourceosctl(monkeypatch):
    calls: list[list[str]] = []

    monkeypatch.setattr("shutil.which", lambda binary: f"/usr/bin/{binary}")

    def fake_run(cmd, check=False):
        calls.append(cmd)
        assert check is False
        return Completed(0)

    monkeypatch.setattr(subprocess, "run", fake_run)

    rc = main(["sourceos", "network", "doctor"])

    assert rc == 0
    assert calls == [["/usr/bin/sourceosctl", "network", "doctor"]]


def test_sourceos_native_assistant_delegates_to_sourceosctl(monkeypatch):
    calls: list[list[str]] = []

    monkeypatch.setattr("shutil.which", lambda binary: f"/usr/bin/{binary}")

    def fake_run(cmd, check=False):
        calls.append(cmd)
        assert check is False
        return Completed(0)

    monkeypatch.setattr(subprocess, "run", fake_run)

    rc = main(["sourceos", "native-assistant", "plan", "--operation", "open-workroom"])

    assert rc == 0
    assert calls == [["/usr/bin/sourceosctl", "native-assistant", "plan", "--operation", "open-workroom"]]


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


def test_agentplane_doctor_delegates_to_sp_run(monkeypatch):
    calls: list[list[str]] = []

    monkeypatch.setattr("shutil.which", lambda binary: f"/usr/bin/{binary}")

    def fake_run(cmd, check=False):
        calls.append(cmd)
        assert check is False
        return Completed(0)

    monkeypatch.setattr(subprocess, "run", fake_run)

    rc = main(["agentplane", "doctor"])

    assert rc == 0
    assert calls == [["/usr/bin/sp-run", "doctor"]]


def test_agentplane_preflight_delegates_to_sp_run(monkeypatch):
    calls: list[list[str]] = []

    monkeypatch.setattr("shutil.which", lambda binary: f"/usr/bin/{binary}")

    def fake_run(cmd, check=False):
        calls.append(cmd)
        assert check is False
        return Completed(0)

    monkeypatch.setattr(subprocess, "run", fake_run)

    rc = main(["agentplane", "preflight", "contract.json"])

    assert rc == 0
    assert calls == [["/usr/bin/sp-run", "preflight", "contract.json"]]


def test_agentplane_admit_delegates_to_sp_run(monkeypatch):
    calls: list[list[str]] = []

    monkeypatch.setattr("shutil.which", lambda binary: f"/usr/bin/{binary}")

    def fake_run(cmd, check=False):
        calls.append(cmd)
        assert check is False
        return Completed(0)

    monkeypatch.setattr(subprocess, "run", fake_run)

    rc = main([
        "agentplane",
        "admit",
        "contract.json",
        "--preflight",
        "preflight.json",
        "--authority-state",
        "authority.json",
    ])

    assert rc == 0
    assert calls == [[
        "/usr/bin/sp-run",
        "admit",
        "contract.json",
        "--preflight",
        "preflight.json",
        "--authority-state",
        "authority.json",
    ]]


def test_governed_runner_alias_delegates_to_sp_run(monkeypatch):
    calls: list[list[str]] = []

    monkeypatch.setattr("shutil.which", lambda binary: f"/usr/bin/{binary}")

    def fake_run(cmd, check=False):
        calls.append(cmd)
        assert check is False
        return Completed(0)

    monkeypatch.setattr(subprocess, "run", fake_run)

    rc = main(["governed-runner", "doctor"])

    assert rc == 0
    assert calls == [["/usr/bin/sp-run", "doctor"]]


def test_governed_runner_dossier_delegates_to_sp_run(monkeypatch):
    calls: list[list[str]] = []

    monkeypatch.setattr("shutil.which", lambda binary: f"/usr/bin/{binary}")

    def fake_run(cmd, check=False):
        calls.append(cmd)
        assert check is False
        return Completed(0)

    monkeypatch.setattr(subprocess, "run", fake_run)

    rc = main(["governed-runner", "dossier", ".socioprophet/runs/demo"])

    assert rc == 0
    assert calls == [["/usr/bin/sp-run", "dossier", ".socioprophet/runs/demo"]]


def test_governed_runner_validate_dossier_delegates_to_sp_run(monkeypatch):
    calls: list[list[str]] = []

    monkeypatch.setattr("shutil.which", lambda binary: f"/usr/bin/{binary}")

    def fake_run(cmd, check=False):
        calls.append(cmd)
        assert check is False
        return Completed(0)

    monkeypatch.setattr(subprocess, "run", fake_run)

    rc = main(["governed-runner", "validate-dossier", "run-dossier.json"])

    assert rc == 0
    assert calls == [["/usr/bin/sp-run", "validate-dossier", "run-dossier.json"]]


def test_missing_delegate_returns_127(monkeypatch, capsys):
    monkeypatch.setattr("shutil.which", lambda binary: None)

    rc = main(["sourceos", "office", "doctor"])

    captured = capsys.readouterr()
    assert rc == 127
    assert "required delegate not found: sourceosctl" in captured.err
    assert "sourceos-devtools" in captured.err


def test_missing_sp_run_delegate_returns_127(monkeypatch, capsys):
    monkeypatch.setattr("shutil.which", lambda binary: None)

    rc = main(["agentplane", "doctor"])

    captured = capsys.readouterr()
    assert rc == 127
    assert "required delegate not found: sp-run" in captured.err
    assert "AgentPlane" in captured.err


def test_missing_sp_run_for_governed_runner_alias_returns_127(monkeypatch, capsys):
    monkeypatch.setattr("shutil.which", lambda binary: None)

    rc = main(["governed-runner", "doctor"])

    captured = capsys.readouterr()
    assert rc == 127
    assert "required delegate not found: sp-run" in captured.err
    assert "AgentPlane" in captured.err
