package app

import (
	"time"

	"github.com/TT-AIXion/github-kanri/internal/repo"
)

type repoStatus struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	Dirty bool   `json:"dirty"`
}

type repoInfo struct {
	Name          string `json:"name"`
	Path          string `json:"path"`
	Origin        string `json:"origin"`
	CurrentBranch string `json:"currentBranch"`
	DefaultBranch string `json:"defaultBranch"`
	Dirty         bool   `json:"dirty"`
}

type repoRecent struct {
	Name       string    `json:"name"`
	Path       string    `json:"path"`
	Unix       int64     `json:"unix"`
	Timestamp  time.Time `json:"timestamp"`
	HasCommits bool      `json:"hasCommits"`
}

type execResult struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	ExitCode int    `json:"exitCode"`
	Duration string `json:"duration"`
	Stdout   string `json:"stdout,omitempty"`
	Stderr   string `json:"stderr,omitempty"`
	Error    string `json:"error,omitempty"`
	Skipped  bool   `json:"skipped,omitempty"`
}

func (a App) handleMultiMatch(result repo.MatchResult) int {
	if a.Out.JSON {
		a.Out.Warn("multiple matches", result.Matches)
		return 2
	}
	for _, line := range repo.Candidates(result) {
		a.Out.Warn(line, nil)
	}
	return 2
}
