package database_dir

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/internal/commandprops"
)

// Handle runs in front of every command: it points sandbox.Config.DatabaseDir
// at the folder --database names, so every database a command opens lives
// under it. The dispatch already refused a value its pattern does not match —
// an absolute path, an upper-case letter, a dot — so the folder the store
// writes to is the one a backup reads from. An empty --database keeps the
// default. It answers nothing, so the command after it runs.
func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	if input.Database != "" {
		sandbox.Config.DatabaseDir = input.Database
	}
	return nil
}
