package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/moduleconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// LoadModuleConf reads and parses the target project's go.mod. A missing file
// is the ordinary "you are not standing in a project" case, so it is reported
// in those words rather than as the raw `open go.mod: …` the filesystem gives.
func LoadModuleConf(deps *deps.Deps, io *smartio.SmartIO) (*moduleconf.ModuleConf, error) {
	content, err := io.ReadFile("go.mod")
	if err != nil {
		return nil, deps.Std.Errorf("no project found: go.mod is missing (run `agnos start` first, or pass --path)")
	}
	conf, err := moduleconf.New(deps, string(content))
	if err != nil {
		return nil, deps.Std.Errorf("go.mod: %w", err)
	}
	return conf, nil
}
