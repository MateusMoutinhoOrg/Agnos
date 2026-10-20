package main

import (
	"os"

	"github.com/MateusMoutinhoOrg/Agnos/adapters/bindings/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The remove-skills example: take the agent skill away again.
//
// It calls the same actions `agnos drop-skills` and `agnos remove-skills`
// call, and writes only inside test-dir.
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

	if err := lib.Actions.RemoveSkills(api.RemoveSkillsProps{Path: "test-dir"}); err != nil {
		panic(err)
	}

	// Both roots are gone with the skill: nothing else was in them.
	for _, root := range []string{"test-dir/.agents", "test-dir/.claude"} {
		if _, err := os.Stat(root); err == nil {
			panic(root + " is still there")
		}
	}

	// What result.yaml records: the same set the cli side copies.
	content, err := os.ReadFile("test-dir/AgnosConfig/extensions.yaml")
	if err != nil {
		panic(err)
	}
	if err := os.MkdirAll("assert-dir/AgnosConfig", 0o755); err != nil {
		panic(err)
	}
	if err := os.WriteFile("assert-dir/AgnosConfig/extensions.yaml", content, 0o644); err != nil {
		panic(err)
	}
}
