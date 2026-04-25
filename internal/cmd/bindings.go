package cmd

import (
	"github.com/socioprophet/prophet-cli/internal/bindings"
	"github.com/spf13/cobra"
)

func newBindingsCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "bindings", Short: "Atomic bindings façade"}
	validate := &cobra.Command{Use: "validate", Short: "Validate binding descriptors"}
	var abdPath string
	abd := &cobra.Command{Use: "abd", Short: "Validate an Agent Binding Descriptor", RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := bindings.ValidateABD(cmd.Context(), abdPath)
		if resp != nil {
			resp["command"] = "prophet bindings validate abd"
			_ = emit(resp)
		}
		return err
	}}
	abd.Flags().StringVar(&abdPath, "abd", "", "path to ABD JSON file relative to the standards repo")
	_ = abd.MarkFlagRequired("abd")
	validate.AddCommand(abd)
	cmd.AddCommand(validate)
	return cmd
}
