package build

import (
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/parsables/moduleconf"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/parsables/projectconf"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/utils"
)

// loadProjectConf reads <ProjectName>Config/project.yaml back through the
// transaction-aware io (so it is visible during `agnos start`, before Persist).
// It never falls back to empty defaults: `agnos start` is a prerequisite for
// `agnos build`, so a missing or unparsable project.yaml is a hard error.
func loadProjectConf(deps *deps.Deps, io *smartio.SmartIO, path string) (*projectconf.ProjectConf, error) {
	rel := config.ProjectName + "Config/project.yaml"

	// The `start` asset group writes project.yaml under its bare relative path,
	// so try that first (visible in the transaction before Persist); fall back
	// to the path-prefixed location for `agnos build <dir>`.
	content, err := io.ReadFile(rel)
	if err != nil {
		content, err = io.ReadFile(path + "/" + rel)
	}
	if err != nil {
		return nil, deps.Std.Errorf("could not read %s: run `agnos start` first (%w)", rel, err)
	}

	return projectconf.New(deps, string(content))
}

// projectNameConst title-cases the configured project name for use as the
// generated config.ProjectName constant (which names the <X>Config/ dir).
func projectNameConst(name string) string {
	if len(name) == 0 {
		return config.ProjectName
	}
	return strings.ToUpper(name[:1]) + name[1:]
}

func BuildInternal(deps *deps.Deps, io *smartio.SmartIO, path string) error {
	deps.Std.Printf("build started with path %s \n", path)

	//Creating the basic dir struct
	io.CreateDir("sandbox/api")
	io.CreateDir("sandbox/internal")

	gomod, err := io.ReadFile(path + "/go.mod")
	if err != nil {
		return err
	}
	module_conf, err := moduleconf.New(deps, string(gomod))
	if err != nil {
		return err
	}

	hasDeps := io.IsDir("sandbox/deps")
	hasCli := io.IsDir("sandbox/internal/cli")

	project_conf, err := loadProjectConf(deps, io, path)
	if err != nil {
		return err
	}

	commands, err := CollectCommands(deps, io)
	if err != nil {
		return err
	}

	vars := map[string]interface{}{
		"Module":       module_conf.Module,
		"Name":         project_conf.Name,
		"Description":  project_conf.Description,
		"Version":      project_conf.Version,
		"ProjectName":  projectNameConst(project_conf.Name),
		"HasDeps":      hasDeps,
		"HasCli":       hasCli,
		"Binds":        CollectBinds(io),
		"Constructors": CollectConstructors(io),
		"DepsLibs":     CollectDepsLibs(io),
		"AdapterLibs":  CollectAdapterLibs(io),
		"Commands":     commands,
	}

	if err := utils.RenderGroup(deps, io, "all", vars); err != nil {
		return err
	}

	if hasDeps {
		if err := utils.RenderGroup(deps, io, "deps", vars); err != nil {
			return err
		}
	}

	if hasCli {
		if err := GenerateCommandEntries(deps, io, commands); err != nil {
			return err
		}
		if err := utils.RenderGroup(deps, io, "cli", vars); err != nil {
			return err
		}
	}

	deps.Std.Printf("successfully rendered template")
	return nil
}
