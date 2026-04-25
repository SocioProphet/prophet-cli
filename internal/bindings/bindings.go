package bindings

import (
	"context"
	"os"

	executil "github.com/socioprophet/prophet-cli/internal/exec"
)

const defaultRepoPath = "../socioprophet-standards-knowledge"

type Response map[string]any

func repoPath() string {
	if v := os.Getenv("PROPHET_VOCAB_REPO"); v != "" {
		return v
	}
	if v := os.Getenv("PROPHET_STANDARDS_REPO"); v != "" {
		return v
	}
	return defaultRepoPath
}

func ValidateABD(ctx context.Context, abdPath string) (Response, error) {
	repo := repoPath()
	argv := []string{"python3", "Lower/tools/validate_abd.py", "--abd", abdPath, "--json"}
	res, err := executil.RunInDir(ctx, repo, argv)
	status := "ok"
	if err != nil {
		status = "failed"
	}
	return Response{
		"delegated_to": "socioprophet-standards-knowledge",
		"operation":    "validate_abd",
		"status":       status,
		"repo_path":    repo,
		"abd_path":     abdPath,
		"validator":    "Lower/tools/validate_abd.py",
		"result":       res,
	}, err
}
