package interview

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	interviewer "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/interviewer"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/databaseconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// commandsDir and routesDir are where the cli and server layers declare their
// units. They are spelled here because utils exposes the path of one command
// and of one route, never the tree they live in.
const (
	commandsDir = "sandbox/internal/commands"
	routesDir   = utils.RoutesDir
)

// The closed vocabularies agnos itself defines. They are values a command
// accepts, not names of anything on disk, so they are listed rather than read.
var (
	fieldTypes     = []string{"string", "boolean", "int", "float"}
	tableTypes     = databaseconf.FieldTypes
	bodyTypes      = []string{"string", "boolean", "int", "float", "object"}
	routeMethods   = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	buildRuntimes  = []string{"go", "none"}
	compileTargets = []string{"linux86", "linuxarm64", "linuxi32", "mac86", "macarm64", "windows86", "windowsi32", "all"}
)

// helpVerb is the one command whose argument names a command of the binary
// running the interview rather than one of the project at --path: help prints
// its own screens, whatever directory it is pointed at.
const helpVerb = "help"

// suggestion is what an interview offers for one field instead of a bare text
// box. Open reports that the list is advisory rather than exhaustive — a new
// category is as valid as an existing one — and adds a row for typing a value
// the list does not hold.
type suggestion struct {
	Options []interviewer.AlternativeOption
	Open    bool
}

// SuggestFor is the whole of what the interview knows about agnos's own
// vocabulary: which field of which command names a thing that already exists,
// and where to read the list of them. It is a table rather than a branch per
// command, and it is the one place a new agnos noun is taught to the
// interview.
//
// A field with no entry here is answered as free text, which is the right
// answer for every name a command is about to create.
//
// answered is what the session has bound so far, and it is what lets a list be
// read off the thing the command is about rather than off the project: the
// query parameters of the route already answered, not of every route there is.
// A field whose list depends on an answer that has not come yet — a name asked
// before its --route, which is what the argument order spells — falls back to
// free text, exactly as it did before there was a list at all.
func SuggestFor(sandbox *api.Sandbox, io *smartio.SmartIO, command api.Command, id string, answered map[string][]any) suggestion {
	verb := verbOf(command)

	switch id {
	case "command":
		if verb == helpVerb {
			return closed(runningCommandOptions(sandbox))
		}
		return closed(commandOptions(sandbox, io))
	case "route":
		return closed(dirOptions(sandbox, io, routesDir))
	case "database":
		return closed(dirOptions(sandbox, io, utils.DatabasesDir))
	case "table":
		return closed(tableOptions(sandbox, io, answered))
	case "parent":
		return closed(nestedTableOptions(sandbox, io, answered))
	case "method":
		return closed(literalOptions(routeMethods))
	case "runtime":
		return closed(literalOptions(buildRuntimes))
	case "target":
		return targetSuggestion(sandbox, io, verb, answered)
	case "theme":
		return closed(themeOptions(sandbox, io))
	case "only":
		return closed(exampleOptions(sandbox, io, ""))
	case "category":
		return open(categorySuggestion(sandbox, io, verb))
	case "format":
		return closed(literalOptions(utils.RouteSchemaFormats))
	case "clear":
		return clearSuggestion(verb)
	case "font":
		return closed(literalOptions(routeconf.ParameterFonts))
	case "trigger-type":
		return closed(literalOptions(routeconf.TriggerTypes))
	case "type":
		return typeSuggestion(verb)
	case "available":
		return availableSuggestion(sandbox, io, verb)
	case "adapter":
		return adapterSuggestion(sandbox, io, verb)
	case "dep":
		return depSuggestion(sandbox, io, verb)
	case "name":
		return nameSuggestion(sandbox, io, verb, answered)
	}

	return suggestion{}
}

// scopeFields are the ids that name the thing the rest of a command's
// questions are about: the route whose parameters are being listed, the
// command whose flags are. They are declared as flags, so the argument order
// puts them after the name they scope — and a list read off a route that has
// not been answered yet is no list at all.
//
// The interview asks them first instead. It changes the order the questions
// come in and nothing about the command line they add up to, and it is what
// turns "which parameter?" from a text box into the parameters that route
// declares.
var scopeFields = []string{"route", "command", "database", "table"}

