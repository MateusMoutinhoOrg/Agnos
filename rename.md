# rename.md — nomes a padronizar

Proposta apenas: nada foi renomeado.

**Como foi testado.** Usei o `release/bootstrap.bin` (agnos v0.14.0, compilado do HEAD `114711e`) num projeto
descartável `shop`, criado fora do repo:

- `start`, depois `backoffice-init`, que instalou cli, server, front, database, deps e 13 deps.
- Backoffice ponta a ponta: `add-backoffice-user --role root`, `start-server --insecure-http`, login por formulário,
  criação de API token e a API `/api/admin/*` (`me`, `list`, `get`, `add`, `edit`, `remove`). Também conferi as
  respostas 401, 404 e 415.
- Comandos: `add-command`, `add-flag`, `add-arg`, `show-command` e `explain-command`.
- Rotas: `add-route`, `add-parameter`, `import-body`, `rename-route`, `rebalance-routes`, `list-routes` e `explain-route`.
- Banco: `add-database`, `add-table`, `add-table-field` e `show-database`.
- Resto: `add-page`, `add-doc`, `add-cli-example` com `exec-test --only`, `list-extensions`, `backoffice-purge` e `interview`.
- No repo: `verify` e `build` numa cópia (`git archive HEAD`), para conferir a idempotência sem tocar na árvore.

**Critérios** (a coluna *Por quê* cita um deles):

| Id | Critério |
|---|---|
| C1 | **Um conceito, um nome**: o mesmo objeto não pode ter dois nomes. |
| C2 | **Um nome, um conceito**: o mesmo nome não pode apontar para coisas diferentes. |
| C3 | **O verbo do cli é o verbo do código**: `add / set / remove / list / show / rename` em comandos, actions, métodos de banco e rotas. |
| C4 | **O caminho e o nome do arquivo dizem a unidade**: `<unit>.yaml`, `<Action>Props`, `<unit>_<aspecto>.go`. |
| C5 | **Vocabulário que um LLM já conhece**: inglês correto e nomes da stdlib Go/HTTP, sem falso cognato. |
| C6 | **Palavras separadas**: `backoffice-user` vira `BackofficeUser`, nunca `Backofficeuser` (o grep por `BackofficeUser` não acha `Backofficeuser`). |

Prio: **A** = confunde de verdade (grep errado, ambiguidade, typo); **M** = inconsistência clara; **B** = cosmético.

---

## 1. Conceitos com dois nomes (onde se ganha mais)

