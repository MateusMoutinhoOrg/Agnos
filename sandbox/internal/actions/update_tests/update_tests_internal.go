package update_tests

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	execTestsAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/exec_tests"
)

// UpdateTestInternal is the run itself: the exec_test action, narrowed to one
// example name and told to write. Nothing about running an example is
// reimplemented here — an update that took a different path through the suite
// would be updating a golden the checking run never produces.
func UpdateTestInternal(deps *deps.Deps, path string, name string) error {
	return execTestsAction.ExecTest(deps, api.ExecTestProps{
		Path:   path,
		Only:   name,
		Update: true,
	})
}
