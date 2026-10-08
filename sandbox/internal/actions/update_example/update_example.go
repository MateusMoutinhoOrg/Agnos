package update_example

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// The package is update_example, not update_example, for the same reason run_examples
// is plural: Go reserves every file whose name ends in _test.go for the
// testing toolchain, so an action directory named after the `update-example`
// command could not hold its own <name>.go. The command, the api field and the
// docs all keep the name.

// UpdateExample rewrites the goldens of one example with what it produces now,
// printing what each write changed. It is `run-examples --update` narrowed to a
// single name, which is what makes an update reviewable: a suite-wide rewrite
// hides the one golden that moved for a reason nobody meant.
//
// The name is required. Without it the command is `run-examples --update` under
// another spelling, and one golden at a time is the whole of it.
func UpdateExample(sandbox *api.Sandbox, props api.UpdateExampleProps) error {
	if sandbox.Deps.StringsDeps.TrimSpace(props.Name) == "" {
		return sandbox.Deps.StdDeps.Errorf("update-example: an example name is required (rewrite every golden with `run-examples --update`)")
	}

	if err := utils.RequireExtension(sandbox, stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName), utils.ExtensionExample); err != nil {
		return err
	}

	sandbox.Deps.StdDeps.Logf("update-example started with path %s \n", props.Path)

	return UpdateExampleInternal(sandbox, props.Path, props.Name)
}
