package cli

import (
	"reflect"
	"testing"

	"github.com/karpikorop/zsh-dcmd-go/internal/store"
)

func TestCopyCommands(t *testing.T) {
	source := &store.Project{
		Name: "source",
		Vars: map[string]string{"project": "src-dev", "config": "src.json"},
		Commands: []store.Command{
			{ID: 1, Run: "build", Description: "Build"},
			{ID: 3, Run: "deploy --project {{project}} --config {{config}}", Description: "Deploy", Confirm: true},
			{ID: 7, Run: "lint", Description: "Lint"},
		},
	}
	target := &store.Project{
		Name:     "target",
		Vars:     map[string]string{"project": "target-dev"},
		Commands: []store.Command{{ID: 5, Run: "build", Description: "Already here"}},
	}

	picked, err := pickCommands(source, []int{3, 1, 3})
	if err != nil {
		t.Fatal(err)
	}
	result := copyCommands(target, picked, source.Vars)

	if len(result.added) != 1 || len(result.skipped) != 2 {
		t.Fatalf("got %d added and %d skipped, want 1 and 2", len(result.added), len(result.skipped))
	}
	if added := result.added[0].command; added.ID != 6 || !added.Confirm {
		t.Errorf("copied command: got %+v, want ID 6 with confirm", added)
	}

	wantVars := map[string]string{"project": "target-dev", "config": "src.json"}
	if !reflect.DeepEqual(target.Vars, wantVars) || !reflect.DeepEqual(result.vars, []string{"config"}) {
		t.Errorf("vars: got %v (copied %v), want %v", target.Vars, result.vars, wantVars)
	}

	if _, err := pickCommands(source, []int{1, 2, 9}); err == nil {
		t.Error("expected an error for missing IDs 2 and 9")
	}
}
