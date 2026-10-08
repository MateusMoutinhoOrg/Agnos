package explain_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	explainRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/explain_route"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	lines, explain_error := explainRouteAction.ExplainRoute(sandbox, api.ExplainRouteProps{
		Path:        props.Path,
		Method:      input.Method,
		RequestPath: input.RequestPath,
		Headers:     input.Header,
		Cookies:     input.Cookie,
	})
	if explain_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", explain_error.Error())
	}

	for _, line := range lines {
		response.Printf("%s\n", line)
	}
	response.SetStatus(api.ExitOk)
	return nil
}
