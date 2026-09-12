package help

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
)

// help is a command like any other — entries.yaml, generated entries.go, and
// this handler.go — except that `agnos build` writes all three instead of the
// user writing two of them. This file is regenerated from every
// sandbox/internal/commands/<name>/entries.yaml: the command metadata is baked
// into helpCommands below, the rendering code is fixed.

const (
	exitOk    = 0
	exitUsage = 2
)

// identifiedBy reports whether name is one of the identifiers a command
// answers to, its aliases included.
func identifiedBy(identifiers []string, name string) bool {
	for _, identifier := range identifiers {
		if identifier == name {
			return true
		}
	}
	return false
}

// binaryName is the executable's name as a user types it: the configured
// project name, lowercased. Usage lines show what to type, not the display
// name of the project.
func binaryName(sandbox *api.Sandbox) string {
	return sandbox.Deps.Stringsdeps.ToLower(config.ProjectName)
}

// ─── ANSI escape sequences ──────────────────────────────────────────────────

const (
	bold    = "\033[1m"
	dim     = "\033[2m"
	italic  = "\033[3m"
	reset   = "\033[0m"
	cyan    = "\033[36m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	magenta = "\033[35m"
	white   = "\033[97m"
	gray    = "\033[90m"
	red     = "\033[31m"
)

// ─── Baked command metadata ─────────────────────────────────────────────────

type helpField struct {
	Identifiers []string // empty for a positional argument
	Name        string   // set for a positional argument
	Description string
	Examples    []string
	Type        string
	Default     string
	Required    bool
}

type helpCommand struct {
	Identifiers     []string
	Category        string
	Description     string
	LongDescription string
	Examples        []string
	Hidden          bool
	Flags           []helpField
	Args            []helpField
}

var helpCommands = []helpCommand{
	{
		Identifiers:     []string{"add-adapter"},
		Category:        "Deps System",
		Description:     "Installs one further adapter for a contract the project already has",
		LongDescription: "Renders assets/adapterlist/<adapter> into the project and writes its declaration to adapters/libs/<adapter>/adapter.yaml. The contract it fills has to be installed already. Installing changes no selection: an available binds one adapter per field, so --available names the one that switches to it.",
		Examples:        []string{"add-adapter reflectsort", "add-adapter reflectsort --available lambda"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--available"}, Description: "the available that should switch to this adapter (installs only when absent)", Examples: []string{"add-adapter reflectsort --available lambda"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"add-adapter reflectsort --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{"add-adapter reflectsort -q"}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "adapter", Description: "the adapter to install from assets/adapterlist", Examples: []string{"add-adapter reflectsort"}, Type: "string", Default: "", Required: true},
		},
	},
	{
		Identifiers:     []string{"add-arg"},
		Category:        "Cli System",
		Description:     "Add a positional arg to a command's entries.yaml",
		LongDescription: "Inserts one positional arg declaration into\nsandbox/internal/commands/<command>/entries.yaml (at --position, else at\nthe end) and runs build so entries.go and the dispatch layer are\nregenerated. Positional args bind by order; an array arg must stay last.\n",
		Examples:        []string{"add-arg file --type string --required --description \"the file to process\" --command exec", "add-arg count --type int --min 1 --position 0 --command exec"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--command", "-c"}, Description: "the command (identifier or package name) that receives the field", Examples: []string{"--command exec"}, Type: "string", Default: "", Required: true},
			{Identifiers: []string{"--type", "-t"}, Description: "the value type: string, boolean, int or float (defaults to string)", Examples: []string{"--type int"}, Type: "string", Default: "string", Required: false},
			{Identifiers: []string{"--description", "-d"}, Description: "help text shown for the field", Examples: []string{"--description \"where the output is written\""}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--example", "-e"}, Description: "an usage example for the field (repeatable)", Examples: []string{"--example \"exec --out result.txt\""}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--default"}, Description: "the literal assigned when the field is absent (cannot be combined with --required)", Examples: []string{"--default ."}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--required", "-r"}, Description: "fail with a usage error when the field is not provided (not for booleans or fields with --default)", Examples: []string{"--required"}, Type: "boolean", Default: "", Required: false},
			{Identifiers: []string{"--array"}, Description: "collect every occurrence into a []T field instead of a single value", Examples: []string{"--array"}, Type: "boolean", Default: "", Required: false},
			{Identifiers: []string{"--min"}, Description: "smallest accepted value (int/float only)", Examples: []string{"--type int --min 1"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--max"}, Description: "largest accepted value (int/float only)", Examples: []string{"--type int --max 10"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--position"}, Description: "zero-based index to insert the field at (defaults to the end)", Examples: []string{"--position 0"}, Type: "int", Default: "-1", Required: false},
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"add-arg file --command exec --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "name", Description: "the arg name (becomes the generated struct field)", Examples: []string{"add-arg file --command exec"}, Type: "string", Default: "", Required: true},
		},
	},
	{
		Identifiers:     []string{"add-available"},
		Category:        "Deps System",
		Description:     "Declares one further available",
		LongDescription: "Creates adapters/availables/<name>/available.yaml as a copy of the standard selection, so it starts filling every field, and build generates its new.go. Point it at another adapter with set-adapter --available.",
		Examples:        []string{"add-available lambda"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"add-available --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{"add-available lambda -q"}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "available", Description: "the name of the new available (it becomes one directory under adapters/availables/)", Examples: []string{"add-available lambda"}, Type: "string", Default: "", Required: true},
		},
	},
	{
		Identifiers:     []string{"add-body-field"},
		Category:        "Server System",
		Description:     "Declare a property of a route's body json-schema",
		LongDescription: "Declares one property of the route's body json-schema at a dotted path, creating the objects it passes through, and runs build so the Body struct and EntriesSchema pick it up. A route that declared no body becomes a json one here. Every keyword the schema subset supports has a flag; ReadBody answers 400 on the first violation, naming the field path.",
		Examples:        []string{"add-body-field email --route create-user --format email --max 254 --required", "add-body-field address.city --route create-user --required", "add-body-field role --route create-user --enum admin --enum member", "add-body-field tags --route create-user --array --unique-items --max-items 10"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--route"}, Description: "the route (identifier or package name) that receives the property", Examples: []string{"add-body-field email --route create-user"}, Type: "string", Default: "", Required: true},
			{Identifiers: []string{"--type"}, Description: "the value type: string, boolean, int, float or object (defaults to string)", Examples: []string{"add-body-field age --route create-user --type int"}, Type: "string", Default: "string", Required: false},
			{Identifiers: []string{"--required"}, Description: "list the property in its parent object's required set", Examples: []string{"add-body-field email --route create-user --required"}, Type: "boolean", Default: "", Required: false},
			{Identifiers: []string{"--array"}, Description: "declare an array of the type instead of a single value", Examples: []string{"add-body-field tags --route create-user --array"}, Type: "boolean", Default: "", Required: false},
			{Identifiers: []string{"--min"}, Description: "minimum for a number, minLength for a string", Examples: []string{"add-body-field age --route create-user --type int --min 0"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--max"}, Description: "maximum for a number, maxLength for a string", Examples: []string{"add-body-field age --route create-user --type int --max 130"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--exclusive-min"}, Description: "exclusiveMinimum for a number property", Examples: []string{"add-body-field score --route create-user --type float --exclusive-min 0"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--exclusive-max"}, Description: "exclusiveMaximum for a number property", Examples: []string{"add-body-field score --route create-user --type float --exclusive-max 1"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--format"}, Description: "json-schema format for a string property: email, uuid, date-time or uri", Examples: []string{"add-body-field email --route create-user --format email"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--pattern"}, Description: "regular expression a string property must match", Examples: []string{"add-body-field zip --route create-user --pattern \"^[0-9]{5}$\""}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--enum"}, Description: "an accepted value of the property (repeatable; declares the enum set)", Examples: []string{"add-body-field role --route create-user --enum admin --enum member"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--const"}, Description: "the single value the property must carry", Examples: []string{"add-body-field kind --route create-user --const user"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--nullable"}, Description: "accept null as well as the declared type", Examples: []string{"add-body-field nickname --route create-user --nullable"}, Type: "boolean", Default: "", Required: false},
			{Identifiers: []string{"--min-items"}, Description: "shortest accepted array (--array only)", Examples: []string{"add-body-field tags --route create-user --array --min-items 1"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--max-items"}, Description: "longest accepted array (--array only)", Examples: []string{"add-body-field tags --route create-user --array --max-items 10"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--unique-items"}, Description: "refuse an array holding the same value twice (--array only)", Examples: []string{"add-body-field tags --route create-user --array --unique-items"}, Type: "boolean", Default: "", Required: false},
			{Identifiers: []string{"--additional-properties"}, Description: "accept undeclared keys inside an object property", Examples: []string{"add-body-field meta --route create-user --type object --additional-properties"}, Type: "boolean", Default: "", Required: false},
			{Identifiers: []string{"--no-additional-properties"}, Description: "refuse undeclared keys inside an object property", Examples: []string{"add-body-field meta --route create-user --type object --no-additional-properties"}, Type: "boolean", Default: "", Required: false},
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"add-body-field email --route create-user --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "name", Description: "the dotted path of the property (address.city); the objects it passes through are created as needed", Examples: []string{"add-body-field address.city --route create-user"}, Type: "string", Default: "", Required: true},
		},
	},
	{
		Identifiers:     []string{"add-cli-example"},
		Category:        "Examples",
		Description:     "Scaffold a new example under examples/cli/",
		LongDescription: "Creates examples/cli/<name>/ with an example.sh stub that already runs and exits 0, then runs build so the example listing of the docs names it. The golden result.yaml is written by the first exec-test, never by hand. Refuses an existing name, and refuses outright in a project with no cli.",
		Examples:        []string{"add-cli-example start"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"add-cli-example start --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{"add-cli-example start -q"}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "name", Description: "the name of the new example (it becomes one directory under examples/cli/)", Examples: []string{"add-cli-example start"}, Type: "string", Default: "", Required: true},
		},
	},
	{
		Identifiers:     []string{"add-command"},
		Category:        "Cli System",
		Description:     "Scaffold a new command package in the project",
		LongDescription: "Creates sandbox/internal/commands/<name>/ with a hand-written\nentries.yaml and a stub handler.go, then runs build so entries.go\nand the dispatch layer are generated for it. Refuses to overwrite\nan existing command.\n",
		Examples:        []string{"add-command my-feature", "add-command my-feature --path ./my-project"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--help"}, Description: "one-line help text for the new command", Examples: []string{"add-command my-feature --help \"does the thing\" --category Misc"}, Type: "string", Default: "", Required: true},
			{Identifiers: []string{"--category"}, Description: "the category the new command is grouped under in help output", Examples: []string{"add-command my-feature --help \"does the thing\" --category Misc"}, Type: "string", Default: "", Required: true},
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"add-command my-feature --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{"add-command my-feature -q"}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "name", Description: "the name of the new command (e.g. my-feature)", Examples: []string{"add-command my-feature"}, Type: "string", Default: "", Required: true},
		},
	},
	{
		Identifiers:     []string{"add-dep"},
		Category:        "Deps System",
		Description:     "Installs one dep of the embedded catalog into the project",
		LongDescription: "Renders the contract of assets/deplist/<dep> and the adapter that fills it, then calls build. The adapter is the dep's default-adapter unless --adapter names another; it is enrolled in every available.",
		Examples:        []string{"add-dep embeddeps", "add-dep embeddeps --path ./my-project"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--adapter"}, Description: "the adapter to fill the dep's contract with (defaults to the dep's default-adapter)", Examples: []string{"add-dep sortdeps --adapter reflectsort"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"add-dep --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{"add-dep embeddeps -q"}, Type: "boolean", Default: "", Required: false},
			{Identifiers: []string{"--as"}, Description: "the name the copied contract takes under sandbox/deps/ (remote deps only; defaults to the last segment of the module path)", Examples: []string{"add-dep github.com/user/MathLib@v1.2.0 --as mathlib"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--remote-available"}, Description: "the available of the remote repo the generated shim builds its sandbox from", Examples: []string{"add-dep github.com/user/MathLib@v1.2.0 --remote-available standard"}, Type: "string", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "dep", Description: "the dep to install from assets/deplist", Examples: []string{"add-dep embeddeps"}, Type: "string", Default: "", Required: true},
		},
	},
	{
		Identifiers:     []string{"add-doc"},
		Category:        "Documentation",
		Description:     "Scaffold a new doc directory under docs/",
		LongDescription: "Creates docs/<name>/ with a doc.md stub and the props.yaml declaring it,\nthen runs build so README.md's index and the parent's Index.md list it.\nA first-level doc needs at least one --theme of themes.yaml; a nested name\n(docs/<Parent>/<Name>) creates a sub-doc, which takes no theme. Refuses to\noverwrite an existing doc.",
		Examples:        []string{"add-doc HandleReports --theme development --description \"How a report is written and regenerated\"", "add-doc PublicApi/api.AddDoc --description \"The AddDoc action of the sandbox api\""},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--theme", "-t"}, Description: "a theme id of themes.yaml the doc belongs to (repeatable; first-level docs only)", Examples: []string{"--theme development --theme cli-usage"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--description", "-d"}, Description: "the one-line summary every index lists the doc with", Examples: []string{"--description \"How a report is written and regenerated\""}, Type: "string", Default: "", Required: true},
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"add-doc HandleReports --theme development -d \"...\" --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{"add-doc HandleReports --theme development -d \"...\" -q"}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "name", Description: "the doc directory under docs/, nested with / for a sub-doc (e.g. PublicApi/api.AddDoc)", Examples: []string{"add-doc HandleReports"}, Type: "string", Default: "", Required: true},
		},
	},
	{
		Identifiers:     []string{"add-flag"},
		Category:        "Cli System",
		Description:     "Add a flag to a command's entries.yaml",
		LongDescription: "Appends one flag declaration to sandbox/internal/commands/<command>/entries.yaml\nand runs build so entries.go and the dispatch layer are regenerated.\nWithout --identifier the flag answers to --<name>. Refuses a name or\nidentifier the command already uses.\n",
		Examples:        []string{"add-flag output --identifier --out --identifier -o --type string --required --command exec", "add-flag verbose --type boolean --description \"print every step\" --command exec", "add-flag retries --type int --min 0 --max 5 --default 1 --command exec"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--identifier", "-i"}, Description: "a cli identifier for the flag, e.g. --out or -o (repeatable; defaults to --<name>)", Examples: []string{"--identifier --out --identifier -o"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--command", "-c"}, Description: "the command (identifier or package name) that receives the field", Examples: []string{"--command exec"}, Type: "string", Default: "", Required: true},
			{Identifiers: []string{"--type", "-t"}, Description: "the value type: string, boolean, int or float (defaults to string)", Examples: []string{"--type int"}, Type: "string", Default: "string", Required: false},
			{Identifiers: []string{"--description", "-d"}, Description: "help text shown for the field", Examples: []string{"--description \"where the output is written\""}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--example", "-e"}, Description: "an usage example for the field (repeatable)", Examples: []string{"--example \"exec --out result.txt\""}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--default"}, Description: "the literal assigned when the field is absent (cannot be combined with --required)", Examples: []string{"--default ."}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--required", "-r"}, Description: "fail with a usage error when the field is not provided (not for booleans or fields with --default)", Examples: []string{"--required"}, Type: "boolean", Default: "", Required: false},
			{Identifiers: []string{"--array"}, Description: "collect every occurrence into a []T field instead of a single value", Examples: []string{"--array"}, Type: "boolean", Default: "", Required: false},
			{Identifiers: []string{"--min"}, Description: "smallest accepted value (int/float only)", Examples: []string{"--type int --min 1"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--max"}, Description: "largest accepted value (int/float only)", Examples: []string{"--type int --max 10"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--position"}, Description: "zero-based index to insert the field at (defaults to the end)", Examples: []string{"--position 0"}, Type: "int", Default: "-1", Required: false},
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"add-flag out --command exec --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "name", Description: "the flag name (becomes the generated struct field, e.g. out-file -> OutFile)", Examples: []string{"add-flag output --command exec"}, Type: "string", Default: "", Required: true},
		},
	},
	{
		Identifiers:     []string{"add-header"},
		Category:        "Server System",
		Description:     "Declare a request header on a route",
		LongDescription: "Declares one request header on a route and runs build so entries.go and the dispatch arm pick it up. The name is the external spelling and is matched without regard to case; the dispatch answers 400 for a missing --required header or one outside --min/--max, before the handler runs.",
		Examples:        []string{"add-header authorization --route create-user --required", "add-header x-retries --route create-user --type int --default 1 --max 5"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--route"}, Description: "the route (identifier or package name) that receives the header", Examples: []string{"add-header authorization --route create-user"}, Type: "string", Default: "", Required: true},
			{Identifiers: []string{"--type"}, Description: "the value type: string, boolean, int or float (defaults to string)", Examples: []string{"add-header x-retries --route create-user --type int"}, Type: "string", Default: "string", Required: false},
			{Identifiers: []string{"--description"}, Description: "help text shown for the header", Examples: []string{"add-header authorization --route create-user --description \"the bearer token\""}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--example"}, Description: "an usage example for the header (repeatable)", Examples: []string{"add-header authorization --route create-user --example \"Bearer abc\""}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--default"}, Description: "the literal assigned when the header is absent (cannot be combined with --required)", Examples: []string{"add-header x-retries --route create-user --type int --default 1"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--required"}, Description: "answer 400 when the header is not provided (not for booleans or headers with --default)", Examples: []string{"add-header authorization --route create-user --required"}, Type: "boolean", Default: "", Required: false},
			{Identifiers: []string{"--min"}, Description: "smallest accepted value (int/float only)", Examples: []string{"add-header x-retries --route create-user --type int --min 0"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--max"}, Description: "largest accepted value (int/float only)", Examples: []string{"add-header x-retries --route create-user --type int --max 5"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--position"}, Description: "zero-based index to insert the header at (defaults to the end)", Examples: []string{"add-header authorization --route create-user --position 0"}, Type: "int", Default: "-1", Required: false},
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"add-header authorization --route create-user --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "name", Description: "the header name, matched without regard to case", Examples: []string{"add-header authorization --route create-user"}, Type: "string", Default: "", Required: true},
		},
	},
	{
		Identifiers:     []string{"add-lib-example"},
		Category:        "Examples",
		Description:     "Scaffold a new example under examples/lib/",
		LongDescription: "Creates examples/lib/<name>/ with an example.go stub (package main) that already runs and exits 0, then runs build so the example listing of the docs names it. The golden result.yaml is written by the first exec-test, never by hand. Refuses an existing name.",
		Examples:        []string{"add-lib-example start"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"add-lib-example start --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{"add-lib-example start -q"}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "name", Description: "the name of the new example (it becomes one directory under examples/lib/)", Examples: []string{"add-lib-example start"}, Type: "string", Default: "", Required: true},
		},
	},
	{
		Identifiers:     []string{"add-page"},
		Category:        "Front System",
		Description:     "Declare a new html page",
		LongDescription: "Declares the route that answers the page and writes the html template it renders under assets/frontend/pages/. The trigger defaults to /<name>; --trigger / declares the home page. An html file already there is kept, which is the way back from a front-purge.",
		Examples:        []string{"add-page home --trigger / --title Home", "add-page about --title About"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{}, Type: "boolean", Default: "", Required: false},
			{Identifiers: []string{"--trigger"}, Description: "the literal segment the page answers on, / included (defaults to /<name>)", Examples: []string{}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--title"}, Description: "the <title> the scaffolded page carries (defaults to the page name)", Examples: []string{}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--help"}, Description: "one-line description of the page, for docs/Routes (defaults to one derived from the name)", Examples: []string{}, Type: "string", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "name", Description: "the page name (becomes the route, its Go package and the html file)", Examples: []string{}, Type: "string", Default: "", Required: true},
		},
	},
	{
		Identifiers:     []string{"add-param"},
		Category:        "Server System",
		Description:     "Declare a query parameter on a route",
		LongDescription: "Declares one query-string parameter on a route and runs build so entries.go and the dispatch arm pick it up. --array collects every occurrence of the key into a []T field; the only other place it is accepted is the last segment of a route's paths, which takes the rest of the path.",
		Examples:        []string{"add-param page --route list-users --type int --default 1 --min 1", "add-param tag --route list-users --array"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--route"}, Description: "the route (identifier or package name) that receives the parameter", Examples: []string{"add-param page --route list-users"}, Type: "string", Default: "", Required: true},
			{Identifiers: []string{"--type"}, Description: "the value type: string, boolean, int or float (defaults to string)", Examples: []string{"add-param page --route list-users --type int"}, Type: "string", Default: "string", Required: false},
			{Identifiers: []string{"--description"}, Description: "help text shown for the parameter", Examples: []string{"add-param page --route list-users --description \"the page to read\""}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--example"}, Description: "an usage example for the parameter (repeatable)", Examples: []string{"add-param page --route list-users --example \"?page=2\""}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--default"}, Description: "the literal assigned when the parameter is absent (cannot be combined with --required)", Examples: []string{"add-param page --route list-users --type int --default 1"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--required"}, Description: "answer 400 when the parameter is not provided (not for booleans or parameters with --default)", Examples: []string{"add-param query --route list-users --required"}, Type: "boolean", Default: "", Required: false},
			{Identifiers: []string{"--array"}, Description: "collect every occurrence into a []T field instead of a single value", Examples: []string{"add-param tag --route list-users --array"}, Type: "boolean", Default: "", Required: false},
			{Identifiers: []string{"--min"}, Description: "smallest accepted value (int/float only)", Examples: []string{"add-param page --route list-users --type int --min 1"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--max"}, Description: "largest accepted value (int/float only)", Examples: []string{"add-param page --route list-users --type int --max 100"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--position"}, Description: "zero-based index to insert the parameter at (defaults to the end)", Examples: []string{"add-param page --route list-users --position 0"}, Type: "int", Default: "-1", Required: false},
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"add-param page --route list-users --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "name", Description: "the query key", Examples: []string{"add-param page --route list-users"}, Type: "string", Default: "", Required: true},
		},
	},
	{
		Identifiers:     []string{"add-route"},
		Category:        "Server System",
		Description:     "Declare a new http route",
		LongDescription: "Writes sandbox/internal/routes/<name>/route.yaml and a stub handler.go, then runs build so entries.go and the dispatch arm are generated. The trigger is normalized to start with /, and defaults to /<name>.",
		Examples:        []string{"add-route create-user --trigger /users --method POST --help \"Create a user\" --category Users"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--trigger"}, Description: "the first literal segment of the path, always starting with / (defaults to /<name>)", Examples: []string{"--trigger /users"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--method", "-m"}, Description: "the http method the route answers: GET, POST, PUT, PATCH, DELETE, HEAD or OPTIONS", Examples: []string{"--method POST"}, Type: "string", Default: "GET", Required: false},
			{Identifiers: []string{"--help"}, Description: "one-line description of the route", Examples: []string{"--help \"Create a user\""}, Type: "string", Default: "", Required: true},
			{Identifiers: []string{"--category"}, Description: "the heading the route is listed under in docs/Routes", Examples: []string{"--category Users"}, Type: "string", Default: "", Required: true},
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"add-route --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "name", Description: "the route name (becomes the directory sandbox/internal/routes/<name> and its Go package)", Examples: []string{"add-route create-user"}, Type: "string", Default: "", Required: true},
		},
	},
	{
		Identifiers:     []string{"add-segment"},
		Category:        "Server System",
		Description:     "Add a segment to a route's path",
		LongDescription: "Appends one segment to the route's paths and runs build so entries.go and the dispatch arm pick it up. With --identifier the segment is a literal, normalized to start with /; with a name it is a capture, which is always required and becomes an Entries field already converted. --array makes that capture take every segment left in the path into a []T field, which only the last segment of a route may do.",
		Examples:        []string{"add-segment --route create-user --identifier /users", "add-segment tenant --route create-user --description \"the tenant the user belongs to\"", "add-segment page --route list-users --type int --min 1", "add-segment rest --route static --array"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--route"}, Description: "the route (identifier or package name) that receives the segment", Examples: []string{"add-segment tenant --route create-user"}, Type: "string", Default: "", Required: true},
			{Identifiers: []string{"--identifier"}, Description: "declare a literal segment instead of a capture, always starting with /", Examples: []string{"add-segment --route create-user --identifier /users"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--type"}, Description: "the value type of a captured segment: string, boolean, int or float", Examples: []string{"add-segment page --route list-users --type int"}, Type: "string", Default: "string", Required: false},
			{Identifiers: []string{"--description"}, Description: "help text shown for the captured segment", Examples: []string{"add-segment tenant --route create-user --description \"the tenant the user belongs to\""}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--example"}, Description: "an usage example for the segment (repeatable)", Examples: []string{"add-segment tenant --route create-user --example /acme"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--array"}, Description: "take every segment left in the path into a []T field (the last segment only)", Examples: []string{"add-segment rest --route static --array"}, Type: "boolean", Default: "", Required: false},
			{Identifiers: []string{"--min"}, Description: "smallest accepted value (int/float only)", Examples: []string{"add-segment page --route list-users --type int --min 1"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--max"}, Description: "largest accepted value (int/float only)", Examples: []string{"add-segment page --route list-users --type int --max 100"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--position"}, Description: "zero-based index to insert the segment at (defaults to the end)", Examples: []string{"add-segment tenant --route create-user --position 0"}, Type: "int", Default: "-1", Required: false},
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"add-segment tenant --route create-user --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "name", Description: "the captured segment's name (omitted when --identifier declares a literal segment)", Examples: []string{"add-segment tenant --route create-user"}, Type: "string", Default: "", Required: false},
		},
	},
	{
		Identifiers:     []string{"build"},
		Category:        "Core Commands",
		Description:     "Build the project in a directory",
		LongDescription: "Re-renders every generated file of the project in the given\ndirectory, then hands the result to the runtime named by\n--runtime (\"go\" resolves the module graph and compiles every\npackage, \"none\" renders only). If no path is provided, the\ncurrent directory is used.\n",
		Examples:        []string{"build", "build --path ./my-project", "build -q"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"build --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{"build -q"}, Type: "boolean", Default: "", Required: false},
			{Identifiers: []string{"--runtime"}, Description: "the toolchain the rendered project is handed to: go (tidy + compile) or none", Examples: []string{"build --runtime none"}, Type: "string", Default: "go", Required: false},
			{Identifiers: []string{"--unsafe"}, Description: "Skips the verify schema gate before building", Examples: []string{"build --unsafe"}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{},
	},
	{
		Identifiers:     []string{"cli-init"},
		Category:        "Cli System",
		Description:     "Initializes the CLI layer for the project",
		LongDescription: "Installs the std and argv deps the CLI layer depends on, renders the\n\"cli\" asset group into the project, and calls build.\n",
		Examples:        []string{"cli-init", "cli-init --path ./my-project"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"cli-init --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{"cli-init -q"}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{},
	},
	{
		Identifiers:     []string{"cli-purge"},
		Category:        "Cli System",
		Description:     "Removes the CLI layer from the project",
		LongDescription: "Removes every file the \"cli\" asset group installs and calls build.\n",
		Examples:        []string{"cli-purge", "cli-purge --path ./my-project"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"cli-purge --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{"cli-purge -q"}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{},
	},
	{
		Identifiers:     []string{"compile"},
		Category:        "Core Commands",
		Description:     "Cross-compile the project's binaries into release/",
		LongDescription: "Runs build over the project and then cross-compiles its ./cmd/main entrypoint once per --target into release/, with CGO disabled. Repeat --target for several targets, or pass --target all to build every one. Targets and their outputs: linux86 -> linux86.out, linuxarm64 -> linuxarm64.out, linuxi32 -> linuxi32.out, mac86 -> mac86.bin, macarm64 -> macarm64.bin, windows86 -> windows86.exe, windowsi32 -> windowsi32.exe.",
		Examples:        []string{"compile --target linux86", "compile --target linux86 --target macarm64", "compile --target all"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--target", "-t"}, Description: "a target to cross-compile (repeatable); one of linux86, linuxarm64, linuxi32, mac86, macarm64, windows86, windowsi32, or all", Examples: []string{"compile --target linux86 --target windows86"}, Type: "string", Default: "", Required: true},
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"compile --target all --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{"compile --target all -q"}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{},
	},
	{
		Identifiers:     []string{"deps-init"},
		Category:        "Deps System",
		Description:     "Initializes the dependency-injection subsystem for the project",
		LongDescription: "Creates the sandbox/deps and adapters directories and calls build.\nRun this once before using add-dep.",
		Examples:        []string{"deps-init", "deps-init --path ./my-project"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"deps-init --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{"deps-init -q"}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{},
	},
	{
		Identifiers:     []string{"deps-purge"},
		Category:        "Deps System",
		Description:     "Removes the dependency-injection subsystem from the project",
		LongDescription: "Removes the sandbox/deps and adapters directories and calls build.\n",
		Examples:        []string{"deps-purge", "deps-purge --path ./my-project"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"deps-purge --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{"deps-purge -q"}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{},
	},
	{
		Identifiers:     []string{"exec-test"},
		Category:        "Examples",
		Description:     "Run the project's examples and check them against their goldens",
		LongDescription: "Runs every example of examples/cli/ and examples/lib/ in alphabetical order, cli side first, each with its own directory as the working directory and the project's own cli in front of the PATH. Every run starts from a removed TestDir and a removed AssertDir, and what it produced - the merged output, the exit status and the sha256 of every file the example copied out of TestDir into AssertDir - is compared against the example's result.yaml, or written there when that golden does not exist yet. An example that copied nothing out fails: it asserted nothing. An example declared on both sides must leave the same tree and exit the same way: the cli is only a wrapper over the lib.",
		Examples:        []string{"exec-test", "exec-test --only start", "exec-test --update"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--only"}, Description: "run a single example by name, both sides (defaults to every example)", Examples: []string{"exec-test --only start"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--update"}, Description: "rewrite every golden result.yaml with what this run produced instead of comparing", Examples: []string{"exec-test --update"}, Type: "boolean", Default: "", Required: false},
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"exec-test --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{"exec-test -q"}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{},
	},
	{
		Identifiers:     []string{"front-init"},
		Category:        "Front System",
		Description:     "Add the html front layer to the project",
		LongDescription: "Installs the deps the front layer needs, renders sandbox/internal/pageio and writes, once, the route serving assets/frontend/static and that tree's skeleton. A project with no server layer is given one first: a page is answered over http.",
		Examples:        []string{"front-init", "front-init --path ./my-project"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{},
	},
	{
		Identifiers:     []string{"front-purge"},
		Category:        "Front System",
		Description:     "Remove the html front layer from the project",
		LongDescription: "Removes sandbox/internal/pageio, the static route and the route of every declared page, then rebuilds. assets/frontend/ is left untouched: pages, styles and scripts are hand-written content, so front-init followed by add-page puts the routes back over it.",
		Examples:        []string{"front-purge"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{},
	},
	{
		Identifiers:     []string{"help", "--help"},
		Category:        "Info",
		Description:     "Display help for a command",
		LongDescription: "When called without arguments, lists every available command\ngrouped by category. When called with a command name, shows\ndetailed usage, arguments, flags, and examples for that command.\n",
		Examples:        []string{"help", "help start"},
		Hidden:          false,
		Flags:           []helpField{},
		Args: []helpField{
			{Name: "command", Description: "The command to describe; omit it to list every command", Examples: []string{"help start"}, Type: "string", Default: "", Required: false},
		},
	},
	{
		Identifiers:     []string{"list-adapters"},
		Category:        "Deps System",
		Description:     "Lists the adapters of the catalog and of the project",
		LongDescription: "One row per adapter: the name, the dep it fills, whether it is installed, the availables binding it, and what backs it.",
		Examples:        []string{"list-adapters"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"list-adapters --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{"list-adapters -q"}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{},
	},
	{
		Identifiers:     []string{"list-deps"},
		Category:        "Deps System",
		Description:     "Lists the deps the embedded catalog can install",
		LongDescription: "One row per dep of the embedded catalog: the name, whether the project has the contract installed, the adapters filling it (or the catalog's default-adapter when none is), and what it provides.",
		Examples:        []string{"list-deps"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"list-deps --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{"list-deps -q"}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{},
	},
	{
		Identifiers:     []string{"local-install"},
		Category:        "Core Commands",
		Description:     "Builds the project and installs it locally",
		LongDescription: "Runs build over the project, then compiles ./cmd/main into /usr/local/bin/<project-name> (~/.local/bin on Windows) so the binary is on PATH.",
		Examples:        []string{},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{},
	},
	{
		Identifiers:     []string{"publish"},
		Category:        "Core Commands",
		Description:     "Builds, compiles and publishes a release via gh",
		LongDescription: "Runs build, then compile (every target by default), and publishes every file of release/ as a gh release named --release-name, defaulting to the version in AgnosConfig/project.yaml.",
		Examples:        []string{},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--path", "-p"}, Description: "The directory holding the project (defaults to the current directory)", Examples: []string{}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--release-name", "-rn"}, Description: "The name of the release", Examples: []string{}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--draft"}, Description: "Create a draft release", Examples: []string{}, Type: "boolean", Default: "", Required: false},
			{Identifiers: []string{"--target", "-t"}, Description: "The target to compile for (defaults to all)", Examples: []string{}, Type: "string", Default: "all", Required: false},
			{Identifiers: []string{"--publisher", "-pub"}, Description: "The publisher to use (defaults to gh)", Examples: []string{}, Type: "string", Default: "gh", Required: false},
		},
		Args: []helpField{},
	},
	{
		Identifiers:     []string{"remove-adapter"},
		Category:        "Deps System",
		Description:     "Uninstalls one adapter, leaving the contract it filled",
		LongDescription: "Removes adapters/libs/<adapter>/ and the require its declaration pins. Refuses an adapter an available still binds — point that available at another adapter first — and refuses one the generator wrote as the shim of a remote dep.",
		Examples:        []string{"remove-adapter reflectsort"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"remove-adapter reflectsort --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{"remove-adapter reflectsort -q"}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "adapter", Description: "the adapter to remove from the project", Examples: []string{"remove-adapter reflectsort"}, Type: "string", Default: "", Required: true},
		},
	},
	{
		Identifiers:     []string{"remove-arg"},
		Category:        "Cli System",
		Description:     "Remove a positional arg from a command's entries.yaml",
		LongDescription: "Drops one positional arg declaration from\nsandbox/internal/commands/<command>/entries.yaml and runs build so\nentries.go and the dispatch layer forget it. Later args shift up.\n",
		Examples:        []string{"remove-arg file --command exec"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--command", "-c"}, Description: "the command (identifier or package name) that owns the arg", Examples: []string{"--command exec"}, Type: "string", Default: "", Required: true},
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"remove-arg file --command exec --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "name", Description: "the arg name", Examples: []string{"remove-arg file --command exec"}, Type: "string", Default: "", Required: true},
		},
	},
	{
		Identifiers:     []string{"remove-available"},
		Category:        "Deps System",
		Description:     "Deletes one available",
		LongDescription: "Removes adapters/availables/<name>/ whole. The standard available is refused: cmd/main/main.go imports it.",
		Examples:        []string{"remove-available lambda"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"remove-available --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{"remove-available lambda -q"}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "available", Description: "the available to remove", Examples: []string{"remove-available lambda"}, Type: "string", Default: "", Required: true},
		},
	},
	{
		Identifiers:     []string{"remove-body-field"},
		Category:        "Server System",
		Description:     "Delete one property from a route's body json-schema",
		LongDescription: "Drops one property of the body json-schema, named by the same dotted path add-body-field declared it with, and unlists it from its parent's required set. The build renders only: dropping a property may leave hand-written code referring to what is gone.",
		Examples:        []string{"remove-body-field address.city --route create-user"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--route"}, Description: "the route (identifier or package name) the property is declared on", Examples: []string{"remove-body-field email --route create-user"}, Type: "string", Default: "", Required: true},
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"remove-body-field email --route create-user --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "name", Description: "the dotted path of the property to drop", Examples: []string{"remove-body-field address.city --route create-user"}, Type: "string", Default: "", Required: true},
		},
	},
	{
		Identifiers:     []string{"remove-cli-example"},
		Category:        "Examples",
		Description:     "Delete an example from examples/cli/",
		LongDescription: "Deletes examples/cli/<name>/ whole - the example.sh, the golden result.yaml and any TestDir or AssertDir the last run left behind - and runs build so the example listing of the docs is rewritten without it.",
		Examples:        []string{"remove-cli-example start"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"remove-cli-example start --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{"remove-cli-example start -q"}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "name", Description: "the example directory under examples/cli/", Examples: []string{"remove-cli-example start"}, Type: "string", Default: "", Required: true},
		},
	},
	{
		Identifiers:     []string{"remove-command"},
		Category:        "Cli System",
		Description:     "Delete a command package from the project",
		LongDescription: "Deletes sandbox/internal/commands/<name>/ (entries.yaml, entries.go,\nhandler.go and anything else inside) and runs build so climain.go and\nhelp stop dispatching to it. The generated help command cannot be removed.\n",
		Examples:        []string{"remove-command my-feature", "remove-command my-feature --path ./my-project"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"remove-command my-feature --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "name", Description: "the command to delete (identifier or package name)", Examples: []string{"remove-command my-feature"}, Type: "string", Default: "", Required: true},
		},
	},
	{
		Identifiers:     []string{"remove-dep"},
		Category:        "Deps System",
		Description:     "Uninstalls one dep from the project",
		LongDescription: "Removes every adapter whose declaration names the dep, its require and its enrollment in every available, then the contract itself, then calls build.",
		Examples:        []string{"remove-dep embeddeps", "remove-dep embeddeps --path ./my-project"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--with-adapters"}, Description: "also remove every adapter that fills the dep (without it, a dep with adapters installed is refused)", Examples: []string{"remove-dep serverdeps --with-adapters"}, Type: "boolean", Default: "", Required: false},
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"remove-dep --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{"remove-dep embeddeps -q"}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "dep", Description: "the dep to remove from the project", Examples: []string{"remove-dep embeddeps"}, Type: "string", Default: "", Required: true},
		},
	},
	{
		Identifiers:     []string{"remove-doc"},
		Category:        "Documentation",
		Description:     "Delete a doc directory from docs/",
		LongDescription: "Deletes docs/<name>/ (doc.md, props.yaml, its assets and every sub-doc\nnested under it) and runs build so the indexes that listed it are rewritten\nwithout it. A theme left with no docs simply stops rendering a section in\nREADME.md; it is not an error, so themes.yaml can keep it.",
		Examples:        []string{"remove-doc HandleReports", "remove-doc PublicApi/api.AddDoc --path ./my-project"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"remove-doc HandleReports --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{"remove-doc HandleReports -q"}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "name", Description: "the doc directory under docs/, nested with / for a sub-doc (e.g. PublicApi/api.AddDoc)", Examples: []string{"remove-doc HandleReports"}, Type: "string", Default: "", Required: true},
		},
	},
	{
		Identifiers:     []string{"remove-flag"},
		Category:        "Cli System",
		Description:     "Remove a flag from a command's entries.yaml",
		LongDescription: "Drops one flag declaration (matched by its name or by one of its\nidentifiers) from sandbox/internal/commands/<command>/entries.yaml and\nruns build so entries.go and the dispatch layer forget it.\n",
		Examples:        []string{"remove-flag output --command exec", "remove-flag --out --command exec"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--command", "-c"}, Description: "the command (identifier or package name) that owns the flag", Examples: []string{"--command exec"}, Type: "string", Default: "", Required: true},
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"remove-flag output --command exec --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "name", Description: "the flag name (or one of its identifiers, e.g. --out)", Examples: []string{"remove-flag output --command exec"}, Type: "string", Default: "", Required: true},
		},
	},
	{
		Identifiers:     []string{"remove-header"},
		Category:        "Server System",
		Description:     "Delete one declared header from a route",
		LongDescription: "Drops one declared header from a route, the exact inverse of add-header. The build renders only: dropping a header may leave hand-written code referring to what is gone.",
		Examples:        []string{"remove-header authorization --route create-user"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--route"}, Description: "the route (identifier or package name) the header is declared on", Examples: []string{"remove-header authorization --route create-user"}, Type: "string", Default: "", Required: true},
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"remove-header authorization --route create-user --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "name", Description: "the header to drop", Examples: []string{"remove-header authorization --route create-user"}, Type: "string", Default: "", Required: true},
		},
	},
	{
		Identifiers:     []string{"remove-lib-example"},
		Category:        "Examples",
		Description:     "Delete an example from examples/lib/",
		LongDescription: "Deletes examples/lib/<name>/ whole - the example.go, the golden result.yaml and any TestDir or AssertDir the last run left behind - and runs build so the example listing of the docs is rewritten without it.",
		Examples:        []string{"remove-lib-example start"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"remove-lib-example start --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{"remove-lib-example start -q"}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "name", Description: "the example directory under examples/lib/", Examples: []string{"remove-lib-example start"}, Type: "string", Default: "", Required: true},
		},
	},
	{
		Identifiers:     []string{"remove-page"},
		Category:        "Front System",
		Description:     "Remove an html page",
		LongDescription: "Deletes the page's route package and its html template both. A route with no html beside it is not a page: remove-route is the editor for those.",
		Examples:        []string{"remove-page about"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "name", Description: "the page to remove", Examples: []string{}, Type: "string", Default: "", Required: true},
		},
	},
	{
		Identifiers:     []string{"remove-param"},
		Category:        "Server System",
		Description:     "Delete one declared query parameter from a route",
		LongDescription: "Drops one declared query parameter from a route, the exact inverse of add-param. The build renders only: dropping a parameter may leave hand-written code referring to what is gone.",
		Examples:        []string{"remove-param page --route list-users"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--route"}, Description: "the route (identifier or package name) the parameter is declared on", Examples: []string{"remove-param page --route list-users"}, Type: "string", Default: "", Required: true},
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"remove-param page --route list-users --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "name", Description: "the query parameter to drop", Examples: []string{"remove-param page --route list-users"}, Type: "string", Default: "", Required: true},
		},
	},
	{
		Identifiers:     []string{"remove-route"},
		Category:        "Server System",
		Description:     "Delete one declared route",
		LongDescription: "Removes sandbox/internal/routes/<name>/ whole and re-renders the dispatch. The build renders only: dropping a route may leave hand-written code referring to what is gone.",
		Examples:        []string{"remove-route create-user"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"remove-route --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "name", Description: "the route to delete (identifier or package name)", Examples: []string{"remove-route create-user"}, Type: "string", Default: "", Required: true},
		},
	},
	{
		Identifiers:     []string{"remove-segment"},
		Category:        "Server System",
		Description:     "Delete one segment from a route's path",
		LongDescription: "Drops one segment from the route's paths, the exact inverse of add-segment: a capture by its name, a literal by the identifier it spells. The build renders only: dropping a segment may leave hand-written code referring to what is gone.",
		Examples:        []string{"remove-segment tenant --route create-user", "remove-segment /users --route create-user"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--route"}, Description: "the route (identifier or package name) the segment is declared on", Examples: []string{"remove-segment tenant --route create-user"}, Type: "string", Default: "", Required: true},
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"remove-segment tenant --route create-user --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "name", Description: "the segment to drop: the captured segment's name, or the identifier of a literal one", Examples: []string{"remove-segment /users --route create-user"}, Type: "string", Default: "", Required: true},
		},
	},
	{
		Identifiers:     []string{"server-init"},
		Category:        "Server System",
		Description:     "Add the http server layer to the project",
		LongDescription: "Installs the deps the server layer needs, renders sandbox/internal/server, the routeio package and the built-in health route, and writes the start-server command. A project with no cli layer is given one first: a server needs a command that starts it.",
		Examples:        []string{"server-init", "server-init --path ./my-project"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"server-init --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{},
	},
	{
		Identifiers:     []string{"server-purge"},
		Category:        "Server System",
		Description:     "Remove the http server layer and every route in it",
		LongDescription: "Drops sandbox/internal/{server,routes,routeio} and the start-server command, then re-renders. The cli layer and the installed deps are left in place.",
		Examples:        []string{"server-purge"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"server-purge --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{},
	},
	{
		Identifiers:     []string{"set-adapter"},
		Category:        "Deps System",
		Description:     "Changes which adapter an available binds for one dep",
		LongDescription: "Rewrites one available.yaml so the named adapter is the one bound for that dep, dropping whichever adapter filled the field before. It is the only editor of that choice.",
		Examples:        []string{"set-adapter sortdeps reflectsort", "set-adapter sortdeps reflectsort --available lambda"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--available"}, Description: "the available to change (defaults to standard)", Examples: []string{"set-adapter sortdeps reflectsort --available lambda"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"set-adapter sortdeps reflectsort --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{"set-adapter sortdeps reflectsort -q"}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "dep", Description: "the dep whose field is being filled", Examples: []string{"set-adapter sortdeps reflectsort"}, Type: "string", Default: "", Required: true},
			{Name: "adapter", Description: "the installed adapter that should fill it", Examples: []string{"set-adapter sortdeps reflectsort"}, Type: "string", Default: "", Required: true},
		},
	},
	{
		Identifiers:     []string{"set-body"},
		Category:        "Server System",
		Description:     "Rewrite the body keys of a route.yaml",
		LongDescription: "Overwrites the body keys of one route.yaml: how the body is read, whether it is required, the longest one accepted and the content-type the dispatch demands. Empty options leave the current value alone. Declaring a json-schema is add-body-field's job; --drop-schema deletes the one already declared.",
		Examples:        []string{"set-body create-user --type json --required --max-bytes 2097152", "set-body upload-avatar --type raw --content-type application/octet-stream", "set-body ping --type none"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--type"}, Description: "how the body is read: none, raw, text or json", Examples: []string{"set-body create-user --type json"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--required"}, Description: "answer 400 when the body is absent or empty", Examples: []string{"set-body create-user --required"}, Type: "boolean", Default: "", Required: false},
			{Identifiers: []string{"--optional"}, Description: "accept an absent body again", Examples: []string{"set-body create-user --optional"}, Type: "boolean", Default: "", Required: false},
			{Identifiers: []string{"--max-bytes"}, Description: "the longest body accepted, in bytes; a longer one is answered 413", Examples: []string{"set-body create-user --max-bytes 2097152"}, Type: "int", Default: "-1", Required: false},
			{Identifiers: []string{"--content-type"}, Description: "the only content-type accepted; a divergent one is answered 415", Examples: []string{"set-body create-user --content-type text/plain"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--drop-schema"}, Description: "delete the declared json-schema, leaving the body unvalidated", Examples: []string{"set-body create-user --drop-schema"}, Type: "boolean", Default: "", Required: false},
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"set-body create-user --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "route", Description: "the route to edit (identifier or package name)", Examples: []string{"set-body create-user --type json"}, Type: "string", Default: "", Required: true},
		},
	},
	{
		Identifiers:     []string{"set-command"},
		Category:        "Cli System",
		Description:     "Update the command-level keys of a command's entries.yaml",
		LongDescription: "Rewrites help, category, long-description and hidden in\nsandbox/internal/commands/<name>/entries.yaml, and appends extra\nidentifiers / examples, then runs build so help output is regenerated.\nKeys not passed are left untouched.\n",
		Examples:        []string{"set-command exec --help \"run the thing\" --category Core", "set-command exec --identifier run --example \"exec file.txt\"", "set-command exec --hidden"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--help"}, Description: "new one-line help text", Examples: []string{"--help \"run the thing\""}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--category"}, Description: "new category the command is grouped under in help output", Examples: []string{"--category Core"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--long-description"}, Description: "new long description shown by help <command>", Examples: []string{"--long-description \"Runs the thing end to end.\""}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--identifier", "-i"}, Description: "an extra verb the command answers to (repeatable)", Examples: []string{"--identifier run"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--example", "-e"}, Description: "an extra usage example (repeatable)", Examples: []string{"--example \"exec file.txt\""}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--hidden"}, Description: "hide the command from help listings", Examples: []string{}, Type: "boolean", Default: "", Required: false},
			{Identifiers: []string{"--visible"}, Description: "show the command in help listings again", Examples: []string{}, Type: "boolean", Default: "", Required: false},
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"set-command exec --hidden --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "name", Description: "the command to update (identifier or package name)", Examples: []string{"set-command exec --help \"run the thing\""}, Type: "string", Default: "", Required: true},
		},
	},
	{
		Identifiers:     []string{"set-dep"},
		Category:        "Deps System",
		Description:     "Moves one remote dep to another version of its module",
		LongDescription: "Re-copies the remote repo's sandbox/api into sandbox/deps/<dep>/ at the given version and regenerates the shim that converts it, then calls build. The module comes from the shim's own adapter.yaml.",
		Examples:        []string{"set-dep mathlib --version v1.3.0"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--version"}, Description: "the module version to copy the contract from", Examples: []string{"set-dep mathlib --version v1.3.0"}, Type: "string", Default: "", Required: true},
			{Identifiers: []string{"--remote-available"}, Description: "the available of the remote repo the regenerated shim builds its sandbox from", Examples: []string{"set-dep mathlib --version v1.3.0 --remote-available standard"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"set-dep mathlib --version v1.3.0 --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{"set-dep mathlib --version v1.3.0 -q"}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "dep", Description: "the remote dep to move", Examples: []string{"set-dep mathlib --version v1.3.0"}, Type: "string", Default: "", Required: true},
		},
	},
	{
		Identifiers:     []string{"set-route"},
		Category:        "Server System",
		Description:     "Rewrite the route-level keys of a route.yaml",
		LongDescription: "Overwrites method, help, category, long-description, hidden and examples on one route. Empty options leave the current value alone; --example appends.",
		Examples:        []string{"set-route create-user --method POST --example \"curl -X POST localhost:8080/users\""},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--method", "-m"}, Description: "the http method the route answers", Examples: []string{"--method POST"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--help"}, Description: "one-line description of the route", Examples: []string{"--help \"Create a user\""}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--category"}, Description: "the heading the route is listed under in docs/Routes", Examples: []string{"--category Users"}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--long-description"}, Description: "the paragraph docs/Routes prints under the route", Examples: []string{"--long-description \"Creates one user under a tenant.\""}, Type: "string", Default: "", Required: false},
			{Identifiers: []string{"--hidden"}, Description: "drop the route from docs/Routes, still dispatched", Examples: []string{"--hidden"}, Type: "boolean", Default: "", Required: false},
			{Identifiers: []string{"--visible"}, Description: "list the route again in docs/Routes", Examples: []string{"--visible"}, Type: "boolean", Default: "", Required: false},
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"set-route --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{}, Type: "boolean", Default: "", Required: false},
			{Identifiers: []string{"--example"}, Description: "an usage example for the route (repeatable)", Examples: []string{"set-route create-user --example \"curl localhost:8080/users\""}, Type: "string", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "route", Description: "the route to edit (identifier or package name)", Examples: []string{"set-route create-user --help \"Create a user\""}, Type: "string", Default: "", Required: true},
		},
	},
	{
		Identifiers:     []string{"start"},
		Category:        "Core Commands",
		Description:     "Initialize a new project in a directory",
		LongDescription: "Scaffolds a new Agnos project in the given directory, creating\nthe required configuration files and folder structure. If no\npath is provided, the current directory is used.\n",
		Examples:        []string{"start -p my-project", "start -p my-project --path ./my-project-dir", "start -p my-project -q"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"start --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--project-name", "-p"}, Description: "the name of the project", Examples: []string{"start -p my-project"}, Type: "string", Default: "", Required: true},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{"start -q"}, Type: "boolean", Default: "", Required: false},
			{Identifiers: []string{"--force", "-f"}, Description: "Forces the creation of the project, overwriting existing files", Examples: []string{"start -f"}, Type: "boolean", Default: "", Required: false},
			{Identifiers: []string{"--module", "-m"}, Description: "the go module path written into go.mod (required when the target dir has no go.mod yet)", Examples: []string{"start -m github.com/user/project"}, Type: "string", Default: "", Required: false},
		},
		Args: []helpField{},
	},
	{
		Identifiers:     []string{"update-test"},
		Category:        "Examples",
		Description:     "Rewrite one example's golden with what it produces now",
		LongDescription: "Runs one example by name, both sides, and writes what it produced over its result.yaml instead of comparing against it. Every write prints what it changes first - the paths that entered, left or changed sha, and the old output against the new one - so a golden is never rewritten unread. It is the normal way one golden is refreshed; exec-test --update rewrites the whole suite at once and hides the one that moved for a reason nobody meant.",
		Examples:        []string{"update-test start", "update-test add-command --path ./my-project"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "name", Description: "the example to update, both sides", Examples: []string{}, Type: "string", Default: "", Required: true},
		},
	},
	{
		Identifiers:     []string{"verify"},
		Category:        "Core Commands",
		Description:     "Checks the project keeps the sandbox/adapter schema",
		LongDescription: "Verifies the structural rules the harness depends on: sandbox/ imports\nstay inside sandbox/, sandbox/ holds only api, binds, deps, internal and\nnew.go, sandbox/api imports nothing but sandbox/deps and sandbox/deps\nimports nothing external, every sandbox/binds file mirrors a sandbox/api\nfile and declares only functions, and adapters/ holds only availables and\nlibs. `agnos build` runs this as a gate unless --unsafe is passed.",
		Examples:        []string{"verify"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"verify --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--runtime"}, Description: "the toolchain the project is handed to after the schema check: go (tidy + compile) or none", Examples: []string{"verify --runtime none"}, Type: "string", Default: "go", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{"verify -q"}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{},
	},
	{
		Identifiers:     []string{"version", "--version"},
		Category:        "Info",
		Description:     "Print the installed version",
		LongDescription: "Prints the current version of the installed binary and exits.\n",
		Examples:        []string{"version"},
		Hidden:          false,
		Flags:           []helpField{},
		Args:            []helpField{},
	},
}

