package remove_flag

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
)

// RemoveFlag drops one flag declaration from
// sandbox/internal/commands/<command>/entries.yaml, then runs build so the
// generated entries.go and dispatch layer forget it.
func RemoveFlag(deps *deps.Deps, path string, command string, name string) error {
	io := smartio.New(deps, path, config.ProjectName)
	if err := RemoveFlagInternal(deps, io, command, name); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(deps, path)
}
