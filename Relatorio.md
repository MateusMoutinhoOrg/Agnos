Relatório de testes de estresse — agnos interview (v0.8.0)
Data: 2026-09-16 · Binário: /usr/local/bin/agnos v0.8.0 · Plataforma: macOS (darwin 22.6.0)

1. Metodologia
O modo entrevista foi exercitado de duas formas, porque ele se comporta de maneira diferente nas duas:

Modo	Como	Para que serviu
TTY real (pty)	driver em Python com pty.fork(), enviando setas/ENTER/q/Ctrl-C/Ctrl-D	reproduzir a experiência real do usuário final
stdin canalizado	printf ... | agnos interview (menus viram number [1]:)	varrer muitos caminhos rapidamente
Aplicações construídas de ponta a ponta, exclusivamente pela entrevista:

taskr — CLI puro: start → cli-init → add-command → add-arg → add-flag
api — servidor HTTP: start → server-init → add-route → add-body-field → front-init → add-page
legacy — projeto iniciado em pasta que já tinha go.mod
Projeto em escala — 30 comandos, para testar os seletores
Cenários destrutivos — server-purge, disable-extension, start --force
Entradas adversariais — unicode/emoji, espaços, !, nomes duplicados, índices fora de faixa, números inválidos
Resultado geral: a entrevista funciona. Os dois aplicativos principais compilaram e rodaram de verdade — o taskr gerou um binário com help, argumentos e flags; a api subiu um servidor que respondeu GET /health → {"status":"ok"}, GET /users → users called e serviu /home em HTML. O conceito está sólido. O que segue são os defeitos e atritos encontrados.

2. O que já está muito bom (não mexer)
Vale registrar, porque várias dessas coisas são raras em CLIs guiados:

