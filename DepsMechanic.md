# DepsMechanic

Plano de implementação para separar **contrato** de **adapter** no sistema de deps, permitir
mais de um adapter por dep, e instalar outro repo agnos como dep. Sucessor do modelo atual,
onde as duas coisas são uma unidade só.

## 1. O defeito atual

A assunção `1 dep = 1 adapter` está gravada em quatro níveis independentes:

| Nível | Assunção | Onde |
|---|---|---|
| Catálogo | um `assets/deplist/<dep>/` empacota contrato **e** adapter | `assets/deplist/argvdeps/` = `sandbox/deps/argvdeps/` + `adapters/libs/verb/` |
| Bind | `New()` liga **todo** dir de `adapters/libs/` | `collect_adapter_libs.go` → `collectLibDirs`, listagem crua |
| Verify | só checa que o campo é *mencionado* em algum fonte | `check_adapters.go` → `checkAdapterCoverage`, um `Contains` sobre a concatenação |
| Versão | o pin do módulo é arquivado sob o nome do **contrato** | `assets/depsversion.yaml`: `argvdeps: .../Verb@v0.0.1`, mas quem importa Verb é o lib `verb` |

Consequência já presente hoje, sem nenhuma mudança: **dois libs que preencham o mesmo campo de
`deps.Deps` compilam, são ambos bindados, e o último em ordem de listagem vence em silêncio.**
Nada no sistema detecta isso. É por isso que `dep-remove serverdeps` não tem resposta correta —
a pergunta não é respondível num modelo em que a dep não sabe quantos adapters tem.

Duas peças já apontam a saída: `adapters/libs/verb` prova que o nome do lib **já** pode diferir
do nome da dep, e `adapters/availables/standard/` prova que o seam de perfil **já** existe — só
nunca teve um segundo perfil nem uma declaração por trás.

## 2. Modelo alvo

Três unidades, com a relação declarada e não inferida:

| Unidade | É | Vive em | Cardinalidade |
|---|---|---|---|
| **dep** | o contrato: `sandbox/deps/<dep>/`, um campo de `deps.Deps` | sandbox (fechado) | 1 |
| **adapter** | uma implementação do contrato, exporta `Bind` | `adapters/libs/<adapter>/` | N por dep |
| **available** | uma seleção: exatamente um adapter por dep | `adapters/availables/<nome>/` | N por projeto |

`adapters/libs/` = o que o projeto **tem**. `adapters/availables/<nome>/` = quem **vence** para
cada campo. `cmd/main/main.go` já importa um available — o conceito só passa a significar algo.

Invariante novo, que hoje não existe: **todo available preenche todo campo de `Deps` exatamente
uma vez.** Zero = nil func que dá panic no primeiro uso; dois = sobrescrita silenciosa. Hoje só
o zero é pego, e mal.

## 3. Layout de arquivos

### 3.1 Catálogo embutido, partido em dois

```
assets/deplist/<dep>/          dep.yaml     + sandbox/deps/<dep>/<dep>.go
assets/adapterlist/<adapter>/  adapter.yaml + adapters/libs/<adapter>/<adapter>.go
```

```yaml
# assets/deplist/serverdeps/dep.yaml
name: serverdeps
field: Serverdeps
help: Http server contract
default-adapter: nethttp
```

```yaml
# assets/adapterlist/awslambda/adapter.yaml
name: awslambda
dep: serverdeps
help: Http server sobre AWS Lambda
module: github.com/aws/aws-lambda-go@v1.47.0   # vazio = só stdlib, nada a pinar
```

`assets/depsversion.yaml` **desaparece**: o pin vai para o `adapter.yaml` de quem de fato
importa o módulo. `assets/deplist/argvdeps/` vira `assets/deplist/argvdeps/` (contrato) +
`assets/adapterlist/verb/` (adapter, com o pin do Verb).

### 3.2 Available declarativo

```
adapters/libs/nethttp/         adapter.yaml + nethttp.go      <- instalados, coexistem
adapters/libs/awslambda/       adapter.yaml + awslambda.go
adapters/availables/standard/  available.yaml + new.go        <- new.go gerado do yaml
adapters/availables/lambda/    available.yaml + new.go
```

