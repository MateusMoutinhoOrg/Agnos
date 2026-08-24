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
			for j := range command.FlagsList {
				flag := &command.FlagsList[j]

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
			for j := range command.ArgsList {
				arg := &command.ArgsList[j]
				if err := collectArg(arg, verb); err != nil {
					sandbox.Deps.Printf("%s\n", err.Error())
					return api.ExitUsage
				}
			}

			flagsRetriver := buildFlagsRetriver(command)
			argsRetriver := buildArgsRetriver(command)

			return command.Handler(sandbox.Deps, &argsRetriver, &flagsRetriver)
		}

		sandbox.Deps.Printf("Unknown Command!\n")
		return api.ExitUsage
	}
}

// printUsage shows all commands and their descriptions.
func printUsage(sandbox *api.SandBox) {
	sandbox.Deps.Printf("Usage: %s <command> [flags] [args]\n\n", sandbox.ProjectName)
	sandbox.Deps.Printf("Commands:\n")
	for _, cmd := range sandbox.Commands {
		if len(cmd.ValidStartIdentifiers) > 0 {
			sandbox.Deps.Printf("  %-20s %s\n", cmd.ValidStartIdentifiers[0], cmd.Description)
		}
	}
}

// collectBoolFlag checks whether a boolean flag is present. If required and
// absent, it returns an error. The flag's Values slice is filled with one
// cliValue holding the result.
func collectBoolFlag(flag *api.Cliflag, verb verbdeps.Lib) error {
	present := verb.IsPresent(flag.ValidIdentifiers)
	if flag.Required && !present {
		return fmt.Errorf("required flag '%s' not provided", flag.Name)
	}
	flag.Values = []api.CliValue{boolValue(present)}
	return nil
}

// collectValueFlag reads the occurrences of a non-bool flag from the argument
// vector. It validates that the number of provided values falls within
// [MinSize, MaxSize] and that required flags have at least one value.
func collectValueFlag(flag *api.Cliflag, verb verbdeps.Lib) error {
	size := verb.GetOptionsSize(flag.ValidIdentifiers)

	if flag.Required && size == 0 {
		return fmt.Errorf("required flag '%s' not provided", flag.Name)
	}

	if size == 0 {
		return nil
	}

	if flag.MinSize > 0 && size < flag.MinSize {
		return fmt.Errorf("flag '%s' requires at least %d value(s), got %d", flag.Name, flag.MinSize, size)
	}

	maxSize := flag.MaxSize
	if maxSize <= 0 {
		maxSize = size
	}
	if size > maxSize {
		return fmt.Errorf("flag '%s' accepts at most %d value(s), got %d", flag.Name, maxSize, size)
	}

	flag.Values = make([]api.CliValue, 0, size)
	for i := 0; i < size; i++ {
		val, err := readFlagValue(flag, verb, i)
		if err != nil {
			return fmt.Errorf("flag '%s': %w", flag.Name, err)
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
	size := arg.Size
	if size <= 0 {
		size = 1
	}

	arg.Values = make([]api.CliValue, 0, size)
	for i := 0; i < size; i++ {
		val, err := readArgValue(arg, verb)
		if err != nil {
			if arg.Required {
				return fmt.Errorf("required arg '%s': %w", arg.Name, err)
			}
			break
		}
		arg.Values = append(arg.Values, val)
	}

	if arg.Required && len(arg.Values) == 0 {
		return fmt.Errorf("required arg '%s' not provided", arg.Name)
	}
	return nil
}

// readArgValue reads one positional value from the next unused argv slot.
func readArgValue(arg *api.CliArg, verb verbdeps.Lib) (api.CliValue, error) {
	switch arg.Type {
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

// buildFlagsRetriver builds a FlagsRetriver from the already-parsed flag
// values stored in each Cliflag's Values slice.
func buildFlagsRetriver(command *api.CliCommand) api.FlagsRetriver {
	flagsByName := make(map[string]*api.Cliflag, len(command.FlagsList))
	for i := range command.FlagsList {
		flagsByName[command.FlagsList[i].Name] = &command.FlagsList[i]
	}

	return api.FlagsRetriver{
		GetStringFlag: func(name string, index int) string {
			f := flagsByName[name]
			if f == nil || index >= len(f.Values) {
				return ""
			}
			return f.Values[index].String()
		},
		GetIntFlag: func(name string, index int) int {
			f := flagsByName[name]
			if f == nil || index >= len(f.Values) {
				return 0
			}
			return f.Values[index].Int()
		},
		GetFloatFlag: func(name string, index int) float64 {
			f := flagsByName[name]
			if f == nil || index >= len(f.Values) {
				return 0
			}
			return f.Values[index].Float()
		},
		GetBoolFlag: func(name string, index int) bool {
			f := flagsByName[name]
			if f == nil || index >= len(f.Values) {
				return false
			}
			return f.Values[index].Bool()
		},
	}
}

// buildArgsRetriver builds an ArgsRetriver from the already-parsed arg values
// stored in each CliArg's Values slice.
func buildArgsRetriver(command *api.CliCommand) api.ArgsRetriver {
	argsByName := make(map[string]*api.CliArg, len(command.ArgsList))
	for i := range command.ArgsList {
		argsByName[command.ArgsList[i].Name] = &command.ArgsList[i]
	}

	return api.ArgsRetriver{
		GetStringArg: func(name string, index int) string {
			a := argsByName[name]
			if a == nil || index >= len(a.Values) {
				return ""
			}
			return a.Values[index].String()
		},
		GetIntArg: func(name string, index int) int {
			a := argsByName[name]
			if a == nil || index >= len(a.Values) {
				return 0
			}
			return a.Values[index].Int()
		},
		GetFloatArg: func(name string, index int) float64 {
			a := argsByName[name]
			if a == nil || index >= len(a.Values) {
				return 0
			}
			return a.Values[index].Float()
		},
		GetBoolArg: func(name string, index int) bool {
			a := argsByName[name]
			if a == nil || index >= len(a.Values) {
				return false
			}
			return a.Values[index].Bool()
		},
	}
}