Tela de confirmação antes de executar. O bloco NOTHING HAS RUN YET com o comando montado é excelente — ensina a CLI enquanto guia.
Editar uma resposta antes de rodar. O change --help (…) na tela de confirmação funciona perfeitamente e volta ao resumo. É a melhor interação da entrevista.
Perguntas condicionais. Ao escolher --default, a pergunta --required some (são incompatíveis). Ao escolher --type boolean, --required também some. Isso previne erros de verdade.
Campos repetíveis. --identifier — one more, or nothing to stop (2 so far) é claro e com contador.
Menu adaptativo. As áreas desligadas ficam fora do menu; o ★ sugere o próximo passo natural. front-init só aparece depois de server-init, o que está correto.
Validação numérica dos menus. 0, -1, 999, 1.5, 1abc são rejeitados; " 2 " é aceito com trim.
Terminal não fica corrompido. As sequências de esconder/mostrar cursor (\e[?25l / \e[?25h) estão balanceadas em todas as saídas testadas (q, Ctrl-C no menu, Ctrl-C no texto).
3. Bugs
🔴 Críticos
B1 — q é engolido como resposta de texto (o cabeçalho promete "go back")
O cabeçalho diz ↑↓ move enter choose q go back. Em qualquer pergunta de texto, q vira o valor literal.

name — the name of the new command (e.g. my-feature)  (text, required)
>  q
--help — one-line help text for the new command  (text, required)
>  q
Resultado: um comando chamado q, com help q. O usuário que quis voltar criou lixo no projeto sem perceber.

B2 — q no menu mata a entrevista inteira, não volta uma etapa
No meio do questionário do add-route, na pergunta --method:

interview ended: interview cancelled
O processo termina. Todas as respostas já dadas são perdidas e é preciso rodar agnos interview de novo. "go back" e "cancela tudo" são coisas muito diferentes — e não existe nenhuma forma de voltar uma pergunta.

B3 — O comando pré-visualizado não é o comando que roda
A entrevista promete: "you could have typed it yourself — that is all I am doing". Não é verdade quando há normalização.

│  $ agnos add-flag --command add "My Flag Name!"     ← o que foi mostrado
add-flag adding my-flag-name! to .../entries.yaml    ← o que aconteceu
O entries.yaml final tem name: my-flag-name! e identifiers: [--my-flag-name!]. A promessa central da tela é quebrada em silêncio.

B4 — Editar rotas embutidas: sucesso reportado, nada é escrito
As rotas health e static são geradas pelo template e recriadas a cada build. A entrevista as oferece no seletor --route (e health é o padrão, por ordem alfabética). Ao adicionar um campo nela:

$ agnos add-body-field --route health token
add-body-field adding token to sandbox/internal/routes/health/route.yaml
build started with path .
successfully rendered template
EXIT=0                                    ← sucesso
E o route.yaml continua idêntico. O campo desapareceu, sem aviso e com código de saída 0. Um usuário que aceite o padrão do seletor cai direto nisso.

B5 — server-purge destrói tudo com um ENTER, sem aviso, e deixa o projeto incoerente
Projeto com 4 rotas (health, home, static, users). A confirmação é exatamente a mesma tela usada para adicionar uma flag, com yes, run it já selecionado:

│  $ agnos server-purge
Run it?
❯ yes, run it            ← um ENTER apaga as 4 rotas
Depois disso, três problemas encadeados:

Assets órfãos — assets/frontend/pages/home.html e about.html continuam lá; só as rotas sumiram.
Estado impossível no cabeçalho — deps on cli on server off front on. O front depende do servidor.
O menu ainda oferece add-page e "Front System" — e usar isso quebra o projeto:
runtime go: `go mod tidy` failed:
go: finding module for package github.com/mateus/api/sandbox/internal/routeio
$ agnos verify
verify found 1 violation(s):
  - sandbox/internal/routes/about has no new.go
$ agnos build ; echo $?
1
A entrevista conduziu o usuário, por caminhos que ela mesma ofereceu, até um projeto que não compila — e ao reabrir a entrevista nesse projeto quebrado nada é dito: o menu aparece normalmente, sem nenhuma menção às violações.

B6 — start --force sobre um projeto existente o destrói (e a entrevista o oferece)
start está listado em "Core Commands" mesmo em projeto já inicializado, e --force é apresentado como um sim/não neutro ("Forces the creation of the project, overwriting existing files").

$ agnos start --project-name hijacked2 --force
sandbox/internal/commands/version/handler.go:11:11:
    sandbox.Deps undefined (type *api.Sandbox has no field or method Deps)
EXIT=1
Saiu com erro e mesmo assim reescreveu README.md, reescreveu AgnosConfig/project.yaml (name: "hijacked2") e derrubou a camada de deps. Todos os handlers de comando existentes pararam de compilar. Sem rollback, sem backup, sem atomicidade. O projeto anterior é irrecuperável.

🟠 Graves
B7 — --module é apresentado como opcional, mas é obrigatório (primeiro passo do primeiro uso)
Este é o primeiro erro que qualquer usuário novo encontra, na primeira tela.

--module — the go module path written into go.mod (required when the target dir has no go.mod yet)  (text, optional)
>
...
number [1]: the module flag (--module) is required when there is no go.mod in the path
✘ start exited 2
A entrevista já sabe se existe go.mod — ela imprime o estado da pasta no topo (project none here yet). Em pasta que já tem go.mod, o campo é de fato opcional e o fluxo funciona. Ou seja: a informação existe e não é usada.

B8 — Toda resposta é perdida quando a execução falha
Em todas as falhas observadas (B7, --position fora de faixa, nome duplicado, nome inválido), o comportamento é o mesmo: mensagem de erro e volta ao menu principal, do zero.

No add-flag isso significa 10 perguntas respondidas de novo por causa de um caractere. E o mais frustrante: a tela de confirmação já tem o mecanismo certo (change --position). Ele simplesmente não é oferecido depois do erro.

B9 — --position não é validado, embora a faixa válida seja conhecida
--position — zero-based index to insert the field at (defaults to the end)  (whole number, "-1" if you skip it)
> 1
...
--position 1 is out of range: this command has 0 flag(s), so the accepted range is 0 to 0
✘ add-flag exited 1
A própria mensagem de erro informa a faixa (0 to 0). Essa faixa poderia estar na pergunta, ou a pergunta poderia ser um seletor "no fim / antes de X / depois de Y". Combinado com B8, o usuário perde o questionário inteiro.

B10 — Validação de nomes inconsistente entre comandos irmãos
Entrada	add-command	add-flag
café-🚀-ação	❌ rejeitado com erro	✅ aceito
My Cool Command	⚠️ aceito com aviso (note: … normalized to "my-cool-command")	⚠️ aceito sem nenhum aviso
My Flag Name!	—	✅ aceito, vira --my-flag-name!
Consequências: o binário publicado ganha uma flag --café-🚀 e outra --my-flag-name! — esta última hostil ao shell (! dispara expansão de histórico em bash/zsh interativos). Além disso, a mensagem de erro do add-command diz "only letters, digits, spaces, dashes and underscores are allowed", mas café e ação são letras: a regra real é ASCII, e o texto não diz isso.

B11 — Nome duplicado só é detectado depois de todo o questionário
│  $ agnos add-command --help "Another add" --category Tasks add
add-command creating sandbox/internal/commands/add
file "sandbox/internal/commands/add/entries.yaml" already exists
✘ add-command exited 1
A entrevista lista os comandos existentes no seletor --command do add-flag — ela sabe quais nomes já existem. Não usa essa informação no add-command. E a mensagem é de desenvolvedor (caminho de arquivo), não de usuário ("já existe um comando chamado add").

B12 — Seletores longos: sem viewport, sem filtro, redesenho da lista inteira a cada tecla
Em projeto com 30 comandos, num terminal 80×24:

a lista sai com 30+ linhas, maior que a tela — o cabeçalho da pergunta rola para fora e o ❯ some de vista;
cada seta redesenha a lista inteira: 3 DOWN produziram ~120 linhas de scrollback;
não há busca/filtro — chegar em cmd9 exige mais de 20 DOWN;
a ordenação é alfabética pura, então cmd10…cmd19 vêm antes de cmd2.
Isso torna a entrevista pouco usável exatamente quando o projeto cresce.

🟡 Médios
B13 — Loop infinito de re-pergunta em campo obrigatório, sem saída indicada
--help — one-line description of the route  (text, required)
>   --help has to be answered
--help — one-line description of the route  (text, required)
>   --help has to be answered
... (indefinidamente)
Sempre a mesma mensagem, sem nunca dizer como sair. Com B1/B2, q também não resolve: ou vira o valor, ou mata tudo.

B14 — Padrão perigoso: um ENTER cego dispara mudanças estruturais pesadas
O item já selecionado (❯) quando o menu abre:

Situação	Padrão sob o ENTER	Por que é ruim
Depois de add-command	server-init	quem faz um CLI ganha um servidor HTTP inteiro
Depois de add-route	front-init	ganha camada HTML + deps embeddeps
Submenu Extensions	disable-extension → e o padrão dela é sandbox (o núcleo)	dois ENTER desligam a geração do núcleo
Submenu Server System	add-body-field → rota padrão health (embutida)	cai direto no B4
Submenu Cli System	add-arg	add-command seria o esperado
Note que depois de add-command e de add-route o ★ desaparece — a entrevista deixa de ter opinião justamente quando o padrão passa a ser perigoso.

B15 — Operações destrutivas misturadas às construtivas, em ordem alfabética
Cli System
  1) add-arg            3) add-flag      5) remove-arg      7) remove-flag
  2) add-command        4) cli-purge     6) remove-command  8) set-command
