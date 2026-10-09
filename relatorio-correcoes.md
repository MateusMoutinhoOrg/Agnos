# Relatório de correções

Base: `relatorio.md` (agnos `v0.17.0`, commit `14019ec`). As correções foram feitas nos templates e catálogos em `assets/`, que é de onde os projetos são gerados, e espelhadas nas cópias que este repositório usa (`sandbox/deps/{serverdeps,embeddeps}`, `adapters/impls/{nethttpserver,goembed}`). Não há serviço externo: tudo usa Go nativo (canais como lock, `html/template`, parser de IP próprio, `sync` no adapter do rate limiter).

**Validado em runtime** = reproduzido num projeto novo (`start` + `backoffice-init`), com o servidor rodando e `curl`. **Código** = implementado e compilado, sem cenário de runtime dedicado.

Um projeto gerado antes desta versão não recebe as correções do backoffice sozinho: os arquivos de `backoffice-init` são escritos uma vez e passam a ser do projeto. Rodar `backoffice-init` de novo reinstala os deps e adapters do catálogo (`ratelimitdeps`, `serverdeps`, `embeddeps`, etc.) mas mantém os arquivos do backoffice já existentes.

## Resumo

| # | Severidade | Status | Validação |
|---|---|---|---|
| 1 | Crítica | Corrigido | Runtime |
| 2 | Crítica | Corrigido | Runtime |
| 3 | Alta | Corrigido | Runtime |
| 4 | Alta | Corrigido | Runtime |
| 5 | Alta | Corrigido, com limitação | Runtime |
| 6 | Alta | Corrigido, com limitação | Runtime |
| 7 | Alta | Corrigido | Código |
| 8 | Média | Corrigido | Runtime |
| 9 | Média | Corrigido | Runtime |
| 10 | Média | Corrigido (aviso) | Código |
| 11 | Média | Corrigido | Runtime |
| 12 | Média | Corrigido | Runtime |
| 13 | Média | Parcial | Código |
| 14 | Média | **Não corrigido** | — |
| 15 | Média | Corrigido | Runtime |
| 16 | Baixa | **Não corrigido** | — |
| 17 | Baixa | Corrigido | Runtime |
| 18 | Baixa | Corrigido | Runtime |
| 19 | Baixa | Corrigido | Código |
| 20 | Baixa | Corrigido | Runtime |
| 21 | Baixa | Corrigido | Código |
| 22 | Baixa | Corrigido | Código |
| 23 | Baixa | Corrigido | Teste do parser |
| 24 | Baixa | Parcial (documentado) | — |
| 25 | Baixa | Corrigido | Runtime |
| 26 | Baixa | Corrigido | Runtime (parcial) |
| 27 | Processo | Corrigido | Novo exemplo |

---

## 1. Esgotamento de memória pelas chaves do rate limiter — corrigido

- A chave por login não guarda mais o que o cliente digitou: é `user:<id>` quando o login resolve um usuário, ou `login:<sha256(lower(trim(login)))>` quando não resolve. Tamanho fixo.
- `POST /admin/login`: `max-bytes` caiu de 1 MiB para 4 KiB; o form-schema limita `username` a 254 e `password` a 1024. O handler também recusa sem checar a senha um login acima de 254 bytes ou uma senha acima de 1024.
- `memoryratelimit` ganhou um teto de 65536 chaves. Passando dele (depois do sweep das janelas expiradas), a janela aberta há mais tempo é descartada.
- **Teste:** um login de 1 MB respondeu `413`, sem chegar ao limiter.
- Arquivos: `backofficethrottle.go`, `backoffice_login/{handler.go,route.yaml}`, `assets/adapter-catalog/memoryratelimit/…/memoryratelimit.go`.

## 2. Rate limit contornado por concorrência — corrigido

- O contrato `ratelimitdeps` ganhou `Undo(key)`. O login agora conta a tentativa **antes** de checar a senha (`Hit` atômico, decidindo pelo valor que ele retorna) e só a desconta (`Undo`) quando o login dá certo. O mesmo vale para o token da API (`ReserveToken` / `TokenSucceeded`).
- Username e email do mesmo usuário caem na mesma chave (`user:<id>`), então alternar entre os dois não dobra mais a cota.
- `CheckPassword` roda no máximo 4 PBKDF2 ao mesmo tempo (semáforo com canal). O resto espera na fila, então uma enxurrada de logins não toma todos os cores.
- **Teste:** 100 logins errados em paralelo para `admin` → `10×401`, `90×429` (antes: 64 senhas checadas). 300 em paralelo de um IP com usernames diferentes → `20×401`, `280×429` em 1,4 s. 50 tokens falsos em paralelo → `20×401`, `30×429`.

