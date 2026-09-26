package interview

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	interviewer "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/interviewer"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// The interview is the one surface of agnos written for a person instead of
// for an llm, and this file is what that costs: the session reads the project
// at --path before it offers anything, so a menu never holds a command that
// project cannot run, and the thing to do next is the first row rather than a
// page to read.
//
// Three tables carry it: the areas, the init that turns each mechanic on, and
// the units an init writes for itself. They are, with suggest.go, the whole of
// the agnos vocabulary this package spells — everything else is still
// generated from the declarations.

// infoCategory answers about the binary rather than about the project, so it
// is the one area offered where there is no project at all.
const infoCategory = "Info"

// area is one category of the command surface as this screen treats it: the
// name the help screen gives it, what it is for in the words of someone who
// has not read the docs, and the mechanic that owns it — "" for an area every
// project has.
type area struct {
	Name      string
	Help      string
	Extension string
}

// areas is every area in the order someone new needs them, which is the one
// place this screen departs from the help screen — that one lists categories
// in the order the declarations were collected. An area whose mechanic is off
// is not an area of this project: it is hidden, and the command that turns it
// on is offered as a step instead. A category missing from this table belongs
// to no mechanic, is always offered, and is listed after these under its own
// name.
var areas = []area{
	{"Core Commands", "build it, check it, install it, publish it", ""},
	{"Cli System", "the commands your program answers to", utils.ExtensionSandboxCli},
	{"Server System", "the http routes your program answers", utils.ExtensionSandboxServer},
	{"Front System", "a website your server serves: any html, css or js you put in assets/frontend", utils.ExtensionSandboxFront},
	{"Database System", "the records your program stores and reads back", utils.ExtensionSandboxDatabase},
	{"Deps System", "the libraries your program is allowed to use", utils.ExtensionSandboxDeps},
	{"Documentation", "the docs/ tree of this project", utils.ExtensionDoc},
	{"Examples", "the examples that guard this project", utils.ExtensionSandboxExample},
	{"Extensions", "what gets generated for this project", ""},
	{infoCategory, "help and version of the tool itself", ""},
}

// areaOf reads one area back by the category it is named after.
func areaOf(category string) (area, bool) {
	for _, one := range areas {
		if one.Name == category {
			return one, true
		}
	}
	return area{}, false
}

// areaHelp is what an area is for, and "" for a category the table does not
// name.
func areaHelp(category string) string {
	one, found := areaOf(category)
	if !found {
		return ""
	}
	return one.Help
}

// extensionInit is the command that turns one mechanic on. It is why an area
// that is off can still be reached: the init is offered as a step, which is
// the only place it is ever offered.
var extensionInit = map[string]string{
	utils.ExtensionSandboxCli:      "cli-init",
	utils.ExtensionSandboxServer:   "server-init",
	utils.ExtensionSandboxFront:    "front-init",
	utils.ExtensionSandboxDeps:     "deps-init",
	utils.ExtensionSandboxDatabase: "database-init",
}

// scaffoldedUnits are the units an init writes for itself: help and version
// from assets/sandbox-cli/, health from assets/sandbox-server/, and the
// frontend route and the index page front-init scaffolds. A layer holding
// nothing else has no unit of its own yet — which is what makes "declare its first one" the step after its
// init, instead of a step no project ever sees.
var scaffoldedUnits = map[string]bool{
	"help":     true,
	"version":  true,
	"health":   true,
	"frontend": true,
	"index":    true,
}

// projectState is the project at --path as the menus need it: whether there is
// one there at all, which mechanics it has on, and whether the units of a
// mechanic have been declared yet. It is read again before every menu, so a
// command that turns a mechanic on changes what the next menu offers.
type projectState struct {
	Started    bool
	Name       string
	Extensions map[string]bool
	Commands   int
	Routes     int
	Pages      int
	Databases  int
}

// readState reads the state off disk. A project with no project.yaml has not
// been started, and nothing else is read — every other command of the surface
// would fail on it. An extensions.yaml that will not parse leaves every
// mechanic off, which offers the inits and hides the areas: the same screen a
// fresh project gets, and the one that leads out of a broken declaration.
func readState(sandbox *api.Sandbox, io *smartio.SmartIO) projectState {
	state := projectState{Extensions: map[string]bool{}}

	if !io.IsFile(utils.ProjectConfPath(sandbox)) {
		return state
	}
	state.Started = true

	if conf, err := utils.LoadProjectConf(sandbox, io); err == nil {
		state.Name = conf.Name
	}

	if conf, err := utils.LoadExtensionsConf(sandbox, io); err == nil {
		utils.NormalizeExtensions(conf)
		for _, spec := range utils.ExtensionCatalog() {
			state.Extensions[spec.Name] = conf.IsEnabled(spec.Name)
		}
	}

	state.Commands = ownUnits(commandOptions(sandbox, io))
	state.Routes = ownUnits(dirOptions(sandbox, io, routesDir))
	state.Pages = ownUnits(pageOptions(sandbox, io))
	state.Databases = ownUnits(dirOptions(sandbox, io, utils.DatabasesDir))

	return state
}

