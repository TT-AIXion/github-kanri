package gitutil

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/TT-AIXion/github-kanri/internal/executil"
)

var symbolicRef = func(ctx context.Context, r executil.Runner, repo string) (executil.Result, error) {
	return r.Run(ctx, repo, "git", "symbolic-ref", "refs/remotes/origin/HEAD")
}

func StatusPorcelain(ctx context.Context, r executil.Runner, repo string) (string, error) {
	res, err := r.Run(ctx, repo, "git", "status", "--porcelain")
	return strings.TrimSpace(res.Stdout), err
}

func IsClean(ctx context.Context, r executil.Runner, repo string) (bool, error) {
	out, err := StatusPorcelain(ctx, r, repo)
	if err != nil {
		return false, err
	}
	return out == "", nil
}

func CurrentBranch(ctx context.Context, r executil.Runner, repo string) (string, error) {
	res, err := r.Run(ctx, repo, "git", "rev-parse", "--abbrev-ref", "HEAD")
	return strings.TrimSpace(res.Stdout), err
}

func DefaultBranch(ctx context.Context, r executil.Runner, repo string) (string, error) {
	res, err := symbolicRef(ctx, r, repo)
	if err != nil {
		return "", err
	}
	out := strings.TrimSpace(res.Stdout)
	parts := strings.Split(out, "/")
	if len(parts) == 0 || parts[len(parts)-1] == "" {
		return "", fmt.Errorf("unexpected origin HEAD")
	}
	return parts[len(parts)-1], nil
}

func OriginURL(ctx context.Context, r executil.Runner, repo string) (string, error) {
	res, err := r.Run(ctx, repo, "git", "remote", "get-url", "origin")
	return strings.TrimSpace(res.Stdout), err
}

func LastCommitUnix(ctx context.Context, r executil.Runner, repo string) (int64, error) {
	res, err := r.Run(ctx, repo, "git", "log", "-1", "--format=%ct")
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(strings.TrimSpace(res.Stdout), 10, 64)
}

func LogOneline(ctx context.Context, r executil.Runner, repo string, limit int) (string, error) {
	args := []string{"log", "--oneline", "--decorate"}
	if limit > 0 {
		args = append(args, "-n", strconv.Itoa(limit))
	}
	res, err := r.Run(ctx, repo, "git", args...)
	return strings.TrimSpace(res.Stdout), err
}

func Clone(ctx context.Context, r executil.Runner, url string, dest string) error {
	_, err := r.Run(ctx, "", "git", "clone", url, dest)
	return err
}

func Pull(ctx context.Context, r executil.Runner, repo string) error {
	_, err := r.Run(ctx, repo, "git", "pull")
	return err
}

func Fetch(ctx context.Context, r executil.Runner, repo string) error {
	_, err := r.Run(ctx, repo, "git", "fetch", "--all", "--tags")
	return err
}

func Checkout(ctx context.Context, r executil.Runner, repo string, ref string) error {
	_, err := r.Run(ctx, repo, "git", "checkout", ref)
	return err
}
