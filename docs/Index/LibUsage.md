# Library Usage

## Description
Index of the documentation for developers consuming Agnos-Cli as a Go library: wiring an adapter into the sandbox, calling the tracker from code, and looking up the public API. Driving the same behavior from a terminal is indexed by [CliUsage.md](/docs/Index/CliUsage.md); changing the library is indexed by [Development.md](/docs/Index/Development.md).

The library is always built the same way: an adapter produces a `deps.Deps`, `lib.New` injects it into the closed sandbox, and the returned `api.Lib` carries every behavior.

---

## Tutorials

- [LibQuickStart.md](/docs/Tutorials/LibQuickStart.md)
  - **description:** Install the module and run the smallest program that uses it
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
  - **description:** Index of every public struct, function, and field, with detail pages
  - [Structs](/docs/References/PublicApi.md#structs)
  - [Functions](/docs/References/PublicApi.md#functions)
  - [Fields](/docs/References/PublicApi.md#fields)
- [Adapters.md](/docs/References/Adapters.md)
  - **description:** Every shipped adapter you can inject, and when to use each one
  - [Available Adapters](/docs/References/Adapters.md#available-adapters)
  - [Embedded Libraries](/docs/References/Adapters.md#embedded-libraries)
- [EmbeddedAssets.md](/docs/References/EmbeddedAssets.md)
  - **description:** Where the text the library displays comes from, and how to serve your own
  - [Why Assets Are a Dependency](/docs/References/EmbeddedAssets.md#why-assets-are-a-dependency)
  - [What Lives in the Assets](/docs/References/EmbeddedAssets.md#what-lives-in-the-assets)
  - [Reading Assets as a Consumer](/docs/References/EmbeddedAssets.md#reading-assets-as-a-consumer)
  - [Serving Assets from Somewhere Else](/docs/References/EmbeddedAssets.md#serving-assets-from-somewhere-else)
  - [Who Needs Assets Filled](/docs/References/EmbeddedAssets.md#who-needs-assets-filled)
- [ApiSamplesList.md](/docs/References/ApiSamplesList.md)
  - **description:** Every Go example shipped in `examples/libraryExamples/`
  - [Examples](/docs/References/ApiSamplesList.md#examples)
- [Bootstrap.md](/docs/References/Bootstrap.md)
  - **description:** Embed this library as a dependency of another library built the same way
  - [The Bootstrap Tree](/docs/References/Bootstrap.md#the-bootstrap-tree)
  - [The Sandbox Wall Applies to Libraries Too](/docs/References/Bootstrap.md#the-sandbox-wall-applies-to-libraries-too)
  - [Declaring the Dependency](/docs/References/Bootstrap.md#declaring-the-dependency)
  - [Filling It from the Adapter](/docs/References/Bootstrap.md#filling-it-from-the-adapter)
  - [Reaching It from the Sandbox](/docs/References/Bootstrap.md#reaching-it-from-the-sandbox)
  - [Embedding This Library in Your Own](/docs/References/Bootstrap.md#embedding-this-library-in-your-own)

The `Deps` contract you fill when injecting your own behavior is documented in [HandleDependencies.md](/docs/Tutorials/HandleDependencies.md).
