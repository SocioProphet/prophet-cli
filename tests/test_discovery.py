"""Zero-config discovery must resolve when a target is present and return a clear
MISS (never a silent guess) when it is not — proven both ways, offline."""
import os
from pathlib import Path

from prophet.integrations import discovery


def test_env_override_wins(monkeypatch, tmp_path):
    monkeypatch.setenv("SOURCEOS_CONTINUUM", "https://continuum.local:8443")
    monkeypatch.setenv("ATLAS_ENDPOINT", "grpc://atlas.local:9090")
    c = discovery.discover_continuum(dev_root=tmp_path)
    a = discovery.discover_tritfabric(dev_root=tmp_path)
    assert c and c.ok and c.kind == "endpoint" and c.source == "env:SOURCEOS_CONTINUUM"
    assert a and a.ok and a.location == "grpc://atlas.local:9090"


def test_sibling_repo_discovery(monkeypatch, tmp_path):
    for name in ("sourceos-continuum", "tritfabric"):
        (tmp_path / name / ".git").mkdir(parents=True)
    monkeypatch.delenv("SOURCEOS_CONTINUUM", raising=False)
    monkeypatch.delenv("ATLAS_ENDPOINT", raising=False)
    monkeypatch.delenv("TRITFABRIC_ENDPOINT", raising=False)
    monkeypatch.setattr(discovery.shutil, "which", lambda _b: None)  # no CLIs on PATH
    monkeypatch.setenv("XDG_CONFIG_HOME", str(tmp_path / "noconfig"))
    c = discovery.discover_continuum(dev_root=tmp_path)
    a = discovery.discover_tritfabric(dev_root=tmp_path)
    assert c and c.kind == "repo" and c.source == "repo:sourceos-continuum"
    assert a and a.kind == "repo" and a.source == "repo:tritfabric"


def test_missing_is_clear_not_guessed(monkeypatch, tmp_path):
    for var in ("SOURCEOS_CONTINUUM", "ATLAS_ENDPOINT", "TRITFABRIC_ENDPOINT"):
        monkeypatch.delenv(var, raising=False)
    monkeypatch.setattr(discovery.shutil, "which", lambda _b: None)
    monkeypatch.setenv("XDG_CONFIG_HOME", str(tmp_path / "noconfig"))
    assert discovery.discover_continuum(dev_root=tmp_path) is None
    assert discovery.discover_tritfabric(dev_root=tmp_path) is None
    both = discovery.discover_all(dev_root=tmp_path)
    assert both["continuum"] is None and both["tritfabric"] is None


if __name__ == "__main__":
    import pytest, sys
    sys.exit(pytest.main([__file__, "-q"]))
