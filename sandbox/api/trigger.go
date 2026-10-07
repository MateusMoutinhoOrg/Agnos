package api

import (
	opinatedagnoscli "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/OpinatedAgnosCli"
)

// TriggerType is how a Trigger compares the text it is handed. It is the
// OpinatedAgnosCli contract's own, shared by the cli and the server layer.
type TriggerType = opinatedagnoscli.TriggerType

// EqualTrigger matches a text that is exactly the trigger's Value.
const EqualTrigger = opinatedagnoscli.EqualTrigger

// PrefixTrigger matches a text that is the Value or continues it with a new
// segment: "/admin" matches "/admin/users", never "/administrator".
const PrefixTrigger = opinatedagnoscli.PrefixTrigger

// TextPrefixTrigger matches a text that begins with the Value, whatever
// follows it.
const TextPrefixTrigger = opinatedagnoscli.TextPrefixTrigger

// SuffixTrigger matches a text that ends with the Value.
const SuffixTrigger = opinatedagnoscli.SuffixTrigger

// RegexTrigger matches a text the Value, a regular expression, matches.
const RegexTrigger = opinatedagnoscli.RegexTrigger

// OneOfTrigger matches a text that is exactly one of the trigger's Values.
const OneOfTrigger = opinatedagnoscli.OneOfTrigger

// Trigger is the condition one declared slice or value has to meet for its
// unit to run at all — the parsed form of one `trigger:` of a route.yaml or a
// command.yaml.
type Trigger = opinatedagnoscli.Trigger
