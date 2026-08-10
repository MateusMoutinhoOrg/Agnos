# Library Usage

## Description
Index of the documentation for developers consuming Agnos as a Go library: wiring an adapter into the sandbox, calling the tracker from code, and looking up the public API. Driving the same behavior from a terminal is indexed by [CliUsage/Index.md](/docs/CliUsage/Index.md); changing the library is indexed by [Development/Index.md](/docs/Development/Index.md).

The library is always built the same way: an adapter produces a `deps.Deps`, `lib.New` injects it into the closed sandbox, and the returned `api.Lib` carries every behavior.

---

## Protocols

| Doc | Description |
| --- | --- |
| [QuickStart.md](/docs/LibUsage/Protocols/QuickStart.md) | Install the module and run the smallest program that uses it |
| [LibInitialization.md](/docs/LibUsage/Protocols/LibInitialization.md) | Install the lib, create deps via an adapter, and run a first program |
| [ManageCategories.md](/docs/LibUsage/Protocols/ManageCategories.md) | Create the categories transactions are tracked under, list them, remove one |
| [TrackTransactions.md](/docs/LibUsage/Protocols/TrackTransactions.md) | Record spend and received transactions, list them, and read a balance |
| [RunApiSample.md](/docs/LibUsage/Protocols/RunApiSample.md) | Run one of the shipped Go examples from the source tree |

---

## References

| Doc | Description |
| --- | --- |
| [PublicApi.md](/docs/LibUsage/References/PublicApi.md) | Index of every public struct, function, and field, with detail pages |
| [Adapters.md](/docs/LibUsage/References/Adapters.md) | Every shipped adapter you can inject, and when to use each one |
| [ApiSamplesList.md](/docs/LibUsage/References/ApiSamplesList.md) | Every Go example shipped in `examples/libraryExamples/` |
| [Bootstrap.md](/docs/LibUsage/References/Bootstrap.md) | Embed this library as a dependency of another library built the same way |

The `Deps` contract you fill when injecting your own behavior is documented in [HandleDependencies.md](/docs/Development/Protocols/HandleDependencies.md).