// ownUnits is how many units of a layer the project declared itself, which is
// every one its init did not scaffold.
func ownUnits(declared []interviewer.AlternativeOption) int {
	count := 0
	for _, one := range declared {
		if !scaffoldedUnits[one.Id] {
			count++
		}
	}
	return count
}

// enabled reports whether one mechanic of the project is on.
func enabled(state projectState, extension string) bool {
	return state.Extensions[extension]
}

// ─── The steps ──────────────────────────────────────────────────────────────

// step is one thing this project has not done yet: the command that does it
// and what it does, in the words of someone meeting agnos for the first time.
// A key step is one the project needs before it does anything at all — those
// come first and the first of them is the suggested one; the rest are offers,
// layers a project may well never want.
type step struct {
	Verb string
	Msg  string
	Key  bool
}

// nextSteps is the ladder a project climbs, kept to the rungs it has not
// climbed: create it, give it a cli, give that cli a command, and so on. The
// list empties itself as the project grows — one with every layer on and a
// unit in each offers no step at all and is driven by the areas alone.
//
// Every <x>-init is here, and that is load-bearing: its own area is hidden
// while the mechanic is off, so the step is the only way in.
func nextSteps(state projectState) []step {
	if !state.Started {
		return []step{{scaffoldVerb, "Create an agnos project here — this folder has none yet", true}}
	}

	steps := []step{}

	steps = appendStep(steps, !enabled(state, utils.ExtensionSandboxCli), true,
		extensionInit[utils.ExtensionSandboxCli], "Give it a command line — the commands it answers to")
	steps = appendStep(steps, enabled(state, utils.ExtensionSandboxCli) && state.Commands == 0, true,
		"add-command", "Declare its first command")
	steps = appendStep(steps, enabled(state, utils.ExtensionSandboxServer) && state.Routes == 0, true,
		"add-route", "Declare its first route")
	steps = appendStep(steps, enabled(state, utils.ExtensionSandboxFront) && state.Pages == 0, true,
		"add-page", "Add its first page")
	steps = appendStep(steps, enabled(state, utils.ExtensionSandboxDatabase) && state.Databases == 0, true,
		"add-database", "Declare its first database")

	steps = appendStep(steps, !enabled(state, utils.ExtensionSandboxServer), false,
		extensionInit[utils.ExtensionSandboxServer], "Give it an http server — the routes it answers")
	steps = appendStep(steps, enabled(state, utils.ExtensionSandboxServer) && !enabled(state, utils.ExtensionSandboxFront), false,
		extensionInit[utils.ExtensionSandboxFront], "Give it html pages, served by that server")
	steps = appendStep(steps, !enabled(state, utils.ExtensionSandboxDatabase), false,
		extensionInit[utils.ExtensionSandboxDatabase], "Give it somewhere to store records — tables it reads and writes")
	steps = appendStep(steps, !enabled(state, utils.ExtensionSandboxDeps), false,
		extensionInit[utils.ExtensionSandboxDeps], "Give it the dependency layer — deps, adapters, availables")

	return steps
}

// appendStep adds one rung when the project has not climbed it, and when the
// command that climbs it is one this binary declares.
func appendStep(steps []step, applies bool, key bool, verb string, msg string) []step {
	if !applies || verb == "" {
		return steps
	}
	return append(steps, step{verb, msg, key})
}

// ─── The gate ───────────────────────────────────────────────────────────────

// applies reports whether a command is worth offering on the project in front
// of the person. Three rules, and no command named in them twice:
//
//   - an <x>-init is a step and never a row of an area: while its mechanic is
//     off the area is hidden and the step is the way in, and once it is on the
//     init has nothing left to do;
//   - with no project, only the Info area answers — every other command of the
//     surface says "run `agnos start` first", and the scaffold that fixes that
//     is the one step offered;
//   - a command of an area belongs to a mechanic, and is offered only while
//     that mechanic is on.
//
// An area is hidden by the same rule rather than by one of its own: a category
// whose mechanic is off has no command that applies, so no row is built for it.
func applies(state projectState, verb string, category string) bool {
	if initExtension(verb) != "" {
		return false
	}
	if !state.Started {
		return category == infoCategory
	}
	if one, found := areaOf(category); found && one.Extension != "" {
		return enabled(state, one.Extension)
	}
	return true
}

// initExtension is the mechanic one <x>-init turns on, and "" for every
// command that is not one.
func initExtension(verb string) string {
	for extension, init := range extensionInit {
		if verb == init {
			return extension
		}
	}
	return ""
}