## 3. Restore interrompido por SIGINT/SIGTERM — corrigido

- O snapshot `pre-restore-*` virou o journal do restore. Ele é criado, marcado com o novo status `rollback` e só então os diretórios são apagados. No fim do restore volta a `ready`.
- No boot, `StartRecover` procura snapshots em `rollback` e os restaura antes de qualquer outra coisa. Enquanto isso, o middleware de manutenção responde `503` a tudo.
- Funciona para SIGINT, SIGTERM, SIGKILL, crash e queda de energia, porque não depende de o processo terminar de forma limpa.
- **Teste:** 3000 arquivos, SIGINT no meio da escrita (2299/3000 gravados). No restart, o log mostrou `the last run stopped in the middle of a restore; data was put back…`, e o dado voltou a 3000 arquivos, idêntico ao estado anterior ao restore. Os usuários continuaram válidos.
- **Limitação:** o shutdown não espera o job terminar, porque o servidor não oferece um gancho para isso sem mudar o lib e o handler do projeto. O journal cobre o caso. Se o processo morrer entre o fim da escrita e a marcação `ready`, o boot desfaz um restore que tinha terminado: o resultado continua consistente, mas o restore precisa ser repetido.

## 4. Restore que falha no meio apaga tudo — corrigido

- A validação prévia passou a recusar caminhos em que um é diretório do outro (`Conflict`). A mesma regra vale no upload (`Import`) e no `close-backup` (novo outcome `conflict`, `400`).
- Se a escrita falha (disco cheio, permissão, um arquivo onde deveria haver uma pasta), o `pre-restore` é restaurado na hora.
- **Teste:** um zip com `data/app/zzz/a` e `data/app/zzz/a/b` → `400` no upload; o mesmo par num snapshot manual → `400` no close. Um restore que falhou com `mkdir data/qqq: not a directory` foi revertido: `data/app` intacto, `/api/admin/me` `200`, log `data was put back as it was`.

## 5. Lockout de qualquer conta por terceiros — corrigido, com limitação

- Um IP que fez login com sucesso numa conta nos últimos 30 dias fica confiável para ela e passa a ter um contador próprio (par conta+IP, 10 falhas). As falhas vindas de outros IPs não o bloqueiam mais.
- O limite global por conta continua valendo para IPs desconhecidos, porque é ele que segura a força bruta distribuída.
- **Teste:** login legítimo de `.5`; 40 tentativas de `.1`–`.4` (`10×401`, `30×429`); `.5` de novo com a senha certa → `303`. Um IP novo (`.6`) → `429`.
- **Limitação:** a confiança fica em memória, então um restart a esquece. Um admin entrando de um IP novo durante um ataque ainda é bloqueado. Não foi usado CAPTCHA, por instrução.

## 6. Restore e snapshot com tráfego ativo — corrigido, com limitação

- Novo middleware `backoffice-maintenance` (prioridade 6, todo caminho):
  - Durante um restore ou rollback: `503` + `Retry-After: 5` em todo request.
  - Durante um snapshot: `503` em todo método que não seja GET, HEAD ou OPTIONS.
- Para o middleware saber o método, o adapter `nethttpserver` passou a expor o header `X-Request-Method`. Como `X-Client-Ip`, ele nunca é lido do cliente.
- **Teste:** durante o snapshot, `GET /health` → `200` e `POST /admin/login` → `503`. Durante o restore, `GET /health` → `503`.
- **Limitação:** um request que já tinha passado do middleware quando o job começou não é interrompido (não existe fase "after" na cadeia para contar requests em andamento). Escritas feitas pela CLI também não são bloqueadas.

## 7. TOCTOU entre restore e remove — corrigido

- `StartRestore` só lê e valida o snapshot **depois** do `acquire`.
- Um restore sem arquivos é recusado (`the snapshot holds no file to restore, nothing was restored`).
- `Remove` recusa um snapshot em `rollback` (`400` / notice `not-ready`).

## 8. Restore ressuscita tokens e credenciais — corrigido

- Por padrão, o restore não toca `backofficedb`. Para incluí-lo é preciso pedir: `include-backoffice: true` no JSON ou o checkbox "users & tokens" na página. A confirmação da página avisa o risco.
- **Teste:** token `t2` revogado depois do snapshot → restore padrão → `401` (continua revogado). Com `include-backoffice: true` → `200`, como esperado quando a opção é pedida.

## 9. Cookie com `Path=/` e front sem headers — corrigido

