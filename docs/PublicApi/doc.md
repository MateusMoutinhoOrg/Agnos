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
| `standard.New` | `func() deps.Deps` (`adapters/bindings/standard`) |

Implementations live under `sandbox/internal` and are unreachable: every contract is a
struct of function fields, filled by a binder.

## The sandbox api

| Page | Declares |
| --- | --- |
| [`sandbox/api/generated.sandbox.go`](api.sandbox.md) | `Sandbox` |
| [`sandbox/api/actions.go`](api.actions.md) | `RuntimeGo`, `RuntimeNone`, `DefaultRoutePriority`, `DefaultMiddlewarePriority`, `BuildProps`, `CompileProps`, `StartProps`, `RunExamplesProps`, `AddDepProps`, `SetDepProps`, `RemoveDepProps`, `AddAdapterProps`, `SetAdapterProps`, `ExtensionInfo`, `DepInfo`, `AdapterInfo`, `AddFlagProps`, `AddArgProps`, `SetArgProps`, `SetFlagProps`, `AddCommandProps`, `RenameCommandProps`, `RebalanceCommandsProps`, `ExplainCommandProps`, `SetCommandProps`, `AddRouteProps`, `SetRouteProps`, `RenameRouteProps`, `RebalanceRoutesProps`, `ExplainRouteProps`, `AddTableFieldProps`, `RemoveTableFieldProps`, `SetTableFieldProps`, `AddPathProps`, `SetPathProps`, `AddParameterProps`, `SetParameterProps`, `SetBodyProps`, `AddBodyFieldProps`, `SetBodyFieldProps`, `ImportBodyProps`, `AddPageProps`, `AddDocProps`, `VerifyProps`, `EnableExtensionProps`, `DisableExtensionProps`, `ListExtensionsProps`, `DepsInitProps`, `DepsPurgeProps`, `ListDepsProps`, `RemoveAdapterProps`, `ListAdaptersProps`, `AddBindingProps`, `RemoveBindingProps`, `CliInitProps`, `CliPurgeProps`, `RemoveCommandProps`, `ListCommandsProps`, `ShowCommandProps`, `RemoveFlagProps`, `RemoveArgProps`, `ServerInitProps`, `ServerPurgeProps`, `RemoveRouteProps`, `RemovePathProps`, `RemoveParameterProps`, `RemoveBodyFieldProps`, `ShowRouteProps`, `ListRoutesProps`, `DatabaseInitProps`, `DatabasePurgeProps`, `AddDatabaseProps`, `RemoveDatabaseProps`, `AddTableProps`, `RemoveTableProps`, `ShowDatabaseProps`, `FrontInitProps`, `FrontPurgeProps`, `BackofficeInitProps`, `BackofficePurgeProps`, `RemovePageProps`, `RemoveDocProps`, `AddCliExampleProps`, `RemoveCliExampleProps`, `AddLibExampleProps`, `RemoveLibExampleProps`, `UpdateExampleProps`, `InterviewProps`, `Actions` |
| [`sandbox/api/generated.cli.go`](api.cli.md) | `ExitOk`, `ExitFailure`, `ExitUsage`, `Cli` |
| [`sandbox/api/generated.clisandbox.go`](api.clisandbox.md) | `CliSandbox` |
| [`sandbox/api/generated.command.go`](api.command.md) | `ArgString`, `ArgInteger`, `ArgNumber`, `ArgUuid`, `FlagString`, `FlagInteger`, `FlagNumber`, `FlagBoolean`, `FlagStringArray`, `FlagIntegerArray`, `FailureHandler`, `FailureNotFound`, `FailureBadUsage`, `FailureUnknownFlag`, `FailureUnexpectedArg`, `ArgType`, `CommandArg`, `FlagType`, `CommandFlag`, `CommandResponse`, `CommandFailureKind`, `CommandFailure`, `Command` |
| [`sandbox/api/generated.config.go`](api.config.md) | `Config` |
| [`sandbox/api/generated.trigger.go`](api.trigger.md) | `TriggerEqual`, `TriggerPrefix`, `TriggerTextPrefix`, `TriggerSuffix`, `TriggerRegex`, `TriggerOneOf`, `TriggerType`, `Trigger` |
| [`sandbox/api/projectconfig.go`](api.projectconfig.md) | `ProjectConfig` |
| [`sandbox/api/projectsandbox.go`](api.projectsandbox.md) | `ProjectSandbox` |

## Dependency contracts

`deps.Deps` has one field per directory of `sandbox/deps/`, named by title-casing it. Each
field is that package's `Sandbox` struct, filled by `adapters/impls/<name>.Bind(&deps)`.

| Page | Declares |
| --- | --- |
| [`deps.OpinionatedAgnosCli`](deps.OpinionatedAgnosCli.md) | `TriggerEqual`, `TriggerPrefix`, `TriggerTextPrefix`, `TriggerSuffix`, `TriggerRegex`, `TriggerOneOf`, `ArgString`, `ArgInteger`, `ArgNumber`, `ArgUuid`, `FlagString`, `FlagInteger`, `FlagNumber`, `FlagBoolean`, `FlagStringArray`, `FlagIntegerArray`, `FailureHandler`, `FailureNotFound`, `FailureBadUsage`, `FailureUnknownFlag`, `FailureUnexpectedArg`, `ExitOk`, `ExitFailure`, `ExitUsage`, `TriggerType`, `Trigger`, `ArgType`, `CommandArg`, `FlagType`, `CommandFlag`, `CommandResponse`, `CommandFailureKind`, `CommandFailure`, `Command`, `Cli`, `MainProps`, `Contract`, `Error` |
| [`deps.ArgvDeps`](deps.argvdeps.md) | `Contract`, `Parser` |
| [`deps.EmbedDeps`](deps.embeddeps.md) | `Contract` |
| [`deps.GoimportsDeps`](deps.goimportsdeps.md) | `Contract`, `File`, `Import`, `Function`, `Param`, `Type`, `Field`, `Value` |
| [`deps.HashDeps`](deps.hashdeps.md) | `Contract` |
| [`deps.InterviewDeps`](deps.interviewdeps.md) | `Option`, `Contract` |
| [`deps.IoDeps`](deps.iodeps.md) | `Contract` |
| [`deps.ReflectDeps`](deps.reflectdeps.md) | `Contract` |
| [`deps.RunDeps`](deps.rundeps.md) | `Contract`, `RunProps`, `Result` |
| [`deps.SerializableDeps`](deps.serializabledeps.md) | `SerializableObject`, `Contract` |
| [`deps.ServerDeps`](deps.serverdeps.md) | `Contract`, `ServerProps`, `Server`, `Request`, `Response` |
| [`deps.SortDeps`](deps.sortdeps.md) | `Contract` |
| [`deps.StdDeps`](deps.stddeps.md) | `Contract` |
| [`deps.StringsDeps`](deps.stringsdeps.md) | `Contract` |
| [`deps.TemplateDeps`](deps.templatedeps.md) | `Contract`, `RenderProps` |
