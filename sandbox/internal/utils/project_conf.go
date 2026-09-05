package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/projectconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// ProjectConfPath is the project-relative path of the file `agnos start`
// writes once and every later command reads.
func ProjectConfPath() string {
	return config.ProjectName + "Config/project.yaml"
}

// LoadProjectConf reads <ProjectName>Config/project.yaml back through the
// transaction-aware io (so it is visible during `agnos start`, before
// Persist). It never falls back to empty defaults: `agnos start` is a
// prerequisite for every other command, so a missing or unparsable
// project.yaml is a hard error.
func LoadProjectConf(deps *deps.Deps, io *smartio.SmartIO) (*projectconf.ProjectConf, error) {
	rel := ProjectConfPath()

	content, err := io.ReadFile(rel)
	if err != nil {
		return nil, deps.Std.Errorf("could not read %s: run `agnos start` first (%w)", rel, err)
	}

	return projectconf.New(deps, string(content))
}
