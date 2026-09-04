package build

import (
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/projectconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// loadProjectConf reads <ProjectName>Config/project.yaml back through the
// transaction-aware io (so it is visible during `agnos start`, before Persist).
// It never falls back to empty defaults: `agnos start` is a prerequisite for
// `agnos build`, so a missing or unparsable project.yaml is a hard error.
func loadProjectConf(deps *deps.Deps, io *smartio.SmartIO, path string) (*projectconf.ProjectConf, error) {
	rel := config.ProjectName + "Config/project.yaml"

	// io resolves rel against the target project root (the --path flag), so
	// this reads the file the `start` asset group wrote — whether or not it
	// has been persisted yet.
	content, err := io.ReadFile(rel)
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
	deps.Std.Log("build started with path %s \n", path)

	//Creating the basic dir struct
	io.CreateDir("sandbox/api")
	io.CreateDir("sandbox/internal")

	module_conf, err := utils.LoadModuleConf(deps, io)
	if err != nil {
		return err
	}

	hasDeps := io.IsDir("sandbox/deps")
	hasCli := io.IsDir("sandbox/internal/cli")

	project_conf, err := loadProjectConf(deps, io, path)
	if err != nil {
		return err
	}

	// The help command's entries.yaml is generated, so it is written before
	// the commands are collected — from there on help is just another entry
	// in the set.
	if hasCli {
		helpVars := map[string]interface{}{
			"Module":      module_conf.Module,
			"ProjectName": projectNameConst(project_conf.Name),
		}
		if err := GenerateHelpEntriesYaml(deps, io, helpVars); err != nil {
			return err
		}
	}

	commands, err := CollectCommands(deps, io)
	if err != nil {
		return err
	}

	themes_conf, err := utils.LoadThemesConf(deps, io)
	if err != nil {
		return err
	}

	// The documentation tree is indexed before anything is rendered, so
	// README.md's index and every sub-doc Index.md come out of the same walk
	// in one build.
	docs, err := CollectDocs(deps, io)
	if err != nil {
		return err
	}

	// docs/PublicApi/doc.md is rendered from the contract sources themselves,
	// so the public surface and its description are always the ones the code
	// declares. `verify` keeps those sources parsable and commented.
	public_api, err := CollectPublicApi(deps, io)
	if err != nil {
		return err
	}

	deps_api, err := CollectDepsApi(deps, io)
	if err != nil {
		return err
	}

	// The docs this build generates are merged in before the index is built:
	// SmartIO listings read disk, so on a project's first build they are not
	// there to be walked yet.
	generated_docs, err := CollectGeneratedDocs(deps, io, docsVars(module_conf.Module, project_conf.Name))
	if err != nil {
		return err
	}
	docs = MergeDocs(docs, generated_docs)

	if err := GenerateSubdocIndexes(deps, io, docs); err != nil {
		return err
	}

	vars := map[string]interface{}{
		"Module":       module_conf.Module,
		"Name":         project_conf.Name,
		"Version":      project_conf.Version,
		"ProjectName":  projectNameConst(project_conf.Name),
		"ConfigDir":    config.ProjectName + "Config",
		"HasDeps":      hasDeps,
		"HasCli":       hasCli,
		"Binds":        CollectBinds(io),
		"Constructors": CollectConstructors(io),
		"DepsLibs":     CollectDepsLibs(io),
		"AdapterLibs":  CollectAdapterLibs(io),
		"Commands":     commands,
		"Themes":       themes_conf.Themes,
		"DocIndex":     CollectDocIndex(docs, themes_conf.Themes),
		"PublicApi":    public_api,
		"DepsApi":      deps_api,
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

	deps.Std.Log("successfully rendered template\n")
	return nil
}
