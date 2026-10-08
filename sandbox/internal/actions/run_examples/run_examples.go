package run_examples

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// The package is run_examples, not run_examples, for one reason only: Go reserves
// every file whose name ends in _test.go for the testing toolchain, so an
// action directory named after the `run-examples` command could not hold its own
// <name>.go. The command, the api field and the docs all keep the name.

// RunExamples runs the project's examples and checks each one against its golden
// result.yaml. Unlike every action that writes into the project tree it opens
// no StagedFS: there is no transaction to persist — each example's test-dir and
// assert-dir are written by a child process, outside any buffer, and the tree
// recorded for it has to be the literal one on disk, unfiltered by
// paths.yaml.
// The project path is therefore joined here rather than at the StagedFS
// boundary.
func RunExamples(sandbox *api.Sandbox, props api.RunExamplesProps) error {
	if err := utils.RequireExtension(sandbox, stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName), utils.ExtensionExample); err != nil {
		return err
	}

	if !sandbox.Deps.IoDeps.IsDir(join(sandbox, props.Path, utils.ExamplesDir)) {
		return sandbox.Deps.StdDeps.Errorf("run-examples: %s has no %s/ directory (create one with add-cli-example / add-lib-example)",
			props.Path, utils.ExamplesDir)
	}

	sandbox.Deps.StdDeps.Logf("run-examples started with path %s \n", props.Path)

	return RunExamplesInternal(sandbox, props.Path, props.Only, props.Update)
}
