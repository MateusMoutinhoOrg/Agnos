package remove_skills

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeSkillsAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_skills"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	remove_error := removeSkillsAction.RemoveSkills(sandbox, api.RemoveSkillsProps{Path: props.Path})
	if remove_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", remove_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