| Atual | Proposto | Onde | Por quê | Prio |
|---|---|---|---|---|
| `admin` (URL `/admin` e `/api/admin`, `routeslist/admin/`, cookie `admin_token`, título "Admin Home") vs `backoffice` (extensão, comandos, pacotes, db) | `backoffice` em todo o código: `routes/backoffice/`, cookie `backoffice_session`. A URL pode continuar `/admin`, mas só a URL | backoffice | C1. O grep por `backoffice` não acha as rotas, e o grep por `admin` não acha os pacotes | A |
| examples vs tests: `sandbox-example`, `add-cli-example`, `examples/` vs `exec-test`, `update-test`, `api.ExecTest`, `api.UpdateTest`, `actions/exec_tests/`, `actions/update_tests/` (no plural, enquanto os comandos estão no singular) | `run-examples` / `update-example`, `RunExamples` / `UpdateExample`, `actions/run_examples/` / `actions/update_example/` | examples | C1. A unidade é *example*: quem procura "como rodar um example" não chega em `exec-test` | A |
| "token" é o JWT de sessão (`backofficeauth.IssueToken`, `SessionOfToken`, cookie `admin_token`) e também o API token (`bo_…`, `backofficetokens`, tabela `apitoken`) | Sessão: `IssueSessionJWT`, `ResolveSession`, cookie `backoffice_session`. API: `api_token` em tudo | backoffice | C2 | A |
| `Font` / `fonts` (`add-parameter --font`, chave `fonts:` do `route.yaml`, `ParameterFont`, `HeaderParam`, `routeDocFonts`) | `Source` / `sources` (`--source`, `ParameterSource`, `SourceHeader`) | server | C5. *font* é falso cognato de "fonte"; em inglês quer dizer tipografia | A |
| `Filtrage` (`<T>Filtrage`, parâmetro `filtrage`) | `Filter` (`<T>Filter`, `filter`) | database | C5. *filtrage* não é inglês corrente | A |
| Texto curto de ajuda: `help:` (yaml), `Command.Help`, `Route.Help` e a flag `--help '<texto>'` de `add-command` / `add-route`. Texto longo: `long-description:` / `LongDescription` | `summary:` / `Summary` / `--summary`, e `description:` / `Description` | cli, server | C2. `--help` também é a flag que abre a ajuda (o CLAUDE.md já documenta a armadilha), e `--description` de `add-flag` é outra coisa | A |
| Verbos de edição: `set-*` (cli), `Update<T><Field>` (banco), rota `edit-backoffice-user`, `backofficerender.EditBackofficeUserForm` | `set` em tudo: `Set<T><Field>`, `set-backoffice-user` | cli, db, backoffice | C3 | A |
| Verbos de criação: `add-*` (cli), `Add<T>` (db) vs rota `create-backoffice-api-token` e `backofficetokens.Create` | `add-backoffice-api-token`, `backofficeapitokens.Add` | backoffice | C3 | M |
| `front` (`front-init`, `sandbox-front`, `OpinatedAgnosFront`, doc `FrontUsage`) vs `frontend` (`assets/frontend/`, rota `frontend`, `frontend_handler.go`, `frontend_route.yaml`, `Root = "frontend"`) | `front` em tudo (`assets/front/`, rota `front`), ou `frontend` em tudo | front | C1 | M |
| Parte "do projeto" chamada de `User` em `api.UserSandbox`, `api.UserConfig`, `usersandbox.go`, `userconfig.go`, e de `project` em `commandprops/project.go`, `routeprops/project.go` | `ProjectSandbox`, `ProjectConfig`, `projectsandbox.go`, `projectconfig.go` | api | C1 + C2. Com backoffice, "user" é `BackofficeUser`, então `UserSandbox` vira ruído no grep | M |
| `deplist` / `adapterlist` (assets) vs "catalog" (`CatalogDeps`, `LoadCatalogAdapterConf`, `origin: catalog`, help "the embedded catalog") | `assets/dep-catalog/`, `assets/adapter-catalog/` | assets | C1 | M |
| Available (`adapters/availables/`, `add-available`, `available.yaml`, `AvailableConf`, `DeclaredAvailables`, `AvailablesBinding`) | Binding (`adapters/bindings/`, `add-binding`, `binding.yaml`) | deps | C5. Adjetivo usado como substantivo; o próprio help diz "which adapter an available **binds**" | M |
| Chave de extensão `sandbox-cli` vs comandos `cli-init` / `cli-purge` (vale também para server, front, database, deps, backoffice e example) | Chaves `cli`, `server`, `front`, `database`, `deps`, `backoffice`, `example` | extensions | C1. A chave e o comando que a liga têm nomes diferentes | M |
| `start --project-name`, `api.Config.ProjectName`, `{{.ProjectName}}` vs `project.yaml: name:` e `{{.Name}}` | `ProjectName` em tudo (yaml `project-name:`) | core | C1 | M |

## 2. Diretórios e arquivos

