package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/moduleconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// LoadModuleConf reads and parses the target project's go.mod. A missing file
// is the ordinary "you are not standing in a project" case, so it is reported
// in those words rather than as the raw `open go.mod: …` the filesystem gives.
func LoadModuleConf(sandbox *api.Sandbox, io *stagedfs.StagedFS) (*moduleconf.ModuleConf, error) {
	content, err := io.ReadFile("go.mod")
	if err != nil {
		return nil, sandbox.Deps.StdDeps.Errorf("no project found: go.mod is missing (run `agnos start` first, or pass --path)")
	}
	conf, err := moduleconf.New(sandbox, string(content))
	if err != nil {
		return nil, sandbox.Deps.StdDeps.Errorf("go.mod: %w", err)
	}
	return conf, nil
}
