package cmd

import "testing"

func TestRootCommandHasExpectedTopLevelCommands(t *testing.T) {
	root := NewRootCommand()
	want := map[string]bool{
		"bootstrap":    false,
		"vocab":        false,
		"bindings":     false,
		"k8s":          false,
		"control-node": false,
		"a2a":          false,
		"ask":          false,
		"plan":         false,
		"agent":        false,
		"mcp":          false,
	}
	for _, c := range root.Commands() {
		if _, ok := want[c.Name()]; ok {
			want[c.Name()] = true
		}
	}
	for name, seen := range want {
		if !seen {
			t.Fatalf("missing top-level command: %s", name)
		}
	}
}