| Atual | Proposto | Onde | Por quê | Prio |
|---|---|---|---|---|
| `sandbox/internal/routeslist/` | `sandbox/internal/routes/` | server | C4. Irmão de `commands/` e `databases/` | A |
| `InternalPureHandler.go` + função `InternalPureHandler` (e os campos `Command.InternalPureHandler` / `Route.InternalPureHandler`, templates `command_internal_pure_handler.go` / `route_internal_pure_handler.go`) | `handler.go` + `Handle` | cli, server | Arquivo Go em PascalCase; "Internal" e "Pure" não informam nada. O CLAUDE.md ainda fala em `handler.go` | A |
| `databases/<db>/specs.yaml` | `database.yaml` | database | C4. Os irmãos são `command.yaml`, `route.yaml`, `adapter.yaml`, `dep.yaml`, `available.yaml` | A |
| `entries.go` / `Entries` | `input.go` / `Input` | cli, server | "Entries" é genérico (entradas de map, do índice de docs). É a entrada declarada do handler | M |
| `docs/<Doc>/props.yaml` (`docpropsconf.DocPropsConf`, `LoadDocProps`) | `doc.yaml` (`docconf.DocConf`) | docs | C4 + C2. "props" já é `RouteProps`, `CommandProps` e os `*Props` das actions | M |
| `sandbox/internal/generated/cli/cli/new.go` e `generated/server/server/new.go` | `generated/cli/new.go` e `generated/server/new.go` | core | Diretório repetido | M |
| `commands/help`, `commands/help_flag`, `commands/version`, `commands/project` e `commands/start_server` ficam sem pasta de categoria, enquanto o agnos usa `commands/core/`, `commands/cli/`… E `add-command --category Demo` não cria `demo/` | Sempre `commands/<categoria>/<nome>/` (`commands/info/help/`, `commands/middleware/help_flag/`) | cli | C4. O caminho deveria sair da categoria | M |
| Middlewares sem marca no nome: `help_flag`, `project`, `backoffice_server` (cli); `client_ip`, `security_headers`, `same_origin`, `authentication`, `root_guard` (rotas) | Sufixo `_middleware` ou pasta `middlewares/` | cli, server | Pelo nome, um LLM não distingue middleware de comando ou rota | M |
| Middleware `project` (lê `--path` / `--quiet`) | `project_flags` | cli | "project" sozinho é o projeto inteiro | M |
| Middleware `backoffice_server` | `backoffice_start_server` | backoffice | O nome sugere um servidor; é um middleware do comando `start-server` | M |
| `sandbox/internal/parsables/` | `sandbox/internal/declarations/` (ou `schemas/`) | core | "parsables" é palavra inventada; cada pacote é o schema de um yaml | B |
| `AgnosConfig/paths.yaml`, lido por `parsables/pathreplacerconf` | Pacote `pathsconf` | config | C4. O arquivo e o parser têm nomes diferentes | B |
| `examples/<side>/<name>/TestDir`, `AssertDir` | `test-dir/`, `assert-dir/` | examples | Todo o resto é minúsculo | B |
| `backofficedb/` (dados, na raiz do projeto) com o mesmo nome do pacote `databases/backofficedb` | `data/backofficedb/` | database | O grep por `backofficedb` mistura dados e código | B |
| `assets/frontend/admin/backoffice.js` | `assets/frontend/backoffice/backoffice.js` | backoffice | C1 (admin vs backoffice) | B |
| `docs/plan/` | `docs/Plan/` | docs | Único doc em minúsculo | B |
| `AgnosTeset/` (untracked) | `AgnosTest/` | raiz | Typo, se a pasta for ficar | B |

## 3. Comandos, args e flags do cli

| Atual | Proposto | Onde | Por quê | Prio |
|---|---|---|---|---|
| O id do arg posicional muda de comando para comando: `Name` (`add-flag`, `add-arg`, `add-route`, `add-parameter`, `add-table`…), `Id` (`add-path`), `Dep` (`add-dep`), `Adapter` (`add-adapter`), `Available` (`add-available`), `Route` (`import-body`, `set-body`) | Sempre `Name`, ou sempre o nome da unidade | cli | C4 | M |
| A flag que aponta a unidade-alvo: `add-flag --command` tem id `Target`, mas `add-parameter --route` tem id `Route` e `add-table --database` tem id `Database` | id = nome da flag (`Command`) | cli | C2. Em `add-table-field`, `Target` é outra coisa: o alvo do `link` | M |
| `add-table-field --type key\|string\|int\|float\|link\|database` vs flags e rotas `string\|integer\|number\|boolean\|string-array` | O mesmo vocabulário: `integer`, `number`, `boolean`, `object` (no lugar de `database`) | database | C1. Testado: `--type number` é recusado em `add-table-field` | A |
| Path id gerado automaticamente: `add-route --trigger /api/products` cria um path chamado `Route`; `--pattern '/api/products/{Id}'` cria `ApiProducts` | Uma regra só, derivada do literal (`ApiProducts`), nunca `Route` | server | C2. Um path chamado `Route` dentro de uma rota | M |
| `add-database --prefix` / yaml `prefix:` | `--key-prefix` / `key-prefix:` | database | "prefix" também é tipo de trigger e campo da tabela `apitoken` | B |
| `start --project-name` | `--name`, como todo `add-*` (ou alinhar com a linha de `ProjectName` da §1) | core | C1 | B |
| Categorias do help ("Deps System", "Cli System", "Server System", "Front System", "Database System", "Backoffice System", "Core Commands", "Documentation") vs pastas `deps/ cli/ server/ front/ database/ backoffice/ core/ docs/` | Categoria = pasta: "Deps", "Cli", "Server", "Front", "Database", "Backoffice", "Core", "Docs" | cli | C4 | B |
| Modo verbal do help misturado: "Installs…", "Initializes…", "Declares…" vs "Add…", "Declare…", "Scaffold…" | Sempre imperativo | cli | Uniformidade facilita a busca | B |

