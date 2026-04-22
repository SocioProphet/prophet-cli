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
	root := &cobra.Command{Use: "prophet", Short: "Prophet façade CLI", SilenceUsage: true}
	root.PersistentFlags().StringVar(&flags.Profile, "profile", "", "named profile")
	root.PersistentFlags().StringVar(&flags.Space, "space", "", "execution space")
	root.PersistentFlags().StringVarP(&flags.Output, "output", "o", "text", "output format")
	root.PersistentFlags().StringVar(&flags.Query, "query", "", "query expression")
	root.PersistentFlags().BoolVarP(&flags.Quiet, "quiet", "q", false, "suppress non-essential output")
	root.PersistentFlags().BoolVar(&flags.Debug, "debug", false, "enable debug output")
	root.PersistentFlags().BoolVar(&flags.NoPager, "no-pager", false, "disable pager")
	root.AddCommand(
		newBootstrapCmd(),
		newVocabCmd(),
		newControlNodeCmd(),
		newA2ACmd(),
		newPlaceholderCmd("ask", "Agent assist: explain or inspect without mutating state"),
		newPlaceholderCmd("plan", "Agent assist: generate a plan over deterministic tools"),
		newPlaceholderCmd("agent", "Agent execute façade (approval-gated scaffold)"),
		newPlaceholderCmd("mcp", "MCP boundary façade"),
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

func newPlaceholderCmd(name, summary string) *cobra.Command {
	return &cobra.Command{Use: name, Short: summary, RunE: func(cmd *cobra.Command, args []string) error {
		return emit(map[string]any{"command": name, "summary": summary, "status": "scaffold"})
	}}
}
