# `deps.Interviewer`

`sandbox/deps/interviewer`

## `AlternativeOption`

AlternativeOption is one choice offered by an alternative question: Id is what the answer returns, Msg is what the person reading it sees.

| Field | Type | Description |
| --- | --- | --- |
| `Id` | `string` | Id is the value the question returns when this option is chosen. It is never displayed. |
| `Msg` | `string` | Msg is the label shown for this option. It is never returned. |

## `Sandbox`

Sandbox is the interview library injected whole as the Deps.Interviewer field. Every field asks one question and blocks until it is answered; an error reports that no answer can be had — the input ended, or the terminal could not be read — and never that the answer was invalid.

| Field | Type | Description |
| --- | --- | --- |
| `IntQuestion` | `func(question string) (int, error)` | IntQuestion asks for a whole number and returns it already converted. |
| `StrQuestion` | `func(question string) (string, error)` | StrQuestion asks for a line of text. An empty answer is a valid one: it is how a caller offering a default learns that the default was taken. |
| `FloatQuestion` | `func(question string) (float64, error)` | FloatQuestion asks for a number and returns it already converted. |
| `BoolQuestion` | `func(question string) (bool, error)` | BoolQuestion asks a yes-or-no question and returns the answer. |
| `SingleAlternativeQuestion` | `func(question string, alternatives []AlternativeOption) (string, error)` | SingleAlternativeQuestion offers a list of options and returns the Id of the one chosen — never its Msg, so the wording of an option may change without moving what the caller matches on. |
| `MultipleAlternativeQuestion` | `func(question string, alternatives []AlternativeOption) ([]string, error)` | MultipleAlternativeQuestion offers a list of options and returns the Ids of every one chosen, in the order the options were declared. Choosing none is a valid answer and returns an empty slice. |
| `Back` | `func(err error) bool` | Back reports whether an error a question returned means the person asked to step back to the question before it, rather than that no answer can be had at all. It is the one way a caller tells the two apart, so that going back can be undone by asking again while every other error ends the session. |

[every contract](doc.md)
