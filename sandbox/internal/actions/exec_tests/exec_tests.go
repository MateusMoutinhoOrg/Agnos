package exec_tests

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// The package is exec_tests, not exec_test, for one reason only: Go reserves
// every file whose name ends in _test.go for the testing toolchain, so an
// action directory named after the `exec-test` command could not hold its own
// <name>.go. The command, the api field and the docs all keep the name.

// ExecTest runs the project's examples and checks each one against its golden
// result.yaml. Unlike every action that writes into the project tree it opens
// no SmartIO: there is no transaction to persist — each example's TestDir is
// written by a child process, outside any buffer, and the tree recorded for it
// has to be the literal one on disk, unfiltered by ignore.yaml / paths.yaml.
// The project path is therefore joined here rather than at the SmartIO
// boundary.
func ExecTest(deps *deps.Deps, props api.ExecTestProps) error {
	if !deps.Iodeps.IsDir(join(props.Path, utils.ExamplesDir)) {
		return deps.Std.Errorf("exec-test: %s has no %s/ directory (create one with add-cli-example / add-lib-example)",
			props.Path, utils.ExamplesDir)
	}

	deps.Std.Log("exec-test started with path %s \n", props.Path)

	return ExecTestInternal(deps, props.Path, props.Only, props.Update)
}
