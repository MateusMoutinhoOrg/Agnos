package cli

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/argvdeps"
	add_adapter "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/add_adapter"
	add_arg "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/add_arg"
	add_available "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/add_available"
	add_body_field "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/add_body_field"
	add_cli_example "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/add_cli_example"
	add_command "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/add_command"
	add_dep "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/add_dep"
	add_doc "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/add_doc"
	add_flag "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/add_flag"
	add_header "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/add_header"
	add_lib_example "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/add_lib_example"
	add_page "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/add_page"
	add_param "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/add_param"
	add_route "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/add_route"
	add_segment "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/add_segment"
	build "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/build"
	cli_init "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/cli_init"
	cli_purge "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/cli_purge"
	compile "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/compile"
	deps_init "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/deps_init"
	deps_purge "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/deps_purge"
	exec_test "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/exec_test"
	front_init "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/front_init"
	front_purge "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/front_purge"
	help "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/help"
	list_adapters "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/list_adapters"
	list_deps "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/list_deps"
	local_install "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/local_install"
	publish "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/publish"
	remove_adapter "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/remove_adapter"
	remove_arg "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/remove_arg"
	remove_available "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/remove_available"
	remove_body_field "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/remove_body_field"
	remove_cli_example "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/remove_cli_example"
	remove_command "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/remove_command"
	remove_dep "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/remove_dep"
	remove_doc "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/remove_doc"
	remove_flag "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/remove_flag"
	remove_header "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/remove_header"
	remove_lib_example "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/remove_lib_example"
	remove_page "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/remove_page"
	remove_param "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/remove_param"
	remove_route "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/remove_route"
	remove_segment "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/remove_segment"
	server_init "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/server_init"
	server_purge "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/server_purge"
	set_adapter "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/set_adapter"
	set_body "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/set_body"
	set_command "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/set_command"
	set_dep "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/set_dep"
	set_route "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/set_route"
	start "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/start"
	update_test "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/update_test"
	verify "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/verify"
	version "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/version"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
)

// Exit codes. Kept here (not in sandbox/api) so the cli layer has no
// dependency on the contract package. They mirror sandbox/api: 0 success,
// 1 a well-formed command that failed, 2 a command line that was wrong.
const (
	ExitOk      = 0
	ExitFailure = 1
	ExitUsage   = 2
)

