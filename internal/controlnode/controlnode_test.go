package controlnode

import (
	"context"
	"errors"
	"testing"
)

func TestStatusIncludesContractsRepo(t *testing.T) {
	resp, err := Status(context.Background())
	if err != nil {
		t.Fatalf("Status returned error: %v", err)
	}
	if got := resp["contracts_repo"]; got != "SourceOS-Linux/sourceos-spec" {
		t.Fatalf("unexpected contracts_repo: %v", got)
	}
}

func TestProcessFailsClosedAndReturnsProbeEnvelope(t *testing.T) {
	resp, err := Process(context.Background(), "input.json", "out")
	if err == nil {
		t.Fatalf("expected Process to fail when processor cannot run")
	}
	var processErr ProcessError
	if !errors.As(err, &processErr) {
		t.Fatalf("expected ProcessError, got %T", err)
	}
	if got := resp["operation"]; got != "process" {
		t.Fatalf("unexpected operation: %v", got)
	}
	if got := resp["status"]; got != "failed" {
		t.Fatalf("expected failed status, got %v", got)
	}
	probe, ok := resp["probe"].(ExecResult)
	if !ok {
		t.Fatalf("probe was not ExecResult: %T", resp["probe"])
	}
	if len(probe.Command) == 0 || probe.Command[0] != "python3" {
		t.Fatalf("unexpected probe command: %#v", probe.Command)
	}
}
