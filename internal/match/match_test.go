package match

import "testing"

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