cli-purge ("Removes the CLI layer") fica entre add-flag e remove-arg, sem marcador, sem cor, sem agrupamento. O mesmo em Server System (server-purge no item 11) e em Core Commands, que traz publish — "Builds, compiles and publishes a release via gh", uma ação externa e irreversível — no item 4, com a mesma tela de confirmação genérica de qualquer outra.

B16 — Depois de add-command / add-route, o próximo passo natural não é sugerido
Acabou de declarar o primeiro comando? O caminho para dar a ele uma flag ou um argumento está enterrado em "Cli System". O mesmo para rotas: add-param, add-header, add-body-field ficam em "Server System". A entrevista guia muito bem até a primeira declaração e então abandona o usuário.

B17 — Seletores inconsistentes: uns têm descrição, outros não
--command (add-flag)          --route (add-body-field)      adapter (add-adapter)
  1) add — Adds a task…         1) health                     1) argvdeps
  2) help — Display help…       2) static                     2) embeddeps
  3) version — Print the…       3) users                      3) goimportsdeps
        ✅ com descrição            ❌ nomes nus                  ❌ nomes nus
No seletor de rotas faltam método, caminho e help — não dá para distinguir uma rota da outra, nem saber que health e static são embutidas. Na lista de adapters, reflectsort ou hashdeps não significam nada para o público-alvo declarado da tela.

B18 — --format é texto livre onde deveria ser um menu
--format — json-schema format for a string property: email, uuid, date-time or uri  (text, optional)
>
O conjunto de valores é fechado e está escrito ali na pergunta. --type e --method, que também são fechados, são menus. Mesmo caso em --min / --max, declarados como text sem validação numérica.

🟢 Menores
B19 — --quiet / -q praticamente não faz nada. Documentado como "Quiets the cli output"; o diff entre com e sem a flag é uma única linha (interview ended: …). O banner, o texto de introdução, os menus e o ruído de build continuam iguais.
B20 — Ctrl-D dá mensagem confusa. interview ended: no input left to answer with parece erro interno de script, não "você cancelou". Ctrl-C sai com um ^C seco, sem mensagem, enquanto q diz "interview cancelled" — três saídas, três comportamentos.
B21 — --path inexistente é aceito em silêncio. agnos interview --path /caminho/que/nao/existe mostra project none here yet — creating one is the first step, idêntico a uma pasta vazia válida. Um erro de digitação no caminho é indistinguível do caso normal.
B22 — Layout de largura fixa (~72 colunas), sem adaptação. A saída é byte a byte idêntica em terminal de 40 e de 200 colunas. E há linhas bem maiores que 80: a pergunta do --identifier tem 167 caracteres numa única linha, quebrando feio já no terminal padrão de 80 colunas.
B23 — Ruído de build exposto ao usuário final. successfully rendered template, runtime go: go mod tidy, build started with path . aparecem repetidos (o cli-init imprimiu esse bloco 4 vezes). Durante os ~40s de go mod tidy não há spinner nem indicação de progresso — parece travado.
B24 — change X não pré-preenche o valor atual. Ao corrigir um texto longo é preciso redigitar tudo.
B25 — Campos repetíveis mostram só o contador, não os valores já digitados, e não há como remover um item errado sem abandonar o questionário.
4. Jargão interno numa tela "feita para pessoas"
A descrição do comando diz: "The one screen of agnos made for a person rather than for a script". Mas os textos das perguntas são reaproveitados literalmente do help da CLI e vazam implementação:

