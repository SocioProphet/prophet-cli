package vocab

import (
	"context"
	"os"

	executil "github.com/socioprophet/prophet-cli/internal/exec"
)

const defaultRepoPath = "../socioprophet-standards-knowledge"

type Response map[string]any

func resolveRepoPath() string {
	if v := os.Getenv("PROPHET_VOCAB_REPO"); v != "" {
		return v
	}
	return defaultRepoPath
}

func normalizeGraphPaths(graphPaths []string) []string {
	if len(graphPaths) == 0 {
		return []string{"."}
	}
	return graphPaths
}

func Fetch(_ context.Context) (Response, error) {
	repoPath := resolveRepoPath()
	return Response{
		"delegated_to": "socioprophet-standards-knowledge",
		"status":       "scaffold",
		"operation":    "fetch",
		"repo_path":    repoPath,
	}, nil
}

func Validate(ctx context.Context, graphPaths []string) (Response, error) {
	repoPath := resolveRepoPath()
	normalized := normalizeGraphPaths(graphPaths)
	cmd := []string{"python3", "policy/tools/validate_all.py", "--data"}
	cmd = append(cmd, normalized...)
	res, err := executil.RunInDir(ctx, repoPath, cmd)
	status := "ok"
	if err != nil {
		status = "failed"
	}
	return Response{
		"delegated_to": "socioprophet-standards-knowledge",
		"status":       status,
		"operation":    "validate",
		"repo_path":    repoPath,
		"graph_paths":  normalized,
		"validator":    "policy/tools/validate_all.py",
		"result":       res,
	}, err
}

func Promote(_ context.Context) (Response, error) {
	repoPath := resolveRepoPath()
	return Response{
		"delegated_to":  "socioprophet-standards-knowledge",
		"status":        "scaffold",
		"operation":     "promote",
		"repo_path":     repoPath,
		"record_target": "/ontology/promotions/<ts>.jsonld",
	}, nil
}

func SRRun(_ context.Context, module string) (Response, error) {
	repoPath := resolveRepoPath()
	return Response{
		"delegated_to": "socioprophet-standards-knowledge",
		"status":       "scaffold",
		"operation":    "sr_run",
		"repo_path":    repoPath,
		"module":       module,
	}, nil
}

func SRGate(_ context.Context) (Response, error) {
	repoPath := resolveRepoPath()
	return Response{
		"delegated_to": "socioprophet-standards-knowledge",
		"status":       "scaffold",
		"operation":    "sr_gate",
		"repo_path":    repoPath,
	}, nil
}
