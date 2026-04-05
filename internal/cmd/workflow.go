package cmd

import "github.com/spf13/cobra"

func newA2ACmd() *cobra.Command {
	cmd := &cobra.Command{}
	cmd.Use = "a2a"
	cmd.Short = "Workflow facade"
	return cmd
}
