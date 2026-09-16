package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// projectNameConst title-cases the configured project name for use as the
// generated Config.ProjectName of the project (which names the <X>Config/ dir).
func projectNameConst(sandbox *api.Sandbox, name string) string {
	if len(name) == 0 {
		return sandbox.Config.ProjectName
	}
	return sandbox.Deps.Stringsdeps.ToUpper(name[:1]) + name[1:]
}

// generatorName is the cli name of the binary running this build — the
// generator, never the project being generated. Docs that spell a command of
// the generator ("agnos add-route") render it from here, while a command of the
// generated project ("<name> start-server") renders from the project's own Name.
func generatorName(sandbox *api.Sandbox) string {
	return sandbox.Deps.Stringsdeps.ToLower(sandbox.Config.ProjectName)
}

func BuildInternal(sandbox *api.Sandbox, io *smartio.SmartIO, path string) error {
	sandbox.Deps.Std.Log("build started with path %s \n", path)

	module_conf, err := utils.LoadModuleConf(sandbox, io)
	if err != nil {
		return err
	}

	// Which mechanics this project wants is declared, never inferred from the
	// directories it happens to carry. A key the catalog gained since the
	// project was scaffolded is filled with its default and written back, so a
	// new extension never breaks an older tree.
	extensions_conf, err := utils.LoadExtensionsConf(sandbox, io)
	if err != nil {
		return err
	}
	if utils.NormalizeExtensions(extensions_conf) {
		if err := utils.SaveExtensionsConf(sandbox, io, extensions_conf); err != nil {
			return err
		}
	}

	hasSandbox := extensions_conf.IsEnabled(utils.ExtensionSandbox)
	hasDeps := extensions_conf.IsEnabled(utils.ExtensionSandboxDeps)
	hasCli := extensions_conf.IsEnabled(utils.ExtensionSandboxCli)
	hasServer := extensions_conf.IsEnabled(utils.ExtensionSandboxServer)
	hasFront := extensions_conf.IsEnabled(utils.ExtensionSandboxFront)
	hasExample := extensions_conf.IsEnabled(utils.ExtensionSandboxExample)
	hasDoc := extensions_conf.IsEnabled(utils.ExtensionDoc)
	hasReadme := extensions_conf.IsEnabled(utils.ExtensionReadme)

	if hasSandbox {
		io.CreateDir("sandbox/api")
		io.CreateDir("sandbox/internal")
	}

	// The contracts of sandbox/api/ and the packages that build them. The
	// first is one field of the Sandbox each; the second is the call list
	// sandbox/new.go is rendered from, and it is wider than the first — a
	// constructor the project wrote itself is in it too.
	constructors := CollectConstructors(sandbox, io)

	// One constructor package per contract, written the first time and never
	// again: from there on how that field is built is the project's.
	if hasSandbox {
		if err := GenerateConstructors(sandbox, io, constructors, module_conf.Module); err != nil {
			return err
		}
	}

	// hasAssets reports that the project carries its own agnos asset groups —
	// it is itself a generator, like agnos. The docs of such a project have to
	// name its templates and its own bootstrap; every other project has no
	// assets/ tree to be told about. It is a property of the tree, not a
	// mechanic the project turns on.
	hasAssets := io.IsDir("assets/" + utils.ExtensionSandbox)

	project_conf, err := utils.LoadProjectConf(sandbox, io)
	if err != nil {
		return err
	}

	// The help command's entries.yaml is generated, so it is written before
	// the commands are collected — from there on help is just another entry
	// in the set.
	if hasCli {
		helpVars := map[string]interface{}{
			"Module":      module_conf.Module,
			"ProjectName": projectNameConst(sandbox, project_conf.Name),
		}
		if err := GenerateHelpEntriesYaml(sandbox, io, helpVars); err != nil {
			return err
		}
	}

	commands, err := CollectCommands(sandbox, io)
	if err != nil {
		return err
	}

	// The server layer's mirror of CollectCommands: one entry per declared
	// route, already ordered for matching so the dispatch only has to range
	// over Server.Routes.
	routes, err := CollectRoutes(sandbox, io)
	if err != nil {
		return err
	}

	themes_conf, err := utils.LoadThemesConf(sandbox, io)
	if err != nil {
		return err
	}

	// The documentation tree is indexed before anything is rendered, so
	// README.md's index and every sub-doc Index.md come out of the same walk
	// in one build.
	docs, err := CollectDocs(sandbox, io)
	if err != nil {
		return err
	}

	// docs/PublicApi/doc.md is rendered from the contract sources themselves,
	// so the public surface and its description are always the ones the code
	// declares. `verify` keeps those sources parsable and commented.
	public_api, err := CollectPublicApi(sandbox, io)
	if err != nil {
		return err
	}

	deps_api, err := CollectDepsApi(sandbox, io)
	if err != nil {
		return err
	}

	// An available is a declared selection, not a directory listing: two
	// adapters may implement the same contract, so which one binds is read
	// from adapters/availables/<name>/available.yaml and nowhere else.
	availables, err := CollectAvailables(sandbox, io)
	if err != nil {
		return err
	}

	// docs/Structure's tree is rendered from the project's own structure.yaml,
	// so the page describes the shape its author declared rather than one
	// typed into the doc by hand. `verify` rejects an item whose path is gone.
	structure, err := CollectStructure(sandbox, io)
	if err != nil {
		return err
	}

	// docs/Commands is rendered from the command declarations themselves, so
	// every visible command, flag, argument and example on the page is the one
	// its entries.yaml declares.
	command_docs, err := CollectCommandDocs(sandbox, io)
	if err != nil {
		return err
	}

	// docs/Routes is rendered from the route declarations themselves, the
	// same way docs/Commands is rendered from the command ones.
	route_docs, err := CollectRouteDocs(sandbox, io)
	if err != nil {
		return err
	}

	// The docs this build generates are merged in before the index is built:
	// SmartIO listings read disk, so on a project's first build they are not
	// there to be walked yet.
	generated_docs, err := CollectGeneratedDocs(sandbox, io, docsVars(module_conf.Module, project_conf.Name, generatorName(sandbox)), utils.DocGroups(sandbox, extensions_conf))
	if err != nil {
		return err
	}
	docs = MergeDocs(sandbox, docs, generated_docs)

	if hasDoc {
		if err := GenerateSubdocIndexes(sandbox, io, docs); err != nil {
			return err
		}
	}

	vars := map[string]interface{}{
		"Module":              module_conf.Module,
		"Name":                project_conf.Name,
		"Version":             project_conf.Version,
		"ProjectName":         projectNameConst(sandbox, project_conf.Name),
		"GeneratorName":       generatorName(sandbox),
		"ConfigDir":           sandbox.Config.ProjectName + "Config",
		"StructureConfFile":   utils.StructureConfFile,
		"HasSandbox":          hasSandbox,
		"HasDeps":             hasDeps,
		"HasCli":              hasCli,
		"HasServer":           hasServer,
		"HasFront":            hasFront,
		"HasExample":          hasExample,
		"HasDoc":              hasDoc,
		"HasReadme":           hasReadme,
		"StaticMount":         CollectFrontMount(sandbox, io),
		"HasAssets":           hasAssets,
		"Constructors":        constructors,
		"ConstructorPackages": CollectConstructorPackages(sandbox, io, constructors),
		"DepsLibs":            CollectDepsLibs(sandbox, io),
		"AdapterLibs":         CollectAdapterLibs(sandbox, io),
		"Availables":          availables,
		"CliExamples":         utils.CollectExamples(sandbox, io, utils.ExampleCliSide),
		"LibExamples":         utils.CollectExamples(sandbox, io, utils.ExampleLibSide),
		"Commands":            commands,
		"CommandDocs":         command_docs,
		"Routes":              routes,
		"RouteDocs":           route_docs,
		"Themes":              themes_conf.Themes,
		"DocIndex":            CollectDocIndex(sandbox, docs, themes_conf.Themes),
		"PublicApi":           public_api,
		"Structure":           structure,
		"DepsApi":             deps_api,
	}

	// The per-unit generators: one new.go per declared available, command and
	// route, each owned by the mechanic that declares the unit.
	if hasDeps {
		if err := GenerateAvailableNews(sandbox, io, availables, module_conf.Module); err != nil {
			return err
		}
	}

	if hasCli {
		if err := GenerateCommandNew(sandbox, io, commands, module_conf.Module); err != nil {
			return err
		}
	}

	if hasServer {
		if err := GenerateRouteNew(sandbox, io, routes, module_conf.Module); err != nil {
			return err
		}
	}

	// Every asset group the declaration turns on, in render order. A mechanic
	// that is off is not rendered at all: what it wrote before stays on disk,
	// untouched, for the project to keep by hand.
	for _, group := range utils.RenderableGroups(extensions_conf) {
		if err := utils.RenderGroup(sandbox, io, group, vars); err != nil {
			return err
		}
	}

	sandbox.Deps.Std.Log("successfully rendered template\n")
	return nil
}
