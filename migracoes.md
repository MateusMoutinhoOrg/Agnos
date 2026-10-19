# Migrações de banco por versão — recomendação

## Contexto

O rascunho anterior deste arquivo propunha: um banco `versions` com a versão de cada banco,
`agnos inc-database <db>` (sobe a versão e cria a função de migração), `agnos update-database
[<db>]` / `update-database-all` (aplica os passos) e um middleware que encerra o programa quando
um banco está desatualizado. A ideia geral está certa. Cinco pontos do rascunho quebram com o
store (Keep) e com as regras do agnos, e é neles que a confiabilidade se ganha ou se perde:

1. **A API tipada não serve para migrar.** `generated.api.go` reflete sempre o `database.yaml`
   *atual*. Os casos que pedem migração (renomear ou remover um campo ou uma tabela, mudar um tipo)
   são justamente aqueles em que o dado antigo já não existe na API atual: o Keep recusa um campo
   que o schema não declara (`InvalidField`), e `Collection(name)` não acha uma tabela removida.
   E um passo escrito contra a API atual deixa de compilar na próxima mudança de schema.
2. **Um banco `versions` separado pode discordar dos dados.** Um snapshot ou restore do backoffice,
   uma pasta copiada ou um `--database` diferente movem os dados sem mover a versão.
3. **`update-database` não pode ser comando do agnos.** O passo é código do projeto e roda sobre
   a pasta `--database` do projeto: é `<projeto> migrate-database` (`{{.ProjectName}}`). No repo do
   agnos os dois nomes renderizam `agnos`, então o erro não apareceria aqui.
4. **O Keep não tem transação entre várias chaves.** Um passo que para no meio deixa o banco meio
   migrado, e reexecutar exige passos idempotentes, o que quem escreve migração raramente garante.
5. **Falta proteção contra downgrade.** Um binário mais antigo que os dados escreveria no formato
   antigo por cima do novo.

## Recomendação

### 1. Versão declarada no `database.yaml`
- Nova chave `version: N` (ausente = 1, então nada quebra em projeto antigo). É a versão que o
  código espera.

### 2. Migração gerada a partir dos comandos, código à mão só para regra de negócio
O `database.yaml` só muda por comando, então o agnos sabe a operação exata. Um diff não distingue
um rename de um remove seguido de add; o comando distingue. Cada edição que quebra dados:
- congela o schema anterior em `sandbox/internal/databases/<db>/migrations/v<N>.yaml`;
- registra a operação em `migrations/v<N+1>.ops.yaml` e sobe para `version: N+1`.

Decidido: **cada edição que quebra dados abre a sua própria versão, sozinha**, e não existe
`inc-database`. A versão mede o formato dos dados, não o schema. Tabela nova, campo opcional novo
e `link` opcional novo nunca abrem versão, porque uma chave ausente lê como vazia ou zero.
`add-database-migration` também abre uma versão, para seed ou regra de negócio.

O `build` renderiza `generated.migrations.go` com os `databasedeps.Props` de cada versão (mesmo
template `items` de `assets/templates/database_new.go`) e o código de cada operação.

| Mudança | Quem escreve a migração |
|---|---|
| `add-table`, `add-table-field` opcional | ninguém: não abre versão (um campo nunca gravado lê como zero) |
| `set-table-field --type key` (de `string`) | gerada: o `Repair` cria o índice e falha num valor duplicado |
| `set-table-field --rename` | gerada: copia o valor antigo para o nome novo |
| `set-table-field --type` (int↔string↔number) | gerada: converte, e um valor que não converte falha a migração (o live fica intacto) |
| `set-table-field --required` (era opcional) / `add-table-field --required` | gerada com `--default <v>`; à mão com `--fill-by-hand` (abre a versão e escreve o stub); sem nenhum dos dois, o comando recusa |
| `set-table-field --required=false` | ninguém: não abre versão |
| `remove-table-field`, `remove-table` | gerada: apaga os valores antigos |
| renomear tabela | gerada, mas precisa de um `set-table --rename` (não existe hoje) |
| regra de negócio (dividir `nome` em `first`/`last`, calcular um campo) | **dev**, em `migrations/v<N>.go`, escrito uma vez por `agnos add-database-migration <db>` |

- O passo à mão recebe `func V<N>(sandbox *api.Sandbox, from databasedeps.Database, to databasedeps.Database) error`:
  o mesmo path aberto com o schema anterior e com o novo. Ele depende só de `databasedeps`, então
  compila para sempre. Roda depois das operações geradas daquela versão.

### 3. Versão gravada dentro do próprio banco
- Uma tabela reservada `agnos-version` entra nos `Props` de todo banco: um registro com `version`
  int. A versão fica sob o mesmo `key-prefix` dos dados, viaja com snapshot e restore e passa pelo
  store, então funciona também com criptografia e com `memstorage`. Fica no lugar do banco
  `versions` do rascunho.
- `verify` passa a recusar uma tabela de usuário com esse nome.

### 4. Copiar, migrar e trocar (atomicidade sem transação)
`<projeto> migrate-database [<db>]` (sem argumento migra todos, o que cobre o `update-database-all`):

