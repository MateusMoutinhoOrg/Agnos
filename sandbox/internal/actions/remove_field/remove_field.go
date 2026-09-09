package remove_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// RemoveField deletes one declared field from
// sandbox/internal/routes/<route>/route.yaml, from the origin named by in,
// then runs build as a follow-up step. The build renders only: dropping a
// field may leave hand-written code referring to what is gone.
func RemoveField(deps *deps.Deps, path string, route string, in string, name string) error {
	io := smartio.New(deps, path, config.ProjectName)
	if err := RemoveFieldInternal(deps, io, route, in, name); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(deps, api.BuildProps{Path: path, Runtime: api.RuntimeNone})
}
