package cmd

import (
	"github.com/socioprophet/prophet-cli/internal/vocab"
	"github.com/spf13/cobra"
)

func newVocabCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "vocab", Short: "Ontogenesis vocabulary and policy-pack façade"}
	cmd.AddCommand(
		&cobra.Command{Use: "fetch", Short: "fetch or refresh Ontogenesis semantic assets", RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := vocab.Fetch(cmd.Context())
			if err != nil { return err }
			resp["command"] = "prophet vocab fetch"
			return emit(resp)
		}},
		&cobra.Command{Use: "validate [graph-path ...]", Short: "Run Ontogenesis semantic-core validation", Args: cobra.ArbitraryArgs, RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := vocab.Validate(cmd.Context(), args)
			if err != nil { return err }
			resp["command"] = "prophet vocab validate"
			return emit(resp)
		}},
		&cobra.Command{Use: "promote", Short: "Promote the current validated context set", RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := vocab.Promote(cmd.Context())
			if err != nil { return err }
			resp["command"] = "prophet vocab promote"
			return emit(resp)
		}},
		newVocabSRCmd(),
	)
	return cmd
}

func newVocabSRCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "sr", Short: "Symbolic regression façade over Ontogenesis runners"}
	cmd.AddCommand(
		&cobra.Command{Use: "run <module>", Short: "Extract, train, and register SR for a module", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := vocab.SRRun(cmd.Context(), args[0])
			if err != nil { return err }
			resp["command"] = "prophet vocab sr run"
			return emit(resp)
		}},
		&cobra.Command{Use: "gate", Short: "Evaluate SR promotion thresholds", RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := vocab.SRGate(cmd.Context())
			if err != nil { return err }
			resp["command"] = "prophet vocab sr gate"
			return emit(resp)
		}},
	)
	return cmd
}