| Passo | Faz |
|---|---|
| 1 | copia `<dir>/<prefix>` para `<dir>/agnos-migrating/<prefix>` |
| 2 | roda `V<k>` de `stored+1` até `N` **na cópia**, gravando a versão na cópia após cada passo |
| 3 | `Collection.Repair` em toda tabela de `to`: refaz índices (um campo que virou `key` ganha o índice) e falha num `KeyConflict`; depois confere que todo registro tem todo campo `required` de `to`, e falha quando falta algum |
| 4 | grava o marcador `done` e renomeia live → `<dir>/agnos-pre-migration/<prefix>-v<old>`, depois cópia → live |

- Falha nos passos 1–3: apaga a cópia, e o live não foi tocado.
- Crash: na próxima execução, uma cópia sem `done` é descartada, e uma com `done` termina a troca.
- Migração é só para frente, sem `down`. Desfazer é devolver a pasta `agnos-pre-migration`.
- O motor fica na lib `OpinionatedAgnosDatabase` (`Migrate`), porque é igual em todo projeto. O
  `generated.migrations.go` só fornece os dados.
- `database-init` passa a instalar `iodeps` para copiar, renomear e remover.

### 5. Guarda em toda linha de comando
- Middleware cli `database_version` (prioridade 7, depois do `database_dir` 6), em
  `assets/database-cli/`. Vale para todo comando menos `migrate-database`, `help` e `version`.

| Estado gravado | Ação |
|---|---|
| `< version` | falha: "rode `<projeto> migrate-database`" |
| `> version` | falha: binário mais antigo que os dados (proteção contra downgrade) |
| sem registro e todas as tabelas vazias | grava `version` (banco novo, nada a migrar) |
| sem registro e com dados | trata como 1 (banco anterior a este recurso) |
| sobra `agnos-migrating/<prefix>` | falha: migração interrompida, rodar `migrate-database` |

- Rode a migração com o servidor parado. A guarda só olha na partida, então um servidor antigo
  ainda no ar não é detectado. Isso fica documentado como regra de deploy: parar, migrar, subir.

### 6. Regras de `verify`
- `version ≥ 1`.
- `migrations/v1..v<N-1>.yaml` e `migrations/v2..v<N>.ops.yaml` contíguos e válidos; um
  `migrations/v<k>.go` só para um `k` que existe.
- O `key-prefix` não muda entre versões nesta primeira entrega.
- Nomes reservados: tabela `agnos-version`, prefixos `agnos-migrating` e `agnos-pre-migration`.

## Arquivos

| Mexe em | O quê |
|---|---|
| `sandbox/internal/declarations/databaseconf/api.go` | campo `Version` |
| `sandbox/internal/actions/build/collect_databases.go`, `generate_database_new.go` | coletar versões e passos, incluir `agnos-version` nos Props, renderizar `generated.migrations.go` |
| `assets/templates/database_migrations.go`, `assets/templates/database_migration_step.go` (novos) | registro gerado e stub escrito uma vez |
| `sandbox/internal/actions/{set_table_field,remove_table_field,remove_table,add_table_field}/*_internal.go` | congelar o schema, registrar a operação e subir a versão numa edição que quebra dados |
| `sandbox/internal/actions/add_database_migration/` + `sandbox/internal/commands/database/add_database_migration/` | novo comando agnos que escreve o `v<N>.go` à mão, copiando `add_table` |
| `assets/database-cli/sandbox/internal/commands/database/migrate_database/` | comando do projeto |
| `assets/database-cli/sandbox/internal/commands/middleware/database_version/` | guarda, copiando `database_dir` |
| `assets/{dep,adapter}-catalog/OpinionatedAgnosDatabase/` | `Migrate`, `StoredVersion`, recuperação de crash; reespelhar a cópia instalada |
| `sandbox/internal/actions/database_init/database_init.go` | `iodeps` em `databaseDeps` |
| `sandbox/internal/actions/verify/check_*.go` (o de databases) | regras do item 6 |
| `assets/doc-database/docs/Databases/doc.md`, `assets/doc/docs/Rules/doc.md`, `docs/Contributing/doc.md`, `AgnosConfig/structure.yaml` | seção "Migrations", regras, padrão novo |

## Verificação

```bash
go build -o release/bootstrap.bin ./cmd/main
./release/bootstrap.bin build && ./release/bootstrap.bin build -q && git diff --quiet && echo idempotent
```

- Exemplo cli novo (`add-cli-example`), num projeto `start` + `database-init`:
  1. `add-database`, `add-table-field nome`, e gravar registros pelo binário;
  2. `set-table-field nome --rename name`, sem escrever código;
  3. `build`;
  4. o binário recusa rodar;
  5. `migrate-database` migra, o binário roda e lê `name`.
- Num segundo exemplo, `add-database-migration` gera um `V2` à mão que retorna erro: a pasta live
  sai byte a byte igual à anterior.
- Banco sem registro de versão e com dados é tratado como v1. Banco com versão maior que a do
  binário é recusado.
- `run-examples --only <nome>` e depois `update-example <nome>`.
