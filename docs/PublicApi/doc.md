# PublicApi

Module `github.com/MateusMoutinhoOrg/Agnos`. Every contract is a struct of function fields filled by binders; implementations live in `sandbox/internal` and are unreachable.

## Entry points

| Symbol | Signature |
|---|---|
| `sandbox.New` | `func(deps *deps.Deps) *api.Sandbox` |
| `standard.New` | `func() deps.Deps` (`adapters/availables/standard`) |
| `api.Sandbox` | `struct { Actions Actions; Cli Cli }` |
| `api.Cli` | `struct { CliMain func(args []string) int }` |
| `api.ExitOk`, `api.ExitFailure`, `api.ExitUsage` | `0`, `1`, `2` |
| `api.RuntimeGo`, `api.RuntimeNone` | `"go"`, `"none"` |

## Actions

`api.Actions` fields. Every `path` is the project dir.

| Field | Signature |
|---|---|
| `Build` | `func(BuildProps) error` |
| `Compile` | `func(CompileProps) error` |
| `Verify` | `func(path string) error` |
| `Start` | `func(StartProps) error` |
| `DepsInit`, `DepsPurge` | `func(path string) error` |
| `DepInstall`, `DepRemove` | `func(path, dep string) error` |
| `DepList` | `func(path string) ([]string, error)` |
| `CliInit`, `CliPurge` | `func(path string) error` |
| `AddCommand` | `func(path, name, help, category string) error` |
| `RemoveCommand` | `func(path, name string) error` |
| `SetCommand` | `func(CommandProps) error` |
| `AddFlag`, `AddArg` | `func(FieldProps) error` |
| `RemoveFlag`, `RemoveArg` | `func(path, command, name string) error` |
| `AddDoc` | `func(DocProps) error` |
| `RemoveDoc` | `func(path, name string) error` |

## Props

| Struct | Fields |
|---|---|
| `BuildProps` | `Path, Runtime string` |
| `CompileProps` | `Path string; Targets []string` (target names or `"all"`) |
| `StartProps` | `Path, ProjectName string; Module *string; Force bool` |
| `CommandProps` | `Path, Command, Help, Category, LongDescription string; Hidden, Visible bool; Identifiers, Examples []string` (empty string = untouched; lists append) |
| `FieldProps` | `Path, Command, Name, Description, Type, Default, Min, Max string; Identifiers, Examples []string; Required, Array bool; Position int` (`Default/Min/Max` raw literals, `""` = unset; `Position < 0` appends) |
| `DocProps` | `Path, Name, Description string; Themes []string` (`Name` is `/`-nested for a sub-doc) |

## Dependency contracts

`deps.Deps` has one field per dir of `sandbox/deps/`, named by title-casing it. Each is a `Lib` struct.

| Field | `Lib` fields |
|---|---|
| `Argvdeps` | `New(args []string) Parser`. `Parser`: `Args []string; Used []bool; IsPresent(flags) bool; GetOptionsSize(flags) int; GetKeyValuesSize(prefixes) int; Get{String,Int,Double,Timestamp}Option(flags, occurrence); Get{String,Int,Double,Timestamp}Arg(index); GetNext{String,Int,Double,Timestamp}Arg(); Get{String,Int,Double,Timestamp}KeyValues(prefixes, occurrence)`, each returning `(T, error)` |
| `Embeddeps` | `ReadFile(path) ([]byte, error); ListFiles(path) ([]string, error); ListFilesRecursively(path) ([]string, error); RenderTemplate(path, vars any) ([]byte, error)` |
| `Goimportsdeps` | `Parse(content) (*File, error); GetPackageName(content) (string, error); GetImports(content) ([]string, error)`. `File{Package, Doc, Imports []Import, Functions []Function, Types []Type, Constants, Variables []Value}` |
| `Iodeps` | `ReadFile(path) ([]byte, error); WriteFile(path, content) error; IsDir/IsFile/Exist(path) bool; CreateDir/RemoveDir(path); ListDirs/ListFiles/ListAll(path) []string` and their `Recursively` variants |
| `Requestdeps` | `NewRequest(url) Request`. `Request{AddHeader(k, v); SetMethod(m); SetBody([]byte); Fetch() (Response, error)}`. `Response{GetStatusCode() int; GetHeader(k) string; ReadBody(size) ([]byte, error); Close() error}` |
| `Rundeps` | `Run(RunProps) (Result, error)`. `RunProps{Dir, Program string; Args, Env []string}`. `Result{Output string; ExitCode int}` |
| `Serializables` | `Create{String,Int,Float,Bool,Null,Object,Array}(...) *SerializibleObject; ParseJson/ParseYaml(data) (*SerializibleObject, error); SerializeToJson/SerializeToYaml(*SerializibleObject) string`. Object: `Is{Int,String,Float,Bool,Null,Object,Array}() bool; Get{Int,Float,String,Bool}(); GetObjectItem(key); HasKey(key); GetKeys(); GetArrayItem(i); GetArraySize(); AddItemToObject/ReplaceItemInObject(key, any); DeleteItemFromObject(key); AddItemToArray(any); DeleteItemFromArray(i)` |
| `Std` | `Now() time.Time; Printf/Log/Error(format, a...) (int, error); Errorf(format, a...) error` |
