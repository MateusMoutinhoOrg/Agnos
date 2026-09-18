# Mecânica de Database

Uma extensão `sandbox-database` que gera mecanicamente os bancos da aplicação — o código, os
métodos e a doc — a partir de uma declaração yaml.

`SampleRepo/sandbox/internal/databases/appdatabase/` é leitura de referência, não contrato: ele
mostra que a ideia fecha sobre o Keep, e a forma abaixo é a que a mecânica gera.

## A unidade

Um **banco** é `sandbox/internal/databases/<db>/`, declarado por `specs.yaml` e gerado inteiro a
partir dele. Mesma relação que `commands/<x>/entries.yaml` e `routes/<x>/route.yaml` têm com o
`new.go` que rendem.

| Arquivo | Escrito por |
|---|---|
| `specs.yaml` | os comandos abaixo, nunca à mão |
| `api.go` | gerado: `<T>Item`, `<T>New`, `<T>Filtrage`, e o struct `<Db>` com um campo func por método |
| `new.go` | gerado: o `database.Props` e o wiring de cada campo func |
| `methods.go` | gerado: o corpo de cada método |
| `methods_custom.go` | **à mão** — o único escape; nenhum build o reescreve |

**Um banco não é uma superfície da api.** `sandbox/api/` não pode importar `sandbox/internal/`,
e os métodos de um banco são tipados por tabela — não há `[]Database` genérico que sirva como
`Cli.Commands` serve para os comandos. Então não há campo em `api.Sandbox`, nem pacote em
`sandbox/constructors/`: quem precisa do banco o constrói na hora, com
`<db>.New(sandbox *api.Sandbox) *<Db>`. Construir é de graça — pelo contrato do Keep, `New`
"touches no key: building one is free and creates nothing until the first record is written".

## O dep

Keep é um repo agnos, então entra como dep **remoto**, não pelo catálogo
([Adapters](docs/Adapters/doc.md)):

```bash
agnos add-dep github.com/MateusMoutinhoOrg/Keep@<version> --as database
```

`database-init` roda isso antes de ligar a chave; o adapter sai `origin: generated`.
`sandbox/deps/database/databases.go` é a superfície inteira: `Props`, `Schema`, `Item`,
`SchemaInstance`, `SchemaItem`, `DatabaseHandle`. O `<db>.New` gerado alcança o dep por
`sandbox.Deps.Database.New(props)` — nenhum `buildDatabase` entra por parâmetro.

## specs.yaml

```yaml
name: appdatabase          # pacote appdatabase, tipo AppDatabase
prefix: app                # Props.Path
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
| toda tabela | `Add<T>(props <T>New) (<T>Item, error)` | `NewItem` |
| toda tabela | `Find<T>ById(id int64) (<T>Item, bool)` | `FindById` |
| campo `key` | `Find<T>By<Field>(v <tipo>) (<T>Item, bool)` | `FindByKey` |
| toda tabela | `List<T>(f <T>Filtrage) ([]<T>Item, error)` | `ListAll` + filtro |
| toda tabela | `Page<T>(position int, chunk int) ([]<T>Item, error)` | `List` |
| toda tabela | `Count<T>() (int, error)` | `ListAll` |
| todo campo plano | `Update<T><Field>(id int64, v <tipo>) error` | `Update` |
| toda tabela | `Remove<T>(id int64) error` | `Remove` |
| campo `link` | `Get<T><Field>(id int64) (<Target>Item, bool)` | `GetLink` |
| campo `database` | `Add<T><Sub>(parentId int64, props <Sub>New) (<Sub>Item, error)` e `List<T><Sub>(parentId int64) ([]<Sub>Item, error)` | `NewSubItem`, `SchemaItem.ListAll` |

**`Find` só nasce de campo `key`.** É o único que o Keep indexa, então é o único que uma busca
direta alcança. Um campo `string`, `int` ou `float` se alcança por `List<T>` e nada mais: gerar
um `Find<T>By<Field>` que varre a tabela inteira seria vender uma varredura com cara de busca
indexada. Por isso `<T>Filtrage` cobre **todo** campo plano — ele é o único caminho até eles.

Três regras de forma, e elas valem para todo método gerado:

- **Busca devolve `(<T>Item, bool)`, escrita devolve `error`** — a mesma convenção que o Keep já
  usa em `FindByKey` e `NewItem`. Um método gerado nunca devolve `nil` tanto para "não achei"
  quanto para "o schema não existe": o que é falha vira `error`, o que é ausência vira `false`.
- **`sandbox *api.Sandbox` primeiro**, em toda função de `methods.go`:
  `func AddUrl(sandbox *api.Sandbox, self *AppDatabase, props UrlNew) (UrlItem, error)`. É a
  regra de `sandbox/internal/`, e o struct `<Db>` guarda o sandbox para fechar sobre ele nos
  campos func.
- **Nenhuma asserção de tipo sem `ok`.** `build<T>Item` lê cada campo por
  `item.Get(name)` e converte na forma vírgula-ok; um valor do tipo errado é `error`, nunca
  panic.

`<T>New` é o struct dos campos de um insert; `<T>Filtrage` leva um campo por campo plano da
tabela — texto vira `<Field>StartsWith` e `<Field>Equals`, numérico vira `<Field>Min` e
`<Field>Max` — e zero value desliga o filtro.

`methods_custom.go` é o escape: mesmo pacote, escrito à mão, para a consulta que o `specs.yaml`
não descreve. O build nunca o lê nem o reescreve; um nome colidindo com um método gerado é erro
de compilação, e `verify` o reporta antes.

## Comandos

Categoria nova `Database System`; todos carregam `--path` e `-q`.

| Comando | Faz |
|---|---|
| `database-init` | instala o dep Keep, escreve `sandbox-database: true`, rende o grupo |
| `database-purge` | remove `utils.ExtensionFiles(sandbox, ExtensionSandboxDatabase)` mais `sandbox/internal/databases/` e `docs/Databases/` inteiros, escreve `false` |
| `add-database <db> [--prefix <p>]` | escreve `sandbox/internal/databases/<db>/specs.yaml` |
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
- `generate_database_pages.go` termina em `removeStaleDocPages`, como `generate_route_pages.go`.
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

## Aceitação

Num projeto de rascunho, com o bootstrap: `database-init`, `add-database`, `add-table`,
`add-table-field` por tipo de campo — inclusive um `link` e um `database` aninhado — e então
`build` compilando e idempotente, `verify` limpo, e a página do banco em `docs/Databases/`
listando o que foi declarado. Fechado isso, o mesmo roteiro vira exemplo nos dois lados de
`examples/`.

## Em aberto

- `Update<T><Field>` é um método por campo: uma tabela de doze campos rende doze. A alternativa é
  um `Update<T>(id int64, props <T>Update) error` só, com os campos opcionais. Um por campo é o
  default até haver motivo contra — é o que dá a assinatura tipada por campo.
