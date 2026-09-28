package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/projectconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// ProjectConfPath is the project-relative path of the file `agnos start`
// writes once and every later command reads.
func ProjectConfPath(sandbox *api.Sandbox) string {
	return sandbox.Config.ProjectName + "Config/project.yaml"
}

// ValidateProjectName reports whether a project name can be what it becomes:
// the name of the binary, the verb its usage line starts with, the alias
// exec-test runs it under and the file local-install writes. So it starts with
// an ASCII letter and holds only ASCII letters, digits, dashes and underscores.
func ValidateProjectName(sandbox *api.Sandbox, name string) error {
	if name == "" {
		return sandbox.Deps.Std.Errorf("a project needs a name (--project-name)")
	}
	first := name[0]
	if !((first >= 'a' && first <= 'z') || (first >= 'A' && first <= 'Z')) {
		return sandbox.Deps.Std.Errorf("invalid project name %q: it must start with an ASCII letter (it becomes the binary's name)", name)
	}
	for i := 0; i < len(name); i++ {
		letter := name[i]
		valid := (letter >= 'a' && letter <= 'z') ||
			(letter >= 'A' && letter <= 'Z') ||
			(letter >= '0' && letter <= '9') ||
			letter == '-' || letter == '_'
		if !valid {
			return sandbox.Deps.Std.Errorf("invalid project name %q: only ASCII letters, digits, dashes and underscores are allowed (it becomes the binary's name)", name)
		}
	}
	return nil
}

// RequireProject fails when io is not rooted at a project: no project.yaml
// means `start` never ran there, which is the one thing worth saying — not that
// some layer of a project that does not exist is missing.
func RequireProject(sandbox *api.Sandbox, io *smartio.SmartIO) error {
	if io.IsFile(ProjectConfPath(sandbox)) {
		return nil
	}
	root := io.Root
	if root == "" {
		root = "."
	}
	return sandbox.Deps.Std.Errorf("%s is not a project (no %s): run `start` there first, or pass --path", root, ProjectConfPath(sandbox))
}

// LoadProjectConf reads <ProjectName>Config/project.yaml back through the
// transaction-aware io (so it is visible during `agnos start`, before
// Persist). It never falls back to empty defaults: `agnos start` is a
// prerequisite for every other command, so a missing or unparsable
// project.yaml is a hard error.
func LoadProjectConf(sandbox *api.Sandbox, io *smartio.SmartIO) (*projectconf.ProjectConf, error) {
	rel := ProjectConfPath(sandbox)

	content, err := io.ReadFile(rel)
	if err != nil {
		return nil, sandbox.Deps.Std.Errorf("could not read %s: run `agnos start` first (%w)", rel, err)
	}

	return projectconf.New(sandbox, string(content))
}