// ─── Entry points ───────────────────────────────────────────────────────────

// CommandHandler backs the `help` / `--help` verb: with no argument it prints
// the general help screen, with a command name it prints that command's
// detailed help.
func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	name := entries.Command
	if name == "" {
		PrintGeneralHelp(sandbox)
		return exitOk
	}

	for i := range helpCommands {
		if identifiedBy(helpCommands[i].Identifiers, name) {
			printCommandHelp(sandbox, &helpCommands[i])
			return exitOk
		}
	}

	e := sandbox.Deps.Std.Error
	e("\n")
	e("  %s%s✘%s Unknown command: %s%s%s\n", bold, red, reset, bold+white, name, reset)
	e("  %sRun '%s help' to see available commands.%s\n", dim, binaryName(sandbox), reset)
	e("\n")
	return exitUsage
}

// ─── General help ──────────────────────────────────────────────────────────

// PrintGeneralHelp lists every command grouped by category. It is also the
// usage screen shown when the binary is run with no arguments.
func PrintGeneralHelp(sandbox *api.Sandbox) {
	p := sandbox.Deps.Std.Printf

	printBanner(sandbox)

	p("  %s%sUSAGE%s\n", bold, cyan, reset)
	p("  %s│%s\n", gray, reset)
	p("  %s│%s  %s$%s %s %s<command>%s %s[flags]%s %s[args]%s\n",
		gray, reset, dim, reset, binaryName(sandbox),
		green, reset, yellow, reset, dim, reset,
	)
	p("  %s│%s\n", gray, reset)
	p("\n")

	categoryOrder := []string{}
	categorized := map[string][]helpCommand{}
	for _, cmd := range helpCommands {
		if cmd.Hidden {
			continue
		}
		cat := cmd.Category
		if cat == "" {
			cat = "Other"
		}
		if _, exists := categorized[cat]; !exists {
			categoryOrder = append(categoryOrder, cat)
		}
		categorized[cat] = append(categorized[cat], cmd)
	}

	maxNameLen := 0
	for _, cmd := range helpCommands {
		if cmd.Hidden || len(cmd.Identifiers) == 0 {
			continue
		}
		if n := len(cmd.Identifiers[0]); n > maxNameLen {
			maxNameLen = n
		}
	}

	for _, cat := range categoryOrder {
		p("  %s%s%s%s\n", bold, cyan, sandbox.Deps.Stringsdeps.ToUpper(cat), reset)
		p("  %s│%s\n", gray, reset)
		for _, cmd := range categorized[cat] {
			if len(cmd.Identifiers) == 0 {
				continue
			}
			name := cmd.Identifiers[0]

			aliasTag := ""
			if len(cmd.Identifiers) > 1 {
				aliasTag = sandbox.Deps.Std.Sprintf("  %s[%s]%s", dim, sandbox.Deps.Stringsdeps.Join(cmd.Identifiers[1:], ", "), reset)
			}

			dotsNeeded := (maxNameLen + 20) - len(name)
			if dotsNeeded < 4 {
				dotsNeeded = 4
			}
			dots := " " + sandbox.Deps.Stringsdeps.Repeat("·", dotsNeeded-2) + " "

			p("  %s│%s  %s%s%s%s%s%s%s%s\n",
				gray, reset, green+bold, name, reset, gray, dots, reset, cmd.Description, aliasTag,
			)
		}
		p("  %s│%s\n", gray, reset)
		p("\n")
	}

	p("  %s%s─── %sTip%s%s ──────────────────────────────%s\n",
		dim, gray, italic, reset+dim+gray, gray, reset,
	)
	p("  %sRun %s%s help <command>%s%s for detailed info on any command.%s\n",
		dim, reset+cyan, binaryName(sandbox), reset, dim, reset,
	)
	p("\n")
}

