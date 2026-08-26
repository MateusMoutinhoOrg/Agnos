package cli

import (
	"fmt"
	"slices"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps/verbdeps"
)

// CliMainFactory fills api.CliApi.CliMain with the dispatch-and-parse loop:
// match a command, collect its declared flags and args from the argument
// vector, validate required fields, and hand the parsed retrivers to the
// command handler.
func CliMainFactory(sandbox *api.SandBox) func(args []string) int {

	return func(args []string) int {
		if len(args) == 0 {
			printUsage(sandbox)
			return api.ExitUsage
		}

		verb := sandbox.Deps.VerbLib

		action, err := verb.GetNextStringArg()
		if err != nil {
			printUsage(sandbox)
			return api.ExitUsage
		}

		for i := range sandbox.Commands {
			command := &sandbox.Commands[i]

			if !slices.Contains(command.ValidStartIdentifiers, action) {
				continue
			}

			// ── Collect flags ──────────────────────────────────────
			for j := range command.Flags {
				flag := &command.Flags[j]

				if flag.Type == api.CliTypeBool {
					if err := collectBoolFlag(flag, verb); err != nil {
						sandbox.Deps.Printf("%s\n", err.Error())
						return api.ExitUsage
					}
					continue
				}

				if err := collectValueFlag(flag, verb); err != nil {
					sandbox.Deps.Printf("%s\n", err.Error())
					return api.ExitUsage
				}
			}

			// ── Collect args (positional, after flags are consumed) ─
			for j := range command.Args {
				arg := &command.Args[j]
				if err := collectArg(arg, verb); err != nil {
					sandbox.Deps.Printf("%s\n", err.Error())
					return api.ExitUsage
				}
			}

			entries := buildCliEntrys(command)

			return command.Handler(sandbox, entries)
		}

		sandbox.Deps.Printf("Unknown Command!\n")
		return api.ExitUsage
	}
}

// printUsage triggers the help command so the user sees the full
// professional help screen when they run the binary with no arguments.
func printUsage(sandbox *api.SandBox) {
	for _, cmd := range sandbox.Commands {
		if slices.Contains(cmd.ValidStartIdentifiers, "help") {
			cmd.Handler(sandbox, buildCliEntrys(&cmd))
			return
		}
	}
	// Fallback — should never happen if help is registered.
	sandbox.Deps.Printf("Usage: %s <command> [flags] [args]\n", sandbox.ProjectName)
}

// collectBoolFlag checks whether a boolean flag is present. If required and
// absent, it returns an error. The flag's Values slice is filled with one
// cliValue holding the result.
func collectBoolFlag(flag *api.Cliflag, verb verbdeps.Lib) error {
	present := verb.IsPresent(flag.ValidIdentifiers)
	if flag.RequiredPresence && !present {
		return fmt.Errorf("required flag '%s' not provided", flag.Id)
	}
	flag.Values = []api.CliValue{boolValue(present)}
	flag.Exist = present
	return nil
}

// collectValueFlag reads the occurrences of a non-bool flag from the argument
// vector. It validates that the number of provided values falls within
// [RequiredMinSize, RequiredMaxSize] and that required flags have at least one value.
func collectValueFlag(flag *api.Cliflag, verb verbdeps.Lib) error {
	size := verb.GetOptionsSize(flag.ValidIdentifiers)

	if flag.RequiredPresence && size == 0 {
		return fmt.Errorf("required flag '%s' not provided", flag.Id)
	}

	if size == 0 {
		return nil
	}

	flag.Exist = true

	if flag.RequiredMinSize > 0 && size < flag.RequiredMinSize {
		return fmt.Errorf("flag '%s' requires at least %d value(s), got %d", flag.Id, flag.RequiredMinSize, size)
	}

	maxSize := flag.RequiredMaxSize
	if maxSize <= 0 {
		maxSize = size
	}
	if size > maxSize {
		return fmt.Errorf("flag '%s' accepts at most %d value(s), got %d", flag.Id, maxSize, size)
	}

	flag.Values = make([]api.CliValue, 0, size)
	for i := 0; i < size; i++ {
		val, err := readFlagValue(flag, verb, i)
		if err != nil {
			return fmt.Errorf("flag '%s': %w", flag.Id, err)
		}
		flag.Values = append(flag.Values, val)
	}
	return nil
}

