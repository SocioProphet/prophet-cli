package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type GlobalFlags struct {
	Profile string
	Space   string
	Output  string
	Query   string
	Quiet   bool
	Debug   bool
	NoPager bool
}

var flags GlobalFlags

func Execute() {
	if err := NewRootCommand().Execute(); err != nil {
		os.Exit(1)
	}
}

func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "prophet",
		Short: "Prophet façade CLI for the SourceOS platform",
		Long: `prophet is a command-line façade for the Prophet / SourceOS platform.

It provides deterministic CLI verbs for bootstrap operations, agent-to-agent
workflow orchestration, and agent-assist scaffolds, while delegating all
sensitive bootstrap business logic to the sourceos-bootstrap engine.

Use --output json for machine-readable output from any command.`,
		SilenceUsage: true,
	}
	root.PersistentFlags().StringVar(&flags.Profile, "profile", "", "named profile")
	root.PersistentFlags().StringVar(&flags.Space, "space", "", "execution space")
	root.PersistentFlags().StringVarP(&flags.Output, "output", "o", "text", "output format")
	root.PersistentFlags().StringVar(&flags.Query, "query", "", "query expression")
	root.PersistentFlags().BoolVarP(&flags.Quiet, "quiet", "q", false, "suppress non-essential output")
	root.PersistentFlags().BoolVar(&flags.Debug, "debug", false, "enable debug output")
	root.PersistentFlags().BoolVar(&flags.NoPager, "no-pager", false, "disable pager")
	root.AddCommand(
		newBootstrapCmd(),
		newA2ACmd(),
		newPlaceholderCmd("ask", "Agent assist: explain or inspect without mutating state",
			"Read-only agent-assist command for explaining or inspecting the system\nwithout mutating any state. Resolves to deterministic tools.\n\nThis command is a scaffold in phase 1.",
			"  prophet ask\n  prophet ask --output json"),
		newPlaceholderCmd("plan", "Agent assist: generate a plan over deterministic tools",
			"Read-only agent-assist command that generates a structured plan over\ndeterministic tools. Does not execute any mutations.\n\nThis command is a scaffold in phase 1.",
			"  prophet plan\n  prophet plan --output json"),
		newPlaceholderCmd("agent", "Agent execute façade (approval-gated scaffold)",
			"Approval-gated agent-execute façade. Mutations require explicit approval\nbefore they are dispatched to deterministic tools.\n\nThis command is a scaffold in phase 1.",
			"  prophet agent\n  prophet agent --output json"),
		newPlaceholderCmd("mcp", "MCP boundary façade",
			"Interface to the Model Context Protocol boundary for agent integrations.\nProject MCP servers are never auto-enabled in untrusted workspaces.\n\nThis command is a scaffold in phase 1.",
			"  prophet mcp\n  prophet mcp --output json"),
	)
	return root
}

func emit(v any) error {
	switch strings.ToLower(flags.Output) {
	case "none":
		return nil
	case "json", "yaml", "table", "tsv":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(v)
	default:
		b, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(b))
		return nil
	}
}

func newPlaceholderCmd(name, summary, long, example string) *cobra.Command {
	return &cobra.Command{
		Use:     name,
		Short:   summary,
		Long:    long,
		Example: example,
		RunE: func(cmd *cobra.Command, args []string) error {
			return emit(map[string]any{"command": name, "summary": summary, "status": "scaffold"})
		},
	}
}
