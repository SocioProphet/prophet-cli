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
			if resp != nil {
				resp["command"] = "prophet vocab fetch"
				_ = emit(resp)
			}
			return err
		}},
		&cobra.Command{Use: "validate [graph-path ...]", Short: "Run Ontogenesis semantic-core validation", Args: cobra.ArbitraryArgs, RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := vocab.Validate(cmd.Context(), args)
			if resp != nil {
				resp["command"] = "prophet vocab validate"
				_ = emit(resp)
			}
			return err
		}},
		&cobra.Command{Use: "promote", Short: "Promote the current validated context set", RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := vocab.Promote(cmd.Context())
			if resp != nil {
				resp["command"] = "prophet vocab promote"
				_ = emit(resp)
			}
			return err
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
			if resp != nil {
				resp["command"] = "prophet vocab sr run"
				_ = emit(resp)
			}
			return err
		}},
		&cobra.Command{Use: "gate", Short: "Evaluate SR promotion thresholds", RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := vocab.SRGate(cmd.Context())
			if resp != nil {
				resp["command"] = "prophet vocab sr gate"
				_ = emit(resp)
			}
			return err
		}},
	)
	return cmd
}
