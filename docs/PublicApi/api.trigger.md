# `sandbox/api/generated.trigger.go`

| Constant | Value | Description |
| --- | --- | --- |
| `TriggerEqual` | `opinionatedagnoscli.TriggerEqual` | TriggerEqual matches a text that is exactly the trigger's Value. |
| `TriggerPrefix` | `opinionatedagnoscli.TriggerPrefix` | TriggerPrefix matches a text that is the Value or continues it with a new segment: "/admin" matches "/admin/users", never "/administrator". |
| `TriggerTextPrefix` | `opinionatedagnoscli.TriggerTextPrefix` | TriggerTextPrefix matches a text that begins with the Value, whatever follows it. |
| `TriggerSuffix` | `opinionatedagnoscli.TriggerSuffix` | TriggerSuffix matches a text that ends with the Value. |
| `TriggerRegex` | `opinionatedagnoscli.TriggerRegex` | TriggerRegex matches a text the Value, a regular expression, matches. |
| `TriggerOneOf` | `opinionatedagnoscli.TriggerOneOf` | TriggerOneOf matches a text that is exactly one of the trigger's Values. |

## `TriggerType`

TriggerType is how a Trigger compares the text it is handed. It is the OpinionatedAgnosCli contract's own, shared by the cli and the server layer.

`type TriggerType = opinionatedagnoscli.TriggerType`

## `Trigger`

Trigger is the condition one declared slice or value has to meet for its unit to run at all — the parsed form of one `trigger:` of a route.yaml or a command.yaml.

`type Trigger = opinionatedagnoscli.Trigger`

[every contract](doc.md)