- O cookie de sessão agora usa `Path=/admin`. As rotas `/api/admin` usam token, não cookie.
- `backoffice-security-headers` passou a rodar em todo caminho: os headers completos em `/admin` e `/api/admin`, e `nosniff`, `X-Frame-Options: SAMEORIGIN` e `Referrer-Policy` nos demais. Não há CSP fora do admin: isso é decisão da aplicação.

## 10. Bind em todas as interfaces com `--insecure-http` — corrigido (aviso)

- `backoffice-start-server` avisa quando `--insecure-http` é usado sem host em `--addr`. A doc passou a mostrar `--addr 127.0.0.1:3000`.
- O padrão `3000:4000` não foi trocado para `127.0.0.1`: isso quebraria todo projeto server que roda em container ou atrás de um proxy em outra máquina.

## 11. Unicidade de username/email não atômica — corrigido

- Add, Set e Remove de usuário rodam sob um lock (canal) do processo. A checagem de unicidade e a escrita acontecem juntas.
- Contra outro processo (um `add-backoffice-user` rodando junto do servidor), há uma checagem depois do insert: entre dois usuários com o mesmo nome ou email, o de id maior se remove e é recusado.
- O nome de token também passou a ser checado e gravado sob um lock.
- **Teste:** 20 `add-backoffice-user --username dup` em paralelo pela CLI (processos separados) → 1 usuário criado e 19 recusados. Antes, eram 18 duplicados.

## 12. Race do "último root" — corrigido

- A contagem de roots e a escrita acontecem sob o mesmo lock.
- Set e Remove conferem, dentro do lock, que quem age ainda é root (novo notice `not-root`, `403` na API). Dois roots rebaixando ou removendo um ao outro em paralelo não zeram mais os roots.
- **Teste:** dois roots se rebaixando em paralelo pela página → um `303` e um `403`, e sobra um root. Antes, sobravam zero.

## 13. Backups — parcial

- **Feito:**
  - O `start-server` avisa quando o `.gitignore` não cobre os stores sob o `--database` da execução.
  - Retenção: o restore mantém só os 5 `pre-restore-*` mais novos.
  - A doc diz explicitamente que os backups não são criptografados, ficam no mesmo disco e não são agendados.
- **Não feito:** criptografia do zip, destino externo e agendamento. São features novas (chave, formato, destino), não correções pontuais.

## 14. Backups não escalam — não corrigido

- Gerar o zip em stream e guardar o manifesto como um único blob exigem mudar `archivedeps` (hoje só em memória), o schema do banco `backup` e o formato dos snapshots já existentes.
- Ficou documentado na página Backups: o zip é montado inteiro em memória e o write timeout corta downloads grandes.

## 15. `build` não é ponto fixo após `*-init` — corrigido

- Causa: `CollectPublicApi` lia `sandbox/api` antes de os agregados (`generated.config.go`, `generated.sandbox.go`) serem renderizados no mesmo build.
- Correção: `GenerateApiSources` renderiza todo arquivo de `sandbox/api` dos grupos ativos antes da coleta, e a coleta inclui os arquivos recém-renderizados.
- **Teste:** num projeto novo, `server-init`, `database-init` e `backoffice-init`, cada um seguido de `build`, não deixaram diff. Neste repo, `build -q && git diff --quiet` continua idempotente.

## 16. Lock de job só dentro do processo — não corrigido

- Exige um primitivo de lock entre processos (flock) no contrato `iodeps`, com implementação por sistema operacional. Fica documentado na página Backups.

## 17. `GET /admin/login` → 404 — corrigido

- Nova rota `backoffice-login-page` (GET `/admin/login`) que responde o formulário com `200`.

## 18. `/openapi.json` expõe `/api/admin` — corrigido

- Nova chave `private: true` em `route.yaml`: a rota sai do `/openapi.json` e continua em `docs/Routes`. É editada por `set-route --private` / `--public`, aparece em `show-route` e está documentada em RouteYaml.
- Todas as rotas não-middleware do backoffice vêm com `private: true`.

## 19. Política de senha e segredo fraca — corrigido

- Senha: 12 a 1024 caracteres, pelo menos 5 caracteres distintos, fora de uma lista de senhas comuns e diferente do username e do email.
- Segredo: além dos 32 caracteres, exige pelo menos 10 caracteres distintos. Um segredo `aaaa…` é recusado.

## 20. `checkBody` aceita `Content-Type` ausente ou por prefixo — corrigido

- O media type é comparado sem os parâmetros (`;charset=…`) e sem diferenciar maiúsculas, nunca por prefixo.
- Um body sem `Content-Type` (`Content-Length > 0` ou `Transfer-Encoding`) recebe `415`. Um request sem body continua aceito em rotas de body opcional.

