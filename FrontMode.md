# FrontMode

A layer `front`: páginas html renderizadas de assets embutidos, servidas por rotas do
projeto. Referência de implementação: `front-model/` (fora do controle de versão).

## Decisão de base: o que é gerado e o que é do usuário

`RenderGroup` escreve com `WriteFileOverwrite`. Tudo que entra num grupo é reescrito por
**todo** `build`. Só entra no grupo `front` o que o agnos garante; rota e conteúdo são
escritos **uma vez** e passam a ser do projeto — `once` no vocabulário de
`docs/GeneratedFiles/doc.md`.

| Arquivo | Escrito por | Reescrita |
|---|---|---|
| `sandbox/internal/pageio/templates.go` | `build` (grupo `front`) | `always` |
| `docs/FrontUsage/` | `build` (grupo `front`) | `always` |
| `sandbox/internal/routes/static/{route.yaml,handler.go}` | `front-init` | `once` |
| `assets/frontend/static/styles/main.css`, `scripts/main.js` | `front-init` | `once` |
| `sandbox/internal/routes/<page>/{route.yaml,handler.go}` | `add-page` | `once` |
| `assets/frontend/pages/<page>.html` | `add-page` | `once` |
| `sandbox/internal/routes/<page>/entries.go` | `build` | `always` |

Contrato com o usuário, dito assim em `docs/FrontUsage`:

- o que o agnos mantém é `pageio` — os helpers e as constantes;
- a rota `static` e as páginas são um ponto de partida, editáveis por conta e risco;
- quem editar `static/handler.go` mantém `safeSegments`: é a única coisa entre o
  `?path` do caller e o resto da árvore de assets;
- restaurar o padrão: `{{.GeneratorName}} remove-route static && {{.GeneratorName}} front-init`
  para a rota, `remove-page <p> && add-page <p>` para uma página (esse apaga o html);
  `front-purge && front-init` devolve a layer inteira sem tocar em `assets/frontend/`.

## Por que é uma layer

Segue `docs/Contributing/doc.md` ("Add a layer"), terceira coluna da tabela cli/server:

| Conceito | CLI | Server | Front |
|---|---|---|---|
| Código compartilhado | — | `sandbox/internal/routeio/` | `sandbox/internal/pageio/` |
| Grupo de assets | `assets/cli/` | `assets/server/` | `assets/front/` |
| Gatilho no build | `io.IsDir("sandbox/internal/cli")` | `io.IsDir("sandbox/internal/server")` | `io.IsDir("sandbox/internal/pageio")` |
| Par init/purge | `cli-init` / `cli-purge` | `server-init` / `server-purge` | `front-init` / `front-purge` |
| Unidade declarada | `commands/<name>/entries.yaml` | `routes/<name>/route.yaml` | `routes/<page>/route.yaml` + `assets/frontend/pages/<page>.html` |

O gatilho é `pageio/` e nunca `assets/frontend/`: essa árvore é conteúdo do usuário e
pode ficar vazia. `hasAssets` continua sendo `assets/all`, então `assets/frontend/` não
faz um projeto se declarar gerador.

Não há coletor novo e não há `docs/Pages`: uma página **é** uma rota, então
`docs/Routes` já a lista.

---

## front-init

```bash
agnos front-init [--path .] [-q]
```

Instala as deps da layer, garante o server, renderiza `pageio` e escreve uma vez a rota
`static` e o esqueleto de `assets/frontend/`.

**Fluxo (`front_init.go`)**

1. `dep-install embeddeps templatedeps hashdeps` — `std`, `stringsdeps`, `sortdeps`,
   `serializables` e `serverdeps` vêm do server. `embeddeps` é quem escreve `assets/asset.go`.
2. abre um `SmartIO`, chama `FrontInitInternal`, `Persist`, `build` com `RuntimeGo`.

**Fluxo (`front_init_internal.go`)**

1. `if !io.IsDir("sandbox/internal/server")` → `serverInitAction.ServerInitInternal(deps, io, path)`
   no **mesmo** SmartIO (padrão de `server_init` com `cli_init`: sem `Persist` nem `build`
   intermediários).
2. `utils.RenderGroup(deps, io, "front", vars)`.
3. `writeStaticRoute`: se `io.IsDir("sandbox/internal/routes/static")` → log
   `"front-init: %s already exists, keeping it"` e segue. Senão renderiza
   `templates/static_route.yaml` e `templates/static_handler.go` para dentro dele.