// AskOrder is the order the interview puts one command's fields to the person:
// every field that scopes the others first, in the order they narrow each
// other — the database before the table it holds — and the declaration's own
// order after them. Only a required flag is hoisted: an optional one may be
// skipped, and an arg is already asked before the flags.
func AskOrder(fields []Field) []Field {
	hoisted := []Field{}
	taken := map[int]bool{}

	for _, id := range scopeFields {
		for index, field := range fields {
			if taken[index] || field.Id != id || !field.IsFlag || !field.Required {
				continue
			}
			taken[index] = true
			hoisted = append(hoisted, field)
			break
		}
	}

	if len(hoisted) == 0 {
		return fields
	}

	for index, field := range fields {
		if !taken[index] {
			hoisted = append(hoisted, field)
		}
	}
	return hoisted
}

// clearSuggestion is the keys --clear takes off, which are the command's own:
// a body property carries the json-schema keywords, a path and a parameter the
// keys of their own entries.
func clearSuggestion(verb string) suggestion {
	switch verb {
	case "set-body-field":
		return closed(literalOptions(utils.RouteBodyFieldClearKeys))
	case "set-table-field":
		return closed(literalOptions(utils.DatabaseFieldClearKeys))
	case "set-path":
		return closed(literalOptions(utils.RoutePathClearKeys))
	}
	return closed(literalOptions(utils.RouteParameterClearKeys))
}

// typeSuggestion picks the vocabulary the declared types come from. A route
// body may also be an object; set-body's --type names a body kind rather than
// a value type, so it is left as free text.
func typeSuggestion(verb string) suggestion {
	switch verb {
	case "set-body":
		return suggestion{}
	case "add-body-field", "set-body-field":
		return closed(literalOptions(bodyTypes))
	case "add-table-field", "set-table-field":
		return closed(literalOptions(tableTypes))
	case "add-parameter", "set-parameter":
		return closed(literalOptions(routeconf.ParameterTypes))
	}
	return closed(literalOptions(fieldTypes))
}

// targetSuggestion tells the two questions about a target apart: a database
// field points at a table of its own database, and `compile` names one of the
// binaries it can cross-build. They are one field id, so the verb decides.
func targetSuggestion(sandbox *api.Sandbox, io *smartio.SmartIO, verb string, answered map[string][]any) suggestion {
	switch verb {
	case "add-table-field", "set-table-field":
		return closed(tableOptions(sandbox, io, answered))
	}
	return closed(literalOptions(compileTargets))
}

// availableSuggestion offers the declared availables, except on add-available,
// where the name is the one being created.
func availableSuggestion(sandbox *api.Sandbox, io *smartio.SmartIO, verb string) suggestion {
	if verb == "add-available" {
		return suggestion{}
	}
	return closed(literalOptions(utils.DeclaredAvailables(sandbox, io)))
}

// adapterSuggestion tells the two questions about an adapter apart: installing
// one picks from the catalog, binding or removing one picks from what the
// project already has.
func adapterSuggestion(sandbox *api.Sandbox, io *smartio.SmartIO, verb string) suggestion {
	switch verb {
	case "add-adapter", "add-dep":
		catalog, err := utils.CatalogAdapters(sandbox)
		if err != nil {
			return suggestion{}
		}
		return closed(literalOptions(catalog))
	}
	return closed(literalOptions(utils.InstalledAdapters(sandbox, io)))
}

// depSuggestion is the same split for a dep: add-dep picks from the catalog —
// and takes a module path besides, which is why its list is open — while every
// other command names one the project has installed.
func depSuggestion(sandbox *api.Sandbox, io *smartio.SmartIO, verb string) suggestion {
	if verb == "add-dep" {
		catalog, err := utils.CatalogDeps(sandbox)
		if err != nil {
			return suggestion{}
		}
		return open(literalOptions(catalog))
	}
	return closed(dirOptions(sandbox, io, utils.ContractsDir))
}

