package repo

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindFilterCandidates(t *testing.T) {
	repos := []Repo{{Name: "alpha"}, {Name: "beta"}, {Name: "gamma"}}
	filtered := Filter(repos, []string{"a*"}, []string{"gamma"})
	if len(filtered) != 1 {
		t.Fatalf("unexpected filter: %v", filtered)
	}
	res := Find(repos, "a")
	if len(res.Matches) == 0 {
		t.Fatalf("expected matches")
	}
	res2 := Find(repos, "*ta")
	if len(res2.Matches) != 1 {
		t.Fatalf("expected one match")
	}
	res3 := Find(repos, "")
	if len(res3.Matches) != 0 {
		t.Fatalf("expected no match")
	}
	candidates := Candidates(res2)
	if len(candidates) != 1 {
		t.Fatalf("expected candidates")
	}
}

func TestPickSuccess(t *testing.T) {
	res := MatchResult{Matches: []Repo{{Name: "alpha"}, {Name: "beta"}}}
	if _, err := Pick(res, 2); err != nil {
		t.Fatalf("unexpected err")
	}
}

func TestScan(t *testing.T) {
	root := t.TempDir()
	repoDir := filepath.Join(root, "repo1", ".git")
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	repos, err := Scan(root)
	if err != nil || len(repos) != 1 {
		t.Fatalf("scan failed: %v %v", err, repos)
	}
}