```yaml
# adapters/availables/standard/available.yaml
adapters: [iodeps, std, verb, nethttp]
```

O `adapter.yaml` copiado para dentro de `adapters/libs/<adapter>/` é o que diz qual dep aquele
adapter preenche — é o que torna o invariante da §2 verificável sem parsear o corpo do `Bind`.

### 3.3 Nada de registro: a árvore é o registro

Não existe `AgnosConfig/deps.yaml` nem manifesto equivalente. Depois da separação da §3.1 cada
unidade é **exatamente um diretório**, e toda pergunta que um registro responderia já está
respondida pela árvore:

| Pergunta | Respondida por |
|---|---|
| Que deps estão instaladas? | os diretórios de `sandbox/deps/` |
| Que adapters estão instalados? | os diretórios de `adapters/libs/` |
| Que dep cada adapter preenche? | `dep:` do `adapters/libs/<adapter>/adapter.yaml` |
| Qual adapter vence? | `available.yaml` de cada `adapters/availables/<nome>/` |
| Que módulo o adapter pina? | `module:` do mesmo `adapter.yaml` |
| De onde veio uma dep copiada de repo remoto? | `origin: generated` + `module:` do shim que a acompanha |

O que hoje força um manifesto é só o straddle: `assets/deplist/<dep>/` cruza duas árvores, então
`dep-remove` precisa da lista de arquivos do pacote para saber o que apagar. Partido o catálogo,
remover uma dep é remover `sandbox/deps/<dep>/` e remover um adapter é remover
`adapters/libs/<adapter>/` — diretório inteiro, sem lista. E o require sai do `go.mod` lendo o
`module:` do `adapter.yaml` que está sendo removido, em vez do `depsversion.yaml` global.

Vale para a dep remota também. `add-dep github.com/user/MathLib@v1.2.0 --as mathlib` escreve
`module: github.com/user/MathLib@v1.2.0` e `origin: generated` no `adapter.yaml` do shim — o
mesmo campo que um adapter embutido já usa para o pin dele. `set-dep mathlib --version v1.3.0`
lê dali, e `list-deps` monta a tabela inteira varrendo os dois diretórios.

Também não há sha256 gravado. O module cache é imutável por versão e o `go.sum` já garante o
conteúdo dele, então a checagem de drift compara `sandbox/deps/<nome>/` byte a byte com
`<cache>/sandbox/api/` diretamente — a mesma regra do `check_deplist`, com o cache no lugar de
`assets/`. Gravar um hash seria uma segunda fonte de verdade para algo que o Go já assina.

O único preço é uma regra de guarda: `remove-adapter` **recusa** um adapter com
`origin: generated`, porque esse shim é metade de uma dep copiada — quem a remove é `remove-dep`,
e quem a re-gera é `set-dep`.

## 4. Superfície de comandos

| Comando | Faz |
|---|---|
| `deps-init` / `deps-purge` | a camada (inalterado; par irmão de `cli-init`, `server-init`, `front-init`) |
| `add-dep <dep> [--adapter <a>]` | contrato + adapter default (ou o escolhido) + inscrição no available |
| `remove-dep <dep> [--with-adapters]` | **recusa** se houver adapter instalado; lista quem a segura |
| `add-adapter <adapter> [--available <a>]` | mais um adapter para uma dep já instalada |
| `remove-adapter <adapter>` | só o adapter; recusa se for o único que preenche o campo num available |
| `set-adapter <dep> <adapter> [--available <a>]` | troca quem o available usa (irmão de `set-route`/`set-command`) |
| `list-deps` / `list-adapters` | catálogo + instalados, com origem e available |
| `add-available <nome>` / `remove-available <nome>` | cria/remove um available |

`remove-dep serverdeps` com `nethttp` e `awslambda` instalados **recusa** e imprime os dois — a
mesma semântica que `remove-route` já tem para uma rota que tem página. `--with-adapters` faz a
cascata. Remover só uma implementação é `remove-adapter awslambda`.