// readFlagValue reads one flag occurrence and returns it as a cliValue of the
// appropriate type.
func readFlagValue(flag *api.Cliflag, verb verbdeps.Lib, occurrence int) (api.CliValue, error) {
	switch flag.Type {
	case api.CliTypeInt:
		v, err := verb.GetIntOption(flag.ValidIdentifiers, occurrence)
		if err != nil {
			return nil, err
		}
		return intValue(v), nil
	case api.CliTypeFloat:
		v, err := verb.GetDoubleOption(flag.ValidIdentifiers, occurrence)
		if err != nil {
			return nil, err
		}
		return floatValue(v), nil
	default: // CliTypeString
		v, err := verb.GetStringOption(flag.ValidIdentifiers, occurrence)
		if err != nil {
			return nil, err
		}
		return stringValue(v), nil
	}
}

// collectArg reads a positional arg from the unused portion of the argument
// vector (via GetNext*Arg). Required args that cannot be read produce an error.
func collectArg(arg *api.CliArg, verb verbdeps.Lib) error {
	minSize := arg.RequiredMinSize
	maxSize := arg.RequiredMaxSize
	if maxSize <= 0 {
		if minSize > 0 {
			maxSize = minSize
		} else {
			maxSize = 1
		}
	}

	arg.Values = make([]api.CliValue, 0, maxSize)
	for i := 0; i < maxSize; i++ {
		val, err := readArgValue(arg, verb)
		if err != nil {
			if i < minSize {
				return fmt.Errorf("required arg '%s': requires at least %d values, got %d", arg.Id, minSize, i)
			}
			break
		}
		arg.Values = append(arg.Values, val)
	}

	if len(arg.Values) == 0 && len(arg.Defaults) > 0 {
		for _, d := range arg.Defaults {
			arg.Values = append(arg.Values, stringValue(d))
		}
	}

	return nil
}

// readArgValue reads one positional value from the next unused argv slot.
func readArgValue(arg *api.CliArg, verb verbdeps.Lib) (api.CliValue, error) {
	switch arg.RequiredType {
	case api.CliTypeInt:
		v, err := verb.GetNextIntArg()
		if err != nil {
			return nil, err
		}
		return intValue(v), nil
	case api.CliTypeFloat:
		v, err := verb.GetNextDoubleArg()
		if err != nil {
			return nil, err
		}
		return floatValue(v), nil
	case api.CliTypeBool:
		v, err := verb.GetNextStringArg()
		if err != nil {
			return nil, err
		}
		return boolValue(v == "true" || v == "1" || v == "yes"), nil
	default: // CliTypeString
		v, err := verb.GetNextStringArg()
		if err != nil {
			return nil, err
		}
		return stringValue(v), nil
	}
}

// buildCliEntrys builds a CliEntrys from the already-parsed flag and arg
// values stored in each Cliflag/CliArg's Values slice.
func buildCliEntrys(command *api.CliCommand) api.CliEntrys {
	flagsById := make(map[string]*api.Cliflag, len(command.Flags))
	for i := range command.Flags {
		flagsById[command.Flags[i].Id] = &command.Flags[i]
	}

	argsById := make(map[string]*api.CliArg, len(command.Args))
	for i := range command.Args {
		argsById[command.Args[i].Id] = &command.Args[i]
	}

	return api.CliEntrys{
		GetFlagById: func(id string) *api.Cliflag {
			return flagsById[id]
		},
		GetArgById: func(id string) *api.CliArg {
			return argsById[id]
		},
	}
}
