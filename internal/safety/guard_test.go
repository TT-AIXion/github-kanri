package safety

import "testing"

func TestGuardCommand(t *testing.T) {
	g := Guard{
		AllowCommands: []string{"git status*"},
		DenyCommands:  []string{"rm -rf*"},
	}
	if err := g.CheckCommand("git status --porcelain"); err != nil {
		t.Fatalf("expected allow")
	}
	if err := g.CheckCommand("rm -rf /tmp"); err == nil {
		t.Fatalf("expected deny")
	}
	if err := g.CheckCommand("echo hi"); err == nil {
		t.Fatalf("expected deny not allowed")
	}
}

func TestGuardPath(t *testing.T) {
	g := Guard{
		AllowPaths: []string{"/tmp/*"},
		DenyPaths:  []string{"/tmp/secret*"},
	}
	if err := g.CheckPath("/tmp/test"); err != nil {
		t.Fatalf("expected allow")
	}
	if err := g.CheckPath("/tmp/secret.txt"); err == nil {
		t.Fatalf("expected deny")
	}
	if err := g.CheckPath("/var/log"); err == nil {
		t.Fatalf("expected deny not allowed")
	}
}

func TestGuardEmptyAllow(t *testing.T) {
	g := Guard{}
	if err := g.CheckCommand("anything"); err != nil {
		t.Fatalf("expected allow when no allow list")
	}
	if err := g.CheckPath("/tmp/anything"); err != nil {
		t.Fatalf("expected allow when no allow list")
	}
}
