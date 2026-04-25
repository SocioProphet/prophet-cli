package k8s

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

func Shapecheck(ctx context.Context, manifestPaths []string) (Response, error) {
	repo := repoPath()
	argv := []string{"python3", "k8s/tools/shapecheck.py", "--json", "--manifest"}
	argv = append(argv, manifestPaths...)
	res, err := executil.RunInDir(ctx, repo, argv)
	status := "ok"
	if err != nil {
		status = "failed"
	}
	return Response{
		"delegated_to":   "socioprophet-standards-knowledge",
		"operation":      "shapecheck",
		"status":         status,
		"repo_path":      repo,
		"manifest_paths": manifestPaths,
		"validator":      "k8s/tools/shapecheck.py",
		"result":         res,
	}, err
}
