package cmd

import (
	"github.com/socioprophet/prophet-cli/internal/controlnode"
	"github.com/spf13/cobra"
)

func newControlNodeCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "control-node", Short: "Local-first control-node façade"}
	cmd.AddCommand(
		&cobra.Command{Use: "status", Short: "Show current control-node lane status", RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := controlnode.Status(cmd.Context())
			if err != nil { return err }
			resp["command"] = "prophet control-node status"
			return emit(resp)
		}},
		&cobra.Command{Use: "process <input.json> <outdir>", Short: "Process a local control-node handoff through agentplane", Args: cobra.ExactArgs(2), RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := controlnode.Process(cmd.Context(), args[0], args[1])
			if err != nil { return err }
			resp["command"] = "prophet control-node process"
			return emit(resp)
		}},
	)
	return cmd
}
