package interview

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	interviewer "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/interviewer"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// commandsDir and routesDir are where the cli and server layers declare their
// units. They are spelled here because utils exposes the path of one command
// and of one route, never the tree they live in.
const (
	commandsDir = "sandbox/internal/commands"
	routesDir   = "sandbox/internal/routes"
)

// The closed vocabularies agnos itself defines. They are values a command
// accepts, not names of anything on disk, so they are listed rather than read.
var (
	fieldTypes     = []string{"string", "boolean", "int", "float"}
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
func SuggestFor(sandbox *api.Sandbox, io *smartio.SmartIO, command api.Command, id string) suggestion {
	verb := verbOf(command)

	switch id {
	case "command":
		if verb == helpVerb {
			return closed(runningCommandOptions(sandbox))
		}
		return closed(commandOptions(sandbox, io))
	case "route":
		return closed(dirOptions(sandbox, io, routesDir))
	case "method":
		return closed(literalOptions(routeMethods))
	case "runtime":
		return closed(literalOptions(buildRuntimes))
	case "target":
		return closed(literalOptions(compileTargets))
	case "theme":
		return closed(themeOptions(sandbox, io))
	case "only":
		return closed(exampleOptions(sandbox, io, ""))
	case "category":
		return open(categoryOptions(sandbox, io))
	case "type":
		return typeSuggestion(verb)
	case "available":
		return availableSuggestion(sandbox, io, verb)
	case "adapter":
		return adapterSuggestion(sandbox, io, verb)
	case "dep":
		return depSuggestion(sandbox, io, verb)
	case "name":
		return nameSuggestion(sandbox, io, verb)
	}

	return suggestion{}
}

// typeSuggestion picks the vocabulary the declared types come from. A route
// body may also be an object; set-body's --type names a body kind rather than
// a value type, so it is left as free text.
func typeSuggestion(verb string) suggestion {
	switch verb {
	case "set-body":
		return suggestion{}
	case "add-body-field":
		return closed(literalOptions(bodyTypes))
	}
	return closed(literalOptions(fieldTypes))
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
func nameSuggestion(sandbox *api.Sandbox, io *smartio.SmartIO, verb string) suggestion {
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
	}

	return suggestion{}
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

// categoryOptions is every category the project's own command surface already
// uses, so a new command joins a heading that exists instead of opening one of
// its own by a typo.
func categoryOptions(sandbox *api.Sandbox, io *smartio.SmartIO) []interviewer.AlternativeOption {
	options := []interviewer.AlternativeOption{}
	seen := map[string]bool{}

	for _, declared := range declaredCommands(sandbox, io) {
		if declared.Category == "" || seen[declared.Category] {
			continue
		}
		seen[declared.Category] = true
		options = append(options, interviewer.AlternativeOption{Id: declared.Category, Msg: declared.Category})
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
