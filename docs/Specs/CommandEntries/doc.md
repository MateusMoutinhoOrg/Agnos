# CommandEntries Specification

## Description
Defines the required shape of `sandbox/internal/commands/<name>/entries.yaml` — the **declaration** of one command, from which `agnos build` generates the typed `Entries` struct and the dispatch arm. It is grown from the command line by `add-command`, `add-flag`, `add-arg` and `set-command`, which re-render it in canonical form; every key and its effect is in [EntriesYaml](/docs/EntriesYaml/doc.md).

### Rules
- The file is YAML with the keys `identifiers`, `category`, `help` (required), and `long-description`, `examples`, `hidden`, `flags`, `args` (optional). No other key.
- `identifiers` is a non-empty list; the first entry is the canonical name. `help` is one line. `long-description` is a block scalar.
- `flags` and `args` are **sequences** of objects, each with a `name`. Order is meaning: positional args bind by their written position; the generated struct lists fields in declaration order. The legacy mapping form is parsed but never written.
- A field's keys are `name`, `identifiers` (flags only), `type` (`string`, `boolean`, `int`, `float`; default `string`), `description`, `examples`, `default` (a string literal), `required`, `array`, `min`, `max` (numbers, `int`/`float` only). No other key.
- `required: true` is not written on a boolean or on a field with `default` — absence is already covered, and the parser drops it.
- An `array: true` arg is the last arg.
- A command taking a project directory declares a `path` flag with `default: "."`, never a positional. A command with progress output declares a boolean `quiet` flag with identifiers `--quiet` and `-q`.
- The file is edited through the CLI, not by hand; the exception is `help`'s declaration, which `agnos build` writes once and the project may then edit.

## Structure
1. **Command keys**: `identifiers`, `category`, `help`, then `long-description`, `examples`, `hidden` as needed.
2. **`args:`** *(optional)*: the positionals, in binding order.
3. **`flags:`** *(optional)*: the options.

The canonical rendering `agnos` writes puts keys in alphabetical order (`args`, `category`, `flags`, `help`, `identifiers`, …); a hand-written file in the order above is equivalent.

> **Note**: For a concrete example, refer to [sample.yaml](/docs/Specs/CommandEntries/sample.yaml).
