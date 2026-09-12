package remove_flag

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// RemoveFlag drops one flag declaration from
// sandbox/internal/commands/<command>/entries.yaml, then runs build so the
// generated entries.go and dispatch layer forget it.
func RemoveFlag(sandbox *api.Sandbox, path string, command string, name string) error {
	io := smartio.New(sandbox, path, config.ProjectName)
	if err := RemoveFlagInternal(sandbox, io, command, name); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: path, Runtime: api.RuntimeNone})
}
