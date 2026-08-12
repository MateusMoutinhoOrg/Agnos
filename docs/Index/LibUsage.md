# Library Usage

## Description
Index of the documentation for developers consuming Agnos-Cli as a Go library: wiring an adapter into the sandbox, calling the tracker from code, and looking up the public API. Driving the same behavior from a terminal is indexed by [CliUsage.md](/docs/Index/CliUsage.md); changing the library is indexed by [Development.md](/docs/Index/Development.md).

The library is always built the same way: an adapter produces a `deps.Deps`, `lib.New` injects it into the closed sandbox, and the returned `api.Lib` carries every behavior.

---

## Tutorials

- [LibInitialization.md](/docs/Tutorials/LibInitialization.md)
  - **description:** Install the lib, create deps via an adapter, and run a first program
- [ManageCategories.md](/docs/Tutorials/ManageCategories.md)
  - **description:** Create the categories transactions are tracked under, list them, remove one
- [TrackTransactions.md](/docs/Tutorials/TrackTransactions.md)
  - **description:** Record spend and received transactions, list them, and read a balance
- [RunApiSample.md](/docs/Tutorials/RunApiSample.md)
  - **description:** Run one of the shipped Go examples from the source tree
  - [Run API Examples](/docs/Tutorials/RunApiSample.md#run-api-examples)

---

## References

- [PublicApi.md](/docs/References/PublicApi.md)
  - **description:** Index of every public-facing entry of the library, logically grouped by their role in the system
  - [Entry Points](/docs/References/PublicApi.md#entry-points)
  - [Core Interface](/docs/References/PublicApi.md#core-interface)
  - [Data Models](/docs/References/PublicApi.md#data-models)
  - [Dependency Contracts](/docs/References/PublicApi.md#dependency-contracts)
- [Adapters.md](/docs/References/Adapters.md)
  - **description:** Every shipped adapter you can inject, and when to use each one
  - [Available Adapters](/docs/References/Adapters.md#available-adapters)
  - [Embedded Libraries](/docs/References/Adapters.md#embedded-libraries)
- [ApiSamplesList.md](/docs/References/ApiSamplesList.md)
  - **description:** Every Go example shipped in `examples/libraryExamples/`
  - [Examples](/docs/References/ApiSamplesList.md#examples)


