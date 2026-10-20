package drop_skills

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// DropSkills turns the skills extension on, then runs build as a follow-up
// step, which writes the skill. It writes no Go, so the build renders only.
func DropSkills(sandbox *api.Sandbox, props api.DropSkillsProps) error {
	io := stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := DropSkillsInternal(sandbox, io, props.Path); err != nil {
		return err
	}
	return buildAction.PersistAndBuild(sandbox, io, api.BuildProps{Path: props.Path, Runtime: api.RuntimeNone})
}
