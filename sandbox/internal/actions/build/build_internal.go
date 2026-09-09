package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// projectNameConst title-cases the configured project name for use as the
// generated config.ProjectName constant (which names the <X>Config/ dir).
func projectNameConst(deps *deps.Deps, name string) string {
	if len(name) == 0 {
		return config.ProjectName
	}
	return deps.Stringsdeps.ToUpper(name[:1]) + name[1:]
}

// generatorName is the cli name of the binary running this build — the
// generator, never the project being generated. Docs that spell a command of
// the generator ("agnos add-route") render it from here, while a command of the
// generated project ("<name> start-server") renders from the project's own Name.
func generatorName(deps *deps.Deps) string {
	return deps.Stringsdeps.ToLower(config.ProjectName)
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
	hasServer := io.IsDir("sandbox/internal/server")

	// hasAssets reports that the project carries its own agnos asset groups —
	// it is itself a generator, like agnos. The docs of such a project have to
	// name its templates and its own bootstrap; every other project has no
	// assets/ tree to be told about.
	hasAssets := io.IsDir("assets/all")

	project_conf, err := utils.LoadProjectConf(deps, io)
	if err != nil {
		return err
	}

	// The help command's entries.yaml is generated, so it is written before
	// the commands are collected — from there on help is just another entry
	// in the set.
	if hasCli {
		helpVars := map[string]interface{}{
			"Module":      module_conf.Module,
			"ProjectName": projectNameConst(deps, project_conf.Name),
		}
		if err := GenerateHelpEntriesYaml(deps, io, helpVars); err != nil {
			return err
		}
	}

	commands, err := CollectCommands(deps, io)
	if err != nil {
		return err
	}

	// The server layer's mirror of CollectCommands: one entry per declared
	// route, already ordered for matching so servermain.go only has to range.
	routes, err := CollectRoutes(deps, io)
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

	// docs/Structure's tree is rendered from the project's own structure.yaml,
	// so the page describes the shape its author declared rather than one
	// typed into the doc by hand. `verify` rejects an item whose path is gone.
	structure, err := CollectStructure(deps, io)
	if err != nil {
		return err
	}

	// docs/Commands is rendered from the command declarations themselves, so
	// every visible command, flag, argument and example on the page is the one
	// its entries.yaml declares.
	command_docs, err := CollectCommandDocs(deps, io)
	if err != nil {
		return err
	}

	// docs/Routes is rendered from the route declarations themselves, the
	// same way docs/Commands is rendered from the command ones.
	route_docs, err := CollectRouteDocs(deps, io)
	if err != nil {
		return err
	}

	// The docs this build generates are merged in before the index is built:
	// SmartIO listings read disk, so on a project's first build they are not
	// there to be walked yet.
	generated_docs, err := CollectGeneratedDocs(deps, io, docsVars(module_conf.Module, project_conf.Name, generatorName(deps)), GeneratedDocsGroups(hasCli, hasServer))
	if err != nil {
		return err
	}
	docs = MergeDocs(deps, docs, generated_docs)

	if err := GenerateSubdocIndexes(deps, io, docs); err != nil {
		return err
	}

	vars := map[string]interface{}{
		"Module":            module_conf.Module,
		"Name":              project_conf.Name,
		"Version":           project_conf.Version,
		"ProjectName":       projectNameConst(deps, project_conf.Name),
		"GeneratorName":     generatorName(deps),
		"ConfigDir":         config.ProjectName + "Config",
		"StructureConfFile": utils.StructureConfFile,
		"HasDeps":           hasDeps,
		"HasCli":            hasCli,
		"HasServer":         hasServer,
		"HasAssets":         hasAssets,
		"Binds":             CollectBinds(deps, io),
		"Constructors":      CollectConstructors(deps, io),
		"DepsLibs":          CollectDepsLibs(deps, io),
		"AdapterLibs":       CollectAdapterLibs(deps, io),
		"CliExamples":       utils.CollectExamples(deps, io, utils.ExampleCliSide),
		"LibExamples":       utils.CollectExamples(deps, io, utils.ExampleLibSide),
		"Commands":          commands,
		"CommandDocs":       command_docs,
		"Routes":            routes,
		"RouteDocs":         route_docs,
		"Themes":            themes_conf.Themes,
		"DocIndex":          CollectDocIndex(deps, docs, themes_conf.Themes),
		"PublicApi":         public_api,
		"Structure":         structure,
		"DepsApi":           deps_api,
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

	if hasServer {
		if err := GenerateRouteEntries(deps, io, routes, module_conf.Module); err != nil {
			return err
		}
		if err := utils.RenderGroup(deps, io, "server", vars); err != nil {
			return err
		}
	}

	deps.Std.Log("successfully rendered template\n")
	return nil
}
