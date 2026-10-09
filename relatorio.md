# Relatório de testes: backoffice, server e backups

Versão testada: agnos `v0.17.0`, commit `14019ec`, binário de bootstrap (`go build -o release/bootstrap.bin ./cmd/main`).

## Como foi testado

- Projeto novo: `agnos start -p testapp` + `agnos backoffice-init` (instala server, front e database). Compilado e executado com `start-server --insecure-http --addr 127.0.0.1:3999`.
- Testes de runtime com `curl` contra `/admin/*`, `/api/admin/*`, a rota `front`, `/openapi.json` e `/health`. IPs de origem diferentes simulados com `--interface 127.0.0.N`.
- Backups: create, download, upload, restore, manual (create-empty/add-file), optimize. Uma segunda instância (`proj2`, 3000 registros) foi usada para medir desempenho e interromper um restore.
- Comandos do agnos: `build` (idempotência), `verify`, `backoffice-purge` → `backoffice-init`, `run-examples --only` para `backoffice-init`, `backoffice-purge`, `server-init`, `front-init` e `database`. Todos os exemplos passaram.
- Revisão do código gerado: auth, throttle, tokens, usuários, snapshots, adapters (`nethttpserver`, `OpinionatedAgnosServer`, `OpinionatedAgnosFront`, `memoryratelimit`, `golangjwt`, `pbkdf2password`, `ziparchive`).

Os caminhos citados são os templates em `assets/`, onde a correção precisa ser feita. O código renderizado num projeto fica no mesmo caminho, sem o prefixo `assets/templates/backoffice/`.

**Confirmado** = reproduzido em runtime. **Código** = identificado na leitura do código, sem reprodução.

## Resumo

| # | Severidade | Problema | Status |
|---|---|---|---|
| 1 | Crítica | Esgotamento de memória sem autenticação pelas chaves do rate limiter | Confirmado |
| 2 | Crítica | Rate limit de login contornado por concorrência: força bruta e DoS de CPU | Confirmado |
| 3 | Alta | Restore interrompido por SIGINT/SIGTERM deixa os bancos pela metade | Confirmado |
| 4 | Alta | Restore que falha no meio apaga todos os bancos, sem rollback | Confirmado |
| 5 | Alta | Lockout de qualquer conta por terceiros, de qualquer IP | Confirmado |
| 6 | Alta | Restore e snapshot rodam com o tráfego ativo: dados inconsistentes | Código |
| 7 | Alta | TOCTOU entre restore e remove: um restore vazio apaga tudo | Código |
| 8 | Média | Restore ressuscita tokens revogados, sessões, usuários removidos e senhas antigas | Confirmado |
| 9 | Média | Cookie de sessão com `Path=/` e front sem headers: um XSS no app vira controle do admin | Confirmado |
| 10 | Média | Bind padrão em todas as interfaces; `--insecure-http` não emite aviso | Código |
| 11 | Média | Unicidade de username/email não é atômica | Confirmado |
| 12 | Média | Race do "último root": o sistema pode ficar com zero roots | Confirmado |
| 13 | Média | Backups sem criptografia, no mesmo disco, sem retenção nem agendamento; `.gitignore` incompleto | Código |
| 14 | Média | Backups não escalam: lentos, inflam o número de arquivos, tudo fica em memória | Confirmado |
| 15 | Média | `build` não é ponto fixo após qualquer `*-init` | Confirmado |
| 16 | Baixa | Lock de job só dentro do processo | Código |
| 17 | Baixa | `GET /admin/login` responde 404, mas a doc manda abrir essa URL | Confirmado |
| 18 | Baixa | `/openapi.json` público expõe toda a superfície de `/api/admin` | Confirmado |
| 19 | Baixa | Senha mínima de 8 caracteres sem nenhuma checagem; segredo validado só pelo tamanho | Confirmado |
| 20 | Baixa | `checkBody` aceita `Content-Type` ausente ou só com o prefixo certo | Confirmado |
| 21 | Baixa | Páginas usam `text/template` com escape manual em vez de `html/template` | Código |
| 22 | Baixa | Cada request da API faz 2 escritas em disco (last-used do token) | Código |
| 23 | Baixa | Allowlist de IP do token compara strings; regex de IPv6 aceita valores inválidos | Código |
| 24 | Baixa | Rate limit em memória: zera no restart e não é compartilhado entre instâncias | Código |
| 25 | Baixa | O front serve dotfiles embutidos (`.env`, `.git/...`) | Código |
| 26 | Baixa | Viewer cria token sem expiração e vê o email de todos os usuários | Confirmado |
| 27 | Processo | Nenhum teste exercita o comportamento HTTP nem os backups | — |

