package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// cliErrorsDir is where the project's own answers to a failed command line
// live — hand-written from the first build on, like a command's
// handler.go.
const cliErrorsDir = "sandbox/internal/cli/errors"

// cliErrorHandlers are the five files of cliErrorsDir, one per
// api.CommandFailureKind, each rendered from assets/templates/cli_<file>.
var cliErrorHandlers = []string{
	"handle_not_found.go",
	"handle_bad_usage.go",
	"handle_unknown_flag.go",
	"handle_unexpected_arg.go",
	"handle_failure.go",
}

// GenerateCliErrorHandlers writes the five handle_*.go of
// sandbox/internal/cli/errors/ that are missing, and leaves every one already
// there alone. They are the cli's mirror of the server's eight: not assets of
// the cli group, since every file of a group is rewritten by every build, and
// these are the project's once written.
func GenerateCliErrorHandlers(sandbox *api.Sandbox, io *stagedfs.StagedFS, module string) error {
	vars := map[string]any{
		"Module":        module,
		"GeneratorName": generatorName(sandbox),
	}
	for _, file := range cliErrorHandlers {
		dest := cliErrorsDir + "/" + file
		if io.IsFile(dest) {
			continue
		}
		if err := utils.RenderTemplateToDest(sandbox, io, "templates/cli_"+file, vars, dest); err != nil {
			return err
		}
	}
	return nil
}
