package repo

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/TT-AIXion/github-kanri/internal/fsutil"
	"github.com/TT-AIXion/github-kanri/internal/match"
)

type Repo struct {
	Name string
	Path string
}

type MatchResult struct {
	Matches []Repo
}

var ErrNoMatch = errors.New("no match")
var ErrMultipleMatches = errors.New("multiple matches")

func Scan(root string) ([]Repo, error) {
	paths, err := fsutil.ListGitRepos(root)
	if err != nil {
		return nil, err
	}
	var repos []Repo
	for _, path := range paths {
		repos = append(repos, Repo{Name: filepath.Base(path), Path: path})
	}
	sort.Slice(repos, func(i, j int) bool { return repos[i].Name < repos[j].Name })
	return repos, nil
}

func Filter(repos []Repo, only []string, exclude []string) []Repo {
	var out []Repo
	for _, r := range repos {
		if len(only) > 0 && !match.Any(only, r.Name) {
			continue
		}
		if len(exclude) > 0 && match.Any(exclude, r.Name) {
			continue
		}
		out = append(out, r)
	}
	return out
}

func Find(repos []Repo, pattern string) MatchResult {
	pattern = strings.TrimSpace(pattern)
	if pattern == "" {
		return MatchResult{Matches: nil}
	}
	var matches []Repo
	for _, r := range repos {
		if matchName(r.Name, pattern) {
			matches = append(matches, r)
		}
	}
	return MatchResult{Matches: matches}
}

func Pick(result MatchResult, pick int) (Repo, error) {
	if len(result.Matches) == 0 {
		return Repo{}, ErrNoMatch
	}
	if len(result.Matches) == 1 {
		return result.Matches[0], nil
	}
	if pick <= 0 || pick > len(result.Matches) {
		return Repo{}, ErrMultipleMatches
	}
	return result.Matches[pick-1], nil
}

func matchName(name, pattern string) bool {
	if strings.ContainsAny(pattern, "*?") {
		return match.Match(pattern, name)
	}
	return strings.Contains(name, pattern)
}

func Candidates(result MatchResult) []string {
	var out []string
	for i, r := range result.Matches {
		out = append(out, fmt.Sprintf("%d: %s", i+1, r.Name))
	}
	return out
}
