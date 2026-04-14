package a2a

import "testing"

func TestDefaultWorkflowHasExpectedPhases(t *testing.T) {
	wf := Default("demo/repo", "DEMO", false)
	if wf.Repo != "demo/repo" {
		t.Fatalf("unexpected repo: %s", wf.Repo)
	}
	if len(wf.Phases) != 6 {
		t.Fatalf("expected 6 phases, got %d", len(wf.Phases))
	}
	if wf.Phases[0].Name != "propose" {
		t.Fatalf("unexpected first phase: %s", wf.Phases[0].Name)
	}
	if wf.Phases[len(wf.Phases)-1].Name != "done" {
		t.Fatalf("unexpected last phase: %s", wf.Phases[len(wf.Phases)-1].Name)
	}
}