// ─── Per-command help ──────────────────────────────────────────────────────

func printCommandHelp(sandbox *api.Sandbox, cmd *helpCommand) {
	p := sandbox.Deps.Std.Printf

	name := cmd.Identifiers[0]

	titleLine := sandbox.Deps.Std.Sprintf("%s %s", binaryName(sandbox), name)
	innerW := len(titleLine) + 4
	if w := len(cmd.Description) + 4; w > innerW {
		innerW = w
	}
	if innerW < 42 {
		innerW = 42
	}

	p("\n")
	p("  %s╭%s╮%s\n", cyan, sandbox.Deps.Stringsdeps.Repeat("─", innerW), reset)
	p("  %s│%s  %s%s%s%s%s│%s\n",
		cyan, reset, bold+white, titleLine, reset,
		sandbox.Deps.Stringsdeps.Repeat(" ", innerW-2-len(titleLine)), cyan, reset,
	)
	p("  %s│%s  %s%s%s%s%s│%s\n",
		cyan, reset, dim, cmd.Description, reset,
		sandbox.Deps.Stringsdeps.Repeat(" ", innerW-2-len(cmd.Description)), cyan, reset,
	)
	p("  %s╰%s╯%s\n", cyan, sandbox.Deps.Stringsdeps.Repeat("─", innerW), reset)
	p("\n")

	if cmd.LongDescription != "" {
		for _, line := range sandbox.Deps.Stringsdeps.Split(cmd.LongDescription, "\n") {
			p("  %s%s%s\n", dim, line, reset)
		}
		p("\n")
	}

	printSection(p, "USAGE")
	usage := sandbox.Deps.Std.Sprintf("  %s$%s %s %s", dim, reset, binaryName(sandbox), name)
	flagPart := ""
	if len(cmd.Flags) > 0 {
		flagPart = sandbox.Deps.Std.Sprintf(" %s[flags]%s", yellow, reset)
	}
	argPart := ""
	for _, arg := range cmd.Args {
		if arg.Required {
			argPart += sandbox.Deps.Std.Sprintf(" %s%s<%s>%s", bold, green, arg.Name, reset)
		} else {
			argPart += sandbox.Deps.Std.Sprintf(" %s[%s]%s", dim, arg.Name, reset)
		}
	}
	p("  %s│%s%s%s%s\n", gray, reset, usage, flagPart, argPart)
	p("  %s│%s\n", gray, reset)
	p("\n")

	if len(cmd.Identifiers) > 1 {
		printSection(p, "ALIASES")
		for _, alias := range cmd.Identifiers {
			bullet := gray + "◦" + reset
			if alias == name {
				bullet = green + "●" + reset
			}
			p("  %s│%s  %s %s%s%s\n", gray, reset, bullet, cyan, alias, reset)
		}
		p("  %s│%s\n", gray, reset)
		p("\n")
	}

	if len(cmd.Args) > 0 {
		printSection(p, "ARGUMENTS")
		for i, arg := range cmd.Args {
			printField(p, arg.Name, arg.Description, arg.Type, arg.Default, arg.Required, arg.Examples)
			if i < len(cmd.Args)-1 {
				p("  %s│%s\n", gray, reset)
			}
		}
		p("  %s│%s\n", gray, reset)
		p("\n")
	}

	if len(cmd.Flags) > 0 {
		printSection(p, "FLAGS")
		for i, flag := range cmd.Flags {
			label := sandbox.Deps.Stringsdeps.Join(flag.Identifiers, gray+", "+reset+yellow+bold)
			printField(p, label, flag.Description, flag.Type, flag.Default, flag.Required, flag.Examples)
			if i < len(cmd.Flags)-1 {
				p("  %s│%s\n", gray, reset)
			}
		}
		p("  %s│%s\n", gray, reset)
		p("\n")
	}

	if len(cmd.Examples) > 0 {
		printSection(p, "EXAMPLES")
		for _, ex := range cmd.Examples {
			p("  %s│%s  %s$%s %s %s\n", gray, reset, dim, reset, binaryName(sandbox), ex)
		}
		p("  %s│%s\n", gray, reset)
		p("\n")
	}
}

