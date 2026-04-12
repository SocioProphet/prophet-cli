package vocab

import (
	"context"
	"path/filepath"

	executil "github.com/socioprophet/prophet-cli/internal/exec"
)

type Response map[string]any

func Fetch(_ context.Context) (Response, error) {
	return Response{
		"delegated_to": "socioprophet-standards-knowledge",
		"status":       "scaffold",
		"operation":    "fetch",
	}, nil
}

func Validate(ctx context.Context, graphPaths []string) (Response, error) {
	cmd := []string{"python3", filepath.ToSlash("policy/tools/validate_all.py"), "--data"}
	cmd = append(cmd, graphPaths...)
	res, _ := executil.Run(ctx, cmd)
	return Response{
		"delegated_to": "socioprophet-standards-knowledge",
		"status":       "scaffold",
		"operation":    "validate",
		"graph_paths":  graphPaths,
		"validator":    "policy/tools/validate_all.py",
		"probe":        res,
	}, nil
}

func Promote(_ context.Context) (Response, error) {
	return Response{
		"delegated_to":  "socioprophet-standards-knowledge",
		"status":        "scaffold",
		"operation":     "promote",
		"record_target": "/ontology/promotions/<ts>.jsonld",
	}, nil
}

func SRRun(_ context.Context, module string) (Response, error) {
	return Response{
		"delegated_to": "socioprophet-standards-knowledge",
		"status":       "scaffold",
		"operation":    "sr_run",
		"module":       module,
	}, nil
}

func SRGate(_ context.Context) (Response, error) {
	return Response{
		"delegated_to": "socioprophet-standards-knowledge",
		"status":       "scaffold",
		"operation":    "sr_gate",
	}, nil
}
