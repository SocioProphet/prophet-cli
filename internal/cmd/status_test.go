package cmd

import "testing"

// TestStatusRegisteredAsTopLevel guards that `prophet status` stays on the surface.
func TestStatusRegisteredAsTopLevel(t *testing.T) {
	root := NewRootCommand()
	found := false
	for _, c := range root.Commands() {
		if c.Name() == "status" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("missing top-level command: status")
	}
}

// TestFacadeSurfacesWellFormed guards the boundary descriptor invariants: every
// surface has a command and note, delegating surfaces name a delegate, and real/
// scaffold surfaces do not.
func TestFacadeSurfacesWellFormed(t *testing.T) {
	for _, s := range facadeSurfaces() {
		if s.Command == "" {
			t.Fatalf("surface with empty command: %+v", s)
		}
		if s.Note == "" {
			t.Fatalf("surface %q has empty note", s.Command)
		}
		switch s.Kind {
		case kindDelegating:
			if s.Delegate == "" {
				t.Fatalf("delegating surface %q names no delegate", s.Command)
			}
		case kindReal, kindScaffold:
			if s.Delegate != "" {
				t.Fatalf("%s surface %q must not name a delegate, got %q", s.Kind, s.Command, s.Delegate)
			}
		default:
			t.Fatalf("surface %q has unknown kind %q", s.Command, s.Kind)
		}
	}
}
