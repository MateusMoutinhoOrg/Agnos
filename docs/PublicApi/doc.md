# PublicApi

Every exported symbol of `github.com/MateusMoutinhoOrg/Agnos`, read straight from the contract sources on
every build: `sandbox/api/` is the surface `sandbox.New` returns, `sandbox/deps/`
the contracts an adapter fills and a caller may replace. Each description below is the
doc comment of the declaration itself — change the comment, run `build`, and this page
follows.

## Entry points

| Symbol | Signature |
| --- | --- |
| `sandbox.New` | `func(deps *deps.Deps) *api.Sandbox` |
| `standard.New` | `func() deps.Deps` (`adapters/availables/standard`) |

Implementations live under `sandbox/internal` and are unreachable: every contract is a
struct of function fields, filled by a binder.

# The sandbox api

## `sandbox/api/sandbox.go`

### `Sandbox`

Sandbox is the whole library: one field per contract declared in sandbox/api/, each filled by its binder. sandbox.New returns it, and nothing callable lives outside of it.

| Field | Type |
| --- | --- |
| `Actions` | `Actions` |
| `Cli` | `Cli` |

## `sandbox/api/actions.go`

| Constant | Value | Description |
| --- | --- | --- |
| `RuntimeGo` | `"go"` | RuntimeGo resolves the module graph and compiles every package after the render. |
| `RuntimeNone` | `"none"` | RuntimeNone renders only, leaving the result unchecked. |

### `BuildProps`

BuildProps describes one (re)render of a project: the directory holding it and the runtime that then checks the result.

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Runtime` | `string` |

### `CompileProps`

CompileProps describes one cross-compile run: the directory holding the project and the target names to build. Each name is one of the keys `agnos compile` accepts (linux86, linuxarm64, linuxi32, mac86, macarm64, windows86, windowsi32) or "all", which expands to every target.

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Targets` | `[]string` |

### `StartProps`

StartProps describes one project to scaffold: the directory to write it into, the name it carries in <Name>Config/project.yaml, the module path for go.mod (nil derives it from the name) and whether an existing directory may be written over.

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `ProjectName` | `string` |
| `Module` | `*string` |
| `Force` | `bool` |

### `FieldProps`

FieldProps describes one flag or positional arg to add to a command's entries.yaml. Default, Min and Max are the raw literals typed on the command line ("" means unset) so the action can tell "not given" from a zero value; Position is the index to insert at (< 0 appends).

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Command` | `string` |
| `Name` | `string` |
| `Identifiers` | `[]string` |
| `Description` | `string` |
| `Examples` | `[]string` |
| `Type` | `string` |
| `Default` | `string` |
| `Required` | `bool` |
| `Array` | `bool` |
| `Min` | `string` |
| `Max` | `string` |
| `Position` | `int` |

### `CommandProps`

CommandProps carries the command-level keys of entries.yaml that set-command may rewrite. Empty strings leave the current value alone; Identifiers / Examples are appended (deduplicated).

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Command` | `string` |
| `Help` | `string` |
| `Category` | `string` |
| `LongDescription` | `string` |
| `Hidden` | `bool` |
| `Visible` | `bool` |
| `Identifiers` | `[]string` |
| `Examples` | `[]string` |

### `DocProps`

