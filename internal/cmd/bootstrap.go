package cmd

import "github.com/spf13/cobra"

func newBootstrapCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bootstrap",
		Short: "Bootstrap façade commands",
		Long: `Façade over the SourceOS bootstrap engine.

All subcommands delegate to the installed sourceos-bootstrap engine.
This CLI owns the command surface, flag grammar, and output formatting.
Bootstrap business logic remains in sourceos-bootstrap.`,
	}
	cmd.AddCommand(
		newBootstrapLeaf("doctor", "diagnose host prerequisites",
			"Diagnose host prerequisites. Checks that the required tools and\nenvironment configuration are in place before a bootstrap session.",
			"  prophet bootstrap doctor\n  prophet bootstrap doctor --output json"),
		newBootstrapLeaf("login", "prepare authenticated bootstrap session",
			"Prepare an authenticated bootstrap session. Establishes credentials\nthat subsequent bootstrap commands will use.",
			"  prophet bootstrap login\n  prophet bootstrap login --profile myprofile"),
		newBootstrapLeaf("build", "submit or describe a bootstrap build request",
			"Submit or describe a bootstrap build request. Delegates build\norchestration to the sourceos-bootstrap engine.",
			"  prophet bootstrap build\n  prophet bootstrap build --output json"),
		newBootstrapLeaf("fetch", "fetch release artifacts",
			"Fetch release artifacts from the bootstrap engine's artifact store.",
			"  prophet bootstrap fetch\n  prophet bootstrap fetch --output json"),
		newBootstrapLeaf("write", "prepare install or recovery media",
			"Prepare install or recovery media using fetched artifacts.",
			"  prophet bootstrap write"),
		newBootstrapLeaf("info", "show bootstrap engine information",
			"Show version and configuration information for the installed\nbootstrap engine.",
			"  prophet bootstrap info\n  prophet bootstrap info --output json"),
		&cobra.Command{
			Use:   "validate <kind> <path>",
			Short: "Validate a frozen object through the bootstrap engine",
			Long: `Validate a frozen object through the bootstrap engine.

<kind> is the object type (e.g. ReleaseSet, ConfigSource) and <path>
is the file system path to the object. Validation logic is owned by
sourceos-bootstrap; this command is a transparent wrapper.`,
			Example: "  prophet bootstrap validate ReleaseSet ./releaseset.json\n  prophet bootstrap validate ConfigSource ./config.yaml --output json",
			Args:    cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				return emit(map[string]any{"command": "prophet bootstrap validate", "kind": args[0], "path": args[1], "delegated_to": "sourceos-bootstrap", "status": "scaffold"})
			},
		},
		&cobra.Command{
			Use:   "verify <kind> <path>",
			Short: "Verify a frozen object or artifact through the bootstrap engine",
			Long: `Cryptographically verify a frozen object or artifact through the
bootstrap engine.

<kind> is the object type (e.g. ReleaseSet, BootReleaseSet) and <path>
is the file system path to the object. Verification logic is owned by
sourceos-bootstrap; this command is a transparent wrapper.`,
			Example: "  prophet bootstrap verify ReleaseSet ./releaseset.json\n  prophet bootstrap verify BootReleaseSet ./boot.json --output json",
			Args:    cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				return emit(map[string]any{"command": "prophet bootstrap verify", "kind": args[0], "path": args[1], "delegated_to": "sourceos-bootstrap", "status": "scaffold"})
			},
		},
	)
	return cmd
}

func newBootstrapLeaf(name, summary, long, example string) *cobra.Command {
	return &cobra.Command{
		Use:     name,
		Short:   summary,
		Long:    long,
		Example: example,
		RunE: func(cmd *cobra.Command, args []string) error {
			return emit(map[string]any{"command": "prophet bootstrap " + name, "summary": summary, "delegated_to": "sourceos-bootstrap", "status": "scaffold"})
		},
	}
}
