package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// errorsDir holds the project's own answer to every way a request can end
// without a route answering it, beside the generated registry
// sandbox/internal/server/generated.new.go.
const errorsDir = "sandbox/internal/server/errors"

// errorHandlerFiles is one file per failure the dispatch can raise, in the
// order sandbox/internal/server/generated.new.go switches on them. Each holds
// one handler with the route handler's own signature.
var errorHandlerFiles = []string{
	"handle_not_found.go",
	"handle_method_not_allowed.go",
	"handle_bad_request.go",
	"handle_unauthorized.go",
	"handle_forbidden.go",
	"handle_payload_too_large.go",
	"handle_unsupported_media_type.go",
	"handle_internal_server_error.go",
}

// GenerateServerErrorHandlers renders assets/templates/server_handle_*.go into
// sandbox/internal/server/errors/ — the eight handlers the generated.new.go
// hands a failure to, one per status.
//
// It is written **once**, the same way a constructor is. A handler already on
// disk is left exactly as it is, however far it has drifted from what the
// template renders: what a 404 says, and whether it is even a 404, is the
// project's. That is also why these are not part of the server asset
// group — every file of a group is rewritten by every build.
//
// Writing them here rather than in server-init is what carries a project that
// ran server-init before they existed: the next build fills in what is
// missing, and the generated.new.go always has something to call.
func GenerateServerErrorHandlers(sandbox *api.Sandbox, io *stagedfs.StagedFS, module string) error {
	for _, handler := range errorHandlerFiles {
		dest := errorsDir + "/" + handler
		if io.IsFile(dest) {
			continue
		}

		vars := map[string]any{
			"Module":        module,
			"GeneratorName": generatorName(sandbox),
		}

		if err := utils.RenderTemplateToDest(sandbox, io, "templates/server_"+handler, vars, dest); err != nil {
			return err
		}
	}
	return nil
}
