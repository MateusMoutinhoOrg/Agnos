# `sandbox/api/trigger.go`

| Constant | Value | Description |
| --- | --- | --- |
| `EqualTrigger` | `iota` | EqualTrigger matches a text that is exactly the trigger's Value. |
| `PrefixTrigger` |  | PrefixTrigger matches a text that is the Value or continues it with a new segment: "/admin" matches "/admin" and "/admin/users", never "/administrator". A Value of "/" matches every path. |
| `TextPrefixTrigger` |  | TextPrefixTrigger matches a text that begins with the Value, whatever follows it: "/admin" matches "/administrator" too. |
| `SuffixTrigger` |  | SuffixTrigger matches a text that ends with the Value. |
| `RegexTrigger` |  | RegexTrigger matches a text the Value, a regular expression, matches. |
| `OneOfTrigger` |  | OneOfTrigger matches a text that is exactly one of the trigger's Values — how a command answers to more than one name. |

## `TriggerType`

TriggerType is how a Trigger compares the text it is handed.

`type TriggerType int`

## `Trigger`

Trigger is the condition one declared slice or value has to meet for its unit to run at all — the parsed form of one `trigger:` of a route.yaml or a command.yaml. Failing it is a non-match, never a usage error or a 400: the input is for some other unit.

| Field | Type | Description |
| --- | --- | --- |
| `Exist` | `bool` | Exist tells a declared trigger from none at all; an entry with none matches whatever the request brought. |
| `Type` | `TriggerType` | Type is how Value is compared. |
| `Value` | `string` | Value is what the text is compared against. |
| `Values` | `[]string` | Values are the texts a OneOfTrigger accepts; nil on every other type. |
| `Negate` | `bool` | Negate inverts the comparison: the trigger holds when the text does not match. |
| `IgnoreCase` | `bool` | IgnoreCase compares without regard to case. |

[every contract](doc.md)