---

## 1. Crítica — Esgotamento de memória sem autenticação pelas chaves do rate limiter

- **Onde:** `assets/templates/backoffice/sandbox/internal/server/backoffice/backofficethrottle/backofficethrottle.go:35` (`loginKey`), `assets/adapter-catalog/memoryratelimit/.../memoryratelimit.go:42` (`hit`).
- **Causa:** a chave do contador por login é `"login:" + lower(trim(username))`, sem limite de tamanho. O body do login aceita até 1 MiB (`max-bytes: 1048576`). Cada tentativa que falha guarda o username inteiro como chave de map por 15 minutos. O `sweep` só remove janelas já expiradas.
- **Evidência:** 150 POSTs em `/admin/login` com username de ~1 MB, vindos de 7 IPs, fizeram o RSS subir de **17 MB para 336 MB**. Depois de 45 s ainda estava em 336 MB.
- **Impacto:** qualquer pessoa na internet derruba o servidor por OOM. O limite por IP não protege por causa do item 2.
- **Correção:** guardar o hash do login na chave (`sha256(lower(login))`), não o login. Truncar ou recusar logins acima de ~254 bytes antes de qualquer processamento. Baixar o `max-bytes` do login para ~4 KiB. Limitar o número de chaves no limiter (LRU).

## 2. Crítica — Rate limit de login contornado por concorrência (TOCTOU)

- **Onde:** `assets/templates/backoffice/sandbox/internal/routes/backoffice/backoffice_login/handler.go:22-32`. O mesmo padrão existe em `backoffice_api_token_auth/handler.go` (`TokenAllowed` / `TokenFailed`).
- **Causa:** o handler chama `LoginAllowed` (`Count`), depois roda o PBKDF2 (~centenas de ms), e só então `LoginFailed` (`Hit`). Requests em paralelo passam todos pela checagem antes de qualquer `Hit` ser registrado.
- **Evidência:**
  - 100 logins errados em paralelo para `admin`: **64 senhas checadas**. O limite documentado é 10 por login.
  - 300 em paralelo de um único IP: **154 checadas**. O limite é 20 por IP.
- **Impacto:**
  - Força bruta várias vezes acima do limite declarado.
  - DoS de CPU sem autenticação: cada request custa 600k iterações de PBKDF2. Com 4 cores, 300 requests ocuparam o servidor por ~5 s.
  - Alternar entre username e email dobra a cota por conta (são duas chaves diferentes).
- **Correção:** registrar o `Hit` **antes** de checar a senha e decidir pelo valor que ele retorna (reserva atômica), zerando no sucesso. Limitar logins concorrentes por IP (semáforo). Usar uma chave de throttle que não dependa de como o login foi escrito (o id do usuário resolvido, quando ele existe).

## 3. Alta — Restore interrompido por SIGINT/SIGTERM deixa os bancos pela metade

- **Onde:** `assets/templates/backoffice/sandbox/internal/snapshots/restore.go:69-82`, `assets/adapter-catalog/OpinionatedAgnosServer/.../main.go:55` (`OnInterrupt` → `Shutdown`).
- **Causa:** o restore roda numa goroutine. O graceful shutdown espera só os requests HTTP em andamento, não os jobs de backup. `main` retorna e o processo sai no meio do `RemoveDir`/`WriteFile`. `StartRecover` só marca como `failed` snapshots em `creating`; não existe recuperação de um restore interrompido.
- **Evidência:** `proj2`, com 3123 arquivos. SIGINT logo depois de `data/backofficedb` ser removido. Resultado: `backofficedb` com **0 arquivos** e `app` com **1223 de 3000**. O processo saiu sem logar erro.
- **Impacto:** um deploy ou restart do systemd durante um restore corrompe todos os bancos e apaga os usuários do backoffice. A recuperação exige criar um root pela CLI e restaurar o `pre-restore-*` manualmente.
- **Correção:**
  - Restaurar num diretório temporário e trocar de forma atômica (`rename`).
  - Gravar um marcador de "restore em andamento" que o boot detecte e use para reverter ao `pre-restore`.
  - Fazer o shutdown esperar o job terminar (ou cancelar antes do ponto de não retorno).

