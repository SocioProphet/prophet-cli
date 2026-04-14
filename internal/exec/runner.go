package executil

import (
	"bytes"
	"context"
	"os/exec"
	"time"
)

type Result struct {
	Command  []string      `json:"command"`
	Dir      string        `json:"dir,omitempty"`
	Stdout   string        `json:"stdout,omitempty"`
	Stderr   string        `json:"stderr,omitempty"`
	ExitCode int           `json:"exit_code"`
	Duration time.Duration `json:"duration_ns"`
}

func Run(ctx context.Context, argv []string) (Result, error) {
	return RunInDir(ctx, "", argv)
}

func RunInDir(ctx context.Context, dir string, argv []string) (Result, error) {
	res := Result{Command: append([]string(nil), argv...), Dir: dir}
	if len(argv) == 0 {
		return res, nil
	}
	start := time.Now()
	cmd := exec.CommandContext(ctx, argv[0], argv[1:]...)
	if dir != "" {
		cmd.Dir = dir
	}
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	res.Stdout = outb.String()
	res.Stderr = errb.String()
	res.Duration = time.Since(start)
	if cmd.ProcessState != nil {
		res.ExitCode = cmd.ProcessState.ExitCode()
	}
	return res, err
}
