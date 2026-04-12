package cmd

import "github.com/spf13/cobra"

func newVocabCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "vocab", Short: "Ontogenesis vocabulary and policy-pack façade"}
	cmd.AddCommand(
		newVocabLeaf("fetch", "fetch or refresh Ontogenesis semantic assets"),
		&cobra.Command{Use: "validate [graph-path ...]", Short: "Run Ontogenesis semantic-core validation", Args: cobra.ArbitraryArgs, RunE: func(cmd *cobra.Command, args []string) error {
			return emit(map[string]any{
				"command":      "prophet vocab validate",
				"graph_paths":  args,
				"delegated_to": "socioprophet-standards-knowledge",
				"validator":    "policy/tools/validate_all.py",
				"status":       "scaffold",
			})
		}},
		&cobra.Command{Use: "promote", Short: "Promote the current validated context set", RunE: func(cmd *cobra.Command, args []string) error {
			return emit(map[string]any{
				"command":       "prophet vocab promote",
				"delegated_to":  "socioprophet-standards-knowledge",
				"record_target": "/ontology/promotions/<ts>.jsonld",
				"status":        "scaffold",
			})
		}},
		newVocabSRCmd(),
	)
	return cmd
}

func newVocabSRCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "sr", Short: "Symbolic regression façade over Ontogenesis runners"}
	cmd.AddCommand(
		&cobra.Command{Use: "run <module>", Short: "Extract, train, and register SR for a module", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
			return emit(map[string]any{
				"command":      "prophet vocab sr run",
				"module":       args[0],
				"delegated_to": "socioprophet-standards-knowledge",
				"status":       "scaffold",
			})
		}},
		&cobra.Command{Use: "gate", Short: "Evaluate SR promotion thresholds", RunE: func(cmd *cobra.Command, args []string) error {
			return emit(map[string]any{
				"command":      "prophet vocab sr gate",
				"delegated_to": "socioprophet-standards-knowledge",
				"status":       "scaffold",
			})
		}},
	)
	return cmd
}

func newVocabLeaf(name, summary string) *cobra.Command {
	return &cobra.Command{Use: name, Short: summary, RunE: func(cmd *cobra.Command, args []string) error {
		return emit(map[string]any{"command": "prophet vocab " + name, "summary": summary, "delegated_to": "socioprophet-standards-knowledge", "status": "scaffold"})
	}}
}
