package remove_skills

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// RemoveSkills deletes the generator's skill and turns the skills extension
// off, then runs build as a follow-up step.
func RemoveSkills(sandbox *api.Sandbox, props api.RemoveSkillsProps) error {
	io := stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := RemoveSkillsInternal(sandbox, io, props.Path); err != nil {
		return err
	}
	return buildAction.PersistAndBuild(sandbox, io, api.BuildProps{Path: props.Path, Runtime: api.RuntimeNone})
}