// ─── Helpers ───────────────────────────────────────────────────────────────

func printField(p func(string, ...any) (int, error), label, description, kind, def string, required bool, examples []string) {
	reqLabel := dim + "optional" + reset
	if required {
		reqLabel = yellow + bold + "required" + reset
	}

	p("  %s│%s  %s%s%s\n", gray, reset, green+bold, label, reset)
	p("  %s│%s    %s\n", gray, reset, description)
	p("  %s│%s    %s%s%s %s│%s %s\n",
		gray, reset, magenta, typeLabel(kind), reset, gray, reset, reqLabel,
	)
	if def != "" {
		p("  %s│%s    %sdefault:%s %s%s%s\n", gray, reset, dim, reset, white+bold, def, reset)
	}
	for _, ex := range examples {
		p("  %s│%s    %s$ %s%s\n", gray, reset, dim, ex, reset)
	}
}

func printBanner(sandbox *api.Sandbox) {
	p := sandbox.Deps.Std.Printf

	titleLine := sandbox.Deps.Std.Sprintf("%s  %s", config.ProjectName, config.Version)
	innerW := len(titleLine) + 4
	if innerW < 42 {
		innerW = 42
	}

	p("\n")
	p("  %s╭%s╮%s\n", cyan, sandbox.Deps.Stringsdeps.Repeat("─", innerW), reset)
	p("  %s│%s  %s%s%s%s%s│%s\n",
		cyan, reset, bold+white, titleLine, reset,
		sandbox.Deps.Stringsdeps.Repeat(" ", innerW-2-len(titleLine)), cyan, reset,
	)
	p("  %s╰%s╯%s\n", cyan, sandbox.Deps.Stringsdeps.Repeat("─", innerW), reset)
	p("\n")
}

func printSection(p func(string, ...any) (int, error), title string) {
	p("  %s%s%s\n", bold+cyan, title, reset)
	p("  %s│%s\n", gray, reset)
}

func typeLabel(kind string) string {
	switch kind {
	case "int":
		return "int"
	case "float":
		return "float"
	case "boolean":
		return "bool"
	default:
		return "string"
	}
}
