package cli

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
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
func CliMain(sandbox *api.Sandbox, args []string) int {

	if len(args) == 0 {
		help.PrintGeneralHelp(sandbox)
		return ExitUsage
	}

	verb := sandbox.Deps.Argvdeps.New(args)

	action, err := verb.GetNextStringArg()
	if err != nil {
		help.PrintGeneralHelp(sandbox)
		return ExitUsage
	}

	switch {
	case action == "add-adapter":
		return dispatchAddAdapter(sandbox, verb)
	case action == "add-arg":
		return dispatchAddArg(sandbox, verb)
	case action == "add-available":
		return dispatchAddAvailable(sandbox, verb)
	case action == "add-body-field":
		return dispatchAddBodyField(sandbox, verb)
	case action == "add-cli-example":
		return dispatchAddCliExample(sandbox, verb)
	case action == "add-command":
		return dispatchAddCommand(sandbox, verb)
	case action == "add-dep":
		return dispatchAddDep(sandbox, verb)
	case action == "add-doc":
		return dispatchAddDoc(sandbox, verb)
	case action == "add-flag":
		return dispatchAddFlag(sandbox, verb)
	case action == "add-header":
		return dispatchAddHeader(sandbox, verb)
	case action == "add-lib-example":
		return dispatchAddLibExample(sandbox, verb)
	case action == "add-page":
		return dispatchAddPage(sandbox, verb)
	case action == "add-param":
		return dispatchAddParam(sandbox, verb)
	case action == "add-route":
		return dispatchAddRoute(sandbox, verb)
	case action == "add-segment":
		return dispatchAddSegment(sandbox, verb)
	case action == "build":
		return dispatchBuild(sandbox, verb)
	case action == "cli-init":
		return dispatchCliInit(sandbox, verb)
	case action == "cli-purge":
		return dispatchCliPurge(sandbox, verb)
	case action == "compile":
		return dispatchCompile(sandbox, verb)
	case action == "deps-init":
		return dispatchDepsInit(sandbox, verb)
	case action == "deps-purge":
		return dispatchDepsPurge(sandbox, verb)
	case action == "exec-test":
		return dispatchExecTest(sandbox, verb)
	case action == "front-init":
		return dispatchFrontInit(sandbox, verb)
	case action == "front-purge":
		return dispatchFrontPurge(sandbox, verb)
	case action == "help" || action == "--help":
		return dispatchHelp(sandbox, verb)
	case action == "list-adapters":
		return dispatchListAdapters(sandbox, verb)
	case action == "list-deps":
		return dispatchListDeps(sandbox, verb)
	case action == "local-install":
		return dispatchLocalInstall(sandbox, verb)
	case action == "publish":
		return dispatchPublish(sandbox, verb)
	case action == "remove-adapter":
		return dispatchRemoveAdapter(sandbox, verb)
	case action == "remove-arg":
		return dispatchRemoveArg(sandbox, verb)
	case action == "remove-available":
		return dispatchRemoveAvailable(sandbox, verb)
	case action == "remove-body-field":
		return dispatchRemoveBodyField(sandbox, verb)
	case action == "remove-cli-example":
		return dispatchRemoveCliExample(sandbox, verb)
	case action == "remove-command":
		return dispatchRemoveCommand(sandbox, verb)
	case action == "remove-dep":
		return dispatchRemoveDep(sandbox, verb)
	case action == "remove-doc":
		return dispatchRemoveDoc(sandbox, verb)
	case action == "remove-flag":
		return dispatchRemoveFlag(sandbox, verb)
	case action == "remove-header":
		return dispatchRemoveHeader(sandbox, verb)
	case action == "remove-lib-example":
		return dispatchRemoveLibExample(sandbox, verb)
	case action == "remove-page":
		return dispatchRemovePage(sandbox, verb)
	case action == "remove-param":
		return dispatchRemoveParam(sandbox, verb)
	case action == "remove-route":
		return dispatchRemoveRoute(sandbox, verb)
	case action == "remove-segment":
		return dispatchRemoveSegment(sandbox, verb)
	case action == "server-init":
		return dispatchServerInit(sandbox, verb)
	case action == "server-purge":
		return dispatchServerPurge(sandbox, verb)
	case action == "set-adapter":
		return dispatchSetAdapter(sandbox, verb)
	case action == "set-body":
		return dispatchSetBody(sandbox, verb)
	case action == "set-command":
		return dispatchSetCommand(sandbox, verb)
	case action == "set-dep":
		return dispatchSetDep(sandbox, verb)
	case action == "set-route":
		return dispatchSetRoute(sandbox, verb)
	case action == "start":
		return dispatchStart(sandbox, verb)
	case action == "update-test":
		return dispatchUpdateTest(sandbox, verb)
	case action == "verify":
		return dispatchVerify(sandbox, verb)
	case action == "version" || action == "--version":
		return dispatchVersion(sandbox, verb)
	}

	sandbox.Deps.Std.Error("unknown command %q — run '%s help' to see the available commands\n", action, binaryName(sandbox))
	return ExitUsage
}

