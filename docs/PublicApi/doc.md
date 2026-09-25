# PublicApi

Every exported symbol of `github.com/MateusMoutinhoOrg/Agnos`, read straight from the contract sources on
every build: `sandbox/api/` is the surface `sandbox.New` returns, `sandbox/deps/`
the contracts an adapter fills and a caller may replace. Each description is the
doc comment of the declaration itself — change the comment, run `build`, and the page
follows.

One page per contract: the tables below say which page declares a symbol, so open that page
rather than reading the whole surface.

## Entry points

| Symbol | Signature |
| --- | --- |
| `sandbox.New` | `func(deps *deps.Deps) *api.Sandbox` |
| `standard.New` | `func() deps.Deps` (`adapters/availables/standard`) |

Implementations live under `sandbox/internal` and are unreachable: every contract is a
struct of function fields, filled by a binder.

## The sandbox api

| Page | Declares |
| --- | --- |
| [`sandbox/api/sandbox.go`](api.sandbox.md) | `Sandbox` |
| [`sandbox/api/actions.go`](api.actions.md) | `RuntimeGo`, `RuntimeNone`, `BuildProps`, `CompileProps`, `StartProps`, `ExecTestProps`, `AddDepProps`, `SetDepProps`, `RemoveDepProps`, `AddAdapterProps`, `SetAdapterProps`, `ExtensionInfo`, `DepInfo`, `AdapterInfo`, `FieldProps`, `CommandProps`, `AddRouteProps`, `RouteProps`, `DatabaseFieldProps`, `DatabaseFieldEditProps`, `RoutePathProps`, `RoutePathEditProps`, `RouteParameterProps`, `RouteParameterEditProps`, `RouteBodyProps`, `RouteBodyFieldProps`, `RouteBodyFieldEditProps`, `RouteBodyImportProps`, `PageProps`, `DocProps`, `Actions` |
| [`sandbox/api/cli.go`](api.cli.md) | `ExitOk`, `ExitFailure`, `ExitUsage`, `Cli` |
| [`sandbox/api/command.go`](api.command.md) | `CommandArg`, `CommandFlag`, `Command`, `NewCommand`, `BindCommand` |
| [`sandbox/api/config.go`](api.config.md) | `Config` |

## Dependency contracts

`deps.Deps` has one field per directory of `sandbox/deps/`, named by title-casing it. Each
field is that package's `Sandbox` struct, filled by `adapters/libs/<name>.Bind(&deps)`.

| Page | Declares |
| --- | --- |
| [`deps.Argvdeps`](deps.argvdeps.md) | `Sandbox`, `Parser` |
| [`deps.Embeddeps`](deps.embeddeps.md) | `Sandbox` |
| [`deps.Goimportsdeps`](deps.goimportsdeps.md) | `Sandbox`, `File`, `Import`, `Function`, `Param`, `Type`, `Field`, `Value` |
| [`deps.Hashdeps`](deps.hashdeps.md) | `Sandbox` |
| [`deps.Interviewer`](deps.interviewer.md) | `AlternativeOption`, `Sandbox` |
| [`deps.Iodeps`](deps.iodeps.md) | `Sandbox` |
| [`deps.Rundeps`](deps.rundeps.md) | `Sandbox`, `RunProps`, `Result` |
| [`deps.Serializables`](deps.serializables.md) | `SerializibleObject`, `Sandbox` |
| [`deps.Serverdeps`](deps.serverdeps.md) | `Sandbox`, `ServerProps`, `Server`, `Request`, `Response` |
| [`deps.Sortdeps`](deps.sortdeps.md) | `Sandbox` |
| [`deps.Std`](deps.std.md) | `Sandbox` |
| [`deps.Stringsdeps`](deps.stringsdeps.md) | `Sandbox` |
| [`deps.Templatedeps`](deps.templatedeps.md) | `Sandbox`, `RenderProps` |
