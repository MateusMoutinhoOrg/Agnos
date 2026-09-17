# `sandbox/api/command.go`

## `CommandArg`

CommandArg is one positional argument a command declares — the parsed form of one entry under `args:` in that command's entries.yaml. It is matched by position, so Id names it for GetItem and for the help screen, never on the command line.

| Field | Type | Description |
| --- | --- | --- |
| `Type` | `string` | Type is the declared type: "string", "boolean", "int" or "float". |
| `Id` | `string` | Id is the name the argument is declared and read back under. |
| `Required` | `bool` | Required reports that the command line is rejected without it. |
| `Array` | `bool` | Array reports that it takes every argument left on the line. |
| `Description` | `string` | Description is the one-line help text. |
| `Examples` | `[]string` | Examples are whole command lines the help screen prints. |
| `Default` | `string` | Default is the value bound when the argument is absent, spelled as it is written in entries.yaml. |
| `HasDefault` | `bool` | HasDefault tells a declared empty default from no default at all. |
| `Min` | `float64` | Min and Max bound a numeric value; HasMin and HasMax tell a bound of zero from no bound. |
| `HasMin` | `bool` |  |
| `Max` | `float64` |  |
| `HasMax` | `bool` |  |

## `CommandFlag`

CommandFlag is one flag a command declares — the parsed form of one entry under `flags:` in that command's entries.yaml. It is matched by Identifiers and read back under Id.

| Field | Type | Description |
| --- | --- | --- |
| `Type` | `string` | Type is the declared type: "string", "boolean", "int" or "float". |
| `Id` | `string` | Id is the name the flag is declared and read back under. |
| `Required` | `bool` | Required reports that the command line is rejected without it. |
| `Array` | `bool` | Array reports that every occurrence is kept, not just the first. |
| `Description` | `string` | Description is the one-line help text. |
| `Examples` | `[]string` | Examples are whole command lines the help screen prints. |
| `Default` | `string` | Default is the value bound when the flag is absent, spelled as it is written in entries.yaml. |
| `HasDefault` | `bool` | HasDefault tells a declared empty default from no default at all. |
| `Identifiers` | `[]string` | Identifiers are the spellings the flag answers to ("--path", "-p"). |
| `Min` | `float64` | Min and Max bound a numeric value; HasMin and HasMax tell a bound of zero from no bound. |
| `HasMin` | `bool` |  |
| `Max` | `float64` |  |
| `HasMax` | `bool` |  |

## `Command`

Command is one command of the project, as the sandbox offers it: the whole of what its entries.yaml declares, plus the handler behind it. Cli.Commands holds one per sandbox/internal/commands/<name>/, each built by that package's generated NewCommand. What Cli.Commands holds is the declaration alone: nothing is ever bound onto it. The cli dispatch copies it with BindCommand, fills that copy's Items from the command line, and hands it to Handler; a caller holding the sandbox can do the same.

| Field | Type | Description |
| --- | --- | --- |
| `Name` | `string` | Name is the package directory of the command, snake_case. |
| `Identifiers` | `[]string` | Identifiers are the verbs it answers to ("add-flag"), the first of which is the one the help screen prints. |
| `Category` | `string` | Category groups it on the general help screen. |
| `Help` | `string` | Help is the one-line description. |
| `LongDescription` | `string` | LongDescription is the paragraph the per-command help screen prints. |
| `Examples` | `[]string` | Examples are whole command lines the help screen prints. |
| `Hidden` | `bool` | Hidden keeps it off the general help screen without disabling it. |
| `Args` | `[]CommandArg` | Args are the positional arguments, in declaration order. |
| `Flags` | `[]CommandFlag` | Flags are the flags, in declaration order. |
| `Items` | `map[string][]any` | Items holds the values bound to the declaration: one entry per flag or arg Id, in declaration order for an array, one element long for a scalar, absent when nothing was bound. The Get* fields read it. |
| `GetItem` | `func(id string) []any` | GetItem returns every value bound under one Id, nil when none was. |
| `GetString` | `func(id string) string` | GetString returns the first string bound under one Id, "" when none was. |
| `GetBool` | `func(id string) bool` | GetBool returns the first boolean bound under one Id, false when none was. |
| `GetInt` | `func(id string) int` | GetInt returns the first int bound under one Id, 0 when none was. |
| `GetFloat` | `func(id string) float64` | GetFloat returns the first float bound under one Id, 0 when none was. |
| `GetStrings` | `func(id string) []string` | GetStrings returns every string bound under one Id, in order. |
| `GetInts` | `func(id string) []int` | GetInts returns every int bound under one Id, in order. |
| `GetFloats` | `func(id string) []float64` | GetFloats returns every float bound under one Id, in order. |
| `Handler` | `func(command *Command) int` | Handler runs the command against one bound copy — the values in its Items — and returns the exit code it answered with. It takes that copy rather than closing over one, so the declaration on Cli.Commands is shared by every caller while nothing bound ever is. It is the command package's own CommandHandler, closed over the sandbox. |

| Function | Description |
| --- | --- |
| `NewCommand() *Command` | NewCommand returns an empty Command with Items open and every Get* reader bound to it. A generated NewCommand fills the declaration and the handler on top of what this returns, so every command reads its values the same way. |
| `BindCommand(command *Command) *Command` | BindCommand copies one declaration into the command a single run binds to: the same declared fields — the slices are read-only and shared — with Items empty, the readers pointed at the copy and the handler carried over. The dispatch calls it once per command line, so what Cli.Commands holds is never written to. |

[every contract](doc.md)
