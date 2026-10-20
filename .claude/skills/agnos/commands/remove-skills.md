# `remove-skills`

Remove the agent skills and turn the skills extension off

```bash
agnos remove-skills [--help] [--path <path>] [--quiet]
```

Deletes the agent skill drop-skills turned on — the generator's directory under .agents/skills/ and .claude/skills/, then each of those left empty — and writes skills: false. Any other skill under the two roots, and the rest of .claude/, is left where it is.

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project-flags](project-flags.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project-flags](project-flags.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project-flags`](project-flags.md) | always |

```bash
agnos remove-skills
agnos remove-skills --path ./my-project
```

Docs · [every command](../SKILL.md)
