# Library Usage

## Description
Index of the documentation for developers consuming Agnos-Cli as a Go library: wiring an adapter into the sandbox, calling the tracker from code, and looking up the public API. Driving the same behavior from a terminal is indexed by [CliUsage.md](/docs/Index/CliUsage.md); changing the library is indexed by [Development.md](/docs/Index/Development.md).

The library is always built the same way: an adapter produces a `deps.Deps`, `lib.New` injects it into the closed sandbox, and the returned `api.Lib` carries every behavior.

---

## Tutorials

| Doc | Description |
| --- | --- |
| [LibQuickStart.md](/docs/Tutorials/LibQuickStart.md) | Install the module and run the smallest program that uses it |
| [LibInitialization.md](/docs/Tutorials/LibInitialization.md) | Install the lib, create deps via an adapter, and run a first program |
| [ManageCategories.md](/docs/Tutorials/ManageCategories.md) | Create the categories transactions are tracked under, list them, remove one |
| [TrackTransactions.md](/docs/Tutorials/TrackTransactions.md) | Record spend and received transactions, list them, and read a balance |
| [RunApiSample.md](/docs/Tutorials/RunApiSample.md) | Run one of the shipped Go examples from the source tree |

---

## References

| Doc | Description |
| --- | --- |
| [PublicApi.md](/docs/References/PublicApi.md) | Index of every public struct, function, and field, with detail pages |
| [Adapters.md](/docs/References/Adapters.md) | Every shipped adapter you can inject, and when to use each one |
| [EmbeddedAssets.md](/docs/References/EmbeddedAssets.md) | Where the text the library displays comes from, and how to serve your own |
| [ApiSamplesList.md](/docs/References/ApiSamplesList.md) | Every Go example shipped in `examples/libraryExamples/` |
| [Bootstrap.md](/docs/References/Bootstrap.md) | Embed this library as a dependency of another library built the same way |

The `Deps` contract you fill when injecting your own behavior is documented in [HandleDependencies.md](/docs/Tutorials/HandleDependencies.md).
