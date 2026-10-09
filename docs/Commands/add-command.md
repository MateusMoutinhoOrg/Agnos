# `add-command`

Scaffold a new command package in the project

```bash
agnos add-command <Name> --summary <summary> [--category <category>] [--trigger <trigger>] [--trigger-type <trigger-type>] [--trigger-negate] [--trigger-ignore-case] [--pattern <pattern>] [--middleware] [--priority <priority>] [--before <before>] [--after <after>] [--dir <dir>] [--help] [--path <path>] [--quiet]
```

Creates <name>/ under the folder of its --category in sandbox/internal/commands — commands/core/<name> for Core — or under the folder --dir names there, with a hand-written command.yaml and a stub handler.go, then runs build so generated.new.go, generated.input.go and the dispatch pick it up. A directory is a command by holding a command.yaml, at any depth; a name is unique across every folder. Its first arg answers to <name> on segment 0 — or to --trigger, or what --pattern compiles to; a --middleware runs in front of every command line and declines. Refuses a name another command already carries.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the name of the new command (e.g. my-feature) |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--summary` | string, required |  | one-line help text for the new command | — |
| `--category` | string |  | the heading the command is listed under in help and docs/Commands (defaults to Commands, or Middleware for a --middleware) | — |
| `--trigger` | string |  | what the words of the command line from segment 0 are compared against (defaults to the name, or every line for a --middleware); a one-of takes its values comma-separated | — |
| `--trigger-type` | string |  | how the trigger is compared: equal, prefix (word by word), text-prefix, suffix, regex or one-of — or starts-with, ends-with, exact, matches, any-of (defaults to equal, or prefix for a --middleware) | — |
| `--trigger-negate` | boolean |  | invert the trigger: the command runs for every line that does not match it | — |
| `--trigger-ignore-case` | boolean |  | compare the trigger without regard to case | — |
| `--pattern` | string |  | the command-line shape the command answers, compiled into its args: literal words, {name}, {name:integer\|number\|uuid} and a last {*rest}; without {*rest} the segment count is fixed (excludes --trigger and --trigger-type) | — |
| `--middleware` | boolean |  | declare a middleware: every command line unless --trigger says otherwise, not strict, priority 10, and a stub handler that hands the line on | — |
| `--priority` | integer | `-1` | the rung this command runs on when several match one command line, lowest first (defaults to 100, or 10 for a --middleware) | — |
| `--before` | string |  | land one rung below the command named, so it runs first (excludes --priority) | — |
| `--after` | string |  | land one rung above the command named, so it runs next (excludes --priority) | — |
| `--dir` | string |  | the folder under sandbox/internal/commands the command lands in, e.g. admin puts it in commands/admin/<name> (defaults to the folder of its category, commands/core/<name> for Core); a directory holding a command.yaml is a command, whatever folder holds it | — |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project-flags](project-flags.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project-flags](project-flags.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project-flags`](project-flags.md) | always |

```bash
agnos add-command my-feature --summary 'Do the feature'
agnos add-command route-add --pattern 'route add {name}' --summary 'Add a route'
agnos add-command profile --middleware --summary 'Read --profile in front of every command'
```

Cli · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
