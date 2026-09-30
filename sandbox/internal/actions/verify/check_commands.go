package verify

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/goimportsdeps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/commandconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// commandsDir is the tree the commands are declared in, at any depth.
const commandsDir = utils.CommandsDir

// legacyCommandFiles are the two files a command was declared with before
// command.yaml: a package still carrying one is an old declaration no build
// collects.
var legacyCommandFiles = []string{"entries.yaml", "handler.go"}

// commandHandlerParams is the canonical InternalPureHandler signature the
// generated new.go closes over: the sandbox, the command line's shared
// CommandProps, the command's own Entries and the response it answers
// through.
var commandHandlerParams = []string{"*api.Sandbox", "*commandprops.CommandProps", "*Entries", "*api.CommandResponse"}

// legacyCommandPropsParam is the props parameter a handler took while
// CommandProps was declared in sandbox/api.
const legacyCommandPropsParam = "*api.CommandProps"

// CheckCommands enforces the shape the cli layer's generators read by
// convention, the way CheckRoutes does for the server's: the four files of a
// command package, a parsable declaration carrying its required keys, args
// and flags that can be matched and bound, Entries ids and keys that do not
// collide, and a hand-written handler with the one signature the dispatch
// calls.
//
// A command is a directory holding a command.yaml, in whatever folder of
// sandbox/internal/commands it sits. A project with no
// sandbox/internal/commands has no cli layer and nothing to check.
func CheckCommands(sandbox *api.Sandbox, io *smartio.SmartIO) []string {
	var violations []string
	if !io.IsDir(commandsDir) {
		return violations
	}

	violations = append(violations, checkUnitTree(sandbox, io, commandsDir, utils.CommandConfFile, "command",
		append([]string{"new.go", "entries.go", utils.CommandHandlerFile}, legacyCommandFiles...))...)

	patterns := map[string]string{}
	for _, unit := range utils.CommandDirs(sandbox, io) {
		dir := unit.Dir

		for _, legacy := range legacyCommandFiles {
			if io.IsFile(dir + "/" + legacy) {
				violations = append(violations, commandViolation(dir,
					"carries "+legacy+", which command.yaml and InternalPureHandler.go replaced: the build no longer reads it"))
			}
		}

		violations = append(violations, checkCommandFiles(sandbox, io, dir)...)

		content, err := io.ReadFile(dir + "/" + utils.CommandConfFile)
		if err != nil {
			continue
		}
		conf, err := commandconf.New(sandbox, string(content))
		if err != nil {
			violations = append(violations, commandViolation(dir, utils.CommandConfFile+" is not parsable: "+err.Error()))
			continue
		}
		violations = append(violations, checkCommandDeclaration(sandbox, dir, conf)...)

		// Two commands may share a pattern as long as they sit on
		// different rungs of the chain — that is what a middleware in front
		// of a command is. Sharing a rung as well is the ambiguity: the two
		// would run in an order nothing declares.
		key := conf.Pattern() + " at priority " + sandbox.Deps.Stringsdeps.FormatInt(int64(conf.Priority), 10)
		if other, taken := patterns[key]; taken {
			violations = append(violations, commandViolation(dir,
				"declares "+key+", which "+other+" already declares"))
		} else {
			patterns[key] = dir
		}
	}

	return violations
}

// checkCommandFiles reports a command package missing any of the four files
// every command has, and a handler with another signature than the one the
// dispatch calls.
func checkCommandFiles(sandbox *api.Sandbox, io *smartio.SmartIO, dir string) []string {
	var violations []string

	for _, file := range []string{utils.CommandConfFile, "new.go", "entries.go", utils.CommandHandlerFile} {
		if !io.IsFile(dir + "/" + file) {
			violations = append(violations, commandViolation(dir, "has no "+file))
		}
	}

	content, err := io.ReadFile(dir + "/" + utils.CommandHandlerFile)
	if err != nil {
		return violations
	}
	parsed, err := sandbox.Deps.Goimportsdeps.Parse(string(content))
	if err != nil {
		return append(violations, commandViolation(dir, utils.CommandHandlerFile+" is not parsable Go: "+err.Error()))
	}

	for _, function := range parsed.Functions {
		if function.Name != routeHandlerName {
			continue
		}
		if isCommandHandler(function) {
			return violations
		}
		return append(violations, commandViolation(dir,
			utils.CommandHandlerFile+" declares "+routeHandlerName+" with another signature; the dispatch calls "+
				routeHandlerName+"(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error"+
				legacyPropsHint(function, legacyCommandPropsParam, utils.CommandPropsDir)))
	}
	return append(violations, commandViolation(dir, utils.CommandHandlerFile+" exports no "+routeHandlerName))
}

// isCommandHandler reports whether one parsed declaration is the command
// handler: a plain exported func of the canonical name, parameters and one
// error result.
func isCommandHandler(function goimportsdeps.Function) bool {
	if function.Receiver != "" || len(function.Params) != len(commandHandlerParams) {
		return false
	}
	for i, kind := range commandHandlerParams {
		if function.Params[i].Type != kind {
			return false
		}
	}
	return len(function.Results) == 1 && function.Results[0].Type == "error"
}