// nameSuggestion answers the one field id that means something different on
// every command. A name a command is about to create has no list, so only the
// commands that name something existing appear here.
func nameSuggestion(sandbox *api.Sandbox, io *smartio.SmartIO, verb string, answered map[string][]any) suggestion {
	switch verb {
	case "enable-extension", "disable-extension":
		return closed(extensionOptions())
	case "remove-command", "set-command":
		return closed(commandOptions(sandbox, io))
	case "remove-route":
		return closed(dirOptions(sandbox, io, routesDir))
	case "remove-page":
		return closed(pageOptions(sandbox, io))
	case "remove-doc":
		return open(docOptions(sandbox, io))
	case "remove-cli-example":
		return closed(exampleOptions(sandbox, io, utils.ExampleCliSide))
	case "remove-lib-example":
		return closed(exampleOptions(sandbox, io, utils.ExampleLibSide))
	case "update-test":
		return closed(exampleOptions(sandbox, io, ""))
	case "set-parameter", "remove-parameter":
		return closed(routeParameterOptions(sandbox, io, answered))
	case "set-path", "remove-path":
		return closed(routePathOptions(sandbox, io, answered))
	case "set-body-field", "remove-body-field":
		return closed(bodyFieldOptions(sandbox, io, answered))
	case "remove-database":
		return closed(dirOptions(sandbox, io, utils.DatabasesDir))
	case "remove-table":
		return closed(tableOptions(sandbox, io, answered))
	case "set-table-field", "remove-table-field":
		return closed(tableFieldOptions(sandbox, io, answered))
	case "remove-flag":
		return closed(commandFieldOptions(sandbox, io, answered, true))
	case "remove-arg":
		return closed(commandFieldOptions(sandbox, io, answered, false))
	}

	return suggestion{}
}

// tableOptions is the tables one database of the project declares, read off
// the --database already answered. A session that has not answered one yet
// gets no list, exactly as it did before there was one.
func tableOptions(sandbox *api.Sandbox, io *smartio.SmartIO, answered map[string][]any) []interviewer.AlternativeOption {
	conf, found := answeredDatabase(sandbox, io, answered)
	if !found {
		return []interviewer.AlternativeOption{}
	}

	options := []interviewer.AlternativeOption{}
	for _, table := range conf.Tables {
		options = append(options, interviewer.AlternativeOption{
			Id:  table.Name,
			Msg: labelled(sandbox, table.Name, tableSummary(sandbox, table)),
		})
	}
	return options
}

// tableFieldOptions is every field one table declares, the fields of its
// nested collections included, by the name the editors of those fields spell.
func tableFieldOptions(sandbox *api.Sandbox, io *smartio.SmartIO, answered map[string][]any) []interviewer.AlternativeOption {
	table, found := answeredTable(sandbox, io, answered)
	if !found {
		return []interviewer.AlternativeOption{}
	}

	options := []interviewer.AlternativeOption{}
	for _, field := range table.Fields {
		options = append(options, interviewer.AlternativeOption{Id: field.Name, Msg: labelled(sandbox, field.Name, field.Type)})
		for _, nested := range field.Fields {
			options = append(options, interviewer.AlternativeOption{
				Id:  nested.Name,
				Msg: labelled(sandbox, field.Name+"."+nested.Name, nested.Type),
			})
		}
	}
	return options
}

// nestedTableOptions is the nested collections one table declares — the only
// values --parent accepts, because a field goes inside a `database` field or
// inside nothing.
func nestedTableOptions(sandbox *api.Sandbox, io *smartio.SmartIO, answered map[string][]any) []interviewer.AlternativeOption {
	table, found := answeredTable(sandbox, io, answered)
	if !found {
		return []interviewer.AlternativeOption{}
	}

	options := []interviewer.AlternativeOption{}
	for _, field := range table.Fields {
		if field.Type != databaseconf.FieldDatabase {
			continue
		}
		options = append(options, interviewer.AlternativeOption{
			Id:  field.Name,
			Msg: labelled(sandbox, field.Name, "a collection nested under each record"),
		})
	}
	return options
}

