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

func TestMatchCommand(t *testing.T) {
	if !MatchCommand("code *", "code /tmp/path") {
		t.Fatalf("expected command match")
	}
	if AnyCommand([]string{"git status*"}, "git status --porcelain") == false {
		t.Fatalf("expected AnyCommand match")
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
