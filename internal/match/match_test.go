package match

import (
	"errors"
	"regexp"
	"testing"
)

func TestMatch(t *testing.T) {
	cases := []struct {
		pattern string
		value   string
		want    bool
	}{
		{"foo*", "foobar", true},
		{"foo?", "fooa", true},
		{"foo?", "foo", false},
		{"**/bar", "foo/bar", true},
		{"**/bar", "bar", true},
		{"bar", "bar", true},
		{"bar", "baz", false},
	}
	for _, tc := range cases {
		if got := Match(tc.pattern, tc.value); got != tc.want {
			t.Fatalf("Match(%q,%q)=%v want %v", tc.pattern, tc.value, got, tc.want)
		}
	}
}

func TestAny(t *testing.T) {
	patterns := []string{"foo*", "bar"}
	if !Any(patterns, "foobar") {
		t.Fatalf("Any should match")
	}
	if Any(patterns, "baz") {
		t.Fatalf("Any should not match")
	}
}

func TestAnySkipsBlank(t *testing.T) {
	if !Any([]string{" ", "foo*"}, "foobar") {
		t.Fatalf("expected blank skip")
	}
}

func TestMatchCommand(t *testing.T) {
	if !MatchCommand("code *", "code /tmp/path") {
		t.Fatalf("expected command match")
	}
	if !MatchCommand("git status", "git status") {
		t.Fatalf("expected exact match")
	}
	if !MatchCommand("git st?tus", "git status") {
		t.Fatalf("expected ? match")
	}
	if !MatchCommand("git (status)*", "git (status) now") {
		t.Fatalf("expected escaped match")
	}
	if AnyCommand([]string{"git status*"}, "git status --porcelain") == false {
		t.Fatalf("expected AnyCommand match")
	}
	if AnyCommand([]string{"git status*"}, "git fetch") {
		t.Fatalf("expected AnyCommand false")
	}
}

func TestAnyCommandSkipsBlank(t *testing.T) {
	if !AnyCommand([]string{" ", "git *"}, "git status") {
		t.Fatalf("expected match")
	}
}

func TestMatchEscapedPath(t *testing.T) {
	if !Match("foo.bar/**", "foo.bar/baz") {
		t.Fatalf("expected escaped path match")
	}
}

func TestMatchCompileError(t *testing.T) {
	orig := compileRegexp
	compileRegexp = func(string) (*regexp.Regexp, error) { return nil, errors.New("bad") }
	defer func() { compileRegexp = orig }()
	if Match("foo", "foo") == false {
		t.Fatalf("expected direct match")
	}
	if Match("foo*", "bar") {
		t.Fatalf("expected false on compile error")
	}
	if MatchCommand("foo*", "bar") {
		t.Fatalf("expected false on compile error")
	}
}
