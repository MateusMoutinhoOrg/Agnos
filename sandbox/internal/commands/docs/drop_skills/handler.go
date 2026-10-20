package drop_skills

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	dropSkillsAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/drop_skills"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	drop_error := dropSkillsAction.DropSkills(sandbox, api.DropSkillsProps{Path: props.Path})
	if drop_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", drop_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
