# Mecânica de Database

Uma extensão `sandbox-database` que gera mecanicamente os bancos da aplicação — o código, os
métodos e a doc — a partir de uma declaração yaml.

## A unidade

Um **banco** é `sandbox/internal/databases/<db>/`, declarado por `specs.yaml` e gerado inteiro a
partir dele. Mesma relação que `commands/<x>/entries.yaml` e `routes/<x>/route.yaml` têm com o
`new.go` que rendem.

| Arquivo | Escrito por |
|---|---|
| `specs.yaml` | os comandos abaixo, nunca à mão |
| `api.go` | gerado: `<Table>Item`, `<Table>Filtrage`, o struct `<Db>` com um campo func por método |
| `new.go` | gerado: o `database.Props` e o wiring de cada campo func |
| `methods.go` | gerado: o corpo de cada método |
| `methods_custom.go` | **à mão** — o único escape; nenhum build o reescreve |

`SampleRepo/sandbox/internal/databases/appdatabase/` é a forma exata que a geração tem que
produzir. O `specs.yaml` que rende aquele diretório byte a byte é o primeiro teste da mecânica.

## O dep

Keep é um repo agnos, então entra como dep **remoto**, não pelo catálogo
([Adapters](docs/Adapters/doc.md)):

```bash
agnos add-dep github.com/MateusMoutinhoOrg/Keep@<version> --as database
```

`database-init` roda isso antes de ligar a chave; o adapter sai `origin: generated`, como em
`SampleRepo/adapters/libs/database/adapter.yaml`. `sandbox/deps/database/databases.go` é a
superfície inteira: `Props`, `Schema`, `Item`, `SchemaInstance`, `SchemaItem`, `DatabaseHandle`.

## specs.yaml

```yaml
name: appdatabase          # pacote appdatabase, tipo AppDatabase
path: app                  # Props.Path
tables:
  - name: url
    fields:
      - {name: alias, type: key, required: true}
      - {name: link, type: string, required: true}
      - {name: creation, type: int, required: true}
      - {name: redirects, type: int, required: true}
```

`type` é um de `key | string | int | float | link | database`, um por um dos `database.Item`.
`link` exige `target: <tabela>`; `database` leva `fields` aninhados e ignora `required`.

## Métodos gerados

Derivados dos campos, cobrindo praticamente toda a superfície do Keep.

| Do quê | Método | Keep por trás |
|---|---|---|
| toda tabela | `Add<T>(props <T>New) (*<T>Item, error)` | `NewItem` |
| toda tabela | `Find<T>ById(id int64) *<T>Item` | `FindById` |
| campo `key` | `Find<T>By<Field>(v) *<T>Item` | `FindByKey` |
| campo `string`/`int`/`float` | `Find<T>By<Field>(v) *<T>Item` | `ListAll` + varredura |
| toda tabela | `List<T>(f <T>Filtrage) []<T>Item` | `ListAll` + filtro |
| toda tabela | `Page<T>(position, chunk int) ([]<T>Item, error)` | `List` |
| toda tabela | `Count<T>() int` | `ListAll` |
| todo campo plano | `Update<T><Field>(id int64, v <tipo>) error` | `Update` |
| toda tabela | `Remove<T>(id int64) error` | `Remove` |
| campo `link` | `Get<T><Field>(id int64) *<Target>Item` | `GetLink` |
| campo `database` | `Add<T><Sub>(parentId int64, props <Sub>New) (*<Sub>Item, error)` e `List<T><Sub>(parentId int64) []<Sub>Item` | `NewSubItem`, `SchemaItem.ListAll` |

Cada tabela rende também o helper privado `build<T>Item(item database.SchemaItem) *<T>Item`,
que é por onde todo método acima devolve um registro.

`<T>Filtrage` leva um campo por campo plano da tabela: texto vira `<Field>StartsWith` e
`<Field>Equals`, numérico vira `<Field>Min` e `<Field>Max`. Zero value desliga o filtro — é o
que `ListUrls` do SampleRepo já faz. `<T>New` é o struct dos campos de um insert, pela regra dos
mais de três valores.

`methods_custom.go` é o escape: mesmo pacote, escrito à mão, para a consulta que o `specs.yaml`
não descreve. O build nunca o lê nem o reescreve; um nome colidindo com um método gerado é erro
de compilação, e `verify` o reporta antes.

## Comandos

Categoria nova `Database System`; todos carregam `--path` e `-q`.

