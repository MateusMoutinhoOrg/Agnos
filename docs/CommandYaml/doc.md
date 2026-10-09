# CommandYaml

`sandbox/internal/commands/[<folder>/]<name>/command.yaml` declares one command — the file is
what makes its directory a command, at any depth. `agnos build`
generates `generated.new.go` (the `api.Command` that lands in `Cli.Commands`) and `generated.input.go` (the `Input`
its handler is handed) from it. Grow it with `add-arg` / `add-flag` / `set-command` and their
`set-` / `remove-` pairs ([Workflow](../Workflow/doc.md#change-the-command-surface)), not by hand:
the editors re-render it with keys in alphabetical order and drop comments.

```yaml
priority: 100
args:
  - { id: Command, start: 0, end: 0, trigger: { type: equal, value: greet } }
  - { id: Name, start: 1, end: 1, required: true, description: who to greet }
flags:
  - { id: Times, keys: [--times, -t], type: integer, default: 1, min: 1 }
category: Demo
summary: Say hello
examples: ["greet bob -t 2"]
```

## Command keys

| Key | Effect |
|---|---|
| `priority` | **Required.** The rung it runs on; the chain runs lowest first. `add-command` writes `100`, `10` for a `--middleware` |
| `segments` | The segment count the command line has to have, `≥ 1`. Absent: any count |
| `strict` | `true` (default, omitted): every token of the line has to be read by it or by a middleware before it. `false`: a middleware |
| `args`, `flags` | What it reads; see below. At least one arg |
| `category`, `summary`, `description`, `examples` | The help screens and `docs/Commands` |
| `hidden` | Dropped from the listings, still dispatches |

## The command line

A command line is its **segments** — every token before the first flag (a token starting with
`-` that is not a number, so `-1` is a segment), plus every token after a bare `--` — and its
**flags**, everything between. `agnos route add foo --force -- a b` has the segments
`route add foo a b`; a segment after a flag is not read, so args come first. A flag takes its
value as the next token or after `=` (`--name bob`, `--name=bob`).

A line whose segments start with a command's verb but do not fit its args is a usage error
naming the arg (exit `2`), never `unknown command`.

## Arg keys

| Key | Effect |
|---|---|
| `id` | The `Input` field, an exported Go name, never `FullCommand` |
| `start`, `end` | The segments read, inclusive; `end: -1` is the last one |
| `type` | `string` (default), `integer`, `number`, `uuid`; anything but string reads one segment (`start == end`). A range binds `[]string` |
| `trigger` | `{type, value, negate, ignore-case}` (`values` for `one-of`), compared against the segments joined by a space. The verb is an arg with a trigger |
| `required` | A line that matches the command without it is a usage error. Excludes `default` |
| `default`, `description` | Bound when absent; help text |

## Flag keys

| Key | Effect |
|---|---|
| `id` | The `Input` field, as for an arg |
| `keys` | The spellings typed (`[--times, -t]`). Absent: `--<id in kebab-case>` |
| `type` | `string` (default), `integer`, `number`, `boolean` (presence, never required), `string-array`, `integer-array` (one element per occurrence) |
| `required`, `default` | As for an arg |
| `min`, `max`, `enum`, `pattern` | Checked on every value before the handler runs |
| `trigger` | What the value has to match for the command to run at all |

`trigger.type` is `equal`, `prefix` (word by word: `add` is `add` or `add …`, never `add-flag`;
an empty value matches every line), `text-prefix`, `suffix`, `regex` or `one-of`. Failing a
trigger, or an arg's type, is a non-match: the line is for another command. A missing `required`,
a value that will not convert or breaks its bounds, and — on a strict command — a token nobody
read are usage errors, answered by `sandbox/internal/cli/errors/handle_*.go` with exit `2`.

`--pattern` on `add-command` compiles to args: literal words are one arg with an `equal` trigger,
`{name}` / `{name:integer|number|uuid}` one segment, a last `{*rest}` `end: -1`; without it the
pattern fixes `segments`.

## The chain

`Main` walks `Cli.Commands` in `priority` order and runs every command whose `Matches`
matches. One answers by `response.SetStatus(code)` or `response.Printf(...)` (which answers
`ExitOk`); one that does neither declined, and the next runs — which is the whole of what a
middleware is. `response.Eprintf` and `response.Logf` answer nothing. Nothing answering is
`handle_not_found.go`: the general help on an empty line, `unknown command` otherwise.

```go
func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	if input.Name == "" {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "Name", "nobody to greet")
	}
	response.Printf("hello %s\n", input.Name)
	return nil
}
```

Every command of one line is handed the same `props`, typed by the project in
`sandbox/internal/commandprops/project.go` — written once by `agnos build`, then the project's.
`commandprops.go` beside it is generated: it embeds every struct of the package.
`help-flag` (priority 5) answers `<command> --help` unless that command declares `--help` itself.
