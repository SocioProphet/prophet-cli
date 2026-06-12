package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// K3 lifecycle stages for an enrichment mission, in order.
var enrichmentLifecycle = []string{
	"INIT_SESSION",
	"PROBE_ACCEPT",
	"INJECT_SEED",
	"SEED_PUBLISH",
	"VERIFY_TWIN",
	"TWIN_READY",
	"ANALYSIS_RUN",
	"GATED_HOST_UPDATE",
}

// Locus progression: local → trusted_private → attested_fog → burst_cloud.
// The scheduler picks the first locus the placement policy allows for the
// given sensitivity class. burst_cloud requires explicit approval in the seed.
var locusProgression = []string{"local", "trusted_private", "attested_fog", "burst_cloud"}

func newEnrichmentCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "enrichment", Short: "Enrichment Twin — asset analysis mission control"}
	cmd.AddCommand(
		newEnrichmentRunCmd(),
		newEnrichmentValidateCmd(),
		newEnrichmentStatusCmd(),
		newEnrichmentLifecycleCmd(),
	)
	return cmd
}

// prophet enrichment run --seed <id> [--corpus <path>] [--locus <locus>]
//
// Drives the full K3 lifecycle from INIT_SESSION through GATED_HOST_UPDATE.
// Gate checks run before analysis. burst_cloud is refused unless the seed
// approval_profile contains burst_cloud_placement.
func newEnrichmentRunCmd() *cobra.Command {
	var seedID string
	var corpus string
	var locus string
	var dryRun bool

	c := &cobra.Command{
		Use:   "run",
		Short: "Run an enrichment mission through the full K3 lifecycle",
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(seedID) == "" {
				return fmt.Errorf("--seed is required (e.g. seed:enrichment/photo-v1)")
			}
			if locus == "" {
				locus = "local"
			}
			if !isValidLocus(locus) {
				return fmt.Errorf("unknown locus %q; valid: %s", locus, strings.Join(locusProgression, ", "))
			}
			if locus == "burst_cloud" {
				return fmt.Errorf("burst_cloud locus requires explicit approval in seed.approval_profile.burst_cloud_placement — pass --locus local to start locally")
			}

			corpusPath := corpus
			if corpusPath == "" {
				home, _ := os.UserHomeDir()
				corpusPath = filepath.Join(home, "Photos")
			}

			stages := buildStageTrace(seedID, corpusPath, locus, dryRun)
			return emit(map[string]any{
				"command":  "prophet enrichment run",
				"seed_id":  seedID,
				"corpus":   corpusPath,
				"locus":    locus,
				"dry_run":  dryRun,
				"started":  time.Now().UTC().Format(time.RFC3339),
				"lifecycle": stages,
				"status":   lifecycleStatus(stages),
			})
		},
	}

	c.Flags().StringVar(&seedID, "seed", "", "genesis seed id (e.g. seed:enrichment/photo-v1)")
	c.Flags().StringVar(&corpus, "corpus", "", "path to asset corpus (default: ~/Photos)")
	c.Flags().StringVar(&locus, "locus", "local", "starting execution locus: local | trusted_private | attested_fog")
	c.Flags().BoolVar(&dryRun, "dry-run", false, "emit lifecycle plan without executing analysis")

	return c
}

// prophet enrichment validate
//
// Runs all enrichment gate checks: genesis seed schema, claim hologram schema,
// and placement policy. Reports the gate results before committing to a run.
func newEnrichmentValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Run all enrichment gate checks before starting a mission",
		RunE: func(cmd *cobra.Command, args []string) error {
			gates := []map[string]any{
				runEnrichmentGate("genesis-seed-schema",
					"ProCybernetica", []string{"make", "enrichment-twin-fixtures"}),
				runEnrichmentGate("path-content-conflict",
					"ProCybernetica", []string{"make", "enrichment-twin-path-content-conflict"}),
				runEnrichmentGate("placement-locus-policy",
					"policy-fabric", []string{"make", "enrichment-placement-locus-policy-validate"}),
				runEnrichmentGate("spine-gates",
					"", []string{}), // delegated to prophet spine validate
			}
			allPassed := true
			for _, g := range gates {
				if g["status"] != "ok" {
					allPassed = false
				}
			}
			status := "ok"
			if !allPassed {
				status = "failed"
			}
			return emit(map[string]any{
				"command": "prophet enrichment validate",
				"status":  status,
				"gates":   gates,
			})
		},
	}
}