// ─── argv helpers ───────────────────────────────────────────────────────────
//
// Every dispatch function below reads the command line through these, so one
// spelling of every usage error is emitted for the whole CLI.

// binaryName is the executable's name as a user types it: the configured
// project name, lowercased.
func binaryName(sandbox *api.Sandbox) string {
	return sandbox.Deps.Stringsdeps.ToLower(config.ProjectName)
}

// silenceLogs turns off the progress channel for the rest of the process. It
// backs the --quiet flag: results (Printf) and errors (Error) still go out,
// only the "… started with path …" notices stop.
func silenceLogs(sandbox *api.Sandbox) {
	sandbox.Deps.Std.Log = func(format string, a ...any) (int, error) {
		return 0, nil
	}
}

// optionValue reads the occurrence-th value of a value flag, reporting a
// clean usage error when the flag was given with nothing after it.
func optionValue(sandbox *api.Sandbox, verb argvdeps.Parser, name string, flags []string, occurrence int) (string, bool) {
	raw, err := verb.GetStringOption(flags, occurrence)
	if err != nil {
		sandbox.Deps.Std.Error("flag '%s': expected a value after %s\n", name, flags[0])
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
func parseStringValue(sandbox *api.Sandbox, subject string, name string, raw string) (string, bool) {
	return raw, true
}

// parseIntValue converts a raw command-line value to an int, reporting the
// failure in the CLI's own words rather than the parser library's.
func parseIntValue(sandbox *api.Sandbox, subject string, name string, raw string) (int, bool) {
	value, err := sandbox.Deps.Stringsdeps.Atoi(raw)
	if err != nil {
		sandbox.Deps.Std.Error("%s '%s': %q is not a valid integer\n", subject, name, raw)
		return 0, false
	}
	return value, true
}

// parseFloatValue converts a raw command-line value to a float64, reporting
// the failure in the CLI's own words rather than the parser library's.
func parseFloatValue(sandbox *api.Sandbox, subject string, name string, raw string) (float64, bool) {
	value, err := sandbox.Deps.Stringsdeps.ParseFloat(raw, 64)
	if err != nil {
		sandbox.Deps.Std.Error("%s '%s': %q is not a valid number\n", subject, name, raw)
		return 0, false
	}
	return value, true
}

// checkUnknownFlags reports the first argument that still looks like a flag
// after every declared flag has been read — a typo such as --pathh, which
// would otherwise be ignored and leave the command running on a default.
func checkUnknownFlags(sandbox *api.Sandbox, verb argvdeps.Parser) bool {
	for i, used := range verb.Used {
		if used || !sandbox.Deps.Stringsdeps.HasPrefix(verb.Args[i], "-") {
			continue
		}
		sandbox.Deps.Std.Error("unknown flag %q — run '%s help' for the accepted flags\n", verb.Args[i], binaryName(sandbox))
		return false
	}
	return true
}

// checkUnusedArgs reports the first argument left over once every declared
// flag and positional arg has been read.
func checkUnusedArgs(sandbox *api.Sandbox, verb argvdeps.Parser) bool {
	for i, used := range verb.Used {
		if used {
			continue
		}
		sandbox.Deps.Std.Error("unexpected argument %q\n", verb.Args[i])
		return false
	}
	return true
}

func dispatchAddAdapter(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &add_adapter.Entries{}
	if verb.GetOptionsSize([]string{"--available"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "available", []string{"--available"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "available", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Available = value
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "adapter", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Adapter = value
	} else {
		sandbox.Deps.Std.Error("required arg 'adapter' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return add_adapter.CommandHandler(sandbox, entries)
}

func dispatchAddArg(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &add_arg.Entries{}
	if verb.GetOptionsSize([]string{"--command", "-c"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "command", []string{"--command", "-c"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "command", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Command = value
	} else {
		sandbox.Deps.Std.Error("required flag 'command' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--type", "-t"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "type", []string{"--type", "-t"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "type", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Type = value
	} else {
		entries.Type = "string"
	}
	if verb.GetOptionsSize([]string{"--description", "-d"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "description", []string{"--description", "-d"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "description", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Description = value
	}
	for occurrence := 0; occurrence < verb.GetOptionsSize([]string{"--example", "-e"}); occurrence++ {
		raw, rawOk := optionValue(sandbox, verb, "example", []string{"--example", "-e"}, occurrence)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "example", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Example = append(entries.Example, value)
	}
	if verb.GetOptionsSize([]string{"--default"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "default", []string{"--default"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "default", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Default = value
	}
	entries.Required = verb.IsPresent([]string{"--required", "-r"})
	entries.Array = verb.IsPresent([]string{"--array"})
	if verb.GetOptionsSize([]string{"--min"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "min", []string{"--min"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "min", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Min = value
	}
	if verb.GetOptionsSize([]string{"--max"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "max", []string{"--max"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "max", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Max = value
	}
	if verb.GetOptionsSize([]string{"--position"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "position", []string{"--position"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseIntValue(sandbox, "flag", "position", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Position = value
	} else {
		entries.Position = -1
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		sandbox.Deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return add_arg.CommandHandler(sandbox, entries)
}

func dispatchAddAvailable(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &add_available.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "available", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Available = value
	} else {
		sandbox.Deps.Std.Error("required arg 'available' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return add_available.CommandHandler(sandbox, entries)
}

func dispatchAddBodyField(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &add_body_field.Entries{}
	if verb.GetOptionsSize([]string{"--route"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "route", []string{"--route"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "route", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Route = value
	} else {
		sandbox.Deps.Std.Error("required flag 'route' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--type"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "type", []string{"--type"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "type", raw)
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
		raw, rawOk := optionValue(sandbox, verb, "min", []string{"--min"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "min", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Min = value
	}
	if verb.GetOptionsSize([]string{"--max"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "max", []string{"--max"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "max", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Max = value
	}
	if verb.GetOptionsSize([]string{"--exclusive-min"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "exclusive-min", []string{"--exclusive-min"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "exclusive-min", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.ExclusiveMin = value
	}
	if verb.GetOptionsSize([]string{"--exclusive-max"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "exclusive-max", []string{"--exclusive-max"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "exclusive-max", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.ExclusiveMax = value
	}
	if verb.GetOptionsSize([]string{"--format"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "format", []string{"--format"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "format", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Format = value
	}
	if verb.GetOptionsSize([]string{"--pattern"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "pattern", []string{"--pattern"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "pattern", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Pattern = value
	}
	for occurrence := 0; occurrence < verb.GetOptionsSize([]string{"--enum"}); occurrence++ {
		raw, rawOk := optionValue(sandbox, verb, "enum", []string{"--enum"}, occurrence)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "enum", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Enum = append(entries.Enum, value)
	}
	if verb.GetOptionsSize([]string{"--const"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "const", []string{"--const"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "const", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Const = value
	}
	entries.Nullable = verb.IsPresent([]string{"--nullable"})
	if verb.GetOptionsSize([]string{"--min-items"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "min-items", []string{"--min-items"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "min-items", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.MinItems = value
	}
	if verb.GetOptionsSize([]string{"--max-items"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "max-items", []string{"--max-items"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "max-items", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.MaxItems = value
	}
	entries.UniqueItems = verb.IsPresent([]string{"--unique-items"})
	entries.AdditionalProperties = verb.IsPresent([]string{"--additional-properties"})
	entries.NoAdditionalProperties = verb.IsPresent([]string{"--no-additional-properties"})
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		sandbox.Deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return add_body_field.CommandHandler(sandbox, entries)
}

func dispatchAddCliExample(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &add_cli_example.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		sandbox.Deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return add_cli_example.CommandHandler(sandbox, entries)
}

func dispatchAddCommand(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &add_command.Entries{}
	if verb.GetOptionsSize([]string{"--help"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "help", []string{"--help"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "help", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Help = value
	} else {
		sandbox.Deps.Std.Error("required flag 'help' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--category"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "category", []string{"--category"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "category", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Category = value
	} else {
		sandbox.Deps.Std.Error("required flag 'category' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		sandbox.Deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return add_command.CommandHandler(sandbox, entries)
}

func dispatchAddDep(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &add_dep.Entries{}
	if verb.GetOptionsSize([]string{"--adapter"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "adapter", []string{"--adapter"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "adapter", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Adapter = value
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if verb.GetOptionsSize([]string{"--as"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "as", []string{"--as"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "as", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.As = value
	}
	if verb.GetOptionsSize([]string{"--remote-available"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "remote-available", []string{"--remote-available"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "remote-available", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.RemoteAvailable = value
	}
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "dep", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Dep = value
	} else {
		sandbox.Deps.Std.Error("required arg 'dep' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return add_dep.CommandHandler(sandbox, entries)
}

func dispatchAddDoc(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &add_doc.Entries{}
	for occurrence := 0; occurrence < verb.GetOptionsSize([]string{"--theme", "-t"}); occurrence++ {
		raw, rawOk := optionValue(sandbox, verb, "theme", []string{"--theme", "-t"}, occurrence)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "theme", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Theme = append(entries.Theme, value)
	}
	if verb.GetOptionsSize([]string{"--description", "-d"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "description", []string{"--description", "-d"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "description", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Description = value
	} else {
		sandbox.Deps.Std.Error("required flag 'description' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		sandbox.Deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return add_doc.CommandHandler(sandbox, entries)
}

func dispatchAddFlag(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &add_flag.Entries{}
	for occurrence := 0; occurrence < verb.GetOptionsSize([]string{"--identifier", "-i"}); occurrence++ {
		raw, rawOk := optionValue(sandbox, verb, "identifier", []string{"--identifier", "-i"}, occurrence)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "identifier", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Identifier = append(entries.Identifier, value)
	}
	if verb.GetOptionsSize([]string{"--command", "-c"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "command", []string{"--command", "-c"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "command", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Command = value
	} else {
		sandbox.Deps.Std.Error("required flag 'command' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--type", "-t"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "type", []string{"--type", "-t"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "type", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Type = value
	} else {
		entries.Type = "string"
	}
	if verb.GetOptionsSize([]string{"--description", "-d"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "description", []string{"--description", "-d"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "description", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Description = value
	}
	for occurrence := 0; occurrence < verb.GetOptionsSize([]string{"--example", "-e"}); occurrence++ {
		raw, rawOk := optionValue(sandbox, verb, "example", []string{"--example", "-e"}, occurrence)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "example", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Example = append(entries.Example, value)
	}
	if verb.GetOptionsSize([]string{"--default"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "default", []string{"--default"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "default", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Default = value
	}
	entries.Required = verb.IsPresent([]string{"--required", "-r"})
	entries.Array = verb.IsPresent([]string{"--array"})
	if verb.GetOptionsSize([]string{"--min"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "min", []string{"--min"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "min", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Min = value
	}
	if verb.GetOptionsSize([]string{"--max"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "max", []string{"--max"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "max", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Max = value
	}
	if verb.GetOptionsSize([]string{"--position"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "position", []string{"--position"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseIntValue(sandbox, "flag", "position", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Position = value
	} else {
		entries.Position = -1
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		sandbox.Deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return add_flag.CommandHandler(sandbox, entries)
}

func dispatchAddHeader(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &add_header.Entries{}
	if verb.GetOptionsSize([]string{"--route"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "route", []string{"--route"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "route", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Route = value
	} else {
		sandbox.Deps.Std.Error("required flag 'route' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--type"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "type", []string{"--type"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "type", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Type = value
	} else {
		entries.Type = "string"
	}
	if verb.GetOptionsSize([]string{"--description"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "description", []string{"--description"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "description", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Description = value
	}
	for occurrence := 0; occurrence < verb.GetOptionsSize([]string{"--example"}); occurrence++ {
		raw, rawOk := optionValue(sandbox, verb, "example", []string{"--example"}, occurrence)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "example", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Example = append(entries.Example, value)
	}
	if verb.GetOptionsSize([]string{"--default"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "default", []string{"--default"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "default", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Default = value
	}
	entries.Required = verb.IsPresent([]string{"--required"})
	if verb.GetOptionsSize([]string{"--min"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "min", []string{"--min"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "min", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Min = value
	}
	if verb.GetOptionsSize([]string{"--max"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "max", []string{"--max"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "max", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Max = value
	}
	if verb.GetOptionsSize([]string{"--position"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "position", []string{"--position"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseIntValue(sandbox, "flag", "position", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Position = value
	} else {
		entries.Position = -1
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		sandbox.Deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return add_header.CommandHandler(sandbox, entries)
}

func dispatchAddLibExample(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &add_lib_example.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		sandbox.Deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return add_lib_example.CommandHandler(sandbox, entries)
}

func dispatchAddPage(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &add_page.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if verb.GetOptionsSize([]string{"--trigger"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "trigger", []string{"--trigger"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "trigger", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Trigger = value
	}
	if verb.GetOptionsSize([]string{"--title"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "title", []string{"--title"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "title", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Title = value
	}
	if verb.GetOptionsSize([]string{"--help"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "help", []string{"--help"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "help", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Help = value
	}
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		sandbox.Deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return add_page.CommandHandler(sandbox, entries)
}

func dispatchAddParam(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &add_param.Entries{}
	if verb.GetOptionsSize([]string{"--route"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "route", []string{"--route"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "route", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Route = value
	} else {
		sandbox.Deps.Std.Error("required flag 'route' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--type"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "type", []string{"--type"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "type", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Type = value
	} else {
		entries.Type = "string"
	}
	if verb.GetOptionsSize([]string{"--description"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "description", []string{"--description"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "description", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Description = value
	}
	for occurrence := 0; occurrence < verb.GetOptionsSize([]string{"--example"}); occurrence++ {
		raw, rawOk := optionValue(sandbox, verb, "example", []string{"--example"}, occurrence)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "example", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Example = append(entries.Example, value)
	}
	if verb.GetOptionsSize([]string{"--default"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "default", []string{"--default"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "default", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Default = value
	}
	entries.Required = verb.IsPresent([]string{"--required"})
	entries.Array = verb.IsPresent([]string{"--array"})
	if verb.GetOptionsSize([]string{"--min"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "min", []string{"--min"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "min", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Min = value
	}
	if verb.GetOptionsSize([]string{"--max"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "max", []string{"--max"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "max", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Max = value
	}
	if verb.GetOptionsSize([]string{"--position"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "position", []string{"--position"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseIntValue(sandbox, "flag", "position", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Position = value
	} else {
		entries.Position = -1
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		sandbox.Deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return add_param.CommandHandler(sandbox, entries)
}

func dispatchAddRoute(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &add_route.Entries{}
	if verb.GetOptionsSize([]string{"--trigger"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "trigger", []string{"--trigger"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "trigger", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Trigger = value
	}
	if verb.GetOptionsSize([]string{"--method", "-m"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "method", []string{"--method", "-m"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "method", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Method = value
	} else {
		entries.Method = "GET"
	}
	if verb.GetOptionsSize([]string{"--help"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "help", []string{"--help"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "help", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Help = value
	} else {
		sandbox.Deps.Std.Error("required flag 'help' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--category"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "category", []string{"--category"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "category", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Category = value
	} else {
		sandbox.Deps.Std.Error("required flag 'category' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		sandbox.Deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return add_route.CommandHandler(sandbox, entries)
}

func dispatchAddSegment(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &add_segment.Entries{}
	if verb.GetOptionsSize([]string{"--route"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "route", []string{"--route"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "route", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Route = value
	} else {
		sandbox.Deps.Std.Error("required flag 'route' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--identifier"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "identifier", []string{"--identifier"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "identifier", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Identifier = value
	}
	if verb.GetOptionsSize([]string{"--type"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "type", []string{"--type"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "type", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Type = value
	} else {
		entries.Type = "string"
	}
	if verb.GetOptionsSize([]string{"--description"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "description", []string{"--description"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "description", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Description = value
	}
	for occurrence := 0; occurrence < verb.GetOptionsSize([]string{"--example"}); occurrence++ {
		raw, rawOk := optionValue(sandbox, verb, "example", []string{"--example"}, occurrence)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "example", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Example = append(entries.Example, value)
	}
	entries.Array = verb.IsPresent([]string{"--array"})
	if verb.GetOptionsSize([]string{"--min"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "min", []string{"--min"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "min", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Min = value
	}
	if verb.GetOptionsSize([]string{"--max"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "max", []string{"--max"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "max", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Max = value
	}
	if verb.GetOptionsSize([]string{"--position"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "position", []string{"--position"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseIntValue(sandbox, "flag", "position", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Position = value
	} else {
		entries.Position = -1
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return add_segment.CommandHandler(sandbox, entries)
}

func dispatchBuild(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &build.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if verb.GetOptionsSize([]string{"--runtime"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "runtime", []string{"--runtime"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "runtime", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Runtime = value
	} else {
		entries.Runtime = "go"
	}
	entries.Unsafe = verb.IsPresent([]string{"--unsafe"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return build.CommandHandler(sandbox, entries)
}

func dispatchCliInit(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &cli_init.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return cli_init.CommandHandler(sandbox, entries)
}

func dispatchCliPurge(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &cli_purge.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return cli_purge.CommandHandler(sandbox, entries)
}

func dispatchCompile(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &compile.Entries{}
	for occurrence := 0; occurrence < verb.GetOptionsSize([]string{"--target", "-t"}); occurrence++ {
		raw, rawOk := optionValue(sandbox, verb, "target", []string{"--target", "-t"}, occurrence)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "target", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Target = append(entries.Target, value)
	}
	if len(entries.Target) == 0 {
		sandbox.Deps.Std.Error("required flag 'target' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return compile.CommandHandler(sandbox, entries)
}

func dispatchDepsInit(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &deps_init.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return deps_init.CommandHandler(sandbox, entries)
}

func dispatchDepsPurge(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &deps_purge.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return deps_purge.CommandHandler(sandbox, entries)
}

func dispatchExecTest(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &exec_test.Entries{}
	if verb.GetOptionsSize([]string{"--only"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "only", []string{"--only"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "only", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Only = value
	}
	entries.Update = verb.IsPresent([]string{"--update"})
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return exec_test.CommandHandler(sandbox, entries)
}

func dispatchFrontInit(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &front_init.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return front_init.CommandHandler(sandbox, entries)
}

func dispatchFrontPurge(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &front_purge.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return front_purge.CommandHandler(sandbox, entries)
}

func dispatchHelp(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &help.Entries{}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "command", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Command = value
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return help.CommandHandler(sandbox, entries)
}

func dispatchListAdapters(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &list_adapters.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return list_adapters.CommandHandler(sandbox, entries)
}

func dispatchListDeps(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &list_deps.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return list_deps.CommandHandler(sandbox, entries)
}

func dispatchLocalInstall(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &local_install.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return local_install.CommandHandler(sandbox, entries)
}

func dispatchPublish(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &publish.Entries{}
	if verb.GetOptionsSize([]string{"--path", "-p"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path", "-p"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	if verb.GetOptionsSize([]string{"--release-name", "-rn"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "release_name", []string{"--release-name", "-rn"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "release_name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.ReleaseName = value
	}
	entries.Draft = verb.IsPresent([]string{"--draft"})
	if verb.GetOptionsSize([]string{"--target", "-t"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "target", []string{"--target", "-t"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "target", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Target = value
	} else {
		entries.Target = "all"
	}
	if verb.GetOptionsSize([]string{"--publisher", "-pub"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "publisher", []string{"--publisher", "-pub"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "publisher", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Publisher = value
	} else {
		entries.Publisher = "gh"
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return publish.CommandHandler(sandbox, entries)
}

func dispatchRemoveAdapter(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &remove_adapter.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "adapter", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Adapter = value
	} else {
		sandbox.Deps.Std.Error("required arg 'adapter' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return remove_adapter.CommandHandler(sandbox, entries)
}

func dispatchRemoveArg(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &remove_arg.Entries{}
	if verb.GetOptionsSize([]string{"--command", "-c"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "command", []string{"--command", "-c"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "command", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Command = value
	} else {
		sandbox.Deps.Std.Error("required flag 'command' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		sandbox.Deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return remove_arg.CommandHandler(sandbox, entries)
}

func dispatchRemoveAvailable(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &remove_available.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "available", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Available = value
	} else {
		sandbox.Deps.Std.Error("required arg 'available' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return remove_available.CommandHandler(sandbox, entries)
}

func dispatchRemoveBodyField(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &remove_body_field.Entries{}
	if verb.GetOptionsSize([]string{"--route"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "route", []string{"--route"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "route", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Route = value
	} else {
		sandbox.Deps.Std.Error("required flag 'route' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		sandbox.Deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return remove_body_field.CommandHandler(sandbox, entries)
}

func dispatchRemoveCliExample(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &remove_cli_example.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		sandbox.Deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return remove_cli_example.CommandHandler(sandbox, entries)
}

func dispatchRemoveCommand(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &remove_command.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		sandbox.Deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return remove_command.CommandHandler(sandbox, entries)
}

func dispatchRemoveDep(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &remove_dep.Entries{}
	entries.WithAdapters = verb.IsPresent([]string{"--with-adapters"})
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "dep", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Dep = value
	} else {
		sandbox.Deps.Std.Error("required arg 'dep' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return remove_dep.CommandHandler(sandbox, entries)
}

func dispatchRemoveDoc(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &remove_doc.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		sandbox.Deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return remove_doc.CommandHandler(sandbox, entries)
}

func dispatchRemoveFlag(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &remove_flag.Entries{}
	if verb.GetOptionsSize([]string{"--command", "-c"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "command", []string{"--command", "-c"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "command", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Command = value
	} else {
		sandbox.Deps.Std.Error("required flag 'command' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		sandbox.Deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return remove_flag.CommandHandler(sandbox, entries)
}

func dispatchRemoveHeader(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &remove_header.Entries{}
	if verb.GetOptionsSize([]string{"--route"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "route", []string{"--route"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "route", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Route = value
	} else {
		sandbox.Deps.Std.Error("required flag 'route' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		sandbox.Deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return remove_header.CommandHandler(sandbox, entries)
}

func dispatchRemoveLibExample(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &remove_lib_example.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		sandbox.Deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return remove_lib_example.CommandHandler(sandbox, entries)
}

func dispatchRemovePage(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &remove_page.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		sandbox.Deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return remove_page.CommandHandler(sandbox, entries)
}

func dispatchRemoveParam(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &remove_param.Entries{}
	if verb.GetOptionsSize([]string{"--route"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "route", []string{"--route"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "route", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Route = value
	} else {
		sandbox.Deps.Std.Error("required flag 'route' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		sandbox.Deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return remove_param.CommandHandler(sandbox, entries)
}

func dispatchRemoveRoute(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &remove_route.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		sandbox.Deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return remove_route.CommandHandler(sandbox, entries)
}

func dispatchRemoveSegment(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &remove_segment.Entries{}
	if verb.GetOptionsSize([]string{"--route"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "route", []string{"--route"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "route", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Route = value
	} else {
		sandbox.Deps.Std.Error("required flag 'route' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		sandbox.Deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return remove_segment.CommandHandler(sandbox, entries)
}

func dispatchServerInit(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &server_init.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return server_init.CommandHandler(sandbox, entries)
}

func dispatchServerPurge(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &server_purge.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return server_purge.CommandHandler(sandbox, entries)
}

func dispatchSetAdapter(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &set_adapter.Entries{}
	if verb.GetOptionsSize([]string{"--available"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "available", []string{"--available"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "available", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Available = value
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "dep", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Dep = value
	} else {
		sandbox.Deps.Std.Error("required arg 'dep' not provided\n")
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "adapter", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Adapter = value
	} else {
		sandbox.Deps.Std.Error("required arg 'adapter' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return set_adapter.CommandHandler(sandbox, entries)
}

func dispatchSetBody(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &set_body.Entries{}
	if verb.GetOptionsSize([]string{"--type"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "type", []string{"--type"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "type", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Type = value
	}
	entries.Required = verb.IsPresent([]string{"--required"})
	entries.Optional = verb.IsPresent([]string{"--optional"})
	if verb.GetOptionsSize([]string{"--max-bytes"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "max-bytes", []string{"--max-bytes"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseIntValue(sandbox, "flag", "max-bytes", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.MaxBytes = value
	} else {
		entries.MaxBytes = -1
	}
	if verb.GetOptionsSize([]string{"--content-type"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "content-type", []string{"--content-type"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "content-type", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.ContentType = value
	}
	entries.DropSchema = verb.IsPresent([]string{"--drop-schema"})
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "route", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Route = value
	} else {
		sandbox.Deps.Std.Error("required arg 'route' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return set_body.CommandHandler(sandbox, entries)
}

func dispatchSetCommand(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &set_command.Entries{}
	if verb.GetOptionsSize([]string{"--help"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "help", []string{"--help"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "help", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Help = value
	}
	if verb.GetOptionsSize([]string{"--category"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "category", []string{"--category"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "category", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Category = value
	}
	if verb.GetOptionsSize([]string{"--long-description"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "long-description", []string{"--long-description"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "long-description", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.LongDescription = value
	}
	for occurrence := 0; occurrence < verb.GetOptionsSize([]string{"--identifier", "-i"}); occurrence++ {
		raw, rawOk := optionValue(sandbox, verb, "identifier", []string{"--identifier", "-i"}, occurrence)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "identifier", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Identifier = append(entries.Identifier, value)
	}
	for occurrence := 0; occurrence < verb.GetOptionsSize([]string{"--example", "-e"}); occurrence++ {
		raw, rawOk := optionValue(sandbox, verb, "example", []string{"--example", "-e"}, occurrence)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "example", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Example = append(entries.Example, value)
	}
	entries.Hidden = verb.IsPresent([]string{"--hidden"})
	entries.Visible = verb.IsPresent([]string{"--visible"})
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		sandbox.Deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return set_command.CommandHandler(sandbox, entries)
}

func dispatchSetDep(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &set_dep.Entries{}
	if verb.GetOptionsSize([]string{"--version"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "version", []string{"--version"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "version", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Version = value
	} else {
		sandbox.Deps.Std.Error("required flag 'version' not provided\n")
		return ExitUsage
	}
	if verb.GetOptionsSize([]string{"--remote-available"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "remote-available", []string{"--remote-available"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "remote-available", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.RemoteAvailable = value
	}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "dep", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Dep = value
	} else {
		sandbox.Deps.Std.Error("required arg 'dep' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return set_dep.CommandHandler(sandbox, entries)
}

func dispatchSetRoute(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &set_route.Entries{}
	if verb.GetOptionsSize([]string{"--method", "-m"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "method", []string{"--method", "-m"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "method", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Method = value
	}
	if verb.GetOptionsSize([]string{"--help"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "help", []string{"--help"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "help", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Help = value
	}
	if verb.GetOptionsSize([]string{"--category"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "category", []string{"--category"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "category", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Category = value
	}
	if verb.GetOptionsSize([]string{"--long-description"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "long-description", []string{"--long-description"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "long-description", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.LongDescription = value
	}
	entries.Hidden = verb.IsPresent([]string{"--hidden"})
	entries.Visible = verb.IsPresent([]string{"--visible"})
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	for occurrence := 0; occurrence < verb.GetOptionsSize([]string{"--example"}); occurrence++ {
		raw, rawOk := optionValue(sandbox, verb, "example", []string{"--example"}, occurrence)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "example", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Example = append(entries.Example, value)
	}
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "route", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Route = value
	} else {
		sandbox.Deps.Std.Error("required arg 'route' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return set_route.CommandHandler(sandbox, entries)
}

func dispatchStart(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &start.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	if verb.GetOptionsSize([]string{"--project-name", "-p"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "project-name", []string{"--project-name", "-p"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "project-name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.ProjectName = value
	} else {
		sandbox.Deps.Std.Error("required flag 'project-name' not provided\n")
		return ExitUsage
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	entries.Force = verb.IsPresent([]string{"--force", "-f"})
	if verb.GetOptionsSize([]string{"--module", "-m"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "module", []string{"--module", "-m"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "module", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Module = value
	}
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return start.CommandHandler(sandbox, entries)
}

func dispatchUpdateTest(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &update_test.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseStringValue(sandbox, "arg", "name", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Name = value
	} else {
		sandbox.Deps.Std.Error("required arg 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return update_test.CommandHandler(sandbox, entries)
}

func dispatchVerify(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &verify.Entries{}
	if verb.GetOptionsSize([]string{"--path"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "path", []string{"--path"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "path", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Path = value
	} else {
		entries.Path = "."
	}
	if verb.GetOptionsSize([]string{"--runtime"}) > 0 {
		raw, rawOk := optionValue(sandbox, verb, "runtime", []string{"--runtime"}, 0)
		if !rawOk {
			return ExitUsage
		}
		value, valueOk := parseStringValue(sandbox, "flag", "runtime", raw)
		if !valueOk {
			return ExitUsage
		}
		entries.Runtime = value
	} else {
		entries.Runtime = "go"
	}
	entries.Quiet = verb.IsPresent([]string{"--quiet", "-q"})
	if entries.Quiet {
		silenceLogs(sandbox)
	}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return verify.CommandHandler(sandbox, entries)
}

func dispatchVersion(sandbox *api.Sandbox, verb argvdeps.Parser) int {
	entries := &version.Entries{}
	if !checkUnknownFlags(sandbox, verb) {
		return ExitUsage
	}
	if !checkUnusedArgs(sandbox, verb) {
		return ExitUsage
	}
	return version.CommandHandler(sandbox, entries)
}