DocProps describes one doc to create under docs/. Name is the doc's directory, optionally nested under its parent ("PublicApi/api.Actions"). Themes are the theme ids of <ProjectName>Config/themes.yaml the doc belongs to: required on a first-level doc, forbidden on a sub-doc.

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Name` | `string` |
| `Description` | `string` |
| `Themes` | `[]string` |

### `Actions`

Actions is the whole set of operations agnos performs on a project. Every field takes the project directory as its first input (`path`, or the Path of a props struct) and reports failure as an error; the ones that change the tree re-render it before returning, so a project is always left in a built state.

| Field | Type | Description |
| --- | --- | --- |
| `Build` | `func(props BuildProps) error` | Build re-renders every generated file of the project and hands the result to the runtime named by the props. |
| `Compile` | `func(props CompileProps) error` | Compile cross-compiles the project's cmd/ binaries into release/, one file per named target. |
| `Verify` | `func(path string) error` | Verify checks the project against the schema every generator assumes and writes nothing; it reports every violation at once. |
| `Start` | `func(props StartProps) error` | Start scaffolds a new project: the config directory, go.mod, the sandbox skeleton and a first build. |
| `DepsInit` | `func(path string) error` | DepsInit adds the dependency layer (sandbox/deps/ and adapters/availables/standard/) to a project that has none. |
| `DepsPurge` | `func(path string) error` | DepsPurge removes the dependency layer and every installed dep with it. |
| `DepInstall` | `func(path string, dep string) error` | DepInstall installs one dep of the built-in list: its contract under sandbox/deps/, its adapter under adapters/libs/ and its go.mod require. |
| `DepRemove` | `func(path string, dep string) error` | DepRemove uninstalls one installed dep, contract, adapter and require. |
| `DepList` | `func(path string) ([]string, error)` | DepList returns the names of the deps installed in the project. |
| `CliInit` | `func(path string) error` | CliInit adds the CLI layer (cmd/main, the dispatcher and the help and version commands) to a project that has none. |
| `CliPurge` | `func(path string) error` | CliPurge removes the CLI layer and every command declared in it. |
| `AddCommand` | `func(path string, name string, help string, category string) error` | AddCommand declares a new command: its entries.yaml, its generated entries.go and a handler.go to fill in. |
| `RemoveCommand` | `func(path string, name string) error` | RemoveCommand deletes one command and unwires it from the dispatcher. |
| `SetCommand` | `func(props CommandProps) error` | SetCommand rewrites the command-level keys of one command's entries.yaml. |
| `AddFlag` | `func(props FieldProps) error` | AddFlag declares one flag on a command. |
| `RemoveFlag` | `func(path string, command string, name string) error` | RemoveFlag deletes one declared flag from a command. |
| `AddArg` | `func(props FieldProps) error` | AddArg declares one positional argument on a command. |
| `RemoveArg` | `func(path string, command string, name string) error` | RemoveArg deletes one declared positional argument from a command. |
| `AddDoc` | `func(props DocProps) error` | AddDoc creates one doc directory under docs/, with its props.yaml and a doc.md to fill in. |
| `RemoveDoc` | `func(path string, name string) error` | RemoveDoc deletes one doc directory and everything under it. |

## `sandbox/api/cli.go`

| Constant | Value | Description |
| --- | --- | --- |
| `ExitOk` | `0` | ExitOk reports that the command did what it was asked to do. |
| `ExitFailure` | `1` | ExitFailure reports that a well-formed command could not be carried out. |
| `ExitUsage` | `2` | ExitUsage reports that the command line itself was wrong — an unknown command or flag, a missing operand, an unparsable amount. Every such error exits with this one code, whichever command produced it. |

### `Cli`

Cli is the CLI surface of the sandbox. CliMain is the generated dispatch-and-parse entry point (see sandbox/internal/cli/climain.go), wired in by sandbox/binds/cli.go.

| Field | Type |
| --- | --- |
| `CliMain` | `func(args []string) int` |

# Dependency contracts

`deps.Deps` has one field per directory of `sandbox/deps/`, named by title-casing it. Each
field is that package's `Lib` struct, filled by `adapters/libs/<name>.Bind(&deps)`.

## `deps.Argvdeps`

`sandbox/deps/argvdeps`

### `Lib`

Lib is the argv-parser constructor injected whole as the Deps.ArgvLib field — the same mechanic as requestdeps.Lib. A parser is bound to one argument vector, so it is created per call rather than injected once: what the sandbox holds is this one-field struct, and the adapter — which lives outside the sandbox — fills New over a concrete argv-parser library.

| Field | Type | Description |
| --- | --- | --- |
| `New` | `func(args []string) Parser` | New builds an argv parser bound to the given arguments. |

### `Parser`

Parser mirrors the concrete argv-parser library's api.Lib — an argument-vector (argv) parser. Every argument starts out unread; calling any Get* field or IsPresent marks the argument(s) it matched as used, so whatever is left over in Args is exactly the positional arguments nothing asked for. The two *Size fields are the exception: they count matches without ever marking anything used. Each getter family (Option, Arg, NextArg, KeyValues) is exposed once per supported value type: String (raw text), Int (base-10), Double (float64), and Timestamp (RFC 3339). A typed getter marks its match as used even when parsing then fails.

| Field | Type | Description |
| --- | --- | --- |
| `Args` | `[]string` | Args is the argument vector being parsed. Every index-based field refers to positions in this slice. Treat it as read-only: mutating it leaves Used out of sync. |
| `Used` | `[]bool` | Used tracks, index for index against Args, which arguments have already been matched by a previous call. Treat it as read-only. |
| `IsPresent` | `func(flags []string) bool` | IsPresent reports whether any of the given flag spellings (e.g. []string{"-q", "--quiet"}) occurs in the unread portion of Args, marking the matched argument used. It never fails: "not present" is a valid outcome. |
| `GetOptionsSize` | `func(flags []string) int` | GetOptionsSize counts how many arguments equal one of the given flag spellings, regardless of Used, and never mutates Used. Pair it with GetStringOption to iterate occurrences 0..size-1. |
| `GetKeyValuesSize` | `func(prefixes []string) int` | GetKeyValuesSize counts how many arguments start with one of the given key=value prefixes (the separator is part of the prefix), regardless of Used, and never mutates Used. |
| `GetStringOption` | `func(flags []string, occurrence int) (string, error)` | GetStringOption returns the argument following the occurrence-th (0-based) match of the given flag spellings, marking both as used. It errors when occurrence is out of range or the flag has no value after it. |
| `GetIntOption` | `func(flags []string, occurrence int) (int, error)` | GetIntOption behaves like GetStringOption, additionally parsing the value as a base-10 integer. |
| `GetDoubleOption` | `func(flags []string, occurrence int) (float64, error)` | GetDoubleOption behaves like GetStringOption, additionally parsing the value as a 64-bit floating-point number. |
| `GetTimestampOption` | `func(flags []string, occurrence int) (time.Time, error)` | GetTimestampOption behaves like GetStringOption, additionally parsing the value as an RFC 3339 timestamp. |
| `GetStringArg` | `func(index int) (string, error)` | GetStringArg returns the argument at the given absolute index of Args and marks it used. It errors when index is out of range. |
| `GetIntArg` | `func(index int) (int, error)` | GetIntArg behaves like GetStringArg, additionally parsing the argument as a base-10 integer. |
| `GetDoubleArg` | `func(index int) (float64, error)` | GetDoubleArg behaves like GetStringArg, additionally parsing the argument as a 64-bit floating-point number. |
| `GetTimestampArg` | `func(index int) (time.Time, error)` | GetTimestampArg behaves like GetStringArg, additionally parsing the argument as an RFC 3339 timestamp. |
| `GetNextStringArg` | `func() (string, error)` | GetNextStringArg returns the first still-unused argument, in order, and marks it used — the leftover positional arguments, drained one call at a time. It errors when every argument has been used. |
| `GetNextIntArg` | `func() (int, error)` | GetNextIntArg behaves like GetNextStringArg, additionally parsing the argument as a base-10 integer. |
| `GetNextDoubleArg` | `func() (float64, error)` | GetNextDoubleArg behaves like GetNextStringArg, additionally parsing the argument as a 64-bit floating-point number. |
| `GetNextTimestampArg` | `func() (time.Time, error)` | GetNextTimestampArg behaves like GetNextStringArg, additionally parsing the argument as an RFC 3339 timestamp. |
| `GetStringKeyValues` | `func(prefixes []string, occurrence int) (string, error)` | GetStringKeyValues returns the text after the matched prefix of the occurrence-th (0-based) argument starting with one of the given key=value prefixes, marking it used. It errors when occurrence is out of range or the value portion is empty. |
| `GetIntKeyValues` | `func(prefixes []string, occurrence int) (int, error)` | GetIntKeyValues behaves like GetStringKeyValues, additionally parsing the value portion as a base-10 integer. |
| `GetDoubleKeyValues` | `func(prefixes []string, occurrence int) (float64, error)` | GetDoubleKeyValues behaves like GetStringKeyValues, additionally parsing the value portion as a 64-bit floating-point number. |
| `GetTimestampKeyValues` | `func(prefixes []string, occurrence int) (time.Time, error)` | GetTimestampKeyValues behaves like GetStringKeyValues, additionally parsing the value portion as an RFC 3339 timestamp. |

## `deps.Embeddeps`

`sandbox/deps/embeddeps`

### `Lib`

Lib is the embedded-asset library injected whole as the Deps.EmbedDeps field. It is read-only by design: assets ship with the program, and nothing in the library ever writes one back. Every path is slash-separated and relative to the root of the asset tree the adapter serves — "report.tmpl", "templates/invoice.tmpl" — never an absolute path and never a path reaching outside that root, so the same call means the same asset whatever the adapter is backed by.

| Field | Type | Description |
| --- | --- | --- |
| `ReadFile` | `func(path string) ([]byte, error)` | ReadFile returns the whole content of one asset. The error reports an asset that does not exist or could not be read; callers inside the sandbox report it rather than assuming the bytes are there, because a missing asset is a packaging mistake, not a user mistake. |
| `ListFiles` | `func(path string) ([]string, error)` | ListFiles returns the names of the assets directly inside the given directory, in lexical order, relative to that directory. Nested directories are not descended into and are not reported. The root itself is addressed as ".". |
| `ListFilesRecursively` | `func(path string) ([]string, error)` | ListFilesRecursively returns every asset at or below the given directory, in lexical order, as slash-separated paths relative to that directory — "templates/invoice.tmpl" and not just "invoice.tmpl". Directories are never reported, only the files inside them. |
| `RenderTemplate` | `func(path string, vars interface{}) ([]byte, error)` | RenderTemplate reads the template from path, renders it using the given variables, and returns the resulting byte slice. |

## `deps.Goimportsdeps`

`sandbox/deps/goimportsdeps`

### `Lib`

Lib is the Go-source parser injected whole as the Deps.Goimportsdeps field. Every field errors when the given content is not parsable Go.

| Field | Type | Description |
| --- | --- | --- |
| `Parse` | `func(content string) (*File, error)` | Parse parses a whole Go source file into a File describing its package clause, imports and top-level declarations (functions, methods, types, constants and variables). |
| `GetPackageName` | `func(content string) (string, error)` | GetPackageName returns the name in the file's package clause. |
| `GetImports` | `func(content string) ([]string, error)` | GetImports returns every import path the file declares, in source order. |
| `Format` | `func(content string) (string, error)` | Format rewrites the given source in the canonical form of the Go toolchain — the same bytes `gofmt` would write: standard indentation, aligned struct fields and comment blocks, no trailing whitespace. Everything agnos generates goes through it, so a regenerated tree diffs to zero against one a formatting editor has saved. |

### `File`

File is the parsed view of one Go source file.

| Field | Type | Description |
| --- | --- | --- |
| `Package` | `string` | Package is the name from the package clause. |
| `Doc` | `string` | Doc is the file's leading doc comment, trimmed, "" when absent. |
| `Imports` | `[]Import` | Imports are the import specs, in source order. |
| `Functions` | `[]Function` | Functions are the top-level function and method declarations, in source order. |
| `Types` | `[]Type` | Types are the top-level type declarations, in source order. |
| `Constants` | `[]Value` | Constants are the names declared in top-level `const` blocks. |
| `Variables` | `[]Value` | Variables are the names declared in top-level `var` blocks. |

### `Import`

Import is a single import spec.

| Field | Type | Description |
| --- | --- | --- |
| `Alias` | `string` | Alias is the explicit local name ("_", "." or an identifier), "" when the import uses its default name. |
| `Path` | `string` | Path is the unquoted import path. |

### `Function`

Function is a top-level function or method declaration, or one method of an interface type.

| Field | Type | Description |
| --- | --- | --- |
| `Name` | `string` | Name is the function or method name. |
| `Doc` | `string` | Doc is the declaration's doc comment, trimmed, "" when absent. |
| `Receiver` | `string` | Receiver is the receiver type name without a leading "*", "" for plain functions and interface methods. |
| `Pointer` | `bool` | Pointer reports whether the receiver is a pointer. |
| `Params` | `[]Param` | Params are the parameters, one entry per name (unnamed params get a single entry with an empty Name). |
| `Results` | `[]Param` | Results are the results, same expansion rule as Params. |
| `Exported` | `bool` | Exported reports whether Name is exported. |

### `Param`

Param is one parameter or result of a function or method.

| Field | Type | Description |
| --- | --- | --- |
| `Name` | `string` | Name is the identifier, "" when the param is unnamed. |
| `Type` | `string` | Type is the type expression, rendered back to source. |

### `Type`

Type is a top-level type declaration.

| Field | Type | Description |
| --- | --- | --- |
| `Name` | `string` | Name is the type name. |
| `Doc` | `string` | Doc is the declaration's doc comment, trimmed, "" when absent. |
| `Kind` | `string` | Kind is one of "struct", "interface", "alias" (`type A = B`) or "other" (any other defined type, e.g. `type ID int`). |
| `Fields` | `[]Field` | Fields are the struct fields, populated only when Kind is "struct". |
| `Methods` | `[]Function` | Methods are the interface methods, populated only when Kind is "interface"; embedded interfaces are skipped. |
| `Underlying` | `string` | Underlying is the underlying type expression, populated only when Kind is "alias" or "other". |
| `Exported` | `bool` | Exported reports whether Name is exported. |

### `Field`

Field is one field of a struct type.

| Field | Type | Description |
| --- | --- | --- |
| `Name` | `string` | Name is the field name, "" for an embedded field. |
| `Type` | `string` | Type is the field type expression, rendered back to source. |
| `Tag` | `string` | Tag is the unquoted struct tag, "" when absent. |
| `Doc` | `string` | Doc is the field's doc comment, trimmed, "" when absent. |
| `Exported` | `bool` | Exported reports whether the field is exported (for an embedded field, whether its type name is exported). |

### `Value`

Value is one name declared in a top-level `const` or `var` block.

| Field | Type | Description |
| --- | --- | --- |
| `Name` | `string` | Name is the identifier. |
| `Doc` | `string` | Doc is the spec's doc comment, trimmed, "" when absent. |
| `Type` | `string` | Type is the declared type expression, "" when the type is inferred. |
| `Value` | `string` | Value is the assigned expression rendered back to source, "" when the spec assigns nothing (a const repeating the previous expression, or a var declared by type alone). |
| `Exported` | `bool` | Exported reports whether Name is exported. |

## `deps.Iodeps`

`sandbox/deps/iodeps`

### `Lib`

Lib is the filesystem library injected whole as the Deps.IoLib field. Paths are whatever the host operating system accepts, resolved by the adapter — unlike embeddeps.Lib, which is always slash-separated and rooted at an asset tree. The listing functions report paths that already include the directory they were given, so a result can be passed straight back in. The predicates report false rather than an error: a path that cannot be stat'd is not a directory and is not a file, which is the answer the caller wanted either way.

| Field | Type | Description |
| --- | --- | --- |
| `ReadFile` | `func(path string) ([]byte, error)` | ReadFile returns the whole content of the file at path. The error reports a file that does not exist or could not be read. |
| `WriteFile` | `func(path string, content []byte) error` | WriteFile writes content to path, creating any missing parent directory first and truncating an existing file. The error reports a directory or a file that could not be written. |
| `IsDir` | `func(path string) bool` | IsDir reports whether path exists and is a directory. |
| `IsFile` | `func(path string) bool` | IsFile reports whether path exists and is not a directory. |
| `Exist` | `func(path string) bool` | Exist reports whether anything exists at path, directory or file. |
| `CreateDir` | `func(path string)` | CreateDir creates the directory at path together with any missing parent. It reports nothing: a directory that already exists and a directory just created are the same outcome to the caller. |
| `RemoveDir` | `func(path string)` | RemoveDir removes the directory or file at path and any children it contains. It reports nothing: a missing path and a path just removed are the same outcome. |
| `ListDirs` | `func(path string) []string` | ListDirs returns the directories directly inside path. Nested directories are not descended into. |
| `ListFiles` | `func(path string) []string` | ListFiles returns the files directly inside path. Directories are not reported. |
| `ListAll` | `func(path string) []string` | ListAll returns every entry directly inside path, directories and files alike. |
| `ListDirsRecursively` | `func(path string) []string` | ListDirsRecursively returns every directory at or below path, excluding path itself. |
| `ListFilesRecursively` | `func(path string) []string` | ListFilesRecursively returns every file at or below path, at any depth. Directories are never reported. |
| `ListAllRecursively` | `func(path string) []string` | ListAllRecursively returns every entry at or below path, directories and files alike, excluding path itself. |

## `deps.Rundeps`

`sandbox/deps/rundeps`

### `Lib`

Lib is the process runner injected whole as the Deps.Rundeps field. It is what the build action reaches for when it has to hand the rendered project to a real toolchain (`go mod tidy`, `go build ./...`) and report whether that toolchain accepted it.

| Field | Type | Description |
| --- | --- | --- |
| `Run` | `func(props RunProps) (Result, error)` | Run executes one program to completion and returns what it wrote. A non-zero exit status is reported in Result.ExitCode, not as the error: the error is reserved for a program that could not be started at all (missing binary, unreadable directory). |

### `RunProps`

RunProps describes one program invocation.

| Field | Type | Description |
| --- | --- | --- |
| `Dir` | `string` | Dir is the working directory the program runs in. "" means the current directory. |
| `Program` | `string` | Program is the executable to run, looked up on PATH. |
| `Args` | `[]string` | Args are the arguments handed to the program, excluding its own name. |
| `Env` | `[]string` | Env is a list of "KEY=VALUE" entries added on top of the current process environment for this one invocation (later entries win). Empty means "inherit the environment unchanged" — the common case. `agnos compile` uses it to set GOOS/GOARCH/CGO_ENABLED per cross-compile. |

### `Result`

Result is what one finished invocation produced.

| Field | Type | Description |
| --- | --- | --- |
| `Output` | `string` | Output is the program's standard output and standard error, merged in the order they were written. |
| `ExitCode` | `int` | ExitCode is the program's exit status; 0 means success. |

## `deps.Serializables`

`sandbox/deps/serializables`

### `SerializibleObject`

SerializibleObject is one node of a parsed document — a scalar, an object or an array — and the whole tree is navigated and edited through its function fields. The same struct is what the Create* constructors of Lib return, so a document can be built in memory and serialized without ever being parsed.

| Field | Type |
| --- | --- |
| `IsInt` | `func() bool` |
| `IsString` | `func() bool` |
| `IsFloat` | `func() bool` |
| `IsBool` | `func() bool` |
| `IsNull` | `func() bool` |
| `IsObject` | `func() bool` |
| `IsArray` | `func() bool` |
| `GetInt` | `func() (int64, error)` |
| `GetFloat` | `func() (float64, error)` |
| `GetString` | `func() (string, error)` |
| `GetBool` | `func() (bool, error)` |
| `GetObjectItem` | `func(key string) (*SerializibleObject, error)` |
| `HasKey` | `func(key string) bool` |
| `GetKeys` | `func() ([]string, error)` |
| `GetArrayItem` | `func(index int) *SerializibleObject` |
| `GetArraySize` | `func() (int, error)` |
| `AddItemToObject` | `func(key string, item any) error` |
| `ReplaceItemInObject` | `func(key string, item any) error` |
| `DeleteItemFromObject` | `func(key string) error` |
| `AddItemToArray` | `func(item any) error` |
| `DeleteItemFromArray` | `func(index int) error` |

### `Lib`

Lib is the JSON/YAML codec injected whole as the Deps.Serializables field: constructors for every node kind, the two parsers and the two serializers.

| Field | Type |
| --- | --- |
| `CreateString` | `func(value string) *SerializibleObject` |
| `CreateInt` | `func(value int64) *SerializibleObject` |
| `CreateFloat` | `func(value float64) *SerializibleObject` |
| `CreateBool` | `func(value bool) *SerializibleObject` |
| `CreateNull` | `func() *SerializibleObject` |
| `CreateObject` | `func() *SerializibleObject` |
| `CreateArray` | `func() *SerializibleObject` |
| `ParseJson` | `func(data string) (*SerializibleObject, error)` |
| `ParseYaml` | `func(data string) (*SerializibleObject, error)` |
| `SerializeToJson` | `func(data *SerializibleObject) string` |
| `SerializeToYaml` | `func(data *SerializibleObject) string` |

## `deps.Std`

`sandbox/deps/std`

### `Lib`

Lib is the runtime library injected whole as the Deps.Std field.

| Field | Type | Description |
| --- | --- | --- |
| `Now` | `func() time.Time` | Now returns the current wall-clock time. |
| `Printf` | `func(format string, a ...any) (n int, err error)` | Printf writes one formatted message to standard output. It carries the command's result — the data a script would read — so it is never silenced. |
| `Log` | `func(format string, a ...any) (n int, err error)` | Log writes one formatted progress message to standard error. It is the channel every "… started with path …" notice goes through, so a caller can keep stdout free of log noise, and it is what --quiet turns off. |
| `Error` | `func(format string, a ...any) (n int, err error)` | Error writes one formatted message to standard error. |
| `Errorf` | `func(format string, a ...any) error` | Errorf formats an error message and returns it as an error. |
