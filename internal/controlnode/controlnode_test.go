package controlnode

import (
	"context"
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

func TestProcessReturnsProbeEnvelope(t *testing.T) {
	resp, err := Process(context.Background(), "input.json", "out")
	if err != nil {
		t.Fatalf("Process returned error: %v", err)
	}
	if got := resp["operation"]; got != "process" {
		t.Fatalf("unexpected operation: %v", got)
	}
	probe, ok := resp["probe"].(ExecResult)
	if !ok {
		t.Fatalf("probe was not ExecResult: %T", resp["probe"])
	}
	if len(probe.Command) == 0 || probe.Command[0] != "python3" {
		t.Fatalf("unexpected probe command: %#v", probe.Command)
	}
}
