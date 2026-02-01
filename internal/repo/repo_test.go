package repo

import "testing"

func TestPickErrors(t *testing.T) {
	result := MatchResult{Matches: []Repo{{Name: "a"}, {Name: "b"}}}
	if _, err := Pick(result, 0); err != ErrMultipleMatches {
		t.Fatalf("expected ErrMultipleMatches")
	}
	result = MatchResult{Matches: nil}
	if _, err := Pick(result, 0); err != ErrNoMatch {
		t.Fatalf("expected ErrNoMatch")
	}
}