4. `writeFrontSkeleton`: `io.WriteFile` (recusa sobrescrita) de
   `assets/frontend/static/styles/main.css` e `assets/frontend/static/scripts/main.js`,
   a partir de `assets/templates/front_main.css` e `front_main.js`.
   Os dois arquivos são obrigatórios, não decoração: `go:embed` não guarda diretório
   vazio, e `dirref "styles"` num diretório inexistente é erro em tempo de render.

`assets/frontend/pages/` fica vazio — quem o preenche é `add-page`.

## front-purge

```bash
agnos front-purge [--path .] [-q]
```

Espelha `server_purge`: remove todo arquivo que o grupo `front` instalaria, mais os
diretórios que a layer possui inteiros, depois `build` com `RuntimeNone`.

```go
var frontDirs = []string{
    "sandbox/internal/pageio",
    "sandbox/internal/routes/static",
}
```

Mais o diretório de cada página — toda rota com `assets/frontend/pages/<name>.html` ao
lado, o mesmo teste que `remove-page` usa. Sem isso o purge deixa handlers importando
um `pageio` que não existe mais, e a árvore para de compilar; `server-purge` derruba
`sandbox/internal/routes` inteiro pela mesma razão.

**Não** remove: `assets/frontend/` (páginas, css e js são conteúdo escrito à mão, não
código da layer), a layer server, nem as deps — mesma razão pela qual `server-purge`
deixa a cli em pé. O log nomeia o que ficou.

O par disso é `add-page` preservando um html já existente: `front-purge` seguido de
`front-init` + `add-page <p>` devolve as rotas por cima do conteúdo intacto.

## add-page

```bash
agnos add-page home --title "Home" [--trigger /] [--path .] [-q]
```

**Fluxo (`add_page_internal.go`)**

1. valida o nome com `utils.ValidateRouteName` e recusa `static`/`health`.
2. se `assets/frontend/pages/<page>.html` já existe, mantém o arquivo e loga
   `"add-page: %s already exists, keeping it"` — é o caminho de volta depois de um
   `front-purge`, que derruba a rota e deixa o html. Uma rota já existente, essa sim,
   é erro: o `io.WriteFile` de `AddRouteInternal` recusa sobrescrita.
3. `addRouteAction.AddRouteInternal(deps, io, name, "GET", trigger, help, "Pages")`
   no mesmo SmartIO — sem duplicar geração de `route.yaml`. `--trigger` cai em
   `/<page>`; `--trigger /` é válido e é o caso da home (`RouteIdentifierSegment`
   devolve `"/"`, `collect_routes.go` ordena a raiz por último, então `/` nunca engole
   `/home`).
4. sobrescreve `handler.go` com `templates/page_handler.go` (`pageio.Render` +
   `pageVars` + 500 em erro de render), via `RenderTemplateToDest`.
5. `io.WriteFile("assets/frontend/pages/<page>.html", ...)` a partir de
   `templates/page_html.html`.

**Armadilha de template:** `templates/page_html.html` passa pelo render do agnos antes
de virar arquivo. Escrito cru, `{{ .Title }}` e `{{ dirref "styles" }}` são executados
no scaffold e chegam vazios ao disco. Tem que ser escapado (`{{ "{{ .Title }}" }}`).
O mesmo vale para qualquer `{{` dentro de `front_main.js`.

## remove-page

```bash
agnos remove-page home [--path .] [-q]
```

1. erro se `assets/frontend/pages/<page>.html` não existe — é isso que distingue uma
   página de uma rota qualquer.
2. `removeRouteAction.RemoveRouteInternal(deps, io, name)` no mesmo SmartIO.
3. remove o html.
4. `build` com `RuntimeNone`.

`remove-route <page>` recusa uma rota que tenha html, apontando para `remove-page`:
um editor por lugar, e sem html órfão.

---

## Mudanças em código existente

| Arquivo | Mudança |
|---|---|
| `sandbox/internal/actions/build/build_internal.go` | `hasFront := io.IsDir("sandbox/internal/pageio")`, `"HasFront": hasFront` no vars map, `RenderGroup(..., "front", vars)` sob `if hasFront` |
| `.../build/collect_generated_docs.go` | `GeneratedDocsGroups(has_cli, has_server, has_front)` |
| `.../build/collect_front_mount.go` (novo) | lê o `identifier` do primeiro segmento de `routes/static/route.yaml`, default `/static`, para o vars `"StaticMount"` |
| `.../remove_route/remove_route_internal.go` | recusa rota com `assets/frontend/pages/<name>.html` |
| `sandbox/api/actions.go`, `sandbox/binds/actions.go` | `FrontInit`, `FrontPurge`, `AddPage`, `RemovePage`, cada campo com doc-comment (`verify` falha sem, e `docs/PublicApi` sai daí) |
| `assets/all/docs/{Rules,GeneratedFiles,Workflow}/doc.md` | blocos `{{ if .HasFront }}` |
| `AgnosConfig/structure.yaml` | entrada para `assets/front/` |
| `docs/Contributing/doc.md`, `CLAUDE.md` | a tabela de layers acima, no mesmo commit |

