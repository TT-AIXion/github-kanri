package repo

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanFilterFind(t *testing.T) {
	root := t.TempDir()
	repo1 := filepath.Join(root, "alpha")
	repo2 := filepath.Join(root, "beta")
	if err := os.MkdirAll(filepath.Join(repo1, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(repo2, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	repos, err := Scan(root)
	if err != nil || len(repos) != 2 {
		t.Fatalf("scan error: %v repos=%v", err, repos)
	}
	filtered := Filter(repos, []string{"*a*"}, []string{"beta"})
	if len(filtered) != 1 || filtered[0].Name != "alpha" {
		t.Fatalf("unexpected filter: %v", filtered)
	}
	filtered = Filter(repos, []string{"nomatch"}, nil)
	if len(filtered) != 0 {
		t.Fatalf("expected empty filter")
	}
	result := Find(repos, "a")
	if len(result.Matches) != 2 {
		t.Fatalf("expected two matches")
	}
	if got := Find(repos, ""); got.Matches != nil {
		t.Fatalf("expected empty pattern")
	}
	wild := Find(repos, "*a*")
	if len(wild.Matches) == 0 {
		t.Fatalf("expected wildcard matches")
	}
	picked, err := Pick(result, 2)
	if err != nil || picked.Name != "beta" {
		t.Fatalf("pick error: %v picked=%v", err, picked)
	}
	one, err := Pick(MatchResult{Matches: []Repo{{Name: "solo"}}}, 1)
	if err != nil || one.Name != "solo" {
		t.Fatalf("expected single pick")
	}
	cands := Candidates(result)
	if len(cands) != 2 {
		t.Fatalf("expected candidates")
	}
}

func TestScanError(t *testing.T) {
	if _, err := Scan(filepath.Join(t.TempDir(), "missing")); err == nil {
		t.Fatalf("expected scan error")
	}
}
