package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// generatorName is the cli name of the binary running this build — the
// generator, never the project being generated. Docs that spell a command of
// the generator ("agnos add-route") render it from here, while a command of the
// generated project ("<name> start-server") renders from the project's own Name.
func generatorName(sandbox *api.Sandbox) string {
	return sandbox.Deps.StringsDeps.ToLower(sandbox.Config.ProjectName)
}

// generatorVersion is the release of the binary running this build. A tree was
// rendered by this version, so it is the version a doc names as the one that
// maintains it — never the project's own Version, which is what the generated
// project releases under.
func generatorVersion(sandbox *api.Sandbox) string {
	return sandbox.Config.Version
}

func BuildInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, path string) error {
	sandbox.Deps.StdDeps.Logf("build started with path %s \n", path)

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
	hasDeps := extensions_conf.IsEnabled(utils.ExtensionDeps)
	hasCli := extensions_conf.IsEnabled(utils.ExtensionCli)
	hasServer := extensions_conf.IsEnabled(utils.ExtensionServer)
	hasFront := extensions_conf.IsEnabled(utils.ExtensionFront)
	hasDatabase := extensions_conf.IsEnabled(utils.ExtensionDatabase)
	hasExample := extensions_conf.IsEnabled(utils.ExtensionExample)
	hasBackoffice := extensions_conf.IsEnabled(utils.ExtensionBackoffice)
	hasDoc := extensions_conf.IsEnabled(utils.ExtensionDoc)
	hasReadme := extensions_conf.IsEnabled(utils.ExtensionReadme)

	if hasSandbox {
		io.CreateDir("sandbox/api")
		io.CreateDir("sandbox/internal")
	}

	// A props struct an older build wrote into sandbox/api/ is moved to its
	// own package first: from here on nothing reads it as a contract.
	if err := MigrateLegacyProps(sandbox, io, hasCli, hasServer); err != nil {
		return err
	}
	// So is a part of api.Sandbox or api.Config named under the old prefix
	// (userconfig_backoffice.go): the aggregates embed by suffix now.
	moved_parts, err := MigrateLegacyApiParts(sandbox, io)
	if err != nil {
		return err
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

	// The help command's command.yaml is generated, so it is written before
	// the commands are collected — from there on help is just another entry
	// in the set.
	if hasCli {
		helpVars := map[string]interface{}{
			"Module":        module_conf.Module,
			"ProjectName":   project_conf.ProjectName,
			"GeneratorName": generatorName(sandbox),
		}
		if err := GenerateHelpCommandYaml(sandbox, io, helpVars); err != nil {
			return err
		}
	}

	commands, err := CollectCommands(sandbox, io)
	if err != nil {
		return err
	}

	// The routes the server group writes itself are declared the same way:
	// their route.yaml is rendered before the routes are collected, so a
	// project whose server group just gained one collects it in this build.
	if hasServer {
		if err := GenerateBuiltinRouteYamls(sandbox, io, map[string]interface{}{
			"Module":        module_conf.Module,
			"ProjectName":   project_conf.ProjectName,
			"GeneratorName": generatorName(sandbox),
		}); err != nil {
			return err
		}
	}

	// The server layer's mirror of CollectCommands: one entry per declared
	// route, already ordered for matching so the dispatch only has to range
	// over Server.Routes.
	routes, err := CollectRoutes(sandbox, io)
	if err != nil {
		return err
	}

	// The database layer's mirror of CollectRoutes: one entry per declared
	// database, already carrying the records and the signatures its three
	// generated files are spelled from.
	databases, err := CollectDatabases(sandbox, io)
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

	dep_contracts, err := CollectDepContracts(sandbox, io)
	if err != nil {
		return err
	}

	// A binding is a declared selection, not a directory listing: two
	// adapters may implement the same contract, so which one binds is read
	// from adapters/bindings/<name>/binding.yaml and nowhere else.
	bindings, err := CollectBindings(sandbox, io)
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
	// its command.yaml declares.
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

	// The OpenAPI document is rendered from the same declarations docs/Routes
	// is, once: docs/Routes/openapi.json holds it and the generated openapi
	// route answers it, so both are always the same bytes.
	open_api := ""
	if hasServer {
		open_api, err = CollectOpenApi(sandbox, io, project_conf.ProjectName, project_conf.Version)
		if err != nil {
			return err
		}
	}

	// docs/Databases is rendered from the database declarations themselves, the
	// same way docs/Routes is rendered from the route ones.
	database_docs, err := CollectDatabaseDocs(sandbox, io)
	if err != nil {
		return err
	}

	// The docs this build generates are merged in before the index is built:
	// StagedFS listings read disk, so on a project's first build they are not
	// there to be walked yet.
	generated_docs, err := CollectGeneratedDocs(sandbox, io, docsVars(module_conf.Module, project_conf.ProjectName, generatorName(sandbox)), utils.DocGroups(sandbox, extensions_conf))
	if err != nil {
		return err
	}
	docs = MergeDocs(sandbox, docs, generated_docs)

	if hasDoc {
		if err := GenerateSubdocIndexes(sandbox, io, docs); err != nil {
			return err
		}
	}

	// The parts api.Config and api.Sandbox embed: every struct of a
	// sandbox/api/*config.go and *sandbox.go file, so each mechanic adds a part
	// of its own instead of editing another's. The parts the enabled groups
	// write are rendered first, and named along with the two start writes and
	// the parts migrated above: start builds before it persists, and the
	// listing reads disk.
	api_parts := append([]string{"sandbox/api/" + utils.ProjectConfigFile, "sandbox/api/" + utils.ProjectSandboxFile}, moved_parts...)
	if hasSandbox {
		group_parts, err := GenerateApiParts(sandbox, io, utils.RenderableGroups(extensions_conf), map[string]interface{}{
			"Module":        module_conf.Module,
			"GeneratorName": generatorName(sandbox),
		})
		if err != nil {
			return err
		}
		api_parts = append(api_parts, group_parts...)
	}
	config_structs, err := utils.CollectEmbeddedStructs(sandbox, io, "sandbox/api", utils.ConfigParts(sandbox), api_parts)
	if err != nil {
		return err
	}
	sandbox_structs, err := utils.CollectEmbeddedStructs(sandbox, io, "sandbox/api", utils.SandboxParts(sandbox), api_parts)
	if err != nil {
		return err
	}

	vars := map[string]interface{}{
		"Module":              module_conf.Module,
		"Version":             project_conf.Version,
		"ProjectName":         project_conf.ProjectName,
		"GeneratorName":       generatorName(sandbox),
		"GeneratorVersion":    generatorVersion(sandbox),
		"ConfigDir":           utils.ConfigDir(sandbox),
		"GoRelease":           utils.GoRelease,
		"GoFloor":             utils.GoFloor,
		"StructureConfFile":   utils.StructureConfFile,
		"HasSandbox":          hasSandbox,
		"HasDeps":             hasDeps,
		"HasCli":              hasCli,
		"HasServer":           hasServer,
		"HasFront":            hasFront,
		"HasDatabase":         hasDatabase,
		"HasExample":          hasExample,
		"HasBackoffice":       hasBackoffice,
		"SecretEnv":           utils.SecretEnvName(sandbox, project_conf.ProjectName),
		"HasDoc":              hasDoc,
		"HasReadme":           hasReadme,
		"HasAssets":           hasAssets,
		"Constructors":        constructors,
		"ConfigStructs":       config_structs,
		"SandboxStructs":      sandbox_structs,
		"ConstructorPackages": CollectConstructorPackages(sandbox, io, constructors),
		"DepLibs":             CollectDepLibs(sandbox, io),
		"AdapterImpls":        CollectAdapterImpls(sandbox, io),
		"Bindings":            bindings,
		"CliExamples":         utils.CollectExamples(sandbox, io, utils.ExampleCliSide),
		"LibExamples":         utils.CollectExamples(sandbox, io, utils.ExampleLibSide),
		"Commands":            commands,
		"CommandDocs":         command_docs,
		"Routes":              routes,
		"RouteDocs":           route_docs,
		"OpenApi":             open_api,
		"Databases":           databases,
		"DatabaseDocs":        database_docs,
		"Themes":              themes_conf.Themes,
		"DocIndex":            CollectDocIndex(sandbox, docs, themes_conf.Themes),
		"PublicApi":           public_api,
		"Structure":           structure,
		"DepContracts":        dep_contracts,
	}

	// The per-unit generators: one new.go per declared binding, command and
	// route, each owned by the mechanic that declares the unit.
	if hasDeps {
		if err := GenerateBindingNewFiles(sandbox, io, bindings, module_conf.Module); err != nil {
			return err
		}
	}

	if hasCli {
		if err := GenerateCommandNew(sandbox, io, commands, module_conf.Module); err != nil {
			return err
		}
		// Written once and never again: the five files that say what this
		// project answers when no command does.
		if err := GenerateCliErrorHandlers(sandbox, io, module_conf.Module); err != nil {
			return err
		}
		// Written once too: what one command line's chain of commands shares.
		if err := GenerateCommandProps(sandbox, io, module_conf.Module); err != nil {
			return err
		}
	}

	// One page per command beside docs/Commands' own index: a lookup costs the
	// page of the command asked about, never every command of the project.
	if hasDoc && hasCli {
		if err := GenerateCommandPages(sandbox, io, command_docs, project_conf.ProjectName); err != nil {
			return err
		}
	}

	// The server layer's mirror of it, one page per declared route.
	if hasDoc && hasServer {
		if err := GenerateRoutePages(sandbox, io, route_docs); err != nil {
			return err
		}
	}

	// The database layer's mirror of it, one page per declared database.
	if hasDoc && hasDatabase {
		if err := GenerateDatabasePages(sandbox, io, database_docs); err != nil {
			return err
		}
	}

	// The same for the contracts: one page per file of sandbox/api and per
	// contract of sandbox/deps, indexed by the symbols each one declares.
	if hasDoc && hasSandbox {
		if err := GeneratePublicApiPages(sandbox, io, public_api, dep_contracts); err != nil {
			return err
		}
	}

	// What an older build wrote under sandbox/internal/generated/ and the
	// OpinionatedAgnos libs replaced, dropped for every mechanic that is on.
	for _, extension := range []string{utils.ExtensionCli, utils.ExtensionServer, utils.ExtensionFront, utils.ExtensionDatabase} {
		if extensions_conf.IsEnabled(extension) {
			utils.RemoveRetiredGenerated(sandbox, io, extension)
		}
	}

	if hasServer {
		if err := GenerateRouteNew(sandbox, io, routes, module_conf.Module); err != nil {
			return err
		}
		// Written once and never again: the eight files that say what this
		// project answers when no route does.
		if err := GenerateServerErrorHandlers(sandbox, io, module_conf.Module); err != nil {
			return err
		}
		// Written once too: what one request's chain of routes shares.
		if err := GenerateRouteProps(sandbox, io, module_conf.Module); err != nil {
			return err
		}
	}

	// A database declares three generated files instead of one: its methods
	// are typed by table, so nothing reads the declaration back at runtime the
	// way the cli and the server dispatches do.
	if hasDatabase {
		if err := GenerateDatabaseNew(sandbox, io, databases, module_conf.Module); err != nil {
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

	sandbox.Deps.StdDeps.Logf("successfully rendered template\n")
	return nil
}
