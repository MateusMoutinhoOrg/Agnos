package verify

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// Verify checks that the project at path keeps the sandbox/adapter schema the
// harness depends on. It performs no filesystem writes, so it never calls
// io.Persist. `agnos build` runs it as a gate before every build unless the
// caller passes --unsafe.
func Verify(sandbox *api.Sandbox, path string) error {
	io := smartio.New(sandbox, path, config.ProjectName)
	return VerifyInternal(sandbox, io, path)
}