// tableSummary is what one table is worth saying in a row: how many fields it
// holds.
func tableSummary(sandbox *api.Sandbox, table databaseconf.Table) string {
	if len(table.Fields) == 1 {
		return "1 field"
	}
	return sandbox.Deps.Std.Sprintf("%d fields", len(table.Fields))
}

// answeredDatabase is the database the session is working on, read off disk,
// or false while the question naming it has not been answered.
func answeredDatabase(sandbox *api.Sandbox, io *smartio.SmartIO, answered map[string][]any) (*databaseconf.DatabaseConf, bool) {
	database := answeredText(sandbox, answered, "database")
	if database == "" {
		return nil, false
	}

	conf, err := utils.LoadDatabaseConf(sandbox, io, database)
	if err != nil {
		return nil, false
	}
	return conf, true
}

// answeredTable is the table the session is working on: the --table answered,
// looked up in the --database answered. The interview asks both before any
// field question, which is what makes this list exist at all.
func answeredTable(sandbox *api.Sandbox, io *smartio.SmartIO, answered map[string][]any) (databaseconf.Table, bool) {
	conf, found := answeredDatabase(sandbox, io, answered)
	if !found {
		return databaseconf.Table{}, false
	}

	index := utils.FindDatabaseTable(sandbox, conf.Tables, answeredText(sandbox, answered, "table"))
	if index < 0 {
		return databaseconf.Table{}, false
	}
	return conf.Tables[index], true
}

// commandFieldOptions is the flags, or the args, one command of the project
// declares — the same answer for a command that routeFieldOptions is for a
// route, read off the --command already answered.
func commandFieldOptions(sandbox *api.Sandbox, io *smartio.SmartIO, answered map[string][]any, flags bool) []interviewer.AlternativeOption {
	name := answeredText(sandbox, answered, "command")
	if name == "" {
		return []interviewer.AlternativeOption{}
	}

	conf, err := utils.LoadCommandConf(sandbox, io, name)
	if err != nil {
		return []interviewer.AlternativeOption{}
	}

	fields := conf.Args
	if flags {
		fields = conf.Flags
	}

	options := []interviewer.AlternativeOption{}
	for _, field := range fields {
		label := field.Key
		if len(field.Identifiers) > 0 {
			label = field.Identifiers[0]
		}
		options = append(options, interviewer.AlternativeOption{Id: field.Key, Msg: labelled(sandbox, label, field.Description)})
	}
	return options
}

// routePathOptions is every entry of one route's `paths`, by the id its
// editors spell. It is read off the route that was answered, so a session that
// has not answered one yet gets no list.
func routePathOptions(sandbox *api.Sandbox, io *smartio.SmartIO, answered map[string][]any) []interviewer.AlternativeOption {
	conf, found := answeredRoute(sandbox, io, answered)
	if !found {
		return []interviewer.AlternativeOption{}
	}

	options := []interviewer.AlternativeOption{}
	for _, path := range conf.Paths {
		label := path.Id
		if path.Trigger.Exists {
			label += "  " + path.Trigger.Value
		}
		options = append(options, interviewer.AlternativeOption{Id: path.Id, Msg: labelled(sandbox, label, path.Description)})
	}
	return options
}

// routeParameterOptions is every entry of one route's `parameters`, by the key
// it is read under — the answer to the question a person editing a
// declaration actually has: which values does this route read.
func routeParameterOptions(sandbox *api.Sandbox, io *smartio.SmartIO, answered map[string][]any) []interviewer.AlternativeOption {
	conf, found := answeredRoute(sandbox, io, answered)
	if !found {
		return []interviewer.AlternativeOption{}
	}

	options := []interviewer.AlternativeOption{}
	for _, parameter := range conf.Parameters {
		options = append(options, interviewer.AlternativeOption{Id: parameter.Key, Msg: labelled(sandbox, parameter.Key, parameter.Description)})
	}
	return options
}

