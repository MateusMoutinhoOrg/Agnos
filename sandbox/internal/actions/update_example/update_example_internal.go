package update_example

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	runExamplesAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/run_examples"
)

// UpdateExampleInternal is the run itself: the run_examples action, narrowed to one
// example name and told to write. Nothing about running an example is
// reimplemented here — an update that took a different path through the suite
// would be updating a golden the checking run never produces.
func UpdateExampleInternal(sandbox *api.Sandbox, path string, name string) error {
	return runExamplesAction.RunExamples(sandbox, api.RunExamplesProps{
		Path:   path,
		Only:   name,
		Update: true,
	})
}