| Comando | Faz |
|---|---|
| `database-init` | instala o dep Keep, escreve `sandbox-database: true`, rende o grupo |
| `database-purge` | remove `utils.ExtensionFiles(sandbox, ExtensionSandboxDatabase)` mais `sandbox/internal/databases/` e `docs/Databases/` inteiros, escreve `false` |
| `add-database <db> [--db-path <prefix>]` | escreve `sandbox/internal/databases/<db>/specs.yaml` |
| `remove-database <db>` | apaga o diretório; recusa enquanto houver `methods_custom.go` |
| `add-table <table> --database <db>` / `remove-table` | uma tabela do `specs.yaml` |
| `add-table-field <field> --database <db> --table <t> --type <tipo> [--required] [--target <tabela>]` | um campo |
| `set-table-field` / `remove-table-field` | o par do anterior |
| `show-database <db>` | lê o `specs.yaml` de volta; não escreve nada (espelha `show-route`) |

Único editor de um `specs.yaml` são esses comandos, e eles re-renderizam o arquivo inteiro.

## A doc gerada

Grupo `doc-database`, exigindo `doc` + `sandbox-database`:

- `assets/doc-database/docs/Databases/{doc.md,props.yaml}` — o índice.
- `assets/templates/database_page.md`, renderizado uma vez por banco em `docs/Databases/<db>.md`:
  tabelas, campos (nome, tipo, required, target) e cada método gerado com sua assinatura.
- `generate_database_pages.go` termina em `removeStaleDocPages`, como
  `generate_route_pages.go`.
- `database-purge` leva `docs/Databases/` inteiro: o grupo instala só `doc.md` e `props.yaml`, e
  deixar as páginas sem `props.yaml` quebra todo build seguinte.

## Onde a mecânica se declara

1. `sandbox/internal/utils/extensions_conf.go`: const `ExtensionSandboxDatabase` e linha em
   `ExtensionCatalog()`, default `false`.
2. `sandbox/internal/utils/asset_groups.go`: `{sandbox-database, [sandbox-database], Code: true}`
   e `{doc-database, [doc, sandbox-database], false}`.
3. `assets/sandbox-database/**` e `assets/doc-database/**`.
4. A chave em `assets/start/AgnosConfig/extensions.yaml` e a linha em
   `assets/doc/docs/Extensions/doc.md`.
5. Build: `collect_databases.go`, `collect_database_docs.go`, `generate_database_new.go`,
   `generate_database_pages.go`, cada um registrado em `build_internal.go`.
6. `sandbox/internal/parsables/databaseconf/` para o `specs.yaml`.
7. `sandbox/internal/actions/verify/check_databases.go`: todo `link` tem `target` de uma tabela
   existente, todo `database` tem `fields`, nome de método gerado não colide com
   `methods_custom.go`, `specs.yaml` presente em todo diretório de `sandbox/internal/databases/`.
8. `interview`: `Database System` em `areas`, `database-init` em `extensionInit`, um degrau em
   `nextSteps`, `databases` em `scaffoldedUnits`, `--target` só com `--type link` em
   `ruled_out.go`, `remove-database` e `database-purge` em `danger.go`.
9. `AgnosConfig/structure.yaml`: entrada para `sandbox/internal/databases`.
10. Um exemplo nos dois lados de `examples/`, e a linha nova em
    `assets/doc/docs/GeneratedFiles/doc.md`.

## Regras que a mecânica herda

- `sandbox/internal/databases/` importa só `sandbox/`: erro por `sandbox.Deps.Std.Errorf`, log
  por `Std.Log`, nunca `fmt`.
- Os comandos acima são do gerador: `{{.GeneratorName}}` os prefixa, `{{.Name}}` nunca.
- Todo exportado que entrar em `sandbox/deps/` ou `sandbox/api/` carrega doc comment — `verify`
  falha sem.
- `sandbox-database` exige `sandbox` ligado, como todo `sandbox-<x>`.
- `false` é parar de gerar, nunca apagar; apagar é o que `database-purge` faz.

## Em aberto

- `Add<T>` quando a tabela tem um ou dois campos: struct `<T>New` sempre, ou posicional até três?
- `Find<T>By<Field>` num campo sem índice é `ListAll` + varredura. Gerar mesmo assim, ou só sob
  um `indexed: true` no campo?
- `<T>Filtrage` cresce com a tabela. Todo campo plano, ou só os marcados no `specs.yaml`?
