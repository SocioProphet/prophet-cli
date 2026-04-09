package cmd

import (
	"fmt"

	"github.com/socioprophet/prophet-cli/internal/a2a"
	"github.com/spf13/cobra"
)

func newA2ACmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "a2a",
		Short: "Agent-to-agent workflow orchestration",
		Long: `Agent-to-agent (A2A) workflow orchestration.

Manages the lifecycle of a structured change workflow across the phases:
  propose → test → review → revise → merge → done

Use 'prophet a2a run' to start or describe a workflow run.`,
	}
	cmd.AddCommand(newA2ARunCmd())
	return cmd
}

func newA2ARunCmd() *cobra.Command {
	var repo, ticket string
	var live bool

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run an A2A workflow for a repository and ticket",
		Long: `Start or describe an agent-to-agent workflow run.

By default this is a dry-run: the workflow plan is printed but nothing
is executed. Pass --live to execute the workflow.

The workflow progresses through six phases:
  1. propose  — author proposes a change
  2. test     — author validates the candidate change
  3. review   — reviewer evaluates and may block or request revision
  4. revise   — author revises after review
  5. merge    — arbiter or maintainer merges the approved change
  6. done     — workflow is complete`,
		Example: `  # Dry-run: print workflow plan
  prophet a2a run --repo owner/repo --ticket TICKET-123

  # Live run
  prophet a2a run --repo owner/repo --ticket TICKET-123 --live

  # JSON output
  prophet a2a run --repo owner/repo --ticket TICKET-123 --output json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if repo == "" {
				return fmt.Errorf("--repo is required")
			}
			if ticket == "" {
				return fmt.Errorf("--ticket is required")
			}
			wf := a2a.Default(repo, ticket, live)
			return emit(wf)
		},
	}
	cmd.Flags().StringVar(&repo, "repo", "", "repository in owner/name format (required)")
	cmd.Flags().StringVar(&ticket, "ticket", "", "ticket or issue identifier (required)")
	cmd.Flags().BoolVar(&live, "live", false, "execute live; default is dry-run")
	return cmd
}
