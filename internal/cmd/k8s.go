package cmd

import (
	"github.com/socioprophet/prophet-cli/internal/k8s"
	"github.com/spf13/cobra"
)

func newK8sCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "k8s", Short: "Kubernetes policy façade"}
	shapecheck := &cobra.Command{Use: "shapecheck <manifest-path>...", Short: "Run K8s manifest shape checks", Args: cobra.MinimumNArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := k8s.Shapecheck(cmd.Context(), args)
		if resp != nil {
			resp["command"] = "prophet k8s shapecheck"
			_ = emit(resp)
		}
		return err
	}}
	cmd.AddCommand(shapecheck)
	return cmd
}