// checkCommandDeclaration enforces the rules that survive parsing: the
// required keys, the args, the flags and the Entries ids and keys they bind.
func checkCommandDeclaration(sandbox *api.Sandbox, dir string, conf *commandconf.CommandConf) []string {
	var violations []string
	strs := sandbox.Deps.Stringsdeps

	for _, key := range conf.Legacy {
		violations = append(violations, commandViolation(dir,
			"declares `"+key+"`, which entries.yaml carried and command.yaml replaced: the verb is an arg's trigger, a flag's spellings are its `keys`, a repeatable flag is a *-array type"))
	}

	if !conf.HasPriority {
		violations = append(violations, commandViolation(dir, "declares no `priority`; every command declares the rung it runs on"))
	} else if conf.Priority < 0 {
		violations = append(violations, commandViolation(dir, "declares a negative `priority`; the chain runs from zero upwards"))
	}
	if conf.HasSegments && conf.Segments < 1 {
		violations = append(violations, commandViolation(dir, "declares `segments` below 1; leave it out for a command that takes any count"))
	}
	if len(conf.Args) == 0 {
		violations = append(violations, commandViolation(dir, "declares no `args`; a command reads one segment at least"))
	}

	ids := map[string]string{}
	claim := func(id string, what string) {
		if !isExportedId(id) {
			violations = append(violations, commandViolation(dir, "declares the "+what+" "+id+", which is not an exported ASCII Go name (an uppercase ASCII letter, then ASCII letters and digits); it names a field of Entries"))
		}
		if contains(utils.CommandReservedIds, id) {
			violations = append(violations, commandViolation(dir, "declares the "+what+" "+id+", which Entries already carries"))
		}
		if other, taken := ids[id]; taken {
			violations = append(violations, commandViolation(dir, "declares the "+what+" "+id+" twice ("+other+" too)"))
		}
		ids[id] = what
	}

	for _, arg := range conf.Args {
		claim(arg.Id, "arg")
		label := "arg " + arg.Id
		if arg.Start < 0 {
			violations = append(violations, commandViolation(dir, "declares the "+label+" with a negative `start`"))
		}
		if arg.End != commandconf.LastSegment && arg.End < arg.Start {
			violations = append(violations, commandViolation(dir, "declares the "+label+" ending before it starts; `end` is -1 or not before `start`"))
		}
		if !contains(commandconf.ArgTypes, arg.Type) {
			violations = append(violations, commandViolation(dir, "declares the "+label+" with the unknown type "+arg.Type+
				" (use "+strs.Join(commandconf.ArgTypes, ", ")+")"))
		} else if arg.Type != commandconf.DefaultArgType && arg.End != arg.Start {
			violations = append(violations, commandViolation(dir, "declares the "+label+" as a "+arg.Type+" reading a range; anything but string reads one segment"))
		}
		if arg.Required && arg.HasDefault {
			violations = append(violations, commandViolation(dir, "declares the "+label+" both `required` and with a `default`"))
		}
		if problem := triggerProblem(sandbox, arg.Trigger, !conf.Strict); problem != "" {
			violations = append(violations, commandViolation(dir, "declares the "+label+" with "+problem))
		}
	}

	keys := map[string]string{}
	for _, flag := range conf.Flags {
		claim(flag.Id, "flag")
		label := "flag " + flag.Id
		if len(flag.Keys) == 0 {
			violations = append(violations, commandViolation(dir, "declares the "+label+" with an empty `keys`"))
		}
		for _, key := range flag.Keys {
			if !strs.HasPrefix(key, "-") || key == "-" || key == "--" {
				violations = append(violations, commandViolation(dir, "declares the "+label+" with the key "+key+"; a key starts with - or --"))
			}
			if other, taken := keys[key]; taken {
				violations = append(violations, commandViolation(dir, "declares the key "+key+" on the "+label+" and on the flag "+other))
			}
			keys[key] = flag.Id
		}
		if !contains(commandconf.FlagTypes, flag.Type) {
			violations = append(violations, commandViolation(dir, "declares the "+label+" with the unknown type "+flag.Type+
				" (use "+strs.Join(commandconf.FlagTypes, ", ")+")"))
		}
		if flag.Required && flag.HasDefault {
			violations = append(violations, commandViolation(dir, "declares the "+label+" both `required` and with a `default`"))
		}
		if flag.Required && flag.Type == "boolean" {
			violations = append(violations, commandViolation(dir, "declares the "+label+" a `required` boolean; its absence already means false"))
		}
		if (flag.HasMin || flag.HasMax) && flag.Type != "integer" && flag.Type != "number" && flag.Type != "integer-array" {
			violations = append(violations, commandViolation(dir, "declares a `min` or a `max` on the "+label+", which is no number"))
		}
		if flag.HasMin && flag.HasMax && flag.Min > flag.Max {
			violations = append(violations, commandViolation(dir, "declares the "+label+" with a `min` above its `max`"))
		}
		if len(flag.Enum) > 0 && flag.Type == "boolean" {
			violations = append(violations, commandViolation(dir, "declares an `enum` on the boolean "+label))
		}
		if flag.Pattern != "" {
			if _, err := strs.MatchPattern(flag.Pattern, ""); err != nil {
				violations = append(violations, commandViolation(dir, "declares the "+label+" with a `pattern` that does not compile: "+err.Error()))
			}
		}
		if problem := triggerProblem(sandbox, flag.Trigger, false); problem != "" {
			violations = append(violations, commandViolation(dir, "declares the "+label+" with "+problem))
		}
	}

	return violations
}

// commandViolation words one finding about one command package.
func commandViolation(dir string, reason string) string {
	return dir + " " + reason
}