// bodyFieldOptions is every property of one route's body json-schema, by the
// dotted path add-body-field declares it at — the nested ones included, which
// is the one list a person cannot read off the route.yaml at a glance.
func bodyFieldOptions(sandbox *api.Sandbox, io *smartio.SmartIO, answered map[string][]any) []interviewer.AlternativeOption {
	conf, found := answeredRoute(sandbox, io, answered)
	if !found || conf.Body.Schema == nil {
		return []interviewer.AlternativeOption{}
	}
	return appendSchemaPaths(sandbox, []interviewer.AlternativeOption{}, conf.Body.Schema, "")
}

// appendSchemaPaths walks one object schema depth-first, so a nested property
// is listed under the object it belongs to.
func appendSchemaPaths(sandbox *api.Sandbox, options []interviewer.AlternativeOption, schema *routeconf.Schema, prefix string) []interviewer.AlternativeOption {
	for _, property := range schema.Properties {
		path := property.Name
		if prefix != "" {
			path = prefix + "." + property.Name
		}

		object := utils.SchemaObjectOf(property.Schema)
		kind := property.Schema.Type
		if object != nil && property.Schema.Type == "array" {
			kind = "[]" + object.Type
		}

		options = append(options, interviewer.AlternativeOption{Id: path, Msg: labelled(sandbox, path, kind)})
		if object != nil && object.Type == "object" {
			options = appendSchemaPaths(sandbox, options, object, path)
		}
	}
	return options
}

// answeredRoute is the route the session is working on, read off disk, or
// false while the question naming it has not been answered — which is what the
// argument order leaves true for a name asked before its --route.
func answeredRoute(sandbox *api.Sandbox, io *smartio.SmartIO, answered map[string][]any) (*routeconf.RouteConf, bool) {
	route := answeredText(sandbox, answered, "route")
	if route == "" {
		return nil, false
	}

	conf, err := utils.LoadRouteConf(sandbox, io, route)
	if err != nil {
		return nil, false
	}
	return conf, true
}

// ─── The lists ──────────────────────────────────────────────────────────────

// commandOptions is every command declared by the project being worked on —
// read off disk, never from sandbox.Cli.Commands. Those are the commands of
// the binary running the interview, and with --path pointing at another
// project the two have nothing to do with each other. Hidden commands are
// listed too: one still dispatches, and a flag may still be declared on it.
func commandOptions(sandbox *api.Sandbox, io *smartio.SmartIO) []interviewer.AlternativeOption {
	options := []interviewer.AlternativeOption{}

	for _, declared := range declaredCommands(sandbox, io) {
		options = append(options, interviewer.AlternativeOption{
			Id:  declared.Identifier,
			Msg: labelled(sandbox, declared.Identifier, declared.Help),
		})
	}

	return options
}

// runningCommandOptions is the command surface of the binary itself — what
// `help` answers about, and the one list that is not read off the project.
func runningCommandOptions(sandbox *api.Sandbox) []interviewer.AlternativeOption {
	options := []interviewer.AlternativeOption{}
	for _, declared := range sandbox.Cli.Commands {
		if len(declared.Identifiers) == 0 {
			continue
		}
		options = append(options, interviewer.AlternativeOption{
			Id:  declared.Identifiers[0],
			Msg: labelled(sandbox, declared.Identifiers[0], declared.Help),
		})
	}
	return options
}

// categorySuggestion offers the headings of the surface the command being
// declared belongs to. A route is listed in docs/Routes and a command in
// docs/Commands, so the two never share a heading — offering a command's
// categories to add-route is what left every route after the first retyping
// its own, one typo away from a heading of its own.
func categorySuggestion(sandbox *api.Sandbox, io *smartio.SmartIO, verb string) []interviewer.AlternativeOption {
	switch verb {
	case "add-route", "set-route":
		return routeCategoryOptions(sandbox, io)
	}
	return commandCategoryOptions(sandbox, io)
}