## 4. `sandbox/api`: Props e Actions

Regra já usada em `AddDepProps`, `AddRouteProps`, `RenameRouteProps` e `ExecTestProps`: **`<Action>Props`**. Estes tipos fogem dela:

| Atual | Proposto | Por quê | Prio |
|---|---|---|---|
| `SetRoute(RouteProps)` | `SetRouteProps` | C2. Colide com `routeprops.RouteProps`, o props de cada request | A |
| `AddFlag(FlagProps)`, `SetFlag(FlagEditProps)` | `AddFlagProps`, `SetFlagProps` | C4 | A |
| `AddArg(ArgProps)`, `SetArg(ArgEditProps)` | `AddArgProps`, `SetArgProps` | C4 | A |
| `AddPath(RoutePathProps)`, `SetPath(RoutePathEditProps)` | `AddPathProps`, `SetPathProps` | C4 | M |
| `AddParameter(RouteParameterProps)`, `SetParameter(RouteParameterEditProps)` | `AddParameterProps`, `SetParameterProps` | C4 | M |
| `SetBody(RouteBodyProps)`, `AddBodyField(RouteBodyFieldProps)`, `SetBodyField(RouteBodyFieldEditProps)`, `ImportBody(RouteBodyImportProps)` | `SetBodyProps`, `AddBodyFieldProps`, `SetBodyFieldProps`, `ImportBodyProps` | C4 | M |
| `AddTableField(DatabaseFieldProps)`, `RemoveTableField(DatabaseFieldProps)`, `SetTableField(DatabaseFieldEditProps)` | `AddTableFieldProps`, `RemoveTableFieldProps`, `SetTableFieldProps` | C4 | M |
| `AddPage(PageProps)`, `AddDoc(DocProps)` | `AddPageProps`, `AddDocProps` | C4 | M |
| Actions com `(path string, name string)` no lugar de props (`Verify`, `RemoveCommand`, `AddTable`, `AddDatabase`, `AddCliExample`…) | Todas recebendo `<Action>Props` | Muda assinatura, não só nome; assim um LLM adivinha a assinatura sem ler | B |

## 5. Contratos de deps (`sandbox/deps/*`)

| Atual | Proposto | Por quê | Prio |
|---|---|---|---|
| `type Sandbox struct` em todo contrato (`argvdeps.Sandbox`, `std.Sandbox`, `serverdeps.Sandbox`…: 26 no catálogo) | `type Contract struct` (ou `<Dep>`: `argvdeps.Argv`) | C2. Neste repo, o grep por `type Sandbox struct` acha 16 arquivos além de `api.Sandbox`, e `Std *std.Sandbox` se lê como "o sandbox" | A |
| Pasta `sandbox/deps/serializables/` com `package serializibles` e tipo `SerializibleObject` | `serializables` / `SerializableObject` | Typo, e pasta diferente do package | A |
| `std.Error` (imprime em stderr) ao lado de `std.Errorf` (retorna um `error`) | `std.Error` → `Eprintf` (ou `PrintError`); `std.Log` → `Logf` | C2 + C5. Só o `f` separa "imprime" de "cria um erro", e em Go `Errorf` sempre cria | A |
| `CommandResponse.Error`, `CommandResponse.Log` | Mesma troca | Espelho do anterior | M |
| `std`, `serializables`, `interviewer`, `database` (remoto, `--as database`) sem o sufixo `deps` | `stddeps`, `serializabledeps`, `interviewdeps`, `databasedeps` | C4. 19 dos 26 deps do catálogo terminam em `deps`, e é o sufixo que evita colisão com a stdlib | M |
| `iodeps.Exist`, `smartio.Exist`, `exist.go` | `Exists`, `exists.go` | C5 | M |
| Campos de `deps.Deps`: `Argvdeps`, `Iodeps`, `Embeddeps`, `Ratelimitdeps`… | `ArgvDeps`, `IoDeps`, `EmbedDeps`, `RateLimitDeps` | C6 | B |
| Adapter com o nome do dep (`adapters/libs/sortdeps` preenche o dep `sortdeps`; a alternativa é `reflectsort`) | `<impl><dep>`: `stdsort`, `reflectsort`, `osargv`, `osio`… | C2. "sortdeps" é ao mesmo tempo o contrato e uma implementação | B |
| `adapters/libs/<adapter>/` | `adapters/impls/<adapter>/` | "libs" colide com as opinated libs | B |
| `std.Goos` | `GOOS` | C5 | B |
| `interviewer.AlternativeOption` | `interviewer.Option` | Redundante | B |

