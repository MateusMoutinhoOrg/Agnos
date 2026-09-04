# EntriesYaml

`sandbox/internal/commands/<name>/entries.yaml` declares one command. `agnos build` generates `entries.go` (the `Entries` struct) and a dispatch arm from it. Grow it with `add-flag`/`add-arg`/`set-command`, not by hand: the editors re-render it with keys in alphabetical order and drop comments.

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
| `examples` | Lines printed as `$ <binary> <example>` |
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

A boolean flag named `quiet` is special: the dispatch replaces `deps.Std.Log` with a no-op before the handler runs.