`StaticMount` merece explicação: a constante vive em `pageio` (gerado) mas o prefixo
real está no `route.yaml` (do usuário). Renderizá-la a partir da declaração é o que
impede que renomear o mount quebre todo `staticref` em silêncio. É a única peça de
maquinário extra do plano — se for cortada, `docs/FrontUsage` tem que dizer que o mount
é fixo em `/static`.

## Declaração dos comandos

Pelo binário de bootstrap, nunca à mão (`docs/Workflow`), categoria `Front System`:

```bash
B=./release/bootstrap.bin
$B add-command front-init  --help "Add the html front layer to the project" --category "Front System"
$B add-command front-purge --help "Remove the html front layer from the project" --category "Front System"
$B add-command add-page     --help "Declare a new html page" --category "Front System"
$B add-command remove-page  --help "Remove an html page" --category "Front System"

# em cada um dos quatro:
$B add-flag path  --command <cmd> --default . --description "the dir holding the project (defaults to the current directory)"
$B add-flag quiet --command <cmd> --identifier --quiet --identifier -q --type boolean --description "Quiets the cli output"

$B add-arg  name    --command add-page --required --description "the page name (becomes the route, its Go package and the html file)"
$B add-flag trigger --command add-page --description "the literal segment the page answers on, / included (defaults to /<name>)"
$B add-flag title   --command add-page --description "the <title> the scaffolded page carries"
$B add-arg  name    --command remove-page --required --description "the page to remove"
```

`add-flag`/`add-arg` leem a própria linha de comando: nenhum valor pode ser exatamente
`--identifier` ou `--example`.

## Docs

- `assets/front/docs/FrontUsage/{doc.md,props.yaml}` — os helpers (`staticref`,
  `cssref`, `jsref`, `dirref`, `inline`, `include`), o contrato de `pageVars`, o mount,
  e o contrato "gerado vs seu" do topo deste plano. `{{.GeneratorName}}` para comando do
  agnos, `{{.Name}}` para comando do projeto — nunca `agnos` hardcoded.
- linhas novas na tabela de `assets/all/docs/GeneratedFiles/doc.md`, na coluna certa
  (`always` só para `pageio` e `docs/FrontUsage`).
- regra nova (páginas/rota são `once`) vai em `assets/all/docs/Rules/doc.md`, em lugar
  nenhum mais.

## Exemplos

Um par cli+lib por comando, como os 27/26 de hoje:

```bash
$B add-cli-example front-init && $B add-lib-example front-init   # idem front-purge, add-page, remove-page
$B update-test front-init                                        # golden, um por vez
```

Cada exemplo copia para `AssertDir` só o que o seu comando toca — `front-init` copia
`sandbox/internal/pageio/`, `routes/static/` e `assets/frontend/`; `add-page` copia
`routes/home/` e `assets/frontend/pages/home.html`.

## Ordem de execução

1. `sandbox/internal/pageio/templates.go` + grupo `assets/front/` + `hasFront` no build.
2. `front-init` (deps, composição com o server, rota `static` `once`, esqueleto).
3. `front-purge`.
4. `add-page` / `remove-page` + guard em `remove-route`.
5. `collect_front_mount.go`.
6. Docs (`FrontUsage`, `GeneratedFiles`, `Rules`, `Contributing`, `CLAUDE.md`, `structure.yaml`).
7. Exemplos + goldens.
8. `bootstrap.bin build` duas vezes, a segunda com `git diff --quiet`; `verify`; `exec-test`.
9. Bump de `version` em `AgnosConfig/project.yaml`.

## Riscos aceitos

- Rota e páginas `once`: correção de bug no `static/handler.go` do agnos não alcança
  projeto já inicializado. O caminho de restauro é apagar e reinicializar.
- `front-purge` apaga `routes/static/` e as rotas de página, edições do usuário
  incluídas. É deliberado — um purge que deixasse handlers órfãos importando `pageio`
  entregaria uma árvore que não compila. O conteúdo de `assets/frontend/` sobrevive.
- Editar `static/handler.go` é editar código de segurança (`safeSegments`, `contentTypeOf`).
- `pageio` gerado ocupa o nome `sandbox/internal/pageio/` no projeto; `utils/` fica livre
  para o código à mão do projeto, que é por isso que `templates.go` não mora lá.
