# `help-flag`

Answer `<command> --help` with that command's help

A middleware: it runs on rung 5, in front of every command line matching
`*`, and hands the line on unless it answers.

Runs in front of every command line. Without --help it hands the line on. With it, it prints the help of the next strict command the line is for — or the general help when there is none — unless that command declares a --help flag of its own, which then reads it.

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | — |

| Runs in front of | When |
| --- | --- |
| [`add-adapter`](add-adapter.md) | always |
| [`add-arg`](add-arg.md) | always |
| [`add-available`](add-available.md) | always |
| [`add-body-field`](add-body-field.md) | always |
| [`add-cli-example`](add-cli-example.md) | always |
| [`add-command`](add-command.md) | always |
| [`add-database`](add-database.md) | always |
| [`add-dep`](add-dep.md) | always |
| [`add-doc`](add-doc.md) | always |
| [`add-flag`](add-flag.md) | always |
| [`add-lib-example`](add-lib-example.md) | always |
| [`add-page`](add-page.md) | always |
| [`add-parameter`](add-parameter.md) | always |
| [`add-path`](add-path.md) | always |
| [`add-route`](add-route.md) | always |
| [`add-table`](add-table.md) | always |
| [`add-table-field`](add-table-field.md) | always |
| [`build`](build.md) | always |
| [`cli-init`](cli-init.md) | always |
| [`cli-purge`](cli-purge.md) | always |
| [`compile`](compile.md) | always |
| [`database-init`](database-init.md) | always |
| [`database-purge`](database-purge.md) | always |
| [`deps-init`](deps-init.md) | always |
| [`deps-purge`](deps-purge.md) | always |
| [`disable-extension`](disable-extension.md) | always |
| [`enable-extension`](enable-extension.md) | always |
| [`exec-test`](exec-test.md) | always |
| [`explain-command`](explain-command.md) | always |
| [`explain-route`](explain-route.md) | always |
| [`front-init`](front-init.md) | always |
| [`front-purge`](front-purge.md) | always |
| [`help`](help.md) | always |
| [`import-body`](import-body.md) | always |
| [`interview`](interview.md) | always |
| [`list-adapters`](list-adapters.md) | always |
| [`list-commands`](list-commands.md) | always |
| [`list-deps`](list-deps.md) | always |
| [`list-extensions`](list-extensions.md) | always |
| [`list-routes`](list-routes.md) | always |
| [`local-install`](local-install.md) | always |
| [`publish`](publish.md) | always |
| [`rebalance-commands`](rebalance-commands.md) | always |
| [`rebalance-routes`](rebalance-routes.md) | always |
| [`remove-adapter`](remove-adapter.md) | always |
| [`remove-arg`](remove-arg.md) | always |
| [`remove-available`](remove-available.md) | always |
| [`remove-body-field`](remove-body-field.md) | always |
| [`remove-cli-example`](remove-cli-example.md) | always |
| [`remove-command`](remove-command.md) | always |
| [`remove-database`](remove-database.md) | always |
| [`remove-dep`](remove-dep.md) | always |
| [`remove-doc`](remove-doc.md) | always |
| [`remove-flag`](remove-flag.md) | always |
| [`remove-lib-example`](remove-lib-example.md) | always |
| [`remove-page`](remove-page.md) | always |
| [`remove-parameter`](remove-parameter.md) | always |
| [`remove-path`](remove-path.md) | always |
| [`remove-route`](remove-route.md) | always |
| [`remove-table`](remove-table.md) | always |
| [`remove-table-field`](remove-table-field.md) | always |
| [`rename-command`](rename-command.md) | always |
| [`rename-route`](rename-route.md) | always |
| [`server-init`](server-init.md) | always |
| [`server-purge`](server-purge.md) | always |
| [`set-adapter`](set-adapter.md) | always |
| [`set-arg`](set-arg.md) | always |
| [`set-body`](set-body.md) | always |
| [`set-body-field`](set-body-field.md) | always |
| [`set-command`](set-command.md) | always |
| [`set-dep`](set-dep.md) | always |
| [`set-flag`](set-flag.md) | always |
| [`set-parameter`](set-parameter.md) | always |
| [`set-path`](set-path.md) | always |
| [`set-route`](set-route.md) | always |
| [`set-table-field`](set-table-field.md) | always |
| [`show-command`](show-command.md) | always |
| [`show-database`](show-database.md) | always |
| [`show-route`](show-route.md) | always |
| [`start`](start.md) | always |
| [`update-test`](update-test.md) | always |
| [`verify`](verify.md) | always |
| [`version`](version.md) | always |

Middlewares · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