## 6. Opinated libs (`OpinatedAgnos{Cli,Server,Front,Database}`)

| Atual | Proposto | Por quê | Prio |
|---|---|---|---|
| `OpinatedAgnos*` (pasta, import alias, campo `Deps.OpinatedAgnosCli`, `utils/opinated_libs.go`, `IsOpinatedLib`, `OpinatedLib`, `isOpinatedContract`, `isInOpinatedContract`) | `OpinionatedAgnos*`, ou mais curto: `agnoscli`, `agnosserver`, `agnosfront`, `agnosdatabase` | C5. *opinated* não existe, e quem busca "opinionated" não acha. Além disso a pasta é PascalCase e o package é minúsculo, ao contrário de todo o resto | A |
| `StatusOk`, `StatusFailure` (500), `StatusUnsupportedMedia`, `StatusUnprocessable`, `StatusUnavailable` | `StatusOK`, `StatusInternalServerError`, `StatusUnsupportedMediaType`, `StatusUnprocessableEntity`, `StatusServiceUnavailable` | C5. Os nomes de `net/http` que todo LLM conhece; `StatusFailure` nem indica que é 500 | A |
| 4 structs `Trigger` (`opinatedagnoscli.Trigger`, `triggerconf.Trigger`, `commandconf.Trigger`, `routeconf.Trigger`) mais o alias `api.Trigger` | Um `triggerconf.Trigger`, usado por `commandconf` e `routeconf`, mais o `api.Trigger` de runtime | C2 | M |
| Constantes `EqualTrigger`, `StringArg`, `StringFlag`, `StringPath`, `StringType` (de ParameterType), `HandlerFailure`, `NotFoundFailure`, `HeaderParam` | Prefixo pelo tipo: `TriggerEqual`, `ArgString`, `FlagString`, `PathString`, `ParameterString`, `FailureHandler`, `FailureNotFound`, `SourceHeader` | C4. O autocomplete e o grep agrupam por tipo; hoje `StringType` não diz de quem é | M |
| `Command.IsActionable`, `Route.IsActionable`, `is_actionable.go` | `Matches`, `matches.go` | Diz o que faz: casa a linha ou o request com a declaração | M |
| `Command.CommandHandler` vs `Route.RequestHandler` | `Run` nos dois | C1. Mesmo papel, dois nomes | M |
| `Cli.CliMain`, `Sandbox.ServerMain` | `Cli.Main`, `Server.Main` | Nome repetido | B |
| `Route.AcceptMethods` | `Methods` | C1 com a chave `methods:` do yaml | B |
| `Trigger.Exist` | `Set` | C5 | B |
| `database.DatabaseHandle` | `database.Handle` | Nome repetido | B |

## 7. Internos: `utils`, `build`, `verify`, `smartio`