## 4. Alta — Restore que falha no meio apaga todos os bancos, sem rollback

- **Onde:** `restore.go:58-82`.
- **Causa:** a validação prévia confere só `SafePath` e se o blob existe. Ela não detecta um conflito arquivo/diretório (`db/a` e `db/a/b`), disco cheio ou erro de permissão. Os diretórios são removidos **antes** da escrita, e um erro no meio não desfaz nada.
- **Evidência:** upload de um zip com `data/app/zzz/a` e `data/app/zzz/a/b`. O upload foi aceito (`201`, `ready`). O restore falhou com `mkdir data/app/zzz/a: not a directory`. Sobraram só `data/app/zzz/a` e `data/backup`. **Todos os usuários e tokens sumiram**: a API passou a responder `401` e só foi possível voltar com `add-backoffice-user` pela CLI.
- **Impacto:** perda total dos dados e lockout do backoffice, causados por um arquivo malformado ou por uma falha comum de I/O.
- **Correção:** a mesma do item 3 (escrever no temporário e depois fazer swap). Rejeitar no `Import`/`Close` caminhos em que um é prefixo de diretório do outro. Em caso de erro, restaurar o `pre-restore` automaticamente.

## 5. Alta — Lockout de qualquer conta por terceiros

- **Onde:** `backofficethrottle.go:53-63`.
- **Causa:** o contador por login é global, contando todos os IPs: 10 falhas bloqueiam aquele login por 15 minutos.
- **Evidência:** depois das tentativas erradas, o login de `admin` **com a senha correta** respondeu `429`, também vindo de outro IP (`127.0.0.2`). Só o login pelo email passou, porque é outra chave.
- **Impacto:** um atacante anônimo, com 10 requests a cada 15 minutos, mantém um admin conhecido fora do sistema por tempo indefinido.
- **Correção:** não bloquear por conta sem considerar o IP. Opções: backoff progressivo por par (IP, conta); CAPTCHA ou atraso em vez de bloqueio; liberar IPs que já tiveram login bem-sucedido naquela conta.

## 6. Alta — Restore e snapshot rodam com o tráfego ativo: dados inconsistentes

- **Onde:** `snapshots.go:254-300` (`fill`, `collect`), `restore.go:69-82`.
- **Causa:**
  - Nada impede escritas nos bancos durante um snapshot. Ele lê arquivo por arquivo (~7 ms cada, ver item 14), então um registro escrito no meio fica com campos de instantes diferentes, e não há consistência entre tabelas.
  - Durante o restore, os requests continuam lendo um diretório vazio ou parcial (sessões somem, `FindBackofficeUserById` falha) e escrevendo em diretórios que serão sobrescritos.
- **Impacto:** backups corrompidos de forma silenciosa, e escritas perdidas ou misturadas durante o restore.
- **Correção:** um lock de escrita global (RW) que o restore tome em modo exclusivo e o snapshot em modo compartilhado. Ou responder `503` nas rotas de escrita durante um job, ou fazer o snapshot a partir de uma cópia consistente.

## 7. Alta — TOCTOU entre restore e remove: um restore vazio apaga tudo

- **Onde:** `restore.go:12-24` e `restore.go:52-56`.
- **Causa:** `StartRestore` lê o snapshot e confere o status **antes** do `acquire`. Um `Remove` que termine nesse intervalo deixa `ListSnapshotContents` vazio, e `restore` não recusa uma lista vazia: remove todos os diretórios e não escreve nada. Um snapshot `ready` criado com o diretório de dados vazio tem o mesmo efeito.
- **Impacto:** todos os bancos apagados. A janela é curta, mas o resultado é catastrófico.
- **Correção:** buscar e validar o snapshot de novo **depois** do `acquire`, e recusar restore de um snapshot sem arquivos (a regra que o `Close` já aplica).

## 8. Média — Restore ressuscita tokens revogados e credenciais antigas

