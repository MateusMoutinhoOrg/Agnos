# `sandbox/api/trigger.go`

| Constant | Value | Description |
| --- | --- | --- |
| `EqualTrigger` | `opinatedagnoscli.EqualTrigger` | EqualTrigger matches a text that is exactly the trigger's Value. |
| `PrefixTrigger` | `opinatedagnoscli.PrefixTrigger` | PrefixTrigger matches a text that is the Value or continues it with a new segment: "/admin" matches "/admin/users", never "/administrator". |
| `TextPrefixTrigger` | `opinatedagnoscli.TextPrefixTrigger` | TextPrefixTrigger matches a text that begins with the Value, whatever follows it. |
| `SuffixTrigger` | `opinatedagnoscli.SuffixTrigger` | SuffixTrigger matches a text that ends with the Value. |
| `RegexTrigger` | `opinatedagnoscli.RegexTrigger` | RegexTrigger matches a text the Value, a regular expression, matches. |
| `OneOfTrigger` | `opinatedagnoscli.OneOfTrigger` | OneOfTrigger matches a text that is exactly one of the trigger's Values. |

## `TriggerType`

TriggerType is how a Trigger compares the text it is handed. It is the OpinatedAgnosCli contract's own, shared by the cli and the server layer.

`type TriggerType = opinatedagnoscli.TriggerType`

## `Trigger`

Trigger is the condition one declared slice or value has to meet for its unit to run at all — the parsed form of one `trigger:` of a route.yaml or a command.yaml.

`type Trigger = opinatedagnoscli.Trigger`

[every contract](doc.md)