| Atual | Proposto | Onde | Por quê | Prio |
|---|---|---|---|---|
| `smartio.WriteFile` (recusa se o arquivo existe) vs `WriteFileOverwrite` | `CreateFile` vs `WriteFile` | smartio | C5. Em Go, `os.WriteFile` sobrescreve, então um LLM vai escolher o errado | A |
| `smartio` / `SmartIO` | `stagedfs` / `StagedFS` | core | "Smart" não informa nada; o essencial é que os writes ficam em buffer até `Persist` | M |
| Ordem verbo/unidade: `CheckCommandLiteral` vs `RouteCheckLiteral` / `RouteCheckParameterLiteral`; `ParseCommandBound` vs `RouteParseBound`; `CheckPosition` (só de command) vs `CheckRoutePosition` | `<Verbo><Unidade><Coisa>`: `CheckRouteLiteral`, `CheckRouteParameterLiteral`, `ParseRouteBound`, `CheckCommandPosition` | utils | C4 | M |
| `utils.FieldType` (só de rota) vs `CommandFlagType`, `DatabaseFieldType` | `RouteFieldType` | utils | C4 | M |
| Funções de slice: `InsertCommandArg`, `InsertCommandFlag`, `RemoveCommandArg`, `RemoveCommandFlag`, `InsertRoutePath`, `RemoveRoutePath`, `InsertRouteParameter`, `RemoveRouteParameter`, `RemoveDatabaseField`, `RemoveDatabaseTable` | Genéricos `InsertAt[T]` / `RemoveAt[T]` | utils | C2. `RemoveCommandArg` parece a action `RemoveArg`, e o grep acha as duas | M |
| Nome exportado calculado em 5+ lugares: `utils.ExportedName`, `exportedName` (build), `CommandEntryId`, `RouteEntryId`, `RouteFieldName` | Um só `GoIdentifier(raw)` | utils | C1 | M |
| `CommandIdentifier` / `RouteIdentifier` / `DatabaseIdentifier` (devolvem o nome kebab, não um identificador Go) | `CommandName` / `RouteName` / `DatabaseName` | utils | C2 com o item anterior | M |
| `CommandArgEdited`, `CommandFlagEdited`, `RoutePathEdited`, `RouteParameterEdited`, `RouteBodyFieldEdited`, `DatabaseFieldEdited` (+ os `…EditEmpty`) | `EditCommandArg`… / `IsCommandArgEditEmpty`… | utils | Verbo primeiro, como o resto | B |
| Arquivos de `utils/` fora do padrão `<unit>_conf.go`: `docs.go`, `examples.go`, `structure.go`, `trigger.go`, `asset_groups.go`, `opinated_libs.go` | `doc_conf.go`, `example_conf.go`, `structure_conf.go`, `trigger_conf.go`, `asset_group.go`, `opinionated_lib.go` | utils | C4 | B |
| `type Reach`, declarado em `command_match.go` e usado também por rotas | `type MiddlewareReach`, em `reach.go` | utils | O nome do arquivo esconde o tipo | B |
| `SecretEnvName` em `render.go` + `backofficeauth.SecretEnv` | Um só, em `backoffice_conf.go` | utils | Arquivo errado e duplicado | B |
| `contains`, `containsPath`, `containsByte`, `AppendUnique` espalhados | `slice.go` | utils | — | B |
| Parâmetros snake_case: `bit_size`, `seq_item`, `project_path`, `old_pkg`, `new_pkg`, `from_dir`, `first_level`, `template_path`, `parent_path`, `doc_path`, `dest_path`, e `parent_id` no código gerado | camelCase | Go todo | C5. Convenção do Go, e o resto do código já usa | B |
| Iniciais em Go: `Id`, `Ip`, `Html`, `Api`, `Http`, `Url` (`PageUrl`, `IsIp`, `ClientIp`, `InsecureHttp`, `backofficerender.Html`) | `ID`, `IP`, `HTML`, `API`, `HTTP`, `URL` | Go todo | C5. Só vale se o `ExportedName` mudar junto; hoje está consistente, então mexer pela metade é pior | B |
| build: `generate_doc_pages.go` só contém `removeStaleDocPages` | `remove_stale_doc_pages.go` | build | C4 | M |
| build: `generate_error_handlers.go` (do server) ao lado de `generate_cli_error_handlers.go` | `generate_server_error_handlers.go` | build | C4 | M |
| build: `GenerateSubdocIndexes` e `generateSubdocIndexes`; `MigrateLegacyProps` e `migrateLegacyProps` | Nomes distintos (`generateSubdocIndex` para o de um item só) | build | Só a caixa distingue as duas | M |
| build: `generate_available_news.go` / `GenerateAvailableNews` | `generate_available_new.go` / `GenerateAvailableNewFiles` | build | "news" parece notícias | B |
| build: `generate_props_aggregate.go` contém `migrateHandWrittenProps` | Mover para `migrate_legacy_props.go` | build | C4 | B |
| build: `collect_deps_api.go`, `collect_deps_libs.go`, `collect_adapter_libs.go` | `collect_dep_contracts.go`, `collect_dep_libs.go`, `collect_adapter_impls.go` | build | C1 com a doc ("contract") | B |
| verify: `check_adapterlist.go` / `CheckAdapterlist`, `check_deplist.go` / `CheckDeplist` | `check_adapter_catalog.go` / `CheckAdapterCatalog`, `check_dep_catalog.go` / `CheckDepCatalog` | verify | C6 + C1 (catalog) | M |
| verify: `check_contracts.go` (dos deps) ao lado de `check_api_shape.go` | `check_dep_contracts.go` | verify | Deixa claro de quais contratos | B |
| verify: `checkUnitTree`, o único check não exportado | `CheckUnitTree` | verify | Segue o padrão `Check<X>` | B |

## 8. Templates (`assets/templates/`) e variáveis de template

