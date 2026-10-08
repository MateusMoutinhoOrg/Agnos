package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/projectconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// ConfigDir is the directory every declaration of a project lives in: the
// generator's own name, title-cased, then Config — AgnosConfig for agnos.
func ConfigDir(sandbox *api.Sandbox) string {
	name := sandbox.Config.ProjectName
	if name == "" {
		return "Config"
	}
	return sandbox.Deps.StringsDeps.ToUpper(name[:1]) + name[1:] + "Config"
}

// ProjectConfPath is the project-relative path of the file `agnos start`
// writes once and every later command reads.
func ProjectConfPath(sandbox *api.Sandbox) string {
	return ConfigDir(sandbox) + "/project.yaml"
}

// ValidateProjectName reports whether a project name can be what it becomes:
// the name of the binary, the verb its usage line starts with, the alias
// run-examples runs it under and the file local-install writes. So it starts with
// an ASCII letter and holds only ASCII letters, digits, dashes and underscores.
func ValidateProjectName(sandbox *api.Sandbox, name string) error {
	if name == "" {
		return sandbox.Deps.StdDeps.Errorf("a project needs a name (--project-name)")
	}
	first := name[0]
	if !((first >= 'a' && first <= 'z') || (first >= 'A' && first <= 'Z')) {
		return sandbox.Deps.StdDeps.Errorf("invalid project name %q: it must start with an ASCII letter (it becomes the binary's name)", name)
	}
	for i := 0; i < len(name); i++ {
		letter := name[i]
		valid := (letter >= 'a' && letter <= 'z') ||
			(letter >= 'A' && letter <= 'Z') ||
			(letter >= '0' && letter <= '9') ||
			letter == '-' || letter == '_'
		if !valid {
			return sandbox.Deps.StdDeps.Errorf("invalid project name %q: only ASCII letters, digits, dashes and underscores are allowed (it becomes the binary's name)", name)
		}
	}
	return nil
}

// RequireProject fails when io is not rooted at a project: no project.yaml
// means `start` never ran there, which is the one thing worth saying — not that
// some layer of a project that does not exist is missing.
func RequireProject(sandbox *api.Sandbox, io *stagedfs.StagedFS) error {
	if io.IsFile(ProjectConfPath(sandbox)) {
		return nil
	}
	root := io.Root
	if root == "" {
		root = "."
	}
	return sandbox.Deps.StdDeps.Errorf("%s is not a project (no %s): run `start` there first, or pass --path", root, ProjectConfPath(sandbox))
}

// LoadProjectConf reads AgnosConfig/project.yaml back through the
// transaction-aware io (so it is visible during `agnos start`, before
// Persist). It never falls back to empty defaults: `agnos start` is a
// prerequisite for every other command, so a missing or unparsable
// project.yaml is a hard error.
func LoadProjectConf(sandbox *api.Sandbox, io *stagedfs.StagedFS) (*projectconf.ProjectConf, error) {
	rel := ProjectConfPath(sandbox)

	content, err := io.ReadFile(rel)
	if err != nil {
		return nil, sandbox.Deps.StdDeps.Errorf("could not read %s: run `agnos start` first (%w)", rel, err)
	}

	return projectconf.New(sandbox, string(content))
}