- **Evidência:** snapshot criado, token do viewer revogado (`401`), snapshot restaurado: o mesmo token voltou a responder **`200`**.
- **Impacto:** depois de um incidente, restaurar um backup anterior reativa tokens vazados, sessões, usuários removidos e senhas trocadas. A doc só avisa que "o root pode ser deslogado".
- **Correção:** por padrão, não restaurar `backofficedb` (ou restaurar sem `api-token` e `session`). Exigir uma opção explícita para incluí-lo e avisar na UI.

## 9. Média — Cookie de sessão com `Path=/` e front sem headers

- **Onde:** `backofficeauth.go:185-191`; rota `front` (`OpinionatedAgnosFront`).
- **Causa:** o cookie `backoffice_session` vai para todo o site. As respostas do front e das rotas da aplicação não têm CSP, `X-Frame-Options` nem `nosniff` (confirmado em `/backoffice/backoffice.js` e `/`). A checagem de Origin aceita requests same-origin.
- **Impacto:** qualquer XSS numa página da aplicação, que é código do projeto e não do backoffice, faz requests autenticados a `/admin/root/*` com o cookie do admin (criar token, baixar backup, restaurar).
- **Correção:** `Path=/admin` no cookie, ou mover o backoffice para um subdomínio/origin separado. Headers de segurança básicos também no front.

## 10. Média — Bind padrão em todas as interfaces; `--insecure-http` não emite aviso

- **Onde:** `assets/templates/start_server_command.yaml` (padrão `3000:4000`, sem host).
- **Causa:** o padrão escuta em `0.0.0.0`. A doc ensina `start-server --insecure-http` para desenvolvimento local, e nesse modo o admin fica em HTTP plano, acessível pela rede, sem nenhum aviso. Só `--allow-x-forwarded-for` gera aviso.
- **Correção:** usar `127.0.0.1` como host padrão, ou emitir aviso quando `--insecure-http` é usado com bind em todas as interfaces.

## 11. Média — Unicidade de username/email não é atômica

- **Onde:** `backofficeusers.go:283` (`validate`) e `:163` (`AddBackofficeUser`).
- **Evidência:** 20 `add-backoffice-user` em paralelo com `username: "dup"` criaram **18 usuários** com o mesmo username.
- **Impacto:** `FindUserByUsernameOrEmail` devolve o primeiro que encontra, então as outras contas ficam inacessíveis e o login fica ambíguo. A mesma race existe no `Set` e no nome de token.
- **Correção:** serializar as escritas de usuário (mutex no processo, no mínimo), ou indexar username e email como `key` normalizada em minúsculas e deixar o store recusar o duplicado.

## 12. Média — Race do "último root"

- **Onde:** `backofficeusers.go:181-200`.
- **Evidência:** dois roots rebaixaram um ao outro em paralelo e a listagem com `role=root` voltou **`total: 0`**.
- **Impacto:** ninguém mais administra o sistema pela web. A recuperação só é possível pela CLI. `Remove` tem a mesma janela (dois roots removendo um ao outro).
- **Correção:** fazer a checagem e a escrita sob o mesmo lock.

## 13. Média — Backups sem criptografia, no mesmo disco, sem retenção nem agendamento

- **Fatos:**
  - O zip e o banco `backup` guardam em texto claro hashes de senha, sessões, hashes de token e todos os dados da aplicação.
  - Ficam dentro do mesmo `--database` que protegem: perder o disco perde os dois.
  - Não há backup agendado, política de retenção nem limite. Cada restore cria mais um `pre-restore-*` completo.
  - O `.gitignore` só cobre `/data/backofficedb` e `/data/backup`. Com `--database var/app`, banco de usuários e backups não são ignorados e podem ir para o git.
- **Correção:**
  - Criptografia opcional do zip, com chave vinda de variável de ambiente.
  - Destino externo (S3 ou diretório configurável) e retenção (`keep N`).
  - Comando ou cron de snapshot.
  - Avisar quando `--database` está fora de um caminho ignorado pelo git.

## 14. Média — Backups não escalam

