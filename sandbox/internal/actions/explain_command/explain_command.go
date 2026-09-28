package explain_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// ExplainCommand runs one command line against the declared commands without
// running any. Like the other readers it opens a SmartIO and never calls
// io.Persist, and it runs no follow-up build.
func ExplainCommand(sandbox *api.Sandbox, props api.ExplainCommandProps) ([]string, error) {
	io := smartio.New(sandbox, props.Path, sandbox.Config.ProjectName)
	return ExplainCommandInternal(sandbox, io, props)
}
