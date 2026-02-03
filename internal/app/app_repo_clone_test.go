package app

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestParseCloneNameFromArgs(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		args    []string
		want    string
		wantErr bool
	}{
		{name: "empty", args: nil, want: ""},
		{name: "name-flag", args: []string{"--name", "alpha"}, want: "alpha"},
		{name: "name-flag-short", args: []string{"-name", "beta"}, want: "beta"},
		{name: "name-eq", args: []string{"--name=gamma"}, want: "gamma"},
		{name: "name-eq-short", args: []string{"-name=delta"}, want: "delta"},
		{name: "missing-arg", args: []string{"--name"}, wantErr: true},
		{name: "duplicate-flag", args: []string{"--name", "a", "--name", "b"}, wantErr: true},
		{name: "duplicate-eq", args: []string{"--name=a", "--name=b"}, wantErr: true},
		{name: "duplicate-eq-short", args: []string{"-name=a", "-name=b"}, wantErr: true},
		{name: "unknown-flag", args: []string{"--bad"}, wantErr: true},
		{name: "unexpected-arg", args: []string{"alpha"}, wantErr: true},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseCloneNameFromArgs(tc.args)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("unexpected name: %q", got)
			}
		})
	}
}

func TestRepoCloneNameFromArgs(t *testing.T) {
	app, cfg := newTestApp(t)
	bare := initBareRepo(t, filepath.Join(t.TempDir(), "remote.git"))
	seedBareRepo(t, bare)
	if code := app.runRepoClone(context.Background(), []string{bare, "--name", "from-args"}); code != 0 {
		t.Fatalf("expected clone")
	}
	if _, err := os.Stat(filepath.Join(cfg.ReposRoot, "from-args")); err != nil {
		t.Fatalf("expected clone dir: %v", err)
	}
}

func TestRepoCloneDuplicateNameFlags(t *testing.T) {
	app, _ := newTestApp(t)
	bare := initBareRepo(t, filepath.Join(t.TempDir(), "remote.git"))
	seedBareRepo(t, bare)
	if code := app.runRepoClone(context.Background(), []string{"--name", "alpha", bare, "--name=beta"}); code == 0 {
		t.Fatalf("expected duplicate name error")
	}
}

func TestRepoCloneNameArgsError(t *testing.T) {
	app, _ := newTestApp(t)
	bare := initBareRepo(t, filepath.Join(t.TempDir(), "remote.git"))
	seedBareRepo(t, bare)
	if code := app.runRepoClone(context.Background(), []string{bare, "--bad"}); code == 0 {
		t.Fatalf("expected parse error")
	}
}
