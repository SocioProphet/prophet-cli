package cmd

import (
	"os/exec"
	"sort"

	"github.com/spf13/cobra"
)

// surfaceKind classifies a façade surface against the engine boundary.
type surfaceKind string

const (
	// kindReal executes locally in the façade with no external engine required.
	kindReal surfaceKind = "real"
	// kindDelegating shells out to a named engine binary (or local-dev fallback).
	kindDelegating surfaceKind = "delegating"
	// kindScaffold is a placeholder surface that mutates nothing.
	kindScaffold surfaceKind = "scaffold"
)

// surfaceDescriptor is the façade's own account of one top-level surface. This is
// curated boundary knowledge (the façade is the source of truth for its own shape),
// kept in one place so it stays reviewable and honest. It is the machine twin of
// docs/IMPLEMENTATION_STATUS.md.
type surfaceDescriptor struct {
	Command  string      `json:"command"`
	Kind     surfaceKind `json:"kind"`
	Delegate string      `json:"delegate,omitempty"`
	Note     string      `json:"note"`
}

// facadeSurfaces returns the boundary map for the top-level command surface.
func facadeSurfaces() []surfaceDescriptor {
	return []surfaceDescriptor{
		{"version", kindReal, "", "suite version, façade-local"},
		{"doctor", kindReal, "", "suite readiness checks"},
		{"self-test", kindReal, "", "lightweight surface self-test"},
		{"emit-evidence", kindReal, "", "suite local evidence"},
		{"status", kindReal, "", "façade boundary legibility (this command)"},
		{"bootstrap", kindDelegating, "sourceos-bootstrap", "engine home: sourceos-sdk/cmd/sourceos-bootstrap"},
		{"vocab", kindDelegating, "ontogenesis", "fetch/gate/promote/sr"},
		{"bindings", kindDelegating, "atomic-bindings", "validate"},
		{"k8s", kindDelegating, "k8s-policy", "scheduling checks"},
		{"control-node", kindDelegating, "control-node", "local-first control-node"},
		{"devtools", kindDelegating, "sourceos-devtools", "profile management; not-yet-wired"},
		{"lab", kindDelegating, "sourceos-devtools", "functional ML labs; not-yet-wired"},
		{"sourceos", kindDelegating, "sourceos-installer", "install; carry delegates to sourceos-ai"},
		{"holmes", kindDelegating, "holmes", "analyze/search/graph/govern"},
		{"model", kindDelegating, "model-router", "route; local-dev python fallback"},
		{"guardrail", kindDelegating, "guardrail-fabric", "test; local-dev python fallback"},
		{"ledger", kindDelegating, "model-governance-ledger", "validate/records; local-dev fallback"},
		{"agent", kindDelegating, "agent-registry", "registry list; local-dev record fallback"},
		{"spine", kindDelegating, "spine-gates", "per-repo validate gates (--repo)"},
		{"enrichment", kindDelegating, "enrichment-twin", "corpus/lifecycle/gate"},
		{"a2a", kindScaffold, "", "workflow façade skeleton"},
		{"ask", kindScaffold, "", "agent assist placeholder"},
		{"plan", kindScaffold, "", "agent assist placeholder"},
		{"mcp", kindScaffold, "", "MCP boundary placeholder"},
	}
}

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show façade boundary: delegation targets, engine presence, stubbed vs real",
		Long: "Report the façade's own boundary: which surfaces execute locally (real), " +
			"which delegate to an external engine, and which are scaffolds. For delegating " +
			"surfaces the delegate engine binary is probed on PATH so an operator can see " +
			"what is actually wired vs stubbed. Read-only; no platform semantics.",
		RunE: func(cmd *cobra.Command, args []string) error {
			surfaces := facadeSurfaces()

			// Probe each distinct delegate engine once for presence on PATH.
			engineState := map[string]any{}
			present, missing := 0, 0
			for _, s := range surfaces {
				if s.Delegate == "" {
					continue
				}
				if _, seen := engineState[s.Delegate]; seen {
					continue
				}
				if path, err := exec.LookPath(s.Delegate); err == nil {
					engineState[s.Delegate] = map[string]any{"status": "present", "path": path}
					present++
				} else {
					engineState[s.Delegate] = map[string]any{"status": "missing"}
					missing++
				}
			}

			counts := map[string]int{"real": 0, "delegating": 0, "scaffold": 0}
			for _, s := range surfaces {
				counts[string(s.Kind)]++
			}

			engines := make([]string, 0, len(engineState))
			for name := range engineState {
				engines = append(engines, name)
			}
			sort.Strings(engines)

			return emit(map[string]any{
				"command":  "prophet status",
				"status":   "ok",
				"repo":     "SocioProphet/prophet-cli",
				"role":     "facade",
				"boundary": "Owns command grammar and docs. Delegates engine logic, transport, and receipt/enrollment semantics to engine repos.",
				"surfaces": surfaces,
				"engines":  engineState,
				"summary": map[string]any{
					"surfaces_total":  len(surfaces),
					"real":            counts["real"],
					"delegating":      counts["delegating"],
					"scaffold":        counts["scaffold"],
					"engines_present": present,
					"engines_missing": missing,
				},
			})
		},
	}
}
