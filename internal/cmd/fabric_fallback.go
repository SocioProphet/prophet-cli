package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type fabricFallback struct {
	Repo       string
	Script     string
	Args       []string
	RecordGlob string
}

func delegateOrFallback(tool string, toolArgs []string, command string, fallback fabricFallback) error {
	if path, err := exec.LookPath(tool); err == nil {
		return runDelegate(path, tool, toolArgs, command)
	}
	root := localRepoRoot(fallback.Repo)
	if fallback.Script != "" {
		script := filepath.Join(root, fallback.Script)
		if fileExists(script) {
			python, err := exec.LookPath("python3")
			if err != nil {
				return emit(map[string]any{"command": command, "status": "not-yet-installed", "delegate": tool, "fallbackRepo": fallback.Repo, "fallbackPath": script, "reason": "python3 not found"})
			}
			args := append([]string{script}, fallback.Args...)
			return runDelegate(python, "python3", args, command)
		}
	}
	if fallback.RecordGlob != "" {
		records, err := loadLocalRecords(filepath.Join(root, fallback.RecordGlob))
		if err == nil && len(records) > 0 {
			return emit(map[string]any{"command": command, "status": "ok", "mode": "local-dev-repo", "repo": fallback.Repo, "records": records})
		}
	}
	return emit(map[string]any{"command": command, "status": "not-yet-installed", "delegate": tool, "fallbackRepo": fallback.Repo, "fallbackPath": root, "reason": "delegate tool and local repo fallback not found"})
}

func runDelegate(path string, tool string, args []string, command string) error {
	run := exec.Command(path, args...)
	var stdout, stderr bytes.Buffer
	run.Stdout = &stdout
	run.Stderr = &stderr
	if err := run.Run(); err != nil {
		return emit(map[string]any{"command": command, "status": "failed", "delegate": tool, "stdout": stdout.String(), "stderr": stderr.String(), "error": err.Error()})
	}
	return emit(map[string]any{"command": command, "status": "ok", "delegate": tool, "stdout": stdout.String(), "stderr": stderr.String()})
}

func localRepoRoot(repo string) string {
	base := os.Getenv("PROPHET_DEV_ROOT")
	if strings.TrimSpace(base) == "" {
		if home, err := os.UserHomeDir(); err == nil {
			base = filepath.Join(home, "dev")
		}
	}
	if strings.TrimSpace(base) == "" {
		base = "."
	}
	return filepath.Join(base, repo)
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func loadLocalRecords(pattern string) ([]any, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	records := make([]any, 0, len(matches))
	for _, match := range matches {
		contents, err := os.ReadFile(match)
		if err != nil {
			return nil, err
		}
		var record any
		if err := json.Unmarshal(contents, &record); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}

func newLedgerCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "ledger", Short: "Model governance ledger facade"}
	cmd.AddCommand(
		&cobra.Command{Use: "validate", Short: "Validate model governance ledger records", RunE: func(cmd *cobra.Command, args []string) error {
			return delegateOrFallback("model-governance-ledger", []string{"validate"}, "prophet ledger validate", fabricFallback{Repo: "model-governance-ledger", Script: "tools/validate_ledger_examples.py"})
		}},
		&cobra.Command{Use: "records", Short: "List local model governance ledger records", RunE: func(cmd *cobra.Command, args []string) error {
			return delegateOrFallback("model-governance-ledger", []string{"records"}, "prophet ledger records", fabricFallback{Repo: "model-governance-ledger", RecordGlob: "examples/*.json"})
		}},
	)
	return cmd
}