- **Evidência** (`proj2`, 3000 registros de ~2 KB):
  - O snapshot levou **21,5 s** para 3123 arquivos, ~7 ms por arquivo. Um banco com 1M de registros (~5M de arquivos) levaria ~10 h.
  - O banco `backup` ficou com **67.066 arquivos / 388 MB** para 5 snapshots de 3123 arquivos. A deduplicação dos blobs funciona, mas cada `(path, sha)` vira um registro com vários arquivos, ~21 arquivos no store por arquivo de dado.
- **Código:**
  - `Export` monta o zip inteiro em memória (`archive.go:53`).
  - `Import` segura o body (até 256 MiB) mais o conteúdo descompactado (até 1 GiB) em RAM (`archive.go:121`).
  - O download fica sujeito ao `--write-timeout-ms` de 10 s e corta em silêncio no meio.
- **Correção:**
  - Gerar o zip em stream direto no response, e processar o upload em stream num arquivo temporário.
  - Guardar o manifesto do snapshot como um único blob (JSON com a lista `path → sha`) em vez de um registro por arquivo.

## 15. Média — `build` não é ponto fixo após qualquer `*-init` (bug do agnos)

- **Onde:** `sandbox/internal/actions/build/build_internal.go:173` (`CollectPublicApi`) roda antes de `:264-268`, que coleta as partes embutidas de `generated.config.go`/`generated.sandbox.go`.
- **Evidência:** num projeto limpo, `start` → `build` não deixou diff. Depois de cada `*-init`, um `build` adicional alterou arquivos:

  | Depois de | Arquivos alterados pelo `build` seguinte |
  |---|---|
  | `server-init` | `docs/PublicApi/api.{cli,clisandbox,sandbox,server,serversandbox}.md` |
  | `database-init` | `docs/PublicApi/api.{config,databaseconfig}.md` |
  | `backoffice-init` | `docs/PublicApi/api.config.md` (falta `BackofficeConfig`) |

  O build seguinte estabiliza.
- **Impacto:** viola a regra "deterministic and idempotent". A doc da API pública fica desatualizada logo depois de cada `-init`, e o check `build -q && git diff --quiet` falha num projeto recém-iniciado. Os exemplos não pegam isso.
- **Correção:** coletar a PublicApi depois de renderizar as partes do `sandbox/api`, ou ler o estado já preparado no StagedFS. Adicionar ao `run-examples` uma checagem de idempotência depois de cada `*-init`.

## 16. Baixa — Lock de job só dentro do processo

- `snapshots.jobs` é um `chan` em memória. Uma segunda instância do servidor sobre o mesmo `--database`, ou um comando da CLI, não respeita o lock. Dois restores em instâncias diferentes rodam ao mesmo tempo.
- **Correção:** um lockfile em `DataDir/backup/.job.lock` (flock).

## 17. Baixa — `GET /admin/login` responde 404

- A doc (`assets/doc-backoffice/docs/Backoffice/doc.md`) diz "Open `/admin/login`". Só existe a rota `POST /admin/login`, e `session-auth` exclui esse caminho, então o GET cai no front e responde `404`. A página de login só aparece ao acessar outra URL de `/admin`.
- **Correção:** adicionar uma rota `GET /admin/login` que renderize o login (ou redirecione para `/admin/home`).

## 18. Baixa — `/openapi.json` público expõe a superfície de `/api/admin`

- Sem autenticação, o `/openapi.json` lista 39 caminhos de `/admin`, inclusive todas as rotas `root/*` de backup.
- **Correção:** omitir as rotas de backoffice do spec público, ou proteger o endpoint.

## 19. Baixa — Política de senha e de segredo fraca

- `12345678` foi aceita (`MinPasswordLength = 8`), sem checagem contra senhas comuns.
- `ReadSecret` aceita `"aaaa…"` com 32 caracteres; só o tamanho é verificado.
- **Correção:** mínimo de 12 caracteres, lista de senhas comuns, e uma checagem mínima de entropia do segredo.

## 20. Baixa — `checkBody` aceita `Content-Type` ausente ou só com o prefixo certo

- **Onde:** `OpinionatedAgnosServer/run.go:204`.
- `content_type != "" && !HasPrefix(...)`: um POST de formulário sem `Content-Type` respondeu `201`. Um tipo como `application/x-www-form-urlencodedXYZ` também passaria.
- Hoje Origin e SameSite cobrem o CSRF, mas essa é uma camada de defesa a menos.
- **Correção:** exigir o header quando há body, e comparar o media type depois de cortar os parâmetros (`;`).

