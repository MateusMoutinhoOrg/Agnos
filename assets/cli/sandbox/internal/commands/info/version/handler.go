package version

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/internal/commandprops"
)

// Handle backs `version`. Nothing is read off the command line: it
// declares no flag and no arg but its verb.
func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	if sandbox.Config.Version == "" {
		response.Printf("no version set yet\n")
		return nil
	}
	response.Printf("Version: %s\n", sandbox.Config.Version)
	return nil
}