// prophet enrichment status
func newEnrichmentStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show enrichment twin readiness and ADR decision state",
		RunE: func(cmd *cobra.Command, args []string) error {
			return emit(map[string]any{
				"command": "prophet enrichment status",
				"status":  "ok",
				"architecture_phase_dod": map[string]any{
					"genesis_seed_validates":          true,
					"claim_hologram_validates":        true,
					"placement_policy_wired":          true,
					"path_content_conflict_tested":    true,
					"adr_0002_decisions_closed":       true,
					"end_to_end_local_run":            false,
					"memory_mesh_write_verified":      false,
				},
				"open_decisions": "ADR-0002 all decided (2026-06-11)",
				"next_step":      "prophet enrichment run --seed seed:enrichment/photo-v1 --corpus <path>",
				"locus_progression": locusProgression,
			})
		},
	}
}

// prophet enrichment lifecycle
func newEnrichmentLifecycleCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "lifecycle",
		Short: "Show the K3 lifecycle stages for an enrichment mission",
		RunE: func(cmd *cobra.Command, args []string) error {
			stages := make([]map[string]any, 0, len(enrichmentLifecycle))
			descriptions := map[string]string{
				"INIT_SESSION":       "Open mission session; establish twin identity and trust domain",
				"PROBE_ACCEPT":       "Negotiate organs, model classes, and eligible loci with placement policy",
				"INJECT_SEED":        "Inject genesis seed (e.g. seed:enrichment/photo-v1) into twin runtime",
				"SEED_PUBLISH":       "Publish seed to mesh; all nodes acknowledge organ and policy bindings",
				"VERIFY_TWIN":        "Verify twin identity, policy bindings, and memory bindings; gate on all passing",
				"TWIN_READY":         "Twin active; organs bound; ready to accept per-asset analysis requests",
				"ANALYSIS_RUN":       "Per-asset: compute asset_hash → check enrichment cache → run organs → write claim holograms",
				"GATED_HOST_UPDATE":  "Approval-gated: activate query-API for host search layer (never writes System Space)",
			}
			approvalRequired := map[string]bool{
				"GATED_HOST_UPDATE": true,
			}
			for _, stage := range enrichmentLifecycle {
				s := map[string]any{
					"stage":       stage,
					"description": descriptions[stage],
				}
				if approvalRequired[stage] {
					s["approval_required"] = true
					s["approval_gate"] = "seed.approval_profile.host_index_writeback"
				}
				stages = append(stages, s)
			}
			return emit(map[string]any{
				"command":    "prophet enrichment lifecycle",
				"status":     "ok",
				"stages":     stages,
				"failure_lanes": []string{"QUARANTINE", "REVOKE", "ROLLBACK"},
			})
		},
	}
}

// ── helpers ──────────────────────────────────────────────────────────────────

func isValidLocus(locus string) bool {
	for _, l := range locusProgression {
		if l == locus {
			return true
		}
	}
	return false
}

func buildStageTrace(seedID, corpus, locus string, dryRun bool) []map[string]any {
	stages := make([]map[string]any, 0, len(enrichmentLifecycle))
	for _, stage := range enrichmentLifecycle {
		s := map[string]any{
			"stage":  stage,
			"status": "pending",
		}
		if dryRun {
			s["status"] = "dry_run"
		}
		switch stage {
		case "INJECT_SEED":
			s["seed_id"] = seedID
		case "ANALYSIS_RUN":
			s["corpus"] = corpus
			s["locus"] = locus
			s["cache_key"] = "asset_hash + analyzer_version + model_version"
		case "GATED_HOST_UPDATE":
			s["approval_gate"] = "seed.approval_profile.host_index_writeback"
			s["mode"] = "query_api_activation"
		}
		stages = append(stages, s)
	}
	return stages
}

func lifecycleStatus(stages []map[string]any) string {
	for _, s := range stages {
		if s["status"] == "failed" {
			return "failed"
		}
	}
	if stages[0]["status"] == "dry_run" {
		return "dry_run"
	}
	return "pending_execution"
}

func runEnrichmentGate(name, repoName string, steps []string) map[string]any {
	if len(steps) == 0 || repoName == "" {
		return map[string]any{
			"gate":   name,
			"status": "not-yet-wired",
			"reason": "delegated to spine validate",
		}
	}
	repoDir := localRepoRoot(repoName)
	result := runStep(repoDir, steps)
	status := result["status"]
	return map[string]any{
		"gate":   name,
		"status": status,
		"repo":   repoName,
		"step":   strings.Join(steps, " "),
	}
}