## 21. Baixa — `text/template` com escape manual nas páginas

- Todas as interpolações conferidas hoje usam `html` ou valores fixos (o teste com `x<script>` saiu escapado).
- Mas um campo novo sem `html`, ou um valor dentro de `href`, `style` ou JS, vira XSS, e nenhum `verify` impede.
- **Correção:** usar `html/template`, com escape por contexto, ou um check no `verify` que recuse `{{.X}}` sem `html` nos templates do backoffice.

## 22. Baixa — Duas escritas em disco por request da API

- `Resolve` (`backofficeapitokens.go:278-282`) grava `last-used-at` e `last-used-ip` em todo request: amplificação de I/O, sem atomicidade entre os dois campos e concorrência de escrita no mesmo registro.
- **Correção:** gravar no máximo uma vez por minuto, ou só quando o IP muda.

## 23. Baixa — Allowlist de IP do token

- A comparação é de strings (`accepts`, `:293`). Um IPv6 escrito de forma não canônica (`2001:0db8::1`) nunca casa com o `net.IP.String()` do adapter: o token fica inutilizável sem nenhum aviso.
- `ipv6Pattern` aceita valores inválidos (`:::::`).
- Não há suporte a CIDR.
- **Correção:** normalizar com um parser de IP no adapter e aceitar CIDR.

## 24. Baixa — Rate limit só em memória

- Um restart zera todos os contadores, e várias instâncias atrás de um load balancer multiplicam o limite. Esse limite está documentado, mas a doc não traz a consequência.
- **Correção:** documentar a consequência e oferecer um adapter persistente ou compartilhado.

## 25. Baixa — O front serve dotfiles

- `//go:embed all:*` inclui arquivos começando com `.`, e o `SafePath` do front só bloqueia `.` e `..`. Um `assets/front/.env` ou `.git/` copiado para a pasta seria servido.
- **Correção:** recusar segmentos que começam com `.` (exceto `.well-known`).

## 26. Baixa — Permissões do viewer

- O viewer cria tokens `never` (sem expiração), sem reautenticação, e vê username e email de todos os usuários (`/api/admin/list-backoffice-users`).
- **Correção:** prazo máximo de token para não-root, pedir a senha de novo ao criar token, e limitar a listagem do viewer.

## 27. Processo — Nenhum teste do comportamento

- Os exemplos `backoffice-init` e `backoffice-purge` conferem só a árvore gerada (34–36 arquivos). Nenhum sobe o servidor nem exercita login, throttle, tokens ou backups.
- Todos os itens 1–8 e 11–12 passariam por qualquer CI atual.
- **Correção:** um exemplo que compile o projeto, suba o server numa porta efêmera e rode um roteiro HTTP com golden: login, 429, tokens, create/restore/upload e o restore com falha do item 4.

---

## O que foi verificado e está correto

- Senhas com PBKDF2-SHA256, 600k iterações e salt próprio, comparadas em tempo constante; arquivos do store com permissão `0600` e diretórios `0700`.
- O JWT aceita só HS256 e exige `exp`. A sessão fica ligada ao IP: um cookie usado de outro IP recebe `401`. Logout e troca de senha encerram as sessões.
- O cliente não consegue forjar `X-Client-Ip` (o adapter sobrescreve). `X-Forwarded-For` só é lido com a flag.
- CSRF: `Origin` de outro host recebe `403`, `text/plain` recebe `415`, e o cookie tem `SameSite=Strict` e `HttpOnly`.
- Os guards root/viewer funcionam nas páginas e na API (`403`). Um viewer não consegue revogar o token de outro usuário (responde `not-found`).
- Path traversal bloqueado em `add-backup-file` (`..`, `%2e%2e`, `%2f`, `\`, `%00`, `backup/` em qualquer caixa), no upload de zip e no front (FS embutido).
- O nome do snapshot é validado por regex, e o filename do download de arquivo é sanitizado.
- O handler de 500 não vaza `Cause` para o cliente.
- O `ziparchive` limita o tamanho descompactado real, não o declarado no header.
- `backoffice-purge` → `build` → `backoffice-init` compila e estabiliza. `verify` passa. Todos os `run-examples` testados passaram.