// CliMain is generated by `agnos build` from every
// sandbox/internal/commands/<name>/entries.yaml. It reads the verb, then hands
// the remaining argv to the matching command's generated dispatch function,
// which fills that command's Entries struct and calls its CommandHandler.
// `help` is dispatched through that same path — it is a declared command whose
// three files agnos happens to write itself — and is reached directly only for
// the empty command line below.
func CliMain(deps *deps.Deps, args []string) int {

	if len(args) == 0 {
		help.PrintGeneralHelp(deps)
		return ExitUsage
	}

	verb := deps.Argvdeps.New(args)

	action, err := verb.GetNextStringArg()
	if err != nil {
		help.PrintGeneralHelp(deps)
		return ExitUsage
	}

	switch {
	case action == "add-adapter":
		return dispatchAddAdapter(deps, verb)
	case action == "add-arg":
		return dispatchAddArg(deps, verb)
	case action == "add-available":
		return dispatchAddAvailable(deps, verb)
	case action == "add-body-field":
		return dispatchAddBodyField(deps, verb)
	case action == "add-cli-example":
		return dispatchAddCliExample(deps, verb)
	case action == "add-command":
		return dispatchAddCommand(deps, verb)
	case action == "add-dep":
		return dispatchAddDep(deps, verb)
	case action == "add-doc":
		return dispatchAddDoc(deps, verb)
	case action == "add-flag":
		return dispatchAddFlag(deps, verb)
	case action == "add-header":
		return dispatchAddHeader(deps, verb)
	case action == "add-lib-example":
		return dispatchAddLibExample(deps, verb)
	case action == "add-page":
		return dispatchAddPage(deps, verb)
	case action == "add-param":
		return dispatchAddParam(deps, verb)
	case action == "add-route":
		return dispatchAddRoute(deps, verb)
	case action == "add-segment":
		return dispatchAddSegment(deps, verb)
	case action == "build":
		return dispatchBuild(deps, verb)
	case action == "cli-init":
		return dispatchCliInit(deps, verb)
	case action == "cli-purge":
		return dispatchCliPurge(deps, verb)
	case action == "compile":
		return dispatchCompile(deps, verb)
	case action == "deps-init":
		return dispatchDepsInit(deps, verb)
	case action == "deps-purge":
		return dispatchDepsPurge(deps, verb)
	case action == "exec-test":
		return dispatchExecTest(deps, verb)
	case action == "front-init":
		return dispatchFrontInit(deps, verb)
	case action == "front-purge":
		return dispatchFrontPurge(deps, verb)
	case action == "help" || action == "--help":
		return dispatchHelp(deps, verb)
	case action == "list-adapters":
		return dispatchListAdapters(deps, verb)
	case action == "list-deps":
		return dispatchListDeps(deps, verb)
	case action == "local-install":
		return dispatchLocalInstall(deps, verb)
	case action == "publish":
		return dispatchPublish(deps, verb)
	case action == "remove-adapter":
		return dispatchRemoveAdapter(deps, verb)
	case action == "remove-arg":
		return dispatchRemoveArg(deps, verb)
	case action == "remove-available":
		return dispatchRemoveAvailable(deps, verb)
	case action == "remove-body-field":
		return dispatchRemoveBodyField(deps, verb)
	case action == "remove-cli-example":
		return dispatchRemoveCliExample(deps, verb)
	case action == "remove-command":
		return dispatchRemoveCommand(deps, verb)
	case action == "remove-dep":
		return dispatchRemoveDep(deps, verb)
	case action == "remove-doc":
		return dispatchRemoveDoc(deps, verb)
	case action == "remove-flag":
		return dispatchRemoveFlag(deps, verb)
	case action == "remove-header":
		return dispatchRemoveHeader(deps, verb)
	case action == "remove-lib-example":
		return dispatchRemoveLibExample(deps, verb)
	case action == "remove-page":
		return dispatchRemovePage(deps, verb)
	case action == "remove-param":
		return dispatchRemoveParam(deps, verb)
	case action == "remove-route":
		return dispatchRemoveRoute(deps, verb)
	case action == "remove-segment":
		return dispatchRemoveSegment(deps, verb)
	case action == "server-init":
		return dispatchServerInit(deps, verb)
	case action == "server-purge":
		return dispatchServerPurge(deps, verb)
	case action == "set-adapter":
		return dispatchSetAdapter(deps, verb)
	case action == "set-body":
		return dispatchSetBody(deps, verb)
	case action == "set-command":
		return dispatchSetCommand(deps, verb)
	case action == "set-dep":
		return dispatchSetDep(deps, verb)
	case action == "set-route":
		return dispatchSetRoute(deps, verb)
	case action == "start":
		return dispatchStart(deps, verb)
	case action == "update-test":
		return dispatchUpdateTest(deps, verb)
	case action == "verify":
		return dispatchVerify(deps, verb)
	case action == "version" || action == "--version":
		return dispatchVersion(deps, verb)
	}

	deps.Std.Error("unknown command %q — run '%s help' to see the available commands\n", action, binaryName(deps))
	return ExitUsage
}

// ─── argv helpers ───────────────────────────────────────────────────────────
//
// Every dispatch function below reads the command line through these, so one
// spelling of every usage error is emitted for the whole CLI.

// binaryName is the executable's name as a user types it: the configured
// project name, lowercased.
func binaryName(deps *deps.Deps) string {
	return deps.Stringsdeps.ToLower(config.ProjectName)
}

// silenceLogs turns off the progress channel for the rest of the process. It
// backs the --quiet flag: results (Printf) and errors (Error) still go out,
// only the "… started with path …" notices stop.
func silenceLogs(deps *deps.Deps) {
	deps.Std.Log = func(format string, a ...any) (int, error) {
		return 0, nil
	}
}

// optionValue reads the occurrence-th value of a value flag, reporting a
// clean usage error when the flag was given with nothing after it.
func optionValue(deps *deps.Deps, verb argvdeps.Parser, name string, flags []string, occurrence int) (string, bool) {
	raw, err := verb.GetStringOption(flags, occurrence)
	if err != nil {
		deps.Std.Error("flag '%s': expected a value after %s\n", name, flags[0])
		return "", false
	}
	return raw, true
}

// nextArgValue drains the next positional argument, reporting whether one was
// left. It never fails on the value itself: parsing is the parse* helpers' job,
// so an unparsable argument is told apart from an exhausted command line.
func nextArgValue(verb argvdeps.Parser) (string, bool) {
	raw, err := verb.GetNextStringArg()
	if err != nil {
		return "", false
	}
	return raw, true
}

