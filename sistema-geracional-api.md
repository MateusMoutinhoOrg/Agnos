# Sistema geracional de api

## Problema

Hoje `sandbox/api/sandbox.go` (template `assets/sandbox/sandbox/api/sandbox.go`) é gerado misturando
fontes diferentes: embute só os `usersandbox*.go` e escreve os campos das extensões (`Cli`, `Actions`,
`Config`, …) direto no template. É difícil de ler e as extensões podem colidir.

## Regra

`api.Sandbox` embute **toda struct `<X>Sandbox` declarada num arquivo `sandbox/api/<x>sandbox.go`**.

| Arquivo | Struct | Dono |
|---|---|---|
| `usersandbox.go` | `UserSandbox` | o projeto |
| `clisandbox.go` | `CliSandbox` | extensão `sandbox-cli` |
| `serversandbox.go` | `ServerSandbox` | extensão `sandbox-server` |
| `backofficesandbox.go` | `BackofficeSandbox` | extensão `sandbox-backoffice` |

- Cada extensão escreve **o seu** arquivo; nenhuma edita o de outra → sem colisão.
- O usuário estende criando outro `<x>sandbox.go` com a struct `<X>Sandbox`; o `build` embute sozinho.

## Config

Mesma regra: `api.Config` embute toda struct `<X>Config` de `sandbox/api/<x>config.go`
(`userconfig.go`, `serverconfig.go`, …).

## Onde mexer

- `sandbox/internal/utils/embedded_structs.go` — `CollectEmbeddedStructs` já faz o parse; trocar o
  filtro `usersandbox*` / `userconfig*` por "termina com `sandbox.go`" / "`config.go`" (exceto os
  próprios `sandbox.go` / `config.go`).
- `assets/sandbox/sandbox/api/sandbox.go` e o template de `config.go` — renderizar só os embeds.
