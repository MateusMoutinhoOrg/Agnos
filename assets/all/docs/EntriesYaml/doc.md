# EntriesYaml

`sandbox/internal/commands/<name>/entries.yaml` declares one command. `agnos build` generates
`entries.go` (the `Entries` struct) and a dispatch arm from it. Grow it with
`add-flag` / `add-arg` / `set-command` ([Workflow](../Workflow/doc.md#change-the-command-surface)),
not by hand: the editors re-render it with keys in alphabetical order and drop comments.

```yaml
identifiers: ["greet", "hello"]
category: Demo
help: Say hello
long-description: |
  Greets someone.
examples: ["greet -n bob 2"]
hidden: false
flags:
  - name: name
    identifiers: ["--name", "-n"]
    type: string
    required: true
    description: who to greet
args:
  - name: times
    type: int
    default: "1"
    min: 1
```

## Command keys

| Key | Effect |
|---|---|
| `identifiers` | Verbs the command answers to; the first is canonical in `help` |
| `category` | Heading in the general help |
| `help` | One-line description |
| `long-description` | Paragraph under `help <command>` |
| `examples` | Lines printed as `$ {{.Name}} <example>` |
| `hidden` | Dropped from the general listing, still dispatches |
| `flags`, `args` | Sequences of fields. Flags are read first in any position; args bind by written order |

## Field keys

| Key | Effect |
|---|---|
| `name` | Go field name (`out-file` -> `OutFile`). For a flag, defaults to the first identifier |
| `identifiers` | Flag spellings (`--name`, `-n`). Flags only. Default `--<name>` |
| `type` | `string` (default), `boolean` (presence only, never required), `int`, `float` |
| `description`, `examples` | Help text |
| `default` | Literal assigned when absent (string, converted to the type). Excludes `required` |
| `required` | Absence is a usage error. Refused on boolean or with `default` |
| `array` | Every occurrence collected into `[]T`. An array arg must be last |
| `min`, `max` | Bounds for `int`/`float`, checked before the handler runs |

Go types: `string`/`[]string`, `bool`, `int`/`[]int`, `float64`/`[]float64`.

A boolean flag named `quiet` is special: the dispatch replaces `deps.Std.Log` with a no-op
before the handler runs.

## Dispatch

`CliMain(args)` matches `args[0]` against every command's identifiers — no match, and an empty
command line, exit `2` with the general help. The arm then reads each declared flag anywhere on
the line, assigns defaults, converts and range-checks numbers, and drains the positionals in
order. An unread `-`-prefixed argument, a leftover positional, a missing required field or a
value out of range all exit `2` before `CommandHandler` runs — which is why a handler never
returns `api.ExitUsage` ([Rules](../Rules/doc.md#exit-codes)).
