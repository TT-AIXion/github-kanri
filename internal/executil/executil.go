package executil

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/TT-AIXion/github-kanri/internal/safety"
)

type Runner struct {
	Guard  safety.Guard
	DryRun bool
}

type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Duration time.Duration
}

func (r Runner) Run(ctx context.Context, dir string, name string, args ...string) (Result, error) {
	cmdline := strings.Join(append([]string{name}, args...), " ")
	if err := r.Guard.CheckCommand(cmdline); err != nil {
		return Result{}, err
	}
	if r.DryRun {
		return Result{ExitCode: 0}, nil
	}
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	start := time.Now()
	err := cmd.Run()
	dur := time.Since(start)
	code := exitCode(err)
	res := Result{Stdout: stdout.String(), Stderr: stderr.String(), ExitCode: code, Duration: dur}
	if err != nil {
		return res, fmt.Errorf("command failed: %s", cmdline)
	}
	return res, nil
}

func (r Runner) RunShell(ctx context.Context, dir string, command string) (Result, error) {
	if err := r.Guard.CheckCommand(command); err != nil {
		return Result{}, err
	}
	if r.DryRun {
		return Result{ExitCode: 0}, nil
	}
	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	cmd.Dir = dir
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	start := time.Now()
	err := cmd.Run()
	dur := time.Since(start)
	code := exitCode(err)
	res := Result{Stdout: stdout.String(), Stderr: stderr.String(), ExitCode: code, Duration: dur}
	if err != nil {
		return res, fmt.Errorf("command failed: %s", command)
	}
	return res, nil
}

func exitCode(err error) int {
	if err == nil {
		return 0
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode()
	}
	return 1
}
