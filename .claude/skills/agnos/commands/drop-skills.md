# `drop-skills`

Turn on the agent skills that teach an LLM to drive this generator

```bash
agnos drop-skills [--help] [--path <path>] [--quiet]
```

Turns the skills extension on and builds, which writes the agent skill into .agents/skills/<generator>/ and .claude/skills/<generator>/: a SKILL.md an agent loads by its description, plus one page per command of the generator. Every later build rewrites it from the binary running that build, so it follows the generator through each release. Run it again on a project that has it on and only the build runs. remove-skills takes it away.

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
agnos drop-skills
agnos drop-skills --path ./my-project
```

Docs · [every command](../SKILL.md)
