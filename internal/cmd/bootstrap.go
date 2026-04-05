package cmd

import "github.com/spf13/cobra"

func newBootstrapCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "bootstrap", Short: "Bootstrap façade commands"}
	cmd.AddCommand(
		newBootstrapLeaf("doctor", "diagnose host prerequisites"),
		newBootstrapLeaf("login", "prepare authenticated bootstrap session"),
		newBootstrapLeaf("build", "submit or describe a bootstrap build request"),
		newBootstrapLeaf("fetch", "fetch release artifacts"),
		newBootstrapLeaf("write", "prepare install or recovery media"),
		newBootstrapLeaf("info", "show bootstrap engine information"),
		&cobra.Command{Use: "validate <kind> <path>", Short: "Validate a frozen object through the bootstrap engine", Args: cobra.ExactArgs(2), RunE: func(cmd *cobra.Command, args []string) error {
			return emit(map[string]any{"command": "prophet bootstrap validate", "kind": args[0], "path": args[1], "delegated_to": "sourceos-bootstrap", "status": "scaffold"})
		}},
		&cobra.Command{Use: "verify <kind> <path>", Short: "Verify a frozen object or artifact through the bootstrap engine", Args: cobra.ExactArgs(2), RunE: func(cmd *cobra.Command, args []string) error {
			return emit(map[string]any{"command": "prophet bootstrap verify", "kind": args[0], "path": args[1], "delegated_to": "sourceos-bootstrap", "status": "scaffold"})
		}},
	)
	return cmd
}

func newBootstrapLeaf(name, summary string) *cobra.Command {
	return &cobra.Command{Use: name, Short: summary, RunE: func(cmd *cobra.Command, args []string) error {
		return emit(map[string]any{"command": "prophet bootstrap " + name, "summary": summary, "delegated_to": "sourceos-bootstrap", "status": "scaffold"})
	}}
}
