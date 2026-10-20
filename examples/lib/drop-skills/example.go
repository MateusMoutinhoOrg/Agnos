package main

import (
	"os"

	"github.com/MateusMoutinhoOrg/Agnos/adapters/bindings/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The drop-skills example: write the agent skill into both roots.
//
// It calls the same action `agnos drop-skills` calls, and writes only inside
// test-dir.
func main() {

	deps := standard.New()    // every adapter lib bound
	lib := sandbox.New(&deps) // *api.Sandbox
	module := "Test"

	if err := lib.Actions.Start(api.StartProps{
		Path:        "test-dir",
		ProjectName: "Test",
		Module:      &module,
	}); err != nil {
		panic(err)
	}

	if err := lib.Actions.DropSkills(api.DropSkillsProps{Path: "test-dir"}); err != nil {
		panic(err)
	}

	// What result.yaml records: the same set the cli side copies.
	for _, dir := range []string{".agents/skills", ".claude/skills"} {
		if err := os.CopyFS("assert-dir/"+dir, os.DirFS("test-dir/"+dir)); err != nil {
			panic(err)
		}
	}
	if err := copyExtensions(); err != nil {
		panic(err)
	}
}

// copyExtensions puts the project's extensions.yaml into assert-dir at the
// place it holds in the tree.
func copyExtensions() error {
	content, err := os.ReadFile("test-dir/AgnosConfig/extensions.yaml")
	if err != nil {
		return err
	}
	if err := os.MkdirAll("assert-dir/AgnosConfig", 0o755); err != nil {
		return err
	}
	return os.WriteFile("assert-dir/AgnosConfig/extensions.yaml", content, 0o644)
}