| Atual | Proposto | Por quê | Prio |
|---|---|---|---|
| `{{ .Name }}` é o projeto (segundo o CLAUDE.md), mas em `example_cli.sh` / `example_lib.go` é o nome do example, e lá o projeto vira `{{ .ProjectName }}` | Nunca `.Name` sozinho: `.ProjectName`, `.GeneratorName`, `.ExampleName`, `.CommandName`, `.RouteName` | C2. Foi assim que nasceu o bug `shop exec-test` (ver Bugs) | A |
| `handle_*.go` (os 8 do server) ao lado de `cli_handle_*.go` | `server_handle_*.go` | C4 | M |
| `handle_server_error.go`, `handle_too_large.go`, `handle_wrong_content_type.go` (nos templates e em `sandbox/internal/server/errors/`) | `handle_internal_server_error.go`, `handle_payload_too_large.go`, `handle_unsupported_media_type.go` | C5. O nome do status que cada um responde (500, 413, 415) | M |
| `doc_doc.md` | `doc_page.md` | Os irmãos são `command_page.md`, `route_page.md`, `database_page.md` e `public_api_page.md` | B |
| `page_html.html`, `page_not_found.html`, `frontend_handler.go`, `frontend_route.yaml` | `front_page.html`, `front_404.html`, `front_handler.go`, `front_route.yaml` | C1 (front) | B |

## 9. Banco: declaração e código gerado

| Atual | Proposto | Por quê | Prio |
|---|---|---|---|
| Nomes grudados em `backofficedb/specs.yaml`: `backofficeuser`, `apitoken`, `passwordhash`, `expiresat`, `tokensha`, `ownerid`, `createdat`, `lastusedat`, `lastusedip` | `backoffice-user`, `api-token`, `password-hash`, `expires-at`, `token-sha256`, `owner-id`, `created-at`, `last-used-at`, `last-used-ip`, que geram `BackofficeUser`, `ApiToken`, `PasswordHash`… | C6. Hoje saem `BackofficeuserItem`, `UpdateApitokenLastusedat` e `FindApitokenByTokensha`: ilegíveis e diferentes do `BackofficeUser` usado no resto do backoffice | A |
| `<T>Item` / `<T>New` / `<T>Filtrage` | `<T>` (ou `<T>Record`) / `<T>Input` / `<T>Filter` | `ProductNew` parece um verbo; `Filtrage` fere C5 | M |
| `Update<T><Field>` | `Set<T><Field>` | C3, igual a `set-table-field` | M |
| `List<T>` (`ListProduct`, `ListBackofficeuser`) | `List<T>s` | Devolve vários; `ListProducts` é o que um LLM digitaria | M |
| `Page<T>(position int, chunk int)` | `List<T>sPage(offset int, limit int)` | Vocabulário padrão de paginação | M |
| Campo aninhado `sessions` de tipo `database`, que gera `SessionsItem`, `SessionsNew`, `AddBackofficeuserSessions` | Campo `session` com tipo `object`, gerando `BackofficeUserSession` | O tipo `database` de campo colide com o conceito de database, e a struct sai no plural | M |
| `RemoveBackofficeuserSessions(sandbox, self, parent_id, id)` (`methods_custom.go`) | `RemoveBackofficeUserSession(sandbox, db, userID, sessionID)` | C6 + camelCase | B |
| Package `store_db` (de `store-db`) vs `backofficedb` | `backoffice_db`, pela mesma regra kebab→snake | C4 | B |

## 10. Backoffice: rotas, pacotes e funções

