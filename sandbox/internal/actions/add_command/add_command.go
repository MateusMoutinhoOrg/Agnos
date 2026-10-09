package add_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// AddCommand scaffolds a new command package under
// sandbox/internal/commands/<name>/ — a hand-written command.yaml and a stub
// handler.go — then runs build as a follow-up step so its generated.new.go
// and generated.input.go — the api.Command that lands in Cli.Commands — are generated
// for it.
func AddCommand(sandbox *api.Sandbox, props api.AddCommandProps) error {
	io := stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := AddCommandInternal(sandbox, io, props); err != nil {
		return err
	}
	return buildAction.PersistAndBuild(sandbox, io, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