// parseStringValue is the string arm of the parse* family. Every raw value is
// already a string, so it only exists to keep the generated dispatch uniform
// across the four field types.
func parseStringValue(deps *deps.Deps, subject string, name string, raw string) (string, bool) {
	return raw, true
}

// parseIntValue converts a raw command-line value to an int, reporting the
// failure in the CLI's own words rather than the parser library's.
func parseIntValue(deps *deps.Deps, subject string, name string, raw string) (int, bool) {
	value, err := deps.Stringsdeps.Atoi(raw)
	if err != nil {
		deps.Std.Error("%s '%s': %q is not a valid integer\n", subject, name, raw)
		return 0, false
	}
	return value, true
}

// parseFloatValue converts a raw command-line value to a float64, reporting
// the failure in the CLI's own words rather than the parser library's.
func parseFloatValue(deps *deps.Deps, subject string, name string, raw string) (float64, bool) {
	value, err := deps.Stringsdeps.ParseFloat(raw, 64)
	if err != nil {
		deps.Std.Error("%s '%s': %q is not a valid number\n", subject, name, raw)
		return 0, false
	}
	return value, true
}

// checkUnknownFlags reports the first argument that still looks like a flag
// after every declared flag has been read — a typo such as --pathh, which
// would otherwise be ignored and leave the command running on a default.
func checkUnknownFlags(deps *deps.Deps, verb argvdeps.Parser) bool {
	for i, used := range verb.Used {
		if used || !deps.Stringsdeps.HasPrefix(verb.Args[i], "-") {
			continue
		}
		deps.Std.Error("unknown flag %q — run '%s help' for the accepted flags\n", verb.Args[i], binaryName(deps))
		return false
	}
	return true
}

// checkUnusedArgs reports the first argument left over once every declared
// flag and positional arg has been read.
func checkUnusedArgs(deps *deps.Deps, verb argvdeps.Parser) bool {
	for i, used := range verb.Used {
		if used {
			continue
		}
		deps.Std.Error("unexpected argument %q\n", verb.Args[i])
		return false
	}
	return true
}