Texto atual	Problema
Add a positional arg to a command's entries.yaml	entries.yaml é arquivo interno
the id the handler reads it back by, e.g. command.GetString("out-file")	código Go numa pergunta para não-programador
becomes the directory sandbox/internal/routes/<name> and its Go package	caminho interno
zero-based index to insert the field at + padrão -1	"zero-based" e padrão negativo se contradizem
nested with / for a sub-doc (e.g. PublicApi/api.AddDoc)	exemplo indecifrável
collect every occurrence into a []T field	sintaxe de genéricos Go
5. Prioridade sugerida
Corrigir primeiro (baixo custo, alto impacto)
B7 — deduzir a obrigatoriedade de --module a partir do go.mod já detectado. É o primeiro erro do primeiro uso.
B8 — em vez de voltar ao menu após falhar, reexibir a tela de confirmação com as respostas preservadas e o campo culpado destacado. O mecanismo (change X) já existe.
B1 + B2 — aceitar q (ou ESC) como "voltar uma pergunta" em todos os prompts, e reservar Ctrl-C para abortar. Se não for viável, corrigir o cabeçalho para dizer o que realmente acontece.
B14 — nunca deixar uma mudança estrutural ou destrutiva como item pré-selecionado. Quando não houver ★, o padrão deveria ser um item inerte (· voltar).
B3 — mostrar a normalização na pré-visualização: add-flag … "My Flag Name!" → criará a flag --my-flag-name!.
Em seguida
B5 + B6 — confirmação reforçada e específica para operações destrutivas, dizendo o que será perdido ("isto apaga 4 rotas: health, home, static, users"), com o padrão em "não". Tornar start/purge atômicos ou fazer backup.
B4 — marcar rotas/comandos embutidos como somente-leitura no seletor e recusar a edição com erro claro, em vez de reportar sucesso falso.
B12 — viewport com rolagem + filtro por digitação nos seletores, e repintura no lugar em vez de reimprimir a lista.
B10 + B11 — unificar a validação de nomes entre add-command, add-flag, add-arg, add-route, checando duplicatas na hora da pergunta.
B5 (menu) — derivar o menu do estado real e avisar quando verify acusar violações, em vez de oferecer caminhos que quebram o projeto.
Polimento
B16 — depois de add-command/add-route, sugerir add-flag/add-arg/add-param com ★.
B15 — separar visualmente as ações destrutivas (seção própria, cor, marcador).
B17 + B18 — descrição em todos os seletores; menu para todo conjunto fechado (--format).
Jargão — textos próprios para a entrevista, separados do help da CLI.
B22 + B23 — quebrar linhas conforme a largura do terminal; esconder o log de build atrás de um indicador de progresso.
B19, B20, B21, B24, B25 — ajustes pontuais.
6. Conclusão
A entrevista cumpre o que promete no essencial: foi possível construir, sem digitar um único comando agnos, um CLI funcional e uma API HTTP com rotas e páginas HTML — e ambos compilaram e rodaram. A tela de confirmação, as perguntas condicionais e o menu adaptativo são acertos de design reais.

Os problemas se concentram em três eixos:

A promessa da tela de confirmação não se sustenta quando há normalização (B3) ou quando a operação falha silenciosamente (B4).
A recuperação de erro é inexistente — qualquer falha joga o usuário de volta ao menu principal (B8), e navegar para trás ou é impossível (B2) ou corrompe dados (B1).
Os padrões não protegem o usuário — ENTER cego instala servidores, desliga o núcleo ou apaga a camada HTTP inteira (B14, B5), e a entrevista chega a oferecer caminhos que deixam o projeto sem compilar (B5) ou a destruir o projeto sem rollback (B6).
Resolvendo os cinco itens da primeira lista de prioridade, a entrevista passa de "funciona se você não errar" para "funciona mesmo quando você erra" — que é exatamente o público que ela declara servir.