Renomeações da superfície atual, para alinhar com a convenção do repo (`add-command`,
`add-route`, `add-page` são verb-first; `<layer>-init`/`<layer>-purge` é o par de camada):

| Hoje | Vira |
|---|---|
| `dep-install <dep>` | `add-dep <dep>` |
| `dep-remove <dep>` | `remove-dep <dep>` |
| `dep-list` | `list-deps` |

Com isso o prefixo `dep-` deixa de existir e a colisão de um `s` entre `deps-init` e
`dep-install` some. Categoria única `Deps System` no lugar de `Dependencies` +
`Dependency System`.

**Decisão em aberto:** os comandos de available se chamam `add-available`/`remove-available`
(nome do diretório) ou renomeamos `adapters/availables/` → `adapters/profiles/` e usamos
`add-profile`/`remove-profile`. Decidir antes da Fase 4; a segunda opção mexe no import de
`cmd/main/main.go` e nos goldens de `start`.

## 5. Dep remota: outro repo agnos

### 5.1 Mecânica

`add-dep github.com/user/MathLib@v1.2.0 --as mathlib`:

1. resolve o módulo pelo `rundeps` (`go mod download -json`) e acha o path no cache;
2. lê `<cache>/sandbox/api/*.go`;
3. **valida** contra a regra de convertibilidade da §5.3, e aborta com a lista de violações;
4. copia cada arquivo para `sandbox/deps/mathlib/`, trocando só a cláusula
   `package api` → `package mathlib`, e passa o resultado por `goimportsdeps.Format`;
5. **gera** `adapters/libs/mathlib/mathlib.go` — o shim da §5.2 — e o `adapter.yaml` com
   `origin: generated`;
6. `AddRequire` do módulo no `go.mod`;
7. inscreve o adapter no `available.yaml` e chama `build`.

Duas propriedades do repo tornam isso quase gratuito:

- **`sandbox/api/` não importa nada.** É a regra do repo (api guarda só contratos), então a
  cópia é auto-contida — nenhum import a reescrever.
- **O nome do tipo já casa.** Depois do commit que renomeou os tipos de lib para `Sandbox`,
  `sandbox/api/sandbox.go` declara `type Sandbox struct` e todo `sandbox/deps/<x>/` declara
  `type Sandbox struct`. A cópia cai no lugar sem renomear nada, e `deps.Mathlib` é do tipo
  `mathlib.Sandbox` pela mesma convenção de sempre (campo = dir title-cased).

`set-dep mathlib --version v1.3.0` re-copia e regenera. `remove-dep mathlib` apaga
`sandbox/deps/mathlib/`, `adapters/libs/mathlib/`, o require e a linha do `available.yaml`.

### 5.2 O adapter gerado: **cast direto não funciona**

Verificado com o compilador, não presumido:

```go
type Props struct{ Addr string }
type Inner struct{ Do func(p Props) int }
type Sandbox struct { Cli func(args []string) int; Inner Inner }
```

- `b.Props(aProps)` → **compila**. Só builtins, tipos idênticos.
- `b.Sandbox(aSandbox)` → **falha**: `cannot convert x (variable of struct type a.Sandbox) to
  type b.Sandbox`.

A regra do Go é identidade de tipo subjacente, e ela **não é recursiva através de tipos
nomeados**: `a.Inner` e `b.Inner` são tipos nomeados distintos, então os structs que os contêm
não são idênticos. Como `api.Sandbox` tem campos `Actions Actions` e `Cli Cli`, o cast de topo
sempre falha para qualquer api real.

A solução é um **shim gerado, um conversor por tipo nomeado, recursivo** — que é exatamente o
tipo de código que o agnos existe para escrever. Forma verificada como compilável:

```go
func convSandbox(v remoteapi.Sandbox) mathlib.Sandbox {
    return mathlib.Sandbox{
        Cli:   v.Cli,                  // func só de builtins: tipo idêntico, atribuição direta
        Inner: convInner(v.Inner),
    }
}

func convInner(v remoteapi.Inner) mathlib.Inner {
    return mathlib.Inner{
        // campo func: closure converte o param na direção **inversa**
        Do: func(p mathlib.Props) int { return v.Do(remoteapi.Props(p)) },
    }
}
```

