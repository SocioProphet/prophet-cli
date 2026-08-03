package receipt

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNewPopulatesSchemaAndDuration(t *testing.T) {
	start := time.Date(2026, 8, 3, 12, 0, 0, 0, time.UTC)
	end := start.Add(1500 * time.Millisecond)
	r := New("prophet holmes analyze", "holmes", "ok", []string{"analyze", "x"}, start, end, nil)
	if r.Schema != Schema {
		t.Fatalf("schema = %q, want %q", r.Schema, Schema)
	}
	if r.DurationMs != 1500 {
		t.Fatalf("duration = %d, want 1500", r.DurationMs)
	}
	if r.Error != "" {
		t.Fatalf("unexpected error: %q", r.Error)
	}
}

func TestWriteToExplicitFile(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "nested", "run.json")
	r := New("prophet ledger validate", "model-governance-ledger", "ok", nil, time.Now(), time.Now(), nil)
	got, err := Write(dest, r)
	if err != nil {
		t.Fatalf("Write: %v", err)
	}
	if got != dest {
		t.Fatalf("path = %q, want %q", got, dest)
	}
	b, err := os.ReadFile(got)
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	var round Receipt
	if err := json.Unmarshal(b, &round); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if round.Command != "prophet ledger validate" || round.Schema != Schema {
		t.Fatalf("round-trip mismatch: %+v", round)
	}
}

func TestWriteToDirectoryMintsTimestampedFile(t *testing.T) {
	dir := t.TempDir()
	r := New("prophet agent registry list", "agent-registry", "ok", nil, time.Now(), time.Now(), nil)
	got, err := Write(dir, r)
	if err != nil {
		t.Fatalf("Write: %v", err)
	}
	if filepath.Dir(got) != dir {
		t.Fatalf("receipt not placed in dir: %q", got)
	}
	if !strings.HasPrefix(filepath.Base(got), "receipt-") {
		t.Fatalf("unexpected receipt name: %q", got)
	}
}

func TestWriteEmptyDestErrors(t *testing.T) {
	if _, err := Write("", Receipt{}); err == nil {
		t.Fatal("expected error for empty dest")
	}
}