// commandCategoryOptions is every category the project's own command surface
// already uses, so a new command joins a heading that exists instead of
// opening one of its own by a typo.
func commandCategoryOptions(sandbox *api.Sandbox, io *smartio.SmartIO) []interviewer.AlternativeOption {
	categories := []string{}
	for _, declared := range declaredCommands(sandbox, io) {
		categories = append(categories, declared.Category)
	}
	return categoryOptions(categories)
}

// routeCategoryOptions is the same list for the server surface, read off the
// `category` of every declared route.yaml. A route whose declaration will not
// parse contributes nothing: it has no heading to join.
func routeCategoryOptions(sandbox *api.Sandbox, io *smartio.SmartIO) []interviewer.AlternativeOption {
	if !io.IsDir(routesDir) {
		return []interviewer.AlternativeOption{}
	}

	names := []string{}
	for _, path := range io.ListDirs(routesDir) {
		if name := utils.LastSegment(sandbox, path); name != "" {
			names = append(names, name)
		}
	}
	sandbox.Deps.Sortdeps.Strings(names)

	categories := []string{}
	for _, name := range names {
		conf, err := utils.LoadRouteConf(sandbox, io, name)
		if err != nil {
			continue
		}
		categories = append(categories, conf.Category)
	}

	return categoryOptions(categories)
}

// categoryOptions is one row per distinct heading, in the order the surface
// declares them, with the blanks dropped: a unit with no category is listed
// under a fallback heading nothing has to be told to reuse.
func categoryOptions(categories []string) []interviewer.AlternativeOption {
	options := []interviewer.AlternativeOption{}
	seen := map[string]bool{}

	for _, category := range categories {
		if category == "" || seen[category] {
			continue
		}
		seen[category] = true
		options = append(options, interviewer.AlternativeOption{Id: category, Msg: category})
	}

	return options
}

// declaredCommand is the little of one project's entries.yaml the suggestions
// need: what the command is called and what to say about it.
type declaredCommand struct {
	Identifier string
	Category   string
	Help       string
}

// declaredCommands reads every sandbox/internal/commands/<name>/entries.yaml of
// the project being worked on. A directory whose declaration will not parse is
// still offered under its own name: the interview is how a person fixes such a
// command, so hiding it would hide the way out.
func declaredCommands(sandbox *api.Sandbox, io *smartio.SmartIO) []declaredCommand {
	if !io.IsDir(commandsDir) {
		return []declaredCommand{}
	}

	names := []string{}
	for _, path := range io.ListDirs(commandsDir) {
		if name := utils.LastSegment(sandbox, path); name != "" {
			names = append(names, name)
		}
	}
	sandbox.Deps.Sortdeps.Strings(names)

	commands := []declaredCommand{}
	for _, name := range names {
		declared := declaredCommand{Identifier: utils.CommandIdentifier(sandbox, name)}

		if conf, err := utils.LoadCommandConf(sandbox, io, name); err == nil {
			if len(conf.Identifiers) > 0 {
				declared.Identifier = conf.Identifiers[0]
			}
			declared.Category = conf.Category
			declared.Help = conf.Help
		}

		commands = append(commands, declared)
	}

	return commands
}

// extensionOptions is the generation-mechanic catalog, with the line
// list-extensions prints as the label.
func extensionOptions() []interviewer.AlternativeOption {
	options := []interviewer.AlternativeOption{}
	for _, spec := range utils.ExtensionCatalog() {
		options = append(options, interviewer.AlternativeOption{Id: spec.Name, Msg: spec.Name + "  —  " + spec.Help})
	}
	return options
}

// themeOptions is the doc themes declared in themes.yaml, by id.
func themeOptions(sandbox *api.Sandbox, io *smartio.SmartIO) []interviewer.AlternativeOption {
	conf, err := utils.LoadThemesConf(sandbox, io)
	if err != nil {
		return []interviewer.AlternativeOption{}
	}

	options := []interviewer.AlternativeOption{}
	for _, theme := range conf.Themes {
		options = append(options, interviewer.AlternativeOption{Id: theme.Id, Msg: labelled(sandbox, theme.Id, theme.Description)})
	}
	return options
}