## 21. `text/template` com escape manual — corrigido

- O contrato `embeddeps` ganhou `RenderHTMLTemplate` (`html/template`). As páginas do backoffice passaram a usá-lo, com escape por contexto.
- Os `{{html …}}` existentes continuam funcionando sem escapar duas vezes (validado executando os 8 templates).

## 22. Duas escritas em disco por request da API — corrigido

- `last-used-ip` só é gravado quando o IP muda, e `last-used-at` no máximo uma vez por minuto (`LastUsedResolution`).

## 23. Allowlist de IP do token — corrigido

- Parser próprio de IPv4/IPv6 no sandbox (`ParseIp`, `FormatIp`, sem `net`), com suporte a CIDR (`IpMatches`).
- As entradas são gravadas na forma canônica e comparadas como endereços: `2001:0DB8::1` casa com `2001:db8::1`, `1.2.3.4` com `::ffff:1.2.3.4`. Valores como `:::::` são recusados.
- **Teste:** 20 endereços comparados com `net.ParseIP`, todos com o mesmo resultado; CIDR v4 e v6 conferidos.

## 24. Rate limit só em memória — parcial

- A consequência (restart zera os contadores, várias instâncias multiplicam o limite) está documentada no contrato `ratelimitdeps` e na página Backoffice.
- Não há adapter persistente ou compartilhado.

## 25. Front serve dotfiles — corrigido

- `SafePath` do `OpinionatedAgnosFront` recusa qualquer segmento que comece com `.`, exceto `.well-known` na raiz.

## 26. Permissões do viewer — corrigido

- Criar token pede a senha de novo, contada no mesmo limite do login.
- Para quem não é root, o token vale no máximo 90 dias (`never` é recusado).
- **Teste:** criar token sem senha → `400`.
- O viewer vê o email dos outros usuários mascarado (`a***@example.com`), e a busca não casa pelo email dos outros. Vale na página, em `list-backoffice-users` e em `get-backoffice-user`.

## 27. Nenhum teste de comportamento — corrigido

- Novo exemplo `examples/cli/backoffice-http`: gera o projeto, compila, sobe o servidor numa porta livre (`127.0.0.1:39300:39399`) e roda um roteiro `curl` com golden.
- O roteiro cobre `GET /admin/login`, o path do cookie, o token sem senha, 30 logins simultâneos, o body grande, o body sem Content-Type, `/.env`, `/openapi.json`, `nosniff`, o restore sem `backofficedb`, o close com conflito e o rollback de um restore que falha.

## Validação final

- `./release/bootstrap.bin build`: compila, e `build -q && git diff --quiet` é idempotente.
- `./release/bootstrap.bin verify`: passou.
- `./release/bootstrap.bin run-examples`: **103 exemplos, 43 cross-checks, 0 falhas**, incluindo o novo `backoffice-http`. A saída gravada no golden dele:

```
GET /admin/login: 200
session cookie: Path=/admin
token without password: 400
GET /api/admin/me: 200
30 wrong sign-ins at once: 401 x10 429 x20
a 64 KiB sign-in: 413
a json body without Content-Type: 415
GET /.env: 404
/admin in /openapi.json: 0
nosniff on /: 1
after restore, the revoked token: 401
after restore, data/app: 1 2 3
close of a file under a file: 400
after a failed restore, GET /api/admin/me: 200
after a failed restore, data/app: 1 2 3
after a failed restore, the log: 1
```

- Goldens regravados com `update-example`:
  - `backoffice-init`, `cli-purge` e `front-init`: mudaram por causa das correções. O `cli-purge` reflete o item 15: a PublicApi agora sai certa no mesmo build.
  - `build`, `start`, `remove-doc` e `interview`: **já falhavam no commit base**. O bump para `v0.17.0` não tinha regravado as páginas `Requirements` nem o banner do `interview`. Confirmado rodando esses exemplos numa worktree do commit `2584aa7`.

## Mudanças de contrato (aditivas)

| Contrato | Adição |
|---|---|
| `ratelimitdeps.Contract` | `Undo(key)` |
| `embeddeps.Contract` | `RenderHTMLTemplate(path, vars)` |
| `serverdeps` (`nethttpserver`) | header `X-Request-Method`, nunca lido do cliente |
| `route.yaml` | chave `private`; `set-route --private` / `--public` |

Nenhum campo, comando ou chave existente foi removido ou renomeado. Um adapter próprio (fora do catálogo) que preencha `ratelimitdeps` ou `embeddeps` precisa preencher o campo novo, ou o `verify` acusa.
