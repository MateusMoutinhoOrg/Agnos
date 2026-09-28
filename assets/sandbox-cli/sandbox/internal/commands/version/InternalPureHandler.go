package version

import (
	"{{.Module}}/sandbox/api"
)

// InternalPureHandler backs `version`. Nothing is read off the command line: it
// declares no flag and no arg but its verb.
func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {
	if sandbox.Config.Version == "" {
		response.Printf("no version set yet\n")
		return nil
	}
	response.Printf("Version: %s\n", sandbox.Config.Version)
	return nil
}
