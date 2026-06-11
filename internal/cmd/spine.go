package cmd

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type spineGate struct {
	Repo  string
	Steps [][]string // each step is a command + args
}

var spineGates = []spineGate{
	{
		Repo:  "prophet-mesh",
		Steps: [][]string{{"make", "validate"}, {"make", "test"}},
	},
	{
		Repo:  "agent-registry",
		Steps: [][]string{{"make", "validate"}, {"make", "test"}},
	},
	{
		Repo:  "model-router",
		Steps: [][]string{{"make", "validate"}, {"make", "test"}},
	},
	{
		Repo:  "agentplane",
		Steps: [][]string{{"make", "validate"}, {"make", "test"}},
	},
	{
		Repo:  "memory-mesh",
		Steps: [][]string{{"make", "validate-prophet-mesh-scope-mirror"}},
	},
}

func newSpineCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "spine", Short: "Prophet Mesh private-preview spine gates"}

	cmd.AddCommand(newSpineValidateCmd())
	cmd.AddCommand(newSpineGateCmd())
	cmd.AddCommand(newSpineListCmd())

	return cmd
}

func newSpineValidateCmd() *cobra.Command {
	var repo string
	c := &cobra.Command{
		Use:   "validate",
		Short: "Run all spine gate checks (or a single repo gate with --repo)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if repo != "" {
				return runSpineGate(repo)
			}
			return runAllSpineGates()
		},
	}
	c.Flags().StringVar(&repo, "repo", "", "run gate for a single repo only")
	return c
}

func newSpineGateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "gate <repo>",
		Short: "Run the spine gate for a specific repo",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSpineGate(args[0])
		},
	}
}

func newSpineListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List spine repos and their gate steps",
		RunE: func(cmd *cobra.Command, args []string) error {
			gates := make([]map[string]any, 0, len(spineGates))
			for _, g := range spineGates {
				steps := make([]string, 0, len(g.Steps))
				for _, s := range g.Steps {
					steps = append(steps, strings.Join(s, " "))
				}
				gates = append(gates, map[string]any{
					"repo":  g.Repo,
					"steps": steps,
				})
			}
			return emit(map[string]any{
				"command": "prophet spine list",
				"status":  "ok",
				"gates":   gates,
			})
		},
	}
}

func runAllSpineGates() error {
	results := make([]map[string]any, 0, len(spineGates))
	allPassed := true

	for _, g := range spineGates {
		result := runGate(g)
		results = append(results, result)
		if result["status"] != "ok" {
			allPassed = false
		}
	}

	status := "ok"
	if !allPassed {
		status = "failed"
	}

	return emit(map[string]any{
		"command": "prophet spine validate",
		"status":  status,
		"gates":   results,
	})
}

func runSpineGate(repo string) error {
	for _, g := range spineGates {
		if g.Repo == repo {
			result := runGate(g)
			result["command"] = "prophet spine gate " + repo
			return emit(result)
		}
	}
	return emit(map[string]any{
		"command": "prophet spine gate " + repo,
		"status":  "error",
		"reason":  "unknown spine repo: " + repo,
		"known":   spineRepoNames(),
	})
}

func runGate(g spineGate) map[string]any {
	repoDir := localRepoRoot(g.Repo)
	stepResults := make([]map[string]any, 0, len(g.Steps))
	gateStatus := "ok"

	for _, step := range g.Steps {
		sr := runStep(repoDir, step)
		stepResults = append(stepResults, sr)
		if sr["status"] != "ok" {
			gateStatus = "failed"
		}
	}

	return map[string]any{
		"repo":   g.Repo,
		"status": gateStatus,
		"steps":  stepResults,
	}
}

func runStep(dir string, args []string) map[string]any {
	if len(args) == 0 {
		return map[string]any{"status": "error", "reason": "empty step"}
	}

	binary, err := exec.LookPath(args[0])
	if err != nil {
		// try resolving relative to dir (e.g. make)
		binary = args[0]
	}

	c := exec.Command(binary, args[1:]...)
	c.Dir = dir

	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr

	start := time.Now()
	runErr := c.Run()
	elapsed := time.Since(start).Milliseconds()

	step := strings.Join(args, " ")
	if runErr != nil {
		return map[string]any{
			"step":      step,
			"status":    "failed",
			"stdout":    strings.TrimSpace(stdout.String()),
			"stderr":    strings.TrimSpace(stderr.String()),
			"error":     runErr.Error(),
			"elapsed_ms": elapsed,
		}
	}

	result := map[string]any{
		"step":      step,
		"status":    "ok",
		"elapsed_ms": elapsed,
	}
	if out := strings.TrimSpace(stdout.String()); out != "" {
		result["stdout"] = out
	}
	return result
}

func spineRepoNames() []string {
	names := make([]string, 0, len(spineGates))
	for _, g := range spineGates {
		names = append(names, g.Repo)
	}
	return names
}

// spineGateDir resolves a spine repo to its absolute local path.
// Exported for use by Makefile validate target.
func spineGateDir(repo string) string {
	return filepath.Join(localRepoRoot(repo), ".")
}
