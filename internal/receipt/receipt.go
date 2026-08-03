// Package receipt emits small, command-scoped, machine-readable records for
// delegated façade actions. These receipts are façade-local: they align
// conceptually with the estate ProofArtifact idea (a signed, replayable account
// of an action) but deliberately do NOT import or implement that spine. They are
// operator/CI convenience only and are not a runtime ledger.
package receipt

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Schema identifies the façade-local receipt shape. Bump on breaking changes.
const Schema = "prophet-cli/receipt/v0"

// Receipt is one delegated-action record.
type Receipt struct {
	Schema     string   `json:"schema"`
	Command    string   `json:"command"`
	Delegate   string   `json:"delegate,omitempty"`
	Status     string   `json:"status"`
	Args       []string `json:"args,omitempty"`
	StartedAt  string   `json:"started_at"`
	FinishedAt string   `json:"finished_at"`
	DurationMs int64    `json:"duration_ms"`
	Error      string   `json:"error,omitempty"`
}

// New builds a Receipt from a delegated action's inputs and outcome.
func New(command, delegate, status string, args []string, started, finished time.Time, execErr error) Receipt {
	r := Receipt{
		Schema:     Schema,
		Command:    command,
		Delegate:   delegate,
		Status:     status,
		Args:       args,
		StartedAt:  started.UTC().Format(time.RFC3339Nano),
		FinishedAt: finished.UTC().Format(time.RFC3339Nano),
		DurationMs: finished.Sub(started).Milliseconds(),
	}
	if execErr != nil {
		r.Error = execErr.Error()
	}
	return r
}

// Write persists r to dest and returns the path written.
//
// If dest is an existing directory, or ends in a path separator, a timestamped
// file is created inside it. Otherwise dest is treated as a file path and its
// parent directory is created as needed. This keeps the flag ergonomic for both
// "one file per run" and "drop receipts in this dir" usage.
func Write(dest string, r Receipt) (string, error) {
	if dest == "" {
		return "", fmt.Errorf("receipt: empty destination")
	}

	target := dest
	asDir := false
	if os.IsPathSeparator(dest[len(dest)-1]) {
		asDir = true
	} else if info, err := os.Stat(dest); err == nil && info.IsDir() {
		asDir = true
	}
	if asDir {
		name := fmt.Sprintf("receipt-%s.json", time.Now().UTC().Format("20060102T150405.000000000Z"))
		target = filepath.Join(dest, name)
	}

	if dir := filepath.Dir(target); dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return "", fmt.Errorf("receipt: create dir %q: %w", dir, err)
		}
	}

	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", fmt.Errorf("receipt: marshal: %w", err)
	}
	b = append(b, '\n')
	if err := os.WriteFile(target, b, 0o644); err != nil {
		return "", fmt.Errorf("receipt: write %q: %w", target, err)
	}
	return target, nil
}