func dispatchAddAdapter(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &add_adapter.Entries{}
	if verb.GetOptionsSize([]string{"--available"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "available", []string{"--available"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "available", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Available = value
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "adapter", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Adapter = value
	} else {
		deps.Std.Error("required arg 'adapter' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return add_adapter.CommandHandler(deps, entries)
}

func dispatchAddArg(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &add_arg.Entries{}
	if verb.GetOptionsSize([]string{"--command", "-c"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "command", []string{"--command", "-c"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "command", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Command = value
	} else {
		deps.Std.Error("required flag 'command' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--type", "-t"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "type", []string{"--type", "-t"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "type", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Type = value
	} else {
		entries.Type = "string"
	}
	if verb.GetOptionsSize([]string{"--description", "-d"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "description", []string{"--description", "-d"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "description", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Description = value
	}
	for occurrence := 0; occurrence < verb.GetOptionsSize([]string{"--example", "-e"}); occurrence++ {
		raw, rawOk := optionValue(deps, verb, "example", []string{"--example", "-e"}, occurrence)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "example", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Example = append(entries.Example, value)
	}
	if verb.GetOptionsSize([]string{"--default"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "default", []string{"--default"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "default", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Default = value
	}
	entries.Required = verb.IsPresent([]string{"--required", "-r"})
	entries.Array = verb.IsPresent([]string{"--array"})
	if verb.GetOptionsSize([]string{"--min"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "min", []string{"--min"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "min", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Min = value
	}
	if verb.GetOptionsSize([]string{"--max"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "max", []string{"--max"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "max", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Max = value
	}
	if verb.GetOptionsSize([]string{"--position"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "position", []string{"--position"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseIntValue(deps, "flag", "position", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Position = value
	} else {
		entries.Position = -1
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return add_arg.CommandHandler(deps, entries)
}

func dispatchAddAvailable(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &add_available.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "available", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Available = value
	} else {
		deps.Std.Error("required arg 'available' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return add_available.CommandHandler(deps, entries)
}

func dispatchAddBodyField(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &add_body_field.Entries{}
	if verb.GetOptionsSize([]string{"--route"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "route", []string{"--route"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "route", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Route = value
	} else {
		deps.Std.Error("required flag 'route' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--type"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "type", []string{"--type"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "type", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Type = value
	} else {
		entries.Type = "string"
	}
	entries.Required = verb.IsPresent([]string{"--required"})
	entries.Array = verb.IsPresent([]string{"--array"})
	if verb.GetOptionsSize([]string{"--min"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "min", []string{"--min"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "min", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Min = value
	}
	if verb.GetOptionsSize([]string{"--max"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "max", []string{"--max"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "max", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Max = value
	}
	if verb.GetOptionsSize([]string{"--exclusive-min"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "exclusive-min", []string{"--exclusive-min"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "exclusive-min", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.ExclusiveMin = value
	}
	if verb.GetOptionsSize([]string{"--exclusive-max"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "exclusive-max", []string{"--exclusive-max"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "exclusive-max", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.ExclusiveMax = value
	}
	if verb.GetOptionsSize([]string{"--format"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "format", []string{"--format"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "format", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Format = value
	}
	if verb.GetOptionsSize([]string{"--pattern"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "pattern", []string{"--pattern"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "pattern", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Pattern = value
	}
	for occurrence := 0; occurrence < verb.GetOptionsSize([]string{"--enum"}); occurrence++ {
		raw, rawOk := optionValue(deps, verb, "enum", []string{"--enum"}, occurrence)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "enum", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Enum = append(entries.Enum, value)
	}
	if verb.GetOptionsSize([]string{"--const"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "const", []string{"--const"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "const", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Const = value
	}
	entries.Nullable = verb.IsPresent([]string{"--nullable"})
	if verb.GetOptionsSize([]string{"--min-items"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "min-items", []string{"--min-items"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "min-items", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.MinItems = value
	}
	if verb.GetOptionsSize([]string{"--max-items"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "max-items", []string{"--max-items"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "max-items", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.MaxItems = value
	}
	entries.UniqueItems = verb.IsPresent([]string{"--unique-items"})
	entries.AdditionalProperties = verb.IsPresent([]string{"--additional-properties"})
	entries.NoAdditionalProperties = verb.IsPresent([]string{"--no-additional-properties"})
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return add_body_field.CommandHandler(deps, entries)
}

func dispatchAddCliExample(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &add_cli_example.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return add_cli_example.CommandHandler(deps, entries)
}

func dispatchAddCommand(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &add_command.Entries{}
	if verb.GetOptionsSize([]string{"--help"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "help", []string{"--help"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "help", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Help = value
	} else {
		deps.Std.Error("required flag 'help' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--category"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "category", []string{"--category"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "category", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Category = value
	} else {
		deps.Std.Error("required flag 'category' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return add_command.CommandHandler(deps, entries)
}

func dispatchAddDep(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &add_dep.Entries{}
	if verb.GetOptionsSize([]string{"--adapter"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "adapter", []string{"--adapter"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "adapter", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Adapter = value
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if verb.GetOptionsSize([]string{"--as"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "as", []string{"--as"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "as", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.As = value
	}
	if verb.GetOptionsSize([]string{"--remote-available"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "remote-available", []string{"--remote-available"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "remote-available", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.RemoteAvailable = value
	}
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "dep", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Dep = value
	} else {
		deps.Std.Error("required arg 'dep' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return add_dep.CommandHandler(deps, entries)
}

func dispatchAddDoc(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &add_doc.Entries{}
	for occurrence := 0; occurrence < verb.GetOptionsSize([]string{"--theme", "-t"}); occurrence++ {
		raw, rawOk := optionValue(deps, verb, "theme", []string{"--theme", "-t"}, occurrence)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "theme", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Theme = append(entries.Theme, value)
	}
	if verb.GetOptionsSize([]string{"--description", "-d"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "description", []string{"--description", "-d"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "description", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Description = value
	} else {
		deps.Std.Error("required flag 'description' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return add_doc.CommandHandler(deps, entries)
}

func dispatchAddFlag(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &add_flag.Entries{}
	for occurrence := 0; occurrence < verb.GetOptionsSize([]string{"--identifier", "-i"}); occurrence++ {
		raw, rawOk := optionValue(deps, verb, "identifier", []string{"--identifier", "-i"}, occurrence)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "identifier", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Identifier = append(entries.Identifier, value)
	}
	if verb.GetOptionsSize([]string{"--command", "-c"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "command", []string{"--command", "-c"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "command", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Command = value
	} else {
		deps.Std.Error("required flag 'command' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--type", "-t"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "type", []string{"--type", "-t"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "type", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Type = value
	} else {
		entries.Type = "string"
	}
	if verb.GetOptionsSize([]string{"--description", "-d"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "description", []string{"--description", "-d"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "description", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Description = value
	}
	for occurrence := 0; occurrence < verb.GetOptionsSize([]string{"--example", "-e"}); occurrence++ {
		raw, rawOk := optionValue(deps, verb, "example", []string{"--example", "-e"}, occurrence)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "example", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Example = append(entries.Example, value)
	}
	if verb.GetOptionsSize([]string{"--default"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "default", []string{"--default"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "default", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Default = value
	}
	entries.Required = verb.IsPresent([]string{"--required", "-r"})
	entries.Array = verb.IsPresent([]string{"--array"})
	if verb.GetOptionsSize([]string{"--min"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "min", []string{"--min"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "min", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Min = value
	}
	if verb.GetOptionsSize([]string{"--max"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "max", []string{"--max"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "max", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Max = value
	}
	if verb.GetOptionsSize([]string{"--position"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "position", []string{"--position"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseIntValue(deps, "flag", "position", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Position = value
	} else {
		entries.Position = -1
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return add_flag.CommandHandler(deps, entries)
}

func dispatchAddHeader(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &add_header.Entries{}
	if verb.GetOptionsSize([]string{"--route"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "route", []string{"--route"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "route", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Route = value
	} else {
		deps.Std.Error("required flag 'route' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--type"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "type", []string{"--type"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "type", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Type = value
	} else {
		entries.Type = "string"
	}
	if verb.GetOptionsSize([]string{"--description"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "description", []string{"--description"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "description", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Description = value
	}
	for occurrence := 0; occurrence < verb.GetOptionsSize([]string{"--example"}); occurrence++ {
		raw, rawOk := optionValue(deps, verb, "example", []string{"--example"}, occurrence)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "example", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Example = append(entries.Example, value)
	}
	if verb.GetOptionsSize([]string{"--default"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "default", []string{"--default"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "default", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Default = value
	}
	entries.Required = verb.IsPresent([]string{"--required"})
	if verb.GetOptionsSize([]string{"--min"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "min", []string{"--min"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "min", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Min = value
	}
	if verb.GetOptionsSize([]string{"--max"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "max", []string{"--max"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "max", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Max = value
	}
	if verb.GetOptionsSize([]string{"--position"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "position", []string{"--position"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseIntValue(deps, "flag", "position", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Position = value
	} else {
		entries.Position = -1
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return add_header.CommandHandler(deps, entries)
}

func dispatchAddLibExample(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &add_lib_example.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return add_lib_example.CommandHandler(deps, entries)
}

func dispatchAddPage(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &add_page.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if verb.GetOptionsSize([]string{"--trigger"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "trigger", []string{"--trigger"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "trigger", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Trigger = value
	}
	if verb.GetOptionsSize([]string{"--title"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "title", []string{"--title"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "title", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Title = value
	}
	if verb.GetOptionsSize([]string{"--help"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "help", []string{"--help"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "help", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Help = value
	}
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return add_page.CommandHandler(deps, entries)
}

func dispatchAddParam(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &add_param.Entries{}
	if verb.GetOptionsSize([]string{"--route"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "route", []string{"--route"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "route", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Route = value
	} else {
		deps.Std.Error("required flag 'route' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--type"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "type", []string{"--type"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "type", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Type = value
	} else {
		entries.Type = "string"
	}
	if verb.GetOptionsSize([]string{"--description"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "description", []string{"--description"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "description", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Description = value
	}
	for occurrence := 0; occurrence < verb.GetOptionsSize([]string{"--example"}); occurrence++ {
		raw, rawOk := optionValue(deps, verb, "example", []string{"--example"}, occurrence)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "example", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Example = append(entries.Example, value)
	}
	if verb.GetOptionsSize([]string{"--default"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "default", []string{"--default"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "default", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Default = value
	}
	entries.Required = verb.IsPresent([]string{"--required"})
	entries.Array = verb.IsPresent([]string{"--array"})
	if verb.GetOptionsSize([]string{"--min"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "min", []string{"--min"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "min", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Min = value
	}
	if verb.GetOptionsSize([]string{"--max"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "max", []string{"--max"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "max", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Max = value
	}
	if verb.GetOptionsSize([]string{"--position"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "position", []string{"--position"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseIntValue(deps, "flag", "position", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Position = value
	} else {
		entries.Position = -1
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return add_param.CommandHandler(deps, entries)
}

func dispatchAddRoute(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &add_route.Entries{}
	if verb.GetOptionsSize([]string{"--trigger"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "trigger", []string{"--trigger"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "trigger", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Trigger = value
	}
	if verb.GetOptionsSize([]string{"--method", "-m"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "method", []string{"--method", "-m"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "method", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Method = value
	} else {
		entries.Method = "GET"
	}
	if verb.GetOptionsSize([]string{"--help"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "help", []string{"--help"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "help", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Help = value
	} else {
		deps.Std.Error("required flag 'help' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--category"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "category", []string{"--category"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "category", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Category = value
	} else {
		deps.Std.Error("required flag 'category' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return add_route.CommandHandler(deps, entries)
}

func dispatchAddSegment(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &add_segment.Entries{}
	if verb.GetOptionsSize([]string{"--route"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "route", []string{"--route"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "route", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Route = value
	} else {
		deps.Std.Error("required flag 'route' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--identifier"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "identifier", []string{"--identifier"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "identifier", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Identifier = value
	}
	if verb.GetOptionsSize([]string{"--type"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "type", []string{"--type"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "type", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Type = value
	} else {
		entries.Type = "string"
	}
	if verb.GetOptionsSize([]string{"--description"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "description", []string{"--description"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "description", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Description = value
	}
	for occurrence := 0; occurrence < verb.GetOptionsSize([]string{"--example"}); occurrence++ {
		raw, rawOk := optionValue(deps, verb, "example", []string{"--example"}, occurrence)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "example", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Example = append(entries.Example, value)
	}
	entries.Array = verb.IsPresent([]string{"--array"})
	if verb.GetOptionsSize([]string{"--min"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "min", []string{"--min"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "min", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Min = value
	}
	if verb.GetOptionsSize([]string{"--max"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "max", []string{"--max"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "max", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Max = value
	}
	if verb.GetOptionsSize([]string{"--position"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "position", []string{"--position"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseIntValue(deps, "flag", "position", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Position = value
	} else {
		entries.Position = -1
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return add_segment.CommandHandler(deps, entries)
}

func dispatchBuild(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &build.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if verb.GetOptionsSize([]string{"--runtime"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "runtime", []string{"--runtime"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "runtime", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Runtime = value
	} else {
		entries.Runtime = "go"
	}
	entries.Unsafe = verb.IsPresent([]string{"--unsafe"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return build.CommandHandler(deps, entries)
}

func dispatchCliInit(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &cli_init.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return cli_init.CommandHandler(deps, entries)
}

func dispatchCliPurge(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &cli_purge.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return cli_purge.CommandHandler(deps, entries)
}

func dispatchCompile(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &compile.Entries{}
	for occurrence := 0; occurrence < verb.GetOptionsSize([]string{"--target", "-t"}); occurrence++ {
		raw, rawOk := optionValue(deps, verb, "target", []string{"--target", "-t"}, occurrence)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "target", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Target = append(entries.Target, value)
	}
	if len(entries.Target) == 0 {
		deps.Std.Error("required flag 'target' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return compile.CommandHandler(deps, entries)
}

func dispatchDepsInit(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &deps_init.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return deps_init.CommandHandler(deps, entries)
}

func dispatchDepsPurge(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &deps_purge.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return deps_purge.CommandHandler(deps, entries)
}

func dispatchExecTest(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &exec_test.Entries{}
	if verb.GetOptionsSize([]string{"--only"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "only", []string{"--only"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "only", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Only = value
	}
	entries.Update = verb.IsPresent([]string{"--update"})
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return exec_test.CommandHandler(deps, entries)
}

func dispatchFrontInit(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &front_init.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return front_init.CommandHandler(deps, entries)
}

func dispatchFrontPurge(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &front_purge.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return front_purge.CommandHandler(deps, entries)
}

func dispatchHelp(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &help.Entries{}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "command", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Command = value
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return help.CommandHandler(deps, entries)
}

func dispatchListAdapters(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &list_adapters.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return list_adapters.CommandHandler(deps, entries)
}

func dispatchListDeps(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &list_deps.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return list_deps.CommandHandler(deps, entries)
}

func dispatchLocalInstall(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &local_install.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return local_install.CommandHandler(deps, entries)
}

func dispatchPublish(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &publish.Entries{}
	if verb.GetOptionsSize([]string{"--path", "-p"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path", "-p"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	if verb.GetOptionsSize([]string{"--release-name", "-rn"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "release_name", []string{"--release-name", "-rn"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "release_name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.ReleaseName = value
	}
	entries.Draft = verb.IsPresent([]string{"--draft"})
	if verb.GetOptionsSize([]string{"--target", "-t"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "target", []string{"--target", "-t"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "target", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Target = value
	} else {
		entries.Target = "all"
	}
	if verb.GetOptionsSize([]string{"--publisher", "-pub"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "publisher", []string{"--publisher", "-pub"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "publisher", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Publisher = value
	} else {
		entries.Publisher = "gh"
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return publish.CommandHandler(deps, entries)
}

func dispatchRemoveAdapter(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &remove_adapter.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "adapter", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Adapter = value
	} else {
		deps.Std.Error("required arg 'adapter' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return remove_adapter.CommandHandler(deps, entries)
}

func dispatchRemoveArg(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &remove_arg.Entries{}
	if verb.GetOptionsSize([]string{"--command", "-c"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "command", []string{"--command", "-c"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "command", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Command = value
	} else {
		deps.Std.Error("required flag 'command' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return remove_arg.CommandHandler(deps, entries)
}

func dispatchRemoveAvailable(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &remove_available.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "available", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Available = value
	} else {
		deps.Std.Error("required arg 'available' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return remove_available.CommandHandler(deps, entries)
}

func dispatchRemoveBodyField(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &remove_body_field.Entries{}
	if verb.GetOptionsSize([]string{"--route"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "route", []string{"--route"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "route", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Route = value
	} else {
		deps.Std.Error("required flag 'route' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return remove_body_field.CommandHandler(deps, entries)
}

func dispatchRemoveCliExample(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &remove_cli_example.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return remove_cli_example.CommandHandler(deps, entries)
}

func dispatchRemoveCommand(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &remove_command.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return remove_command.CommandHandler(deps, entries)
}

func dispatchRemoveDep(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &remove_dep.Entries{}
	entries.WithAdapters = verb.IsPresent([]string{"--with-adapters"})
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "dep", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Dep = value
	} else {
		deps.Std.Error("required arg 'dep' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return remove_dep.CommandHandler(deps, entries)
}

func dispatchRemoveDoc(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &remove_doc.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return remove_doc.CommandHandler(deps, entries)
}

func dispatchRemoveFlag(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &remove_flag.Entries{}
	if verb.GetOptionsSize([]string{"--command", "-c"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "command", []string{"--command", "-c"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "command", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Command = value
	} else {
		deps.Std.Error("required flag 'command' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return remove_flag.CommandHandler(deps, entries)
}

func dispatchRemoveHeader(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &remove_header.Entries{}
	if verb.GetOptionsSize([]string{"--route"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "route", []string{"--route"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "route", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Route = value
	} else {
		deps.Std.Error("required flag 'route' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return remove_header.CommandHandler(deps, entries)
}

func dispatchRemoveLibExample(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &remove_lib_example.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return remove_lib_example.CommandHandler(deps, entries)
}

func dispatchRemovePage(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &remove_page.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return remove_page.CommandHandler(deps, entries)
}

func dispatchRemoveParam(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &remove_param.Entries{}
	if verb.GetOptionsSize([]string{"--route"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "route", []string{"--route"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "route", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Route = value
	} else {
		deps.Std.Error("required flag 'route' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return remove_param.CommandHandler(deps, entries)
}

func dispatchRemoveRoute(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &remove_route.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return remove_route.CommandHandler(deps, entries)
}

func dispatchRemoveSegment(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &remove_segment.Entries{}
	if verb.GetOptionsSize([]string{"--route"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "route", []string{"--route"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "route", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Route = value
	} else {
		deps.Std.Error("required flag 'route' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return remove_segment.CommandHandler(deps, entries)
}

func dispatchServerInit(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &server_init.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return server_init.CommandHandler(deps, entries)
}

func dispatchServerPurge(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &server_purge.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return server_purge.CommandHandler(deps, entries)
}

func dispatchSetAdapter(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &set_adapter.Entries{}
	if verb.GetOptionsSize([]string{"--available"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "available", []string{"--available"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "available", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Available = value
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "dep", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Dep = value
	} else {
		deps.Std.Error("required arg 'dep' not provided\n")
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "adapter", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Adapter = value
	} else {
		deps.Std.Error("required arg 'adapter' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return set_adapter.CommandHandler(deps, entries)
}

func dispatchSetBody(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &set_body.Entries{}
	if verb.GetOptionsSize([]string{"--type"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "type", []string{"--type"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "type", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Type = value
	}
	entries.Required = verb.IsPresent([]string{"--required"})
	entries.Optional = verb.IsPresent([]string{"--optional"})
	if verb.GetOptionsSize([]string{"--max-bytes"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "max-bytes", []string{"--max-bytes"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseIntValue(deps, "flag", "max-bytes", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.MaxBytes = value
	} else {
		entries.MaxBytes = -1
	}
	if verb.GetOptionsSize([]string{"--content-type"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "content-type", []string{"--content-type"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "content-type", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.ContentType = value
	}
	entries.DropSchema = verb.IsPresent([]string{"--drop-schema"})
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "route", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Route = value
	} else {
		deps.Std.Error("required arg 'route' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return set_body.CommandHandler(deps, entries)
}

func dispatchSetCommand(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &set_command.Entries{}
	if verb.GetOptionsSize([]string{"--help"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "help", []string{"--help"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "help", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Help = value
	}
	if verb.GetOptionsSize([]string{"--category"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "category", []string{"--category"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "category", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Category = value
	}
	if verb.GetOptionsSize([]string{"--long-description"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "long-description", []string{"--long-description"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "long-description", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.LongDescription = value
	}
	for occurrence := 0; occurrence < verb.GetOptionsSize([]string{"--identifier", "-i"}); occurrence++ {
		raw, rawOk := optionValue(deps, verb, "identifier", []string{"--identifier", "-i"}, occurrence)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "identifier", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Identifier = append(entries.Identifier, value)
	}
	for occurrence := 0; occurrence < verb.GetOptionsSize([]string{"--example", "-e"}); occurrence++ {
		raw, rawOk := optionValue(deps, verb, "example", []string{"--example", "-e"}, occurrence)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "example", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Example = append(entries.Example, value)
	}
	entries.Hidden = verb.IsPresent([]string{"--hidden"})
	entries.Visible = verb.IsPresent([]string{"--visible"})
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return set_command.CommandHandler(deps, entries)
}

func dispatchSetDep(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &set_dep.Entries{}
	if verb.GetOptionsSize([]string{"--version"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "version", []string{"--version"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "version", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Version = value
	} else {
		deps.Std.Error("required flag 'version' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--remote-available"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "remote-available", []string{"--remote-available"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "remote-available", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.RemoteAvailable = value
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "dep", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Dep = value
	} else {
		deps.Std.Error("required arg 'dep' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return set_dep.CommandHandler(deps, entries)
}

func dispatchSetRoute(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &set_route.Entries{}
	if verb.GetOptionsSize([]string{"--method", "-m"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "method", []string{"--method", "-m"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "method", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Method = value
	}
	if verb.GetOptionsSize([]string{"--help"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "help", []string{"--help"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "help", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Help = value
	}
	if verb.GetOptionsSize([]string{"--category"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "category", []string{"--category"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "category", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Category = value
	}
	if verb.GetOptionsSize([]string{"--long-description"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "long-description", []string{"--long-description"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "long-description", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.LongDescription = value
	}
	entries.Hidden = verb.IsPresent([]string{"--hidden"})
	entries.Visible = verb.IsPresent([]string{"--visible"})
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	for occurrence := 0; occurrence < verb.GetOptionsSize([]string{"--example"}); occurrence++ {
		raw, rawOk := optionValue(deps, verb, "example", []string{"--example"}, occurrence)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "example", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Example = append(entries.Example, value)
	}
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "route", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Route = value
	} else {
		deps.Std.Error("required arg 'route' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return set_route.CommandHandler(deps, entries)
}

func dispatchStart(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &start.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	if verb.GetOptionsSize([]string{"--project-name", "-p"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "project-name", []string{"--project-name", "-p"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "project-name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.ProjectName = value
	} else {
		deps.Std.Error("required flag 'project-name' not provided\n")
		return ExitUsage
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	entries.Force = verb.IsPresent([]string{"--force", "-f"})
	if verb.GetOptionsSize([]string{"--module", "-m"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "module", []string{"--module", "-m"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "module", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Module = value
	}
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return start.CommandHandler(deps, entries)
}

func dispatchUpdateTest(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &update_test.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(deps, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return update_test.CommandHandler(deps, entries)
}

func dispatchVerify(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &verify.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	if verb.GetOptionsSize([]string{"--runtime"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "runtime", []string{"--runtime"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(deps, "flag", "runtime", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Runtime = value
	} else {
		entries.Runtime = "go"
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(deps)
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return verify.CommandHandler(deps, entries)
}

func dispatchVersion(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &version.Entries{}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return version.CommandHandler(deps, entries)
}
