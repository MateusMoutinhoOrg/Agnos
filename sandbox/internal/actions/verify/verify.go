package verify

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
)

// Verify checks that the project at path keeps the sandbox/adapter schema the
// harness depends on. It performs no filesystem writes, so it never calls
// io.Persist. `agnos build` runs it as a gate before every build unless the
// caller passes --unsafe.
func Verify(deps *deps.Deps, path string) error {
	io := smartio.New(deps, path, config.ProjectName)
	return VerifyInternal(deps, io, path)
}
