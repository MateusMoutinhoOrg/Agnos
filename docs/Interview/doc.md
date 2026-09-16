# Interview

`agnos interview` drives the whole command surface by asking questions. It declares nothing of
its own: every menu and every question is generated from `Cli.Commands`, so a command added
tomorrow is covered without this feature changing.

```bash
agnos interview                       # the current directory
agnos interview --path ./my-project   # another project
```

## The session

| Step | Comes from |
| --- | --- |
| What do you want to do? | the `category` of every declared command |
| Which command? | the commands of that category, with their `help` |
| One question per field | each `CommandArg` and `CommandFlag` of that command |
| About to run | the command line the answers add up to |
| Run it / change one answer / back | — |

Running a command returns to the first menu. `· exit` ends the session; so does ctrl-c, `q`, or
the input running out — none of them is a failure, and all exit `0`.

## Questions

Each field is asked as what it declares: `boolean` as yes or no, `int` and `float` as a number
checked against `min` and `max`, `array` as one value at a time until an empty answer, everything
else as text. An answer that will not convert or falls outside the bounds is asked again — the
same rules the dispatch applies to a command line, applied before a handler runs.

| Field | Answered by |
| --- | --- |
| `--path` | the session's own `--path`, never asked |
| `--quiet` | always off, so the build a command runs stays visible |
| everything else | asked, defaulting to what `entries.yaml` declares |

A field that names something already in the project is offered as a list instead of a text box —
the commands declared, the deps installed, the adapters, the availables, the routes, the pages,
the docs, the examples, the themes of `themes.yaml`, the extensions of the catalog, and the
closed vocabularies (`string`/`boolean`/`int`/`float`, the http methods, the compile targets).
Those lists are read from the project at `--path`, not from the binary running the interview.

## Running

The interview binds the answers onto a copy of the declaration with `api.BindCommand` and calls
that command's own handler — the same call the dispatch makes, with the values coming from
questions instead of from argv. Each handler runs its own action, so each command persists and
builds for itself; the interview writes nothing.

## The contract behind it

Questions are asked through [`deps.Interviewer`](../PublicApi/doc.md), which says nothing about
how one is presented. The adapter answers with an arrow-key menu when stdin is a terminal it can
put in raw mode, and with a numbered list read line by line when it is not — a pipe, a file, CI
or Windows. It is installable in any project:

```bash
agnos add-dep interviewer
```