// pageOptions is every route with an html template beside it — the one test
// that tells a page from any other route.
func pageOptions(sandbox *api.Sandbox, io *smartio.SmartIO) []interviewer.AlternativeOption {
	options := []interviewer.AlternativeOption{}
	for _, option := range dirOptions(sandbox, io, routesDir) {
		if utils.IsPage(sandbox, io, option.Id) {
			options = append(options, option)
		}
	}
	return options
}

// docOptions is every doc of the tree, sub-docs included, named the way
// remove-doc spells them ("PublicApi", "PublicApi/api.Actions").
func docOptions(sandbox *api.Sandbox, io *smartio.SmartIO) []interviewer.AlternativeOption {
	tree, err := utils.CollectDocTree(sandbox, io)
	if err != nil {
		return []interviewer.AlternativeOption{}
	}

	options := []interviewer.AlternativeOption{}
	return appendDocs(sandbox, options, tree)
}

// appendDocs walks the doc tree depth-first, so a sub-doc is listed under the
// parent it belongs to.
func appendDocs(sandbox *api.Sandbox, options []interviewer.AlternativeOption, docs []utils.Doc) []interviewer.AlternativeOption {
	for _, doc := range docs {
		name := sandbox.Deps.Stringsdeps.TrimPrefix(doc.Path, utils.DocsDir+"/")
		options = append(options, interviewer.AlternativeOption{Id: name, Msg: labelled(sandbox, name, doc.Description)})
		options = appendDocs(sandbox, options, doc.Subdocs)
	}
	return options
}

// exampleOptions is the examples of one side, or of both when side is empty —
// exec-test --only and update-test name an example on both sides at once.
func exampleOptions(sandbox *api.Sandbox, io *smartio.SmartIO, side string) []interviewer.AlternativeOption {
	sides := []string{side}
	if side == "" {
		sides = utils.ExampleSides
	}

	names := []string{}
	seen := map[string]bool{}
	for _, one := range sides {
		for _, example := range utils.CollectExamples(sandbox, io, one) {
			if seen[example.Name] {
				continue
			}
			seen[example.Name] = true
			names = append(names, example.Name)
		}
	}

	sandbox.Deps.Sortdeps.Strings(names)
	return literalOptions(names)
}

// dirOptions is one option per directory of a tree, by its own name. It is how
// every list read off disk is built, so an absent tree yields no options
// rather than an error.
func dirOptions(sandbox *api.Sandbox, io *smartio.SmartIO, dir string) []interviewer.AlternativeOption {
	if !io.IsDir(dir) {
		return []interviewer.AlternativeOption{}
	}

	names := []string{}
	for _, path := range io.ListDirs(dir) {
		name := utils.LastSegment(sandbox, path)
		if name != "" {
			names = append(names, name)
		}
	}

	sandbox.Deps.Sortdeps.Strings(names)
	return literalOptions(names)
}

// ─── Helpers ────────────────────────────────────────────────────────────────

// literalOptions is a list of values that are their own labels.
func literalOptions(values []string) []interviewer.AlternativeOption {
	options := []interviewer.AlternativeOption{}
	for _, value := range values {
		options = append(options, interviewer.AlternativeOption{Id: value, Msg: value})
	}
	return options
}

// labelled writes one row as the value plus what it is, when there is
// something to say about it.
func labelled(sandbox *api.Sandbox, id string, description string) string {
	if sandbox.Deps.Stringsdeps.TrimSpace(description) == "" {
		return id
	}
	return sandbox.Deps.Std.Sprintf("%s  —  %s", id, description)
}

// verbOf is the canonical identifier of a command, the one every table above
// is keyed by.
func verbOf(command api.Command) string {
	if len(command.Identifiers) == 0 {
		return ""
	}
	return command.Identifiers[0]
}

// closed and open are the two shapes a suggestion takes: a list that is the
// whole of what the field accepts, and one that is only a head start.
func closed(options []interviewer.AlternativeOption) suggestion {
	return suggestion{Options: options, Open: false}
}

func open(options []interviewer.AlternativeOption) suggestion {
	return suggestion{Options: options, Open: true}
}
