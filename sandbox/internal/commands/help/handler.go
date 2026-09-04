package help

import (
	"fmt"
	"slices"
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
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

// binaryName is the executable's name as a user types it: the configured
// project name, lowercased. Usage lines show what to type, not the display
// name of the project.
func binaryName() string {
	return strings.ToLower(config.ProjectName)
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
		Identifiers:     []string{"dep-install"},
		Category:        "Dependencies",
		Description:     "Installs an embedded dep into the project",
		LongDescription: "Renders every file under assets/deplist/<dep> into the project\nat the path it holds inside that dep, then calls build.\n",
		Examples:        []string{"dep-install embeddeps", "dep-install embeddeps --path ./my-project"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"dep-install --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{"dep-install embeddeps -q"}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "dep", Description: "the dep to install from assets/deplist", Examples: []string{"dep-install embeddeps"}, Type: "string", Default: "", Required: true},
		},
	},
	{
		Identifiers:     []string{"dep-list"},
		Category:        "Dependencies",
		Description:     "Lists the embedded deps available to install",
		LongDescription: "Lists the name of every dep under assets/deplist that dep-install\ncan render into a project.\n",
		Examples:        []string{"dep-list"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"dep-list --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{"dep-list -q"}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{},
	},
	{
		Identifiers:     []string{"dep-remove"},
		Category:        "Dependencies",
		Description:     "Removes an embedded dep from the project",
		LongDescription: "Removes every file that assets/deplist/<dep> installs into the\nproject, then calls build.\n",
		Examples:        []string{"dep-remove embeddeps", "dep-remove embeddeps --path ./my-project"},
		Hidden:          false,
		Flags: []helpField{
			{Identifiers: []string{"--path"}, Description: "the dir holding the project (defaults to the current directory)", Examples: []string{"dep-remove --path ./my-project"}, Type: "string", Default: ".", Required: false},
			{Identifiers: []string{"--quiet", "-q"}, Description: "Quiets the cli output", Examples: []string{"dep-remove embeddeps -q"}, Type: "boolean", Default: "", Required: false},
		},
		Args: []helpField{
			{Name: "dep", Description: "the dep to remove from the project", Examples: []string{"dep-remove embeddeps"}, Type: "string", Default: "", Required: true},
		},
	},
	{
		Identifiers:     []string{"deps-init"},
		Category:        "Dependency System",
		Description:     "Initializes the dependency-injection subsystem for the project",
		LongDescription: "Creates the sandbox/deps and adapters directories and calls build.\nRun this once before using dep-install.\n",
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
		Category:        "Dependency System",
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
		Identifiers:     []string{"verify"},
		Category:        "Core Commands",
		Description:     "Checks the project keeps the sandbox/adapter schema",
		LongDescription: "Verifies the structural rules the harness depends on: sandbox/ imports\nstay inside sandbox/, sandbox/ holds only api, binds, deps, internal and\nnew.go, sandbox/api and sandbox/deps import nothing external, every\nsandbox/binds file mirrors a sandbox/api file and declares only functions,\nand adapters/ holds only availables and libs. `agnos build` runs this as a\ngate unless --unsafe is passed.\n",
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
func CommandHandler(deps *deps.Deps, entries *Entries) int {
	name := entries.Command
	if name == "" {
		PrintGeneralHelp(deps)
		return exitOk
	}

	for i := range helpCommands {
		if slices.Contains(helpCommands[i].Identifiers, name) {
			printCommandHelp(deps, &helpCommands[i])
			return exitOk
		}
	}

	e := deps.Std.Error
	e("\n")
	e("  %s%s✘%s Unknown command: %s%s%s\n", bold, red, reset, bold+white, name, reset)
	e("  %sRun '%s help' to see available commands.%s\n", dim, binaryName(), reset)
	e("\n")
	return exitUsage
}

// ─── General help ──────────────────────────────────────────────────────────

// PrintGeneralHelp lists every command grouped by category. It is also the
// usage screen shown when the binary is run with no arguments.
func PrintGeneralHelp(deps *deps.Deps) {
	p := deps.Std.Printf

	printBanner(deps)

	p("  %s%sUSAGE%s\n", bold, cyan, reset)
	p("  %s│%s\n", gray, reset)
	p("  %s│%s  %s$%s %s %s<command>%s %s[flags]%s %s[args]%s\n",
		gray, reset, dim, reset, binaryName(),
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
		p("  %s%s%s%s\n", bold, cyan, strings.ToUpper(cat), reset)
		p("  %s│%s\n", gray, reset)
		for _, cmd := range categorized[cat] {
			if len(cmd.Identifiers) == 0 {
				continue
			}
			name := cmd.Identifiers[0]

			aliasTag := ""
			if len(cmd.Identifiers) > 1 {
				aliasTag = fmt.Sprintf("  %s[%s]%s", dim, strings.Join(cmd.Identifiers[1:], ", "), reset)
			}

			dotsNeeded := (maxNameLen + 20) - len(name)
			if dotsNeeded < 4 {
				dotsNeeded = 4
			}
			dots := " " + strings.Repeat("·", dotsNeeded-2) + " "

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
		dim, reset+cyan, binaryName(), reset, dim, reset,
	)
	p("\n")
}

// ─── Per-command help ──────────────────────────────────────────────────────

func printCommandHelp(deps *deps.Deps, cmd *helpCommand) {
	p := deps.Std.Printf

	name := cmd.Identifiers[0]

	titleLine := fmt.Sprintf("%s %s", binaryName(), name)
	innerW := len(titleLine) + 4
	if w := len(cmd.Description) + 4; w > innerW {
		innerW = w
	}
	if innerW < 42 {
		innerW = 42
	}

	p("\n")
	p("  %s╭%s╮%s\n", cyan, strings.Repeat("─", innerW), reset)
	p("  %s│%s  %s%s%s%s%s│%s\n",
		cyan, reset, bold+white, titleLine, reset,
		strings.Repeat(" ", innerW-2-len(titleLine)), cyan, reset,
	)
	p("  %s│%s  %s%s%s%s%s│%s\n",
		cyan, reset, dim, cmd.Description, reset,
		strings.Repeat(" ", innerW-2-len(cmd.Description)), cyan, reset,
	)
	p("  %s╰%s╯%s\n", cyan, strings.Repeat("─", innerW), reset)
	p("\n")

	if cmd.LongDescription != "" {
		for _, line := range strings.Split(cmd.LongDescription, "\n") {
			p("  %s%s%s\n", dim, line, reset)
		}
		p("\n")
	}

	printSection(p, "USAGE")
	usage := fmt.Sprintf("  %s$%s %s %s", dim, reset, binaryName(), name)
	flagPart := ""
	if len(cmd.Flags) > 0 {
		flagPart = fmt.Sprintf(" %s[flags]%s", yellow, reset)
	}
	argPart := ""
	for _, arg := range cmd.Args {
		if arg.Required {
			argPart += fmt.Sprintf(" %s%s<%s>%s", bold, green, arg.Name, reset)
		} else {
			argPart += fmt.Sprintf(" %s[%s]%s", dim, arg.Name, reset)
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
			label := strings.Join(flag.Identifiers, gray+", "+reset+yellow+bold)
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
			p("  %s│%s  %s$%s %s %s\n", gray, reset, dim, reset, binaryName(), ex)
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

func printBanner(deps *deps.Deps) {
	p := deps.Std.Printf

	titleLine := fmt.Sprintf("%s  %s", config.ProjectName, config.Version)
	innerW := len(titleLine) + 4
	if innerW < 42 {
		innerW = 42
	}

	p("\n")
	p("  %s╭%s╮%s\n", cyan, strings.Repeat("─", innerW), reset)
	p("  %s│%s  %s%s%s%s%s│%s\n",
		cyan, reset, bold+white, titleLine, reset,
		strings.Repeat(" ", innerW-2-len(titleLine)), cyan, reset,
	)
	p("  %s╰%s╯%s\n", cyan, strings.Repeat("─", innerW), reset)
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
