# Interview

`agnos interview` is the one surface of agnos written for a person instead of for an llm — an
llm drives agnos through the plain cli ([Commands](../Commands/doc.md)). It asks in plain words,
offers only what the project in front of it can run, and leads with the step that project needs
next. Every other rule of this repo still binds its code; the exception is only about who its
screens are written for.

```bash
agnos interview                       # the current directory
agnos interview --path ./my-project   # another project
```

It declares nothing of its own: every menu and every question is generated from `Cli.Commands`,
so a command added tomorrow is covered without this feature changing.

## The session

| Step | Comes from |
| --- | --- |
| What do you want to do? | the steps this project has not taken, then the areas it has |
| Which command? | the commands of that area, with their `help` |
| One question per field | each `CommandArg` and `CommandFlag` of that command |
| Nothing has run yet | the command line the answers add up to |
| Run it / change one answer / back | — |

Running a command returns to the first menu, rebuilt from the project on disk — a command that
turned a layer on opens that layer's area on the next pass, and the step that turned it on is
gone. `· exit` ends the session; so does ctrl-c, `q`, or the input running out — none of them is
a failure, and all exit `0`.

## The first menu

The project at `--path` is read before every menu: whether `AgnosConfig/project.yaml` is there,
which extensions are on, and how many units of its own each layer declares.

**Steps** — what this project has not done yet, in the order a project grows. The first key one
is marked `★`; the rest are offers.

| Step | Offered when | Key |
| --- | --- | --- |
| `start` | there is no project in the folder | ★ |
| `cli-init` | `sandbox-cli` is off | ★ |
| `add-command` | the cli is on and declares no command of its own | ★ |
| `add-route` | the server is on and declares no route of its own | ★ |
| `add-page` | the front is on and has no page | ★ |
| `server-init` | `sandbox-server` is off | |
| `front-init` | the server is on and `sandbox-front` is off | |
| `deps-init` | `sandbox-deps` is off | |

`help`, `version`, `health` and `static` are what an init scaffolds, so they never count as
units the project declared itself.

**Areas** — one row per category of the command surface, offered only while the mechanic that
owns it is on.

| Area | Needs |
| --- | --- |
| `Cli System` | `sandbox-cli` |
| `Server System` | `sandbox-server` |
| `Front System` | `sandbox-front` |
| `Deps System` | `sandbox-deps` |
| `Examples` | `sandbox-example` |
| `Documentation` | `doc` |
| `Core Commands`, `Extensions`, `Info` | — |

With no project in the folder only `start` and `Info` are offered: every other command answers
*run `agnos start` first*. An `<x>-init` is only ever a step — while its mechanic is off the area
is hidden and the step is the way in, and once it is on the init is offered nowhere.

## Questions

Each field is asked as what it declares: `boolean` as yes or no, `int` and `float` as a number
checked against `min` and `max`, `array` as one value at a time until an empty answer, everything
else as text. An answer that will not convert or falls outside the bounds is asked again — the
same rules the dispatch applies to a command line, applied before a handler runs.

| Field | Answered by |
| --- | --- |
| `--path` | the session's own `--path`, never asked |
| `--path` of `start` | asked: it is the folder being created, with the session's path as its default |
| `--quiet` | always off, so the build a command runs stays visible |
| everything else | asked, defaulting to what `entries.yaml` declares |

A field an answer already given decides is not asked at all, and does not appear on the confirm
screen: `--required` after a `--default` (and on a boolean, whose absence already means false),
`--default` after a `--required`, `--min`/`--max` on a field that is not a number, `set-body`'s
`--optional` after its `--required`, and each `add-body-field` keyword outside the `--type` it
applies to. Every one of them is a combination the command itself refuses, so the question only
led to an error after the confirm screen. Changing an answer there drops the answers it has just
ruled out.

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