| Atual | Proposto | Por quê | Prio |
|---|---|---|---|
| Rota `add-backoffice-user` (POST html) tem o mesmo nome do comando cli `add-backoffice-user`; também existem `add-backoffice-user-page` e `api-add-backoffice-user` | Sufixo por superfície: `…-form` (POST html), `…-page` (GET html), `api-…` (json) | C2. Comando e rota com o mesmo nome | M |
| `home`, `login`, `logout`, `authentication`, `root-guard` (nomes de rota são globais) | `backoffice-home`, `backoffice-login`, `backoffice-logout`, `backoffice-session-auth`, `backoffice-root-guard` | Um projeto que declare `login` colide | M |
| `api-authentication`, `api-root-guard` | `backoffice-api-token-auth`, `backoffice-api-root-guard` | Diz o que cada um autentica | M |
| `api-me` | `api-get-current-backoffice-user` | Segue `api-get-backoffice-user` | B |
| `backofficeguard` (contém `IsIp`, `ClientIp`, `SecurityHeaders`, `SameOrigin`, `ListensEverywhere`) | `backofficehttp` | Não é um guard; o guard é o `root_guard` | M |
| `backofficetokens` | `backofficeapitokens` | C2 (sessão vs API) | M |
| `backofficerender.Html`, `Login`, `Home`, `BackofficeApiTokens`, `CreateBackofficeApiTokenForm`, `AddBackofficeUserForm`, `EditBackofficeUserForm` | `RenderHTML`, `RenderLoginPage`, `RenderHomePage`, `RenderApiTokensPage`, `RenderAddApiTokenPage`, `RenderAddUserPage`, `RenderSetUserPage` | `Login` e `Home` colidem com as rotas no grep; `CreateBackofficeApiTokenForm` parece uma action | M |
| `<NAME>_SECRET` (ex.: `SHOP_SECRET`) e `api.BackofficeConfig.Secret` | `<NAME>_BACKOFFICE_SECRET` e `SessionSecret` | Genérico demais: colide com qualquer outro segredo do projeto | M |
| `backofficeapi.User`, `UserDocument`, `Listing`, `Ok`, `Role` | `UserJSON`, `UserResponseJSON`, `UserListJSON`, `OkJSON`; `Role` sai, porque duplica `backofficeauth.ParseRole` | Diz que é um encoder json | B |
| `backofficetokens.Listed` | `ApiTokenRow` | — | B |
| `backofficeauth.FindByLogin`, `SessionOfToken` | `FindUserByUsernameOrEmail`, `ResolveSession` | "login" também é uma rota; o par `Resolve` já existe em tokens | B |
| URLs RPC: `/admin/list-backoffice-users`, `/admin/list-backoffice-api-tokens`, `/api/admin/list-backoffice-users` (POST) | `/backoffice/users`, `/backoffice/api-tokens` | Mais curtas e previsíveis (opcional) | B |

## 11. Nomes velhos que ainda aparecem nas docs

| Onde | Diz | Hoje é |
|---|---|---|
| `CLAUDE.md` (tabela de editores, Architecture) | `entries.yaml`, `docs/EntriesYaml/doc.md` | `command.yaml`, `docs/CommandYaml/doc.md` |
| `CLAUDE.md` (Naming is load-bearing, Two layers) | `handler.go`, `CommandHandler(sandbox *api.Sandbox, command *api.Command) int` | `InternalPureHandler.go`, `InternalPureHandler(...)` |
| comentário de `api.Sandbox.Config` (template) | `<ProjectName>Config/project.yaml` | `AgnosConfig/project.yaml` |

---

## Bugs encontrados nos testes (não são renomes)

| # | O que | Reproduz | Prio |
|---|---|---|---|
| 1 | `shop help start-server` mostra a ajuda do middleware `backoffice-server` ("Read the backoffice secret…") e esconde as flags do próprio `start-server` (`--read-timeout-ms`, `--write-timeout-ms`, `--shutdown-timeout-ms`) | `backoffice-init`, depois `shop help start-server` | A |
| 2 | O `example.sh` gerado diz "run by `shop exec-test`", mas `exec-test` é do agnos. É a armadilha do CLAUDE.md: `assets/templates/example_cli.sh:1` usa `{{ .ProjectName }}` onde deveria ser `{{ .GeneratorName }}` | `add-cli-example greet` num projeto que não é o agnos | M |
| 3 | `explain-command greet bob --loud` falha com `unknown flag "--loud"` e só funciona com `--`. O comando não tem `examples:` mostrando isso | — | M |
| 4 | `add-table-field --type number` é recusado (o vocabulário de tipos do banco é `int` / `float`) | Ver §3 | M |
| 5 | Num clone limpo, o `verify` falha: `AgnosConfig/structure.yaml` descreve `release`, mas `release/` está no `.gitignore` ("ghost spec") | `git archive HEAD` → `build` | M |
| 6 | O HEAD não foi rebuildado depois do bump: um `build` no HEAD muda `docs/Requirements/doc.md` de v0.13.0 para v0.14.0 (fora isso, o build é idempotente) | `git archive HEAD` → `build` → `git diff` | B |
| 7 | A API `/api/admin/root/add-backoffice-user` exige `password`, enquanto o cli gera a senha; `edit-backoffice-user` exige todos os campos (substitui, não faz patch) | curl | B |
