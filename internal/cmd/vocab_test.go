package cmd

import "testing"

func TestVocabCommandHasExpectedSubcommands(t *testing.T) {
	cmd := newVocabCmd()
	want := map[string]bool{
		"fetch":    false,
		"validate": false,
		"promote":  false,
		"sr":       false,
	}
	for _, c := range cmd.Commands() {
		if _, ok := want[c.Name()]; ok {
			want[c.Name()] = true
		}
	}
	for name, seen := range want {
		if !seen {
			t.Fatalf("missing vocab subcommand: %s", name)
		}
	}
}
