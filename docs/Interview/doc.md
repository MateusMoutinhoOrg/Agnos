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
| Nothing has run yet | the command line the answers add up to, plus what it rewrites and what it removes |
| Run it / change one answer / back | — |

Running a command returns to the first menu, rebuilt from the project on disk — a command that
turned a layer on opens that layer's area on the next pass, and the step that turned it on is
gone. A command that **fails** does not: the confirm screen comes back with every answer still on
it, so one rejected character costs one `change`, not the whole questionnaire.

`· exit` ends the session; so does ctrl-c, ctrl-d, or the input running out — none of them is a
failure, and all exit `0`.

## Going back

| Where | Key | Line mode |
| --- | --- | --- |
| any question | `esc` (and `q` on a menu) | a line that reads `:back` |
| ends the session | `ctrl-c`, `ctrl-d` | the input running out |

Going back re-asks the question before this one and drops every answer after it, so what is asked
again is asked against the state it was asked against the first time. Going back past the first
question of a command leaves that command unrun and brings the menu back; going back at the menu
is the way out.

Nothing is answered by going back — a text question reads the keys itself rather than a whole
line, which is what keeps every printable character an answer and leaves escape meaning escape.
`deps.Interviewer.Back` is how a session tells that error from one that means no answer can be
had at all.

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
| `--module` of `start` | asked as **required** when the target folder has no `go.mod`, optional when it has one |
| everything else | asked, defaulting to what `entries.yaml` declares |

A declaration is written once and read everywhere, so it can only call a field optional; the
folder being worked on is what decides. `--module` is the one field this applies to today, and
the rule is the one `start`'s handler rejects on — asking it as required only moves that error
from after the confirm screen to the question itself.

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

## Nothing destructive under a blind enter

A guided screen tells the person the pre-selected row is the safe one, so it has to be.

| Menu | Leads with |
| --- | --- |
| an area whose first command removes or overwrites (`Extensions`, headed by `disable-extension`) | `· back` |
| the first menu when every step left on it is an offer to install a whole layer | `· exit` |
| the confirm screen of a command that takes something away | `· no, back to the menu` |

A purge's confirm screen names what goes with it — `this removes 3 commands: greet, help,
version` — read off the project, not off the command line.

## The command line is the command that runs

The confirm screen promises *you could have typed it yourself*. A name typed with capitals,
spaces or underscores is written down normalized, so the screen says so as well:

```
│  $ agnos add-flag --command greet "My Flag Name!"
│  "My Flag Name!" is written down as the flag --my-flag-name!
```

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
