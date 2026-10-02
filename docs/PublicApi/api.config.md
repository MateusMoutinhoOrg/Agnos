# `sandbox/api/config.go`

## `Config`

Config is what the project knows about itself: the values of <ProjectName>Config/project.yaml, rendered into the sandbox by the build that read them. It is a field of the Sandbox like any other contract, so a caller may replace it — a test that runs the cli under another name, say — and every reader of it follows.

| Field | Type | Description |
| --- | --- | --- |
| `UserConfig` | `UserConfig` | UserConfig is a part of the Config, declared in sandbox/api/userconfig.go. Embedded, so each of its fields is read as sandbox.Config.<Field>. Every struct of a sandbox/api/<x>config.go file is one, so a mechanic adds its own part beside the project's rather than editing it. |
| `ProjectName` | `string` | ProjectName is the project's name, title-cased. It prefixes the <ProjectName>Config/ directory that holds every declaration, and lower-cased it is the name the cli answers to. |
| `Version` | `string` | Version is the release the project is at, the `version` key of <ProjectName>Config/project.yaml. |

[every contract](doc.md)
