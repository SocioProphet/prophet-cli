package vocab

import "testing"

func TestResolveRepoPathPrefersEnv(t *testing.T) {
	t.Setenv("PROPHET_VOCAB_REPO", "/tmp/ontogenesis")
	if got := resolveRepoPath(); got != "/tmp/ontogenesis" {
		t.Fatalf("unexpected repo path: %s", got)
	}
}

func TestNormalizeGraphPathsDefaultsToDot(t *testing.T) {
	got := normalizeGraphPaths(nil)
	if len(got) != 1 || got[0] != "." {
		t.Fatalf("unexpected default graph paths: %#v", got)
	}
}
