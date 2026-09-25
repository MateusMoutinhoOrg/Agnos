# Commands

`agnos <command> [flags] [args]`. `agnos help <command>`, or `agnos <command> --help`,
prints the same for one command; an empty command line prints the general help and exits 2.
A command declaring a `--help` flag of its own keeps it, and is described through `help` alone.

One page per command, each rendered from that command's `entries.yaml`
([EntriesYaml](../EntriesYaml/doc.md)) on each build — open the one you need rather than this
whole page. Hidden commands are not listed. Flags may appear anywhere on the command line;
positionals bind in order after them. A `repeatable` field is given once per value.

## Deps System

| Command | Does |
| --- | --- |
| [`add-adapter`](add-adapter.md) | Installs one further adapter for a contract the project already has |
| [`add-available`](add-available.md) | Declares one further available |
| [`add-dep`](add-dep.md) | Installs one dep of the embedded catalog into the project |
| [`deps-init`](deps-init.md) | Initializes the dependency-injection subsystem for the project |
| [`deps-purge`](deps-purge.md) | Removes the dependency-injection subsystem from the project |
| [`list-adapters`](list-adapters.md) | Lists the adapters of the catalog and of the project |
| [`list-deps`](list-deps.md) | Lists the deps the embedded catalog can install |
| [`remove-adapter`](remove-adapter.md) | Uninstalls one adapter, leaving the contract it filled |
| [`remove-available`](remove-available.md) | Deletes one available |
| [`remove-dep`](remove-dep.md) | Uninstalls one dep from the project |
| [`set-adapter`](set-adapter.md) | Changes which adapter an available binds for one dep |
| [`set-dep`](set-dep.md) | Moves one remote dep to another version of its module |

## Cli System

| Command | Does |
| --- | --- |
| [`add-arg`](add-arg.md) | Add a positional arg to a command's entries.yaml |
| [`add-command`](add-command.md) | Scaffold a new command package in the project |
| [`add-flag`](add-flag.md) | Add a flag to a command's entries.yaml |
| [`cli-init`](cli-init.md) | Initializes the CLI layer for the project |
| [`cli-purge`](cli-purge.md) | Removes the CLI layer from the project |
| [`remove-arg`](remove-arg.md) | Remove a positional arg from a command's entries.yaml |
| [`remove-command`](remove-command.md) | Delete a command package from the project |
| [`remove-flag`](remove-flag.md) | Remove a flag from a command's entries.yaml |
| [`set-command`](set-command.md) | Update the command-level keys of a command's entries.yaml |

## Server System

| Command | Does |
| --- | --- |
| [`add-body-field`](add-body-field.md) | Declare a property of a route's body json-schema |
| [`add-parameter`](add-parameter.md) | Declare one value a route reads from the query string or the headers |
| [`add-path`](add-path.md) | Add one slice of the request path to a route |
| [`add-route`](add-route.md) | Declare a new http route |
| [`import-body`](import-body.md) | Infer a route's body json-schema from an example payload |
| [`remove-body-field`](remove-body-field.md) | Delete one property from a route's body json-schema |
| [`remove-parameter`](remove-parameter.md) | Delete one entry of a route's parameters |
| [`remove-path`](remove-path.md) | Delete one entry of a route's paths |
| [`remove-route`](remove-route.md) | Delete one declared route |
| [`server-init`](server-init.md) | Add the http server layer to the project |
| [`server-purge`](server-purge.md) | Remove the http server layer and every route in it |
| [`set-body`](set-body.md) | Rewrite the body keys of a route.yaml |
| [`set-body-field`](set-body-field.md) | Rewrite one property of a route's body json-schema |
| [`set-parameter`](set-parameter.md) | Rewrite one entry of a route's parameters |
| [`set-path`](set-path.md) | Rewrite one entry of a route's paths |
| [`set-route`](set-route.md) | Rewrite the route-level keys of a route.yaml |
| [`show-route`](show-route.md) | Print one route's whole declaration as a tree |

## Examples

| Command | Does |
| --- | --- |
| [`add-cli-example`](add-cli-example.md) | Scaffold a new example under examples/cli/ |
| [`add-lib-example`](add-lib-example.md) | Scaffold a new example under examples/lib/ |
| [`exec-test`](exec-test.md) | Run the project's examples and check them against their goldens |
| [`remove-cli-example`](remove-cli-example.md) | Delete an example from examples/cli/ |
| [`remove-lib-example`](remove-lib-example.md) | Delete an example from examples/lib/ |
| [`update-test`](update-test.md) | Rewrite one example's golden with what it produces now |

## Database System

| Command | Does |
| --- | --- |
| [`add-database`](add-database.md) | Declare a new database in the project |
| [`add-table`](add-table.md) | Declare one collection of records on a database |
| [`add-table-field`](add-table-field.md) | Declare one field on a table of a database |
| [`database-init`](database-init.md) | Add the database layer to the project |
| [`database-purge`](database-purge.md) | Remove the database layer and every declared database |
| [`remove-database`](remove-database.md) | Delete one database package whole |
| [`remove-table`](remove-table.md) | Delete one collection from a database |
| [`remove-table-field`](remove-table-field.md) | Delete one declared field from a table |
| [`set-table-field`](set-table-field.md) | Rewrite one declared field of a table |
| [`show-database`](show-database.md) | Print one database's whole declaration as a tree |

## Documentation

| Command | Does |
| --- | --- |
| [`add-doc`](add-doc.md) | Scaffold a new doc directory under docs/ |
| [`remove-doc`](remove-doc.md) | Delete a doc directory from docs/ |

## Front System

| Command | Does |
| --- | --- |
| [`add-page`](add-page.md) | Declare a new html page |
| [`front-init`](front-init.md) | Add the html front layer to the project |
| [`front-purge`](front-purge.md) | Remove the html front layer from the project |
| [`remove-page`](remove-page.md) | Remove an html page |

## Core Commands

| Command | Does |
| --- | --- |
| [`build`](build.md) | Build the project in a directory |
| [`compile`](compile.md) | Cross-compile the project's binaries into release/ |
| [`local-install`](local-install.md) | Builds the project and installs it locally |
| [`publish`](publish.md) | Builds, compiles and publishes a release via gh |
| [`start`](start.md) | Initialize a new project in a directory |
| [`verify`](verify.md) | Checks the project keeps the sandbox/adapter schema |

## Extensions

| Command | Does |
| --- | --- |
| [`disable-extension`](disable-extension.md) | Turn one generation mechanic off |
| [`enable-extension`](enable-extension.md) | Turn one generation mechanic on |
| [`list-extensions`](list-extensions.md) | Lists the generation mechanics and which are on |

## Info

| Command | Does |
| --- | --- |
| [`help`](help.md) | Display help for a command |
| [`interview`](interview.md) | Guided mode: answer questions instead of typing commands |
| [`version`](version.md) | Print the installed version |

Output channels and exit codes are in [Rules](../Rules/doc.md#output-channels).