Note a contravariância: um campo `func` precisa do conversor no sentido remoto→local para os
resultados e local→remoto para os parâmetros. O gerador emite o par `conv<T>`/`rev<T>` só para
os tipos que cada direção alcança.

O `Bind` fecha o ciclo construindo o sandbox remoto com os adapters **dele**:

```go
func Bind(deps *deps.Deps) {
    rdeps := remotestd.New()                 // github.com/user/MathLib/adapters/availables/standard
    remote := remotelib.New(&rdeps)          // github.com/user/MathLib/sandbox
    deps.Mathlib = convSandbox(*remote)
}
```

Qual available do repo remoto usar é o flag `--remote-available` (default `standard`).

O gerador tem tudo que precisa em `goimportsdeps`: `File.Types` já traz `Kind`
(`struct`/`interface`/`alias`/`other`), `Fields` com `Name`/`Type`, e `Underlying`. Nenhum dep
novo.

### 5.3 Regra de convertibilidade

| Forma | Tratamento |
|---|---|
| builtin, e slice/map/ponteiro de builtin | atribuição direta |
| struct nomeada do mesmo pacote | conversor gerado, campo a campo |
| campo `func` com params e results conversíveis | closure, params na direção inversa |
| slice/map de tipo nomeado conversível | loop gerado |
| `interface`, `chan`, generics, campo embedded | **rejeitado** |
| tipo importado de fora (`time.Time`, `io.Reader`) | **rejeitado** (impossível: api não importa nada) |

