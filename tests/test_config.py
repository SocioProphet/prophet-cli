"""Tests for the cross-platform config core: fail-closed doctor + native-rc emit."""
import os
import stat
import tempfile
from pathlib import Path

from prophet_cli import config


def _cfg(**over):
    base = {
        "profile": {"name": "box", "endpoint": "box"},
        "compute": {"kube_context": "sourceos-box", "kube_namespace": "ns"},
        "identity": {"tenant": "you", "user": "dev"},
        "hellgraph": {"http_port": 8787},
        "secrets": {},
    }
    base.update(over)
    return base


def _levels(results):
    return {r["name"]: r["level"] for r in results}


def test_clean_config_has_no_required_errors():
    r = config.doctor(_cfg(), platform="posix")
    assert [x for x in r if x["level"] == "error" and x["name"].startswith("required")] == []


def test_prod_kube_context_is_an_error():
    for ctx in ("gke_prod-cluster", "eks-production", "my-PROD-ctx"):
        r = config.doctor(_cfg(compute={"kube_context": ctx, "kube_namespace": "ns"}), platform="posix")
        assert _levels(r)["compute.kube_context"] == "error", ctx


def test_sovereign_context_is_ok():
    assert _levels(config.doctor(_cfg(), platform="posix"))["compute.kube_context"] == "ok"


def test_missing_required_key_is_error():
    c = _cfg()
    del c["identity"]["tenant"]
    assert _levels(config.doctor(c, platform="posix"))["required:identity.tenant"] == "error"


def test_secret_file_existence_and_0600_enforced_on_posix():
    with tempfile.TemporaryDirectory() as td:
        good = Path(td) / "k.key"; good.write_text("x"); os.chmod(good, 0o600)
        bad = Path(td) / "b.key"; bad.write_text("x"); os.chmod(bad, 0o700)  # not 0600 -> refused; no group/world bits at all
        c = _cfg(secrets={
            "good": {"file": str(good)},
            "loose": {"file": str(bad)},
            "gone": {"file": str(Path(td) / "nope.key")},
        })
        lv = _levels(config.doctor(c, platform="posix"))
        assert lv["secret:good"] == "ok"
        assert lv["secret:loose"] == "error"   # 0700 → refused (anything but 0600)
        assert lv["secret:gone"] == "error"    # missing → refused


def test_windows_uses_cred_ref_not_file_perms():
    c = _cfg(secrets={"k": {"file": "/nonexistent", "cred": "sourceos-signing"}})
    assert _levels(config.doctor(c, platform="windows"))["secret:k"] == "ok"  # cred ref; no 0600 on win


def test_emit_posix_generates_prophetrc():
    out = config.emit(_cfg(secrets={"signing_key": {"file": "${NOETICA_HOME}/sovereign-root.key",
                                                     "cred": "sourceos-signing"}}), target="posix")
    assert 'export PROPHET_PROFILE="box"' in out
    assert 'export SOURCEOS_KUBE_CONTEXT="sourceos-box"' in out
    assert 'export SOURCEOS_SIGNING_KEY_FILE="${NOETICA_HOME}/sovereign-root.key"' in out


def test_emit_powershell_generates_psm1():
    out = config.emit(_cfg(secrets={"signing_key": {"file": "x", "cred": "sourceos-signing"}}),
                      target="powershell")
    assert '$env:PROPHET_PROFILE = "box"' in out
    assert '$env:SOURCEOS_KUBE_CONTEXT = "sourceos-box"' in out
    assert '$env:SOURCEOS_SIGNING_KEY_REF = "cred:sourceos-signing"' in out
    assert "$env:APPDATA" in out  # native Windows idiom, not XDG


def test_doctor_never_echoes_a_secret_path():
    """CodeQL flagged this as clear-text logging of sensitive information, and it was right.

    `doctor` is a REPORTER: its output goes to terminals, CI logs and pasted transcripts. The path
    of a signing key is attack surface even though it is not the key. The estate rule is that
    scanners do not echo their matches, and this is the same rule.
    """
    with tempfile.TemporaryDirectory() as td:
        secret = Path(td) / "sovereign-root.key"
        secret.write_text("x")
        os.chmod(secret, 0o700)                      # wrong mode, so it is reported at all
        results = config.doctor(_cfg(secrets={"signing_key": {"file": str(secret)}}), platform="posix")
        for r in results:
            assert td not in r["message"], f"leaked the directory in {r['name']}"
            assert "sovereign-root" not in r["message"], f"leaked the filename in {r['name']}"
        msg = _levels(results)["secret:signing_key"]
        assert msg == "error", "the wrong mode must still be refused — redaction is not leniency"


def test_emit_still_carries_the_reference_because_the_rc_needs_it():
    """The counterpart: `emit` is a GENERATOR, not a reporter. Its stdout IS the rc file, which must
    carry the reference for the shell to resolve at load. Redacting there would break the feature."""
    out = config.emit(_cfg(secrets={"signing_key": {"file": "${NOETICA_HOME}/sovereign-root.key"}}),
                      target="posix")
    assert 'SOURCEOS_SIGNING_KEY_FILE="${NOETICA_HOME}/sovereign-root.key"' in out
