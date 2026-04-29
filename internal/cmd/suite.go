package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func newSuiteDoctorCmd() *cobra.Command {
	return &cobra.Command{Use: "doctor", Short: "Check Prophet suite readiness", RunE: func(cmd *cobra.Command, args []string) error {
		tools := []string{"sourceos-ai", "holmes", "model-router", "guardrail-fabric", "model-governance-ledger", "agent-registry"}
		checks := make([]map[string]any, 0, len(tools))
		for _, tool := range tools {
			checks = append(checks, toolCheck(tool))
		}
		return emit(map[string]any{"command": "prophet doctor", "status": "ok", "checks": checks})
	}}
}

func newSuiteVersionCmd() *cobra.Command {
	return &cobra.Command{Use: "version", Short: "Show Prophet suite version", RunE: func(cmd *cobra.Command, args []string) error {
		return emit(map[string]any{"command": "prophet version", "tool": "prophet", "status": "ok", "version": "0.1.0-dev"})
	}}
}

func newSuiteSelfTestCmd() *cobra.Command {
	return &cobra.Command{Use: "self-test", Short: "Run lightweight Prophet suite self-test", RunE: func(cmd *cobra.Command, args []string) error {
		return emit(map[string]any{"command": "prophet self-test", "status": "ok", "checks": []string{"command-surface"}})
	}}
}

func newSuiteEvidenceCmd() *cobra.Command {
	return &cobra.Command{Use: "emit-evidence", Short: "Emit Prophet suite local evidence", RunE: func(cmd *cobra.Command, args []string) error {
		return emit(map[string]any{"command": "prophet emit-evidence", "status": "ok", "repo": "SocioProphet/prophet-cli", "surface": "suite-facade"})
	}}
}

func newDevtoolsCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "devtools", Short: "SourceOS developer tools facade"}
	profile := &cobra.Command{Use: "profile", Short: "Manage SourceOS devtools profiles"}
	profile.AddCommand(
		&cobra.Command{Use: "list", Short: "List devtools profiles", RunE: func(cmd *cobra.Command, args []string) error {
			return emit(map[string]any{"command": "prophet devtools profile list", "status": "ok", "profiles": []string{"core", "build", "containers", "k8s", "ai-core", "labs", "security"}})
		}},
		&cobra.Command{Use: "apply <profiles>", Short: "Apply devtools profiles", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
			return emit(map[string]any{"command": "prophet devtools profile apply", "status": "not-yet-wired", "profiles": strings.Split(args[0], ","), "delegate": "sourceos-devtools"})
		}},
		&cobra.Command{Use: "current", Short: "Show current devtools profiles", RunE: func(cmd *cobra.Command, args []string) error {
			return emit(map[string]any{"command": "prophet devtools profile current", "status": "not-yet-wired", "delegate": "sourceos-devtools"})
		}},
	)
	cmd.AddCommand(profile)
	cmd.AddCommand(
		&cobra.Command{Use: "doctor", Short: "Check devtools readiness", RunE: func(cmd *cobra.Command, args []string) error {
			return emit(map[string]any{"command": "prophet devtools doctor", "status": "not-yet-wired", "delegate": "sourceos-devtools"})
		}},
		&cobra.Command{Use: "emit-evidence", Short: "Emit devtools evidence", RunE: func(cmd *cobra.Command, args []string) error {
			return emit(map[string]any{"command": "prophet devtools emit-evidence", "status": "not-yet-wired", "delegate": "sourceos-devtools"})
		}},
	)
	return cmd
}

func newLabCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "lab", Short: "Functional ML lab facade"}
	labs := []string{"nlplab", "translationlab", "embeddinglab", "speechlab", "ocrlab", "imagelab", "videolab", "timeserieslab", "graphlab"}
	cmd.AddCommand(
		&cobra.Command{Use: "list", Short: "List functional labs", RunE: func(cmd *cobra.Command, args []string) error {
			return emit(map[string]any{"command": "prophet lab list", "status": "ok", "labs": labs})
		}},
		&cobra.Command{Use: "enable <labs...>", Short: "Enable functional labs", Args: cobra.MinimumNArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
			return emit(map[string]any{"command": "prophet lab enable", "status": "not-yet-wired", "labs": args, "delegate": "sourceos-devtools"})
		}},
		&cobra.Command{Use: "disable <labs...>", Short: "Disable functional labs", Args: cobra.MinimumNArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
			return emit(map[string]any{"command": "prophet lab disable", "status": "not-yet-wired", "labs": args, "delegate": "sourceos-devtools"})
		}},
		&cobra.Command{Use: "status", Short: "Show lab status", RunE: func(cmd *cobra.Command, args []string) error {
			return emit(map[string]any{"command": "prophet lab status", "status": "not-yet-wired", "delegate": "sourceos-devtools"})
		}},
		&cobra.Command{Use: "doctor", Short: "Check lab readiness", RunE: func(cmd *cobra.Command, args []string) error {
			return emit(map[string]any{"command": "prophet lab doctor", "status": "not-yet-wired", "delegate": "sourceos-devtools"})
		}},
	)
	return cmd
}

func newSourceOSCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "sourceos", Short: "SourceOS installer and carry facade"}
	installCmd := &cobra.Command{Use: "install", Short: "Run SourceOS install facade", RunE: func(cmd *cobra.Command, args []string) error {
		target, _ := cmd.Flags().GetString("target")
		channel, _ := cmd.Flags().GetString("channel")
		return emit(map[string]any{"command": "prophet sourceos install", "status": "not-yet-wired", "target": target, "channel": channel, "delegate": "sourceos-installer"})
	}}
	installCmd.Flags().String("target", "", "target platform")
	installCmd.Flags().String("channel", "dev", "release channel")
	carry := &cobra.Command{Use: "carry", Short: "SourceOS AI carry facade"}
	carry.AddCommand(
		delegateCommand("list", "List SourceOS carry refs", "sourceos-ai", []string{"carry", "list"}),
		delegateCommand("validate", "Validate SourceOS carry refs", "sourceos-ai", []string{"carry", "validate"}),
		delegateCommand("doctor", "Check SourceOS carry readiness", "sourceos-ai", []string{"carry", "doctor"}),
		delegateCommand("emit-evidence", "Emit SourceOS carry evidence", "sourceos-ai", []string{"emit-evidence"}),
	)
	cmd.AddCommand(installCmd, carry)
	return cmd
}

func newHolmesCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "holmes", Short: "Holmes language intelligence facade"}
	cmd.AddCommand(
		delegateArgsCommand("analyze <path>", "Analyze a document", "holmes", []string{"analyze"}, cobra.ExactArgs(1)),
		delegateArgsCommand("search <query>", "Search with Holmes/Sherlock", "holmes", []string{"search"}, cobra.MinimumNArgs(1)),
		delegateArgsCommand("graph <path>", "Build semantic graph", "holmes", []string{"graph"}, cobra.ExactArgs(1)),
		delegateArgsCommand("govern <path>", "Evaluate language governance", "holmes", []string{"govern"}, cobra.ExactArgs(1)),
	)
	return cmd
}

func newModelCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "model", Short: "Model routing facade"}
	route := &cobra.Command{Use: "route", Short: "Route a model/service request", RunE: func(cmd *cobra.Command, args []string) error {
		task, _ := cmd.Flags().GetString("task")
		privacy, _ := cmd.Flags().GetString("privacy")
		return delegateOrFallback("model-router", []string{"route", "--task", task, "--privacy", privacy}, "prophet model route", fabricFallback{Repo: "model-router", Script: "tools/model_router.py", Args: []string{"emit-demo-decision"}})
	}}
	route.Flags().String("task", "", "task name")
	route.Flags().String("privacy", "standard", "privacy policy")
	cmd.AddCommand(route)
	return cmd
}

func newGuardrailCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "guardrail", Short: "Guardrail fabric facade"}
	cmd.AddCommand(&cobra.Command{Use: "test <policy.json> <input.json>", Short: "Test guardrail policy", Args: cobra.ExactArgs(2), RunE: func(cmd *cobra.Command, args []string) error {
		return delegateOrFallback("guardrail-fabric", append([]string{"test"}, args...), "prophet guardrail test", fabricFallback{Repo: "guardrail-fabric", Script: "tools/guardrail_fabric.py", Args: []string{"emit-demo-decision"}})
	}})
	return cmd
}

func newAgentSuiteCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "agent", Short: "Agent registry and execution facade"}
	registry := &cobra.Command{Use: "registry", Short: "Agent registry facade"}
	registry.AddCommand(&cobra.Command{Use: "list", Short: "List agent registry entries", RunE: func(cmd *cobra.Command, args []string) error {
		return delegateOrFallback("agent-registry", []string{"list"}, "prophet agent registry list", fabricFallback{Repo: "agent-registry", RecordGlob: "examples/*.json"})
	}})
	cmd.AddCommand(registry)
	return cmd
}

func delegateCommand(use, short, tool string, baseArgs []string) *cobra.Command {
	return &cobra.Command{Use: use, Short: short, RunE: func(cmd *cobra.Command, args []string) error {
		return delegateOrReport(tool, append(baseArgs, args...), "prophet "+cmd.CommandPath())
	}}
}

func delegateArgsCommand(use, short, tool string, baseArgs []string, argRule cobra.PositionalArgs) *cobra.Command {
	return &cobra.Command{Use: use, Short: short, Args: argRule, RunE: func(cmd *cobra.Command, args []string) error {
		return delegateOrReport(tool, append(baseArgs, args...), "prophet "+cmd.CommandPath())
	}}
}

func delegateOrReport(tool string, args []string, command string) error {
	path, err := exec.LookPath(tool)
	if err != nil {
		return emit(map[string]any{"command": command, "status": "not-yet-wired", "delegate": tool, "reason": "delegate tool not found"})
	}
	run := exec.Command(path, args...)
	var stdout, stderr bytes.Buffer
	run.Stdout = &stdout
	run.Stderr = &stderr
	if err := run.Run(); err != nil {
		return emit(map[string]any{"command": command, "status": "failed", "delegate": tool, "stdout": stdout.String(), "stderr": stderr.String(), "error": err.Error()})
	}
	return emit(map[string]any{"command": command, "status": "ok", "delegate": tool, "stdout": stdout.String(), "stderr": stderr.String()})
}

func toolCheck(tool string) map[string]any {
	path, err := exec.LookPath(tool)
	if err != nil {
		return map[string]any{"tool": tool, "status": "missing"}
	}
	return map[string]any{"tool": tool, "status": "present", "path": path}
}

func commandRequiredFlag(name string, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("missing required --%s", name)
	}
	return nil
}
