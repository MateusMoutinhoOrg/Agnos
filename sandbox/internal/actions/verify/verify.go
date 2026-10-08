package verify

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// Verify checks that the project at path keeps the sandbox/adapter schema the
// harness depends on. It performs no filesystem writes, so it never calls
// io.Persist. `agnos build` runs it as a gate before every build unless the
// caller passes --unsafe.
func Verify(sandbox *api.Sandbox, props api.VerifyProps) error {
	io := stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName)
	return VerifyInternal(sandbox, io, props.Path)
}
