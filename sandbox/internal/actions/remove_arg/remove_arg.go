package remove_arg

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
)

// RemoveArg drops one positional arg declaration from
// sandbox/internal/commands/<command>/entries.yaml, then runs build so the
// generated entries.go and dispatch layer forget it.
func RemoveArg(deps *deps.Deps, path string, command string, name string) error {
	io := smartio.New(deps, path, config.ProjectName)
	if err := RemoveArgInternal(deps, io, command, name); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(deps, api.BuildProps{Path: path, Runtime: api.RuntimeNone})
}
