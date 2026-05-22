package controlnode

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type Response map[string]any

type ExecResult struct {
	Command  []string      `json:"command"`
	Stdout   string        `json:"stdout,omitempty"`
	Stderr   string        `json:"stderr,omitempty"`
	ExitCode int           `json:"exit_code"`
	Duration time.Duration `json:"duration_ns"`
	WorkDir  string        `json:"work_dir,omitempty"`
}

type ProcessError struct {
	Probe ExecResult
}

func (e ProcessError) Error() string {
	return fmt.Sprintf("control-node processor failed with exit code %d", e.Probe.ExitCode)
}

func repoPath() string {
	if p := os.Getenv("PROPHET_AGENTPLANE_REPO"); p != "" {
		return p
	}
	return filepath.Clean("../agentplane")
}

func runInDir(ctx context.Context, dir string, argv []string) (ExecResult, error) {
	res := ExecResult{Command: append([]string(nil), argv...), WorkDir: dir}
	if len(argv) == 0 {
		return res, nil
	}
	start := time.Now()
	cmd := exec.CommandContext(ctx, argv[0], argv[1:]...)
	cmd.Dir = dir
	out, err := cmd.Output()
	res.Stdout = string(out)
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			res.Stderr = string(ee.Stderr)
			res.ExitCode = ee.ExitCode()
		} else {
			res.Stderr = err.Error()
			res.ExitCode = 1
		}
	} else if cmd.ProcessState != nil {
		res.ExitCode = cmd.ProcessState.ExitCode()
	}
	res.Duration = time.Since(start)
	return res, err
}

func Status(_ context.Context) (Response, error) {
	return Response{
		"delegated_to": "agentplane + prophet-platform + source-os",
		"status": "scaffold",
		"operation": "status",
		"agentplane_repo": repoPath(),
		"runtime_service": "prophet-platform/apps/node-commander",
		"contracts_repo": "SourceOS-Linux/sourceos-spec",
	}, nil
}

func Process(ctx context.Context, inputPath, outDir string) (Response, error) {
	argv := []string{"python3", filepath.ToSlash("scripts/process_local_control_node_input.py"), inputPath, outDir}
	res, err := runInDir(ctx, repoPath(), argv)
	resp := Response{
		"delegated_to": "SocioProphet/agentplane",
		"status": "completed",
		"operation": "process",
		"input_path": inputPath,
		"out_dir": outDir,
		"processor": "scripts/process_local_control_node_input.py",
		"probe": res,
	}
	if err != nil || res.ExitCode != 0 {
		resp["status"] = "failed"
		return resp, ProcessError{Probe: res}
	}
	return resp, nil
}