Isto **não é uma regra nova** — é a disciplina que os contratos de `sandbox/deps/` já enunciam
("Only builtin types cross this boundary — no `time.Time`, no `io.Reader`, no type of the
concrete library", em `serverdeps.go`). O plano só a estende a `sandbox/api/` de um repo que
queira ser instalável, e a torna verificável.

Um repo se declara instalável com `installable: true` em `AgnosConfig/project.yaml`; um
`check_installable_api.go` novo em `verify` aplica a tabela acima, para que o autor descubra a
violação **antes** de publicar, e não o consumidor na hora de instalar.

## 6. Fases

Cada fase termina com `verify` limpo, `build` compilando e idempotente, e `exec-test` verde.
Sempre por bootstrap (`go build -o release/bootstrap.bin ./cmd/main`), nunca com um agnos
instalado — as fases mexem em templates, collectors e deps.

### Fase 1 — partir o catálogo

- criar `assets/adapterlist/<adapter>/` para os 14 adapters; deixar em `assets/deplist/<dep>/`
  só o contrato;
- escrever os `dep.yaml` e `adapter.yaml`, movendo os pins de `assets/depsversion.yaml` para os
  `adapter.yaml` que importam o módulo; apagar `depsversion.yaml` e
  `parsables/depsversionconf/`;
- `parsables/depconf/` e `parsables/adapterconf/` novos (padrão de `commandconf`/`routeconf`);
- `dep_install` passa a instalar contrato + `default-adapter`, com flag `--adapter`;
- `check_deplist.go` vira `check_deplist.go` + `check_adapterlist.go`, mesma regra byte-a-byte;
- entrada em `AgnosConfig/structure.yaml` para `assets/adapterlist/`.

Sem mudança de superfície de comando. Nenhum golden se move além dos de `dep-install`.

### Fase 2 — available declarativo

- `available.yaml` em `adapters/availables/standard/`, escrito por `deps-init`;
- `collect_adapter_libs.go` → `collect_available_adapters.go`: lê o `available.yaml` em vez de
  listar `adapters/libs/`; `new.go` passa a ser gerado da declaração;
- `checkAdapterCoverage` reescrito: para cada available, cada campo de `Deps` é preenchido
  **exatamente uma vez**, resolvido pelo `dep` de cada `adapter.yaml` — zero e dois viram
  violações distintas e nomeadas;
- `adaptersAllowedDirs` continua `{availables, libs}`.

É aqui que o bug latente da §1 deixa de existir.

### Fase 3 — renomear a superfície

- `dep_install`→`add_dep`, `dep_remove`→`remove_dep`, `dep_list`→`list_deps`, em
  `commands/` e `actions/`; `identifiers` e `category` nos `entries.yaml`;
- `api.Actions`: `DepInstall`/`DepRemove`/`DepList` → `AddDep`/`RemoveDep`/`ListDeps`, e os
  binds correspondentes;
- exemplos: `remove-cli-example dep-install` etc., `add-cli-example add-dep` etc. — nunca por
  edição manual de `result.yaml`;
- docs e `README.md` regenerados; `docs/DepList/doc.md` reescrito.

Superfície limpa **antes** de crescer, para não misturar renomeação com comando novo no mesmo
diff.

### Fase 4 — multi-adapter

- comandos novos: `add-adapter`, `remove-adapter`, `set-adapter`, `list-adapters`,
  `add-available`, `remove-available`;
- `remove-dep` recusa com adapter instalado, `--with-adapters` cascateia;
- `remove-adapter` recusa quando é o único a preencher o campo num available;
- um exemplo por comando novo.

### Fase 5 — dep remota

- `add-dep <module>@<version> --as <nome>` (o argumento com `/` é module path; sem `/` é nome
  do catálogo — mesma desambiguação do `go get`);
- copiador de `sandbox/api/` → `sandbox/deps/<nome>/`;
- validador da §5.3 e gerador de shim da §5.2, em
  `sandbox/internal/actions/add_dep/generate_shim.go`;
- `set-dep <nome> --version`;
- `installable: true` + `check_installable_api.go`;
- `check_remote_deps.go`: compara `sandbox/deps/<nome>/` com `<cache>/sandbox/api/` byte a byte
  quando o cache está disponível, espelhando a regra do `check_deplist`;
- exemplo com um repo agnos mínimo fixado por versão.

## 7. Docs a atualizar

| Doc | Mudança |
|---|---|
| `assets/all/docs/DepList/doc.md` | vira a tabela de deps; ganha a coluna de adapters |
| `assets/all/docs/Adapters/doc.md` | novo (`agnos add-doc`), tema `development` |
| `assets/all/docs/Workflow/doc.md` | "Add a dependency" reescrito com o par dep/adapter |
| `assets/all/docs/Rules/doc.md` | regra do invariante "um adapter por campo por available"; regra da convertibilidade |
| `assets/all/docs/GeneratedFiles/doc.md` | `new.go` gerado do yaml; `adapter.yaml` e `available.yaml` são `once` |
| `assets/all/docs/Structure/doc.md` | via `AgnosConfig/structure.yaml` |
| `docs/Contributing/doc.md` | espelhar o padrão, no mesmo commit |
| `CLAUDE.md` | seção Architecture: as três unidades |
| `docs/plan/doc.md` | riscar o item 2 do roadmap quando a Fase 5 fechar |

Cuidado de sempre nos templates: `{{.GeneratorName}}` para comando do agnos, `{{.Name}}` para
comando do projeto gerado — em `assets/all/docs/` os dois renderizam `agnos` neste repo e a
troca só aparece num projeto scaffoldado.

## 8. Riscos

| Risco | Mitigação |
|---|---|
| Fase 1 move 14 deps de uma vez; um erro de byte quebra `check_deplist` | fazer um dep por commit, `verify` entre cada um |
| A cópia da api remota fica desatualizada em relação ao módulo | `check_remote_deps.go` compara com o module cache, que o `go.sum` já assina |
| Repo remoto com api não conversível | rejeitado na instalação com lista de violações; `installable: true` pega antes de publicar |
| Diamante de versões no `go.mod` | `go mod tidy` do próprio `build` resolve; conflito real vira erro do toolchain, não silêncio |
| Ordem alfabética de `libs/` sumindo muda a ordem de bind em `new.go` | ordem passa a ser a do `available.yaml`; goldens de `start` e `deps-init` se movem uma vez, na Fase 2 |
