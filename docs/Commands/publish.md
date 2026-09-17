# `publish`

Builds, compiles and publishes a release via gh

```bash
agnos publish [--path <path>] [--release-name <release_name>] [--draft] [--target <target>] [--publisher <publisher>]
```

Runs build, then compile (every target by default), and publishes every file of release/ as a gh release named --release-name, defaulting to the version in AgnosConfig/project.yaml.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path`, `-p` | string | `.` | The directory holding the project (defaults to the current directory) |
| `--release-name`, `-rn` | string |  | The name of the release |
| `--draft` | boolean |  | Create a draft release |
| `--target`, `-t` | string | `all` | The target to compile for (defaults to all) |
| `--publisher`, `-pub` | string | `gh` | The publisher to use (defaults to gh) |

Core Commands · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
