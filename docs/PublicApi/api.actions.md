# `sandbox/api/actions.go`

| Constant | Value | Description |
| --- | --- | --- |
| `RuntimeGo` | `"go"` | RuntimeGo resolves the module graph and compiles every package after the render. |
| `RuntimeNone` | `"none"` | RuntimeNone renders only, leaving the result unchecked. |
| `DefaultRoutePriority` | `100` | DefaultRoutePriority is the rung a route lands on when it names none — high enough to leave the rungs below it to the middlewares in front. |
| `DefaultMiddlewarePriority` | `10` | DefaultMiddlewarePriority is the rung `add-route --middleware` lands on when it names none: below every route left on DefaultRoutePriority. |

## `BuildProps`

BuildProps describes one (re)render of a project: the directory holding it and the runtime that then checks the result.

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Runtime` | `string` |

## `CompileProps`

CompileProps describes one cross-compile run: the directory holding the project and the target names to build. Each name is one of the keys `agnos compile` accepts (linux86, linuxarm64, linuxi32, mac86, macarm64, windows86, windowsi32) or "all", which expands to every target.

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Targets` | `[]string` |

## `StartProps`

StartProps describes one project to scaffold: the directory to write it into, the name it carries in <Name>Config/project.yaml, the module path for go.mod (nil derives it from the name) and whether an existing directory may be written over.

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `ProjectName` | `string` |
| `Module` | `*string` |
| `Force` | `bool` |

## `ExecTestProps`

ExecTestProps describes one run of the project's example suite: the directory holding the project, the single example name to run (empty runs every one, both sides) and whether the goldens are rewritten with what the run produced instead of compared against it.

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Only` | `string` |
| `Update` | `bool` |

## `AddDepProps`

AddDepProps describes one dep to install. Dep is either a name of the embedded catalog or, when it holds a "/", the module path of another agnos repo — the same disambiguation `go get` makes. Adapter picks which implementation of a catalog dep fills the contract ("" takes the dep's declared default-adapter); As names the copied contract of a remote one ("" takes the last segment of the module path); RemoteAvailable is the available of the remote repo the generated shim builds its sandbox from ("" is standard).

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Dep` | `string` |
| `Adapter` | `string` |
| `As` | `string` |
| `RemoteAvailable` | `string` |

## `SetDepProps`

SetDepProps describes one remote dep to move to another version of the module it was copied from. RemoteAvailable is the available of the remote repo the regenerated shim builds its sandbox from ("" is standard).

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Dep` | `string` |
| `Version` | `string` |
| `RemoteAvailable` | `string` |

## `RemoveDepProps`

RemoveDepProps describes one dep to uninstall. A dep with adapters installed is refused, because which of them was meant is not a question the tree can answer; WithAdapters is the caller saying "all of them".

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Dep` | `string` |
| `WithAdapters` | `bool` |

## `AddAdapterProps`

AddAdapterProps describes one further implementation to install for a contract the project already has. Available names the available that should switch to it; "" installs the package and leaves every selection alone.

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Adapter` | `string` |
| `Available` | `string` |

## `SetAdapterProps`

SetAdapterProps describes one selection to change: which adapter fills a dep's field in one available. Available is the standard one when empty.

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Dep` | `string` |
| `Adapter` | `string` |
| `Available` | `string` |

## `ExtensionInfo`

ExtensionInfo is one row of ListExtensions: one generation mechanic of the catalog, what it looks after, and whether this project turned it on.

| Field | Type |
| --- | --- |
| `Name` | `string` |
| `Enabled` | `bool` |
| `Help` | `string` |

## `DepInfo`

DepInfo is one row of ListDeps: a dep of the embedded catalog, and what the project holds of it.

| Field | Type |
| --- | --- |
| `Name` | `string` |
| `Field` | `string` |
| `Help` | `string` |
| `DefaultAdapter` | `string` |
| `Installed` | `bool` |
| `Adapters` | `[]string` |

## `AdapterInfo`

AdapterInfo is one row of ListAdapters: an adapter of the embedded catalog or one installed in the project, and which availables bind it. Origin is "catalog" for one the catalog installs and "generated" for the shim of a dep copied from a remote repo.

| Field | Type |
| --- | --- |
| `Name` | `string` |
| `Dep` | `string` |
| `Help` | `string` |
| `Module` | `string` |
| `Origin` | `string` |
| `Installed` | `bool` |
| `Availables` | `[]string` |

## `FlagProps`

FlagProps describes one flag to declare in a command's command.yaml. Name is the flag's name, whose exported Go form is the Entries field it binds to ("out-file" -> OutFile); Keys are the spellings a user types, --<name> when none is given. Type is string, integer, number, boolean, string-array or integer-array, and Array asks for the repeatable form of a scalar Type. Default, Min and Max are the raw literals typed on the command line ("" means unset); Enum is every value accepted, Pattern a regular expression every value matches, and the Trigger* fields what the value has to match for the command to run at all. Position is the index to insert at (< 0 appends).

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Command` | `string` |
| `Name` | `string` |
| `Keys` | `[]string` |
| `Type` | `string` |
| `Array` | `bool` |
| `Required` | `bool` |
| `Default` | `string` |
| `Min` | `string` |
| `Max` | `string` |
| `Enum` | `[]string` |
| `Pattern` | `string` |
| `Trigger` | `string` |
| `TriggerType` | `string` |
| `TriggerNegate` | `bool` |
| `TriggerIgnoreCase` | `bool` |
| `Description` | `string` |
| `Position` | `int` |

## `ArgProps`

ArgProps describes one arg to declare in a command's command.yaml: the segments Start to End of the command line — the raw indexes typed on the command line, "" being the first segment no arg reads yet and Start again, "-1" the last segment — bound to the Entries field Name becomes. Type is string, integer, number or uuid, anything but string reading one segment alone; Trigger and TriggerType are what the segments, joined by a space, have to match for the command to run at all, and TriggerNegate / TriggerIgnoreCase the two switches on it. Position is the index to insert at (< 0 appends).

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Command` | `string` |
| `Name` | `string` |
| `Start` | `string` |
| `End` | `string` |
| `Type` | `string` |
| `Required` | `bool` |
| `Default` | `string` |
| `Trigger` | `string` |
| `TriggerType` | `string` |
| `TriggerNegate` | `bool` |
| `TriggerIgnoreCase` | `bool` |
| `Description` | `string` |
| `Position` | `int` |

## `ArgEditProps`

ArgEditProps describes the change set-arg applies to one arg a command declares. Name is the arg as it is declared now and Rename the name it takes on ("" leaves it alone); every other key overwrites what is there when it is given. Clear takes "trigger", "trigger-negate", "trigger-ignore-case", "type", "required", "default" or "description" off again.

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Command` | `string` |
| `Name` | `string` |
| `Rename` | `string` |
| `Start` | `string` |
| `End` | `string` |
| `Type` | `string` |
| `Required` | `bool` |
| `Default` | `string` |
| `Trigger` | `string` |
| `TriggerType` | `string` |
| `TriggerNegate` | `bool` |
| `TriggerIgnoreCase` | `bool` |
| `Description` | `string` |
| `Clear` | `[]string` |

## `FlagEditProps`

FlagEditProps describes the change set-flag applies to one flag a command declares. Name is the flag as it is declared now — its name, its id or one of its keys — and Rename the name it takes on ("" leaves it alone); Keys replace the spellings when any is given, Enum the accepted values, and every other key overwrites what is there when it is given. Clear takes "keys", "type", "required", "default", "min", "max", "enum", "pattern", "trigger", "trigger-negate", "trigger-ignore-case" or "description" off again.

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Command` | `string` |
| `Name` | `string` |
| `Rename` | `string` |
| `Keys` | `[]string` |
| `Type` | `string` |
| `Required` | `bool` |
| `Default` | `string` |
| `Min` | `string` |
| `Max` | `string` |
| `Enum` | `[]string` |
| `Pattern` | `string` |
| `Trigger` | `string` |
| `TriggerType` | `string` |
| `TriggerNegate` | `bool` |
| `TriggerIgnoreCase` | `bool` |
| `Description` | `string` |
| `Clear` | `[]string` |

## `AddCommandProps`

AddCommandProps describes one command to scaffold. Name is its package and the verb its first arg answers to on segment 0; Trigger is another value that arg compares against ("" is the name, or every command line for a Middleware), TriggerType how ("" is equal — prefix for a Middleware), and TriggerNegate / TriggerIgnoreCase the two switches on it. Pattern declares the args from one command-line shape instead ("route add {name} {*rest}") and excludes Trigger and TriggerType. A Middleware is not strict, runs on DefaultMiddlewarePriority and answers nothing. Priority is the rung it runs on, used only when HasPriority is set; Before and After name another command to land one rung below or above instead, and exclude Priority. Dir is the folder under sandbox/internal/commands the command lands in ("" is the top; "admin" puts it in commands/admin/<name>): its command.yaml is what makes it a command, whatever folder holds it.

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Name` | `string` |
| `Trigger` | `string` |
| `TriggerType` | `string` |
| `TriggerNegate` | `bool` |
| `TriggerIgnoreCase` | `bool` |
| `Pattern` | `string` |
| `Middleware` | `bool` |
| `Priority` | `int` |
| `HasPriority` | `bool` |
| `Before` | `string` |
| `After` | `string` |
| `Help` | `string` |
| `Category` | `string` |
| `Dir` | `string` |

## `RenameCommandProps`

RenameCommandProps describes one command to rename: Command as it is declared now, Name the name it takes on. Dir, when HasDir is set, is the folder under sandbox/internal/commands it moves to ("" is the top); without it the command stays in the folder it sits in. Name may be its current one when only the folder changes.

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Command` | `string` |
| `Name` | `string` |
| `Dir` | `string` |
| `HasDir` | `bool` |

## `RebalanceCommandsProps`

RebalanceCommandsProps describes one rebalance of the cli chain: every command is laid down again Step rungs apart, in the order it runs now, the first one on Step.

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Step` | `int` |

## `ExplainCommandProps`

ExplainCommandProps describes one command line to run against the declared commands without running any: Argv is the line, as it would be typed after the binary name.

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Argv` | `[]string` |

## `SetCommandProps`

SetCommandProps carries the command-level keys of command.yaml that set-command may rewrite. Empty strings leave the current value alone; Identifiers — further verbs the command answers to — and Examples are appended (deduplicated), and Hidden / Visible, Strict / Loose are the two sides of one switch each. Priority is the rung the command runs on, read when HasPriority is set; Before and After name another command to land one rung below or above instead. Segments is the segment count the command line has to have, read when HasSegments is set. Clear takes "segments" off again.

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Command` | `string` |
| `Help` | `string` |
| `Category` | `string` |
| `LongDescription` | `string` |
| `Hidden` | `bool` |
| `Visible` | `bool` |
| `Strict` | `bool` |
| `Loose` | `bool` |
| `Priority` | `int` |
| `HasPriority` | `bool` |
| `Before` | `string` |
| `After` | `string` |
| `Segments` | `int` |
| `HasSegments` | `bool` |
| `Identifiers` | `[]string` |
| `Examples` | `[]string` |
| `Clear` | `[]string` |

## `AddRouteProps`

AddRouteProps describes one route to scaffold. Trigger is the whole-path value its first path compares against ("" is "/" followed by the name), TriggerType how ("equal", "prefix", "text-prefix", "suffix" or "regex", or the aliases starts-with, ends-with, exact, equals and matches; "" is equal — prefix for a Middleware), and TriggerNegate / TriggerIgnoreCase the two switches on it. Pattern declares the paths from one url shape instead ("/users/{id:integer}/{*rest}") and excludes Trigger and TriggerType. Methods are the http methods it answers to ([] is GET — ANY for a Middleware) and ResponseType the Content-Type its responses carry ("" is application/json — text/plain for a Middleware). Priority is the rung it runs on, used only when HasPriority is set; without it the route lands on DefaultRoutePriority, or DefaultMiddlewarePriority for a Middleware. Before and After name another route to land one rung below or above instead, and exclude Priority. Dir is the folder under sandbox/internal/routeslist the route lands in ("" is the top; "admin" puts it in routeslist/admin/<name>): its route.yaml is what makes it a route, whatever folder holds it.

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Name` | `string` |
| `Methods` | `[]string` |
| `Trigger` | `string` |
| `TriggerType` | `string` |
| `TriggerNegate` | `bool` |
| `TriggerIgnoreCase` | `bool` |
| `Pattern` | `string` |
| `Middleware` | `bool` |
| `Priority` | `int` |
| `HasPriority` | `bool` |
| `Before` | `string` |
| `After` | `string` |
| `ResponseType` | `string` |
| `Help` | `string` |
| `Category` | `string` |
| `Dir` | `string` |

## `RouteProps`

RouteProps carries the route-level keys of route.yaml that set-route may rewrite. Empty strings leave the current value alone; Methods replace the whole list when any is given; Examples are appended (deduplicated), and Hidden / Visible are the two sides of one switch. Priority is the rung the route runs on, and HasPriority is what tells a priority declared as zero from one not given at all; Before and After name another route to land one rung below or above instead. Segments is the segment count the request path has to have, read when HasSegments is set. Clear takes "segments" off again.

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Route` | `string` |
| `Methods` | `[]string` |
| `ResponseType` | `string` |
| `Help` | `string` |
| `Category` | `string` |
| `LongDescription` | `string` |
| `Hidden` | `bool` |
| `Visible` | `bool` |
| `Priority` | `int` |
| `HasPriority` | `bool` |
| `Before` | `string` |
| `After` | `string` |
| `Segments` | `int` |
| `HasSegments` | `bool` |
| `Examples` | `[]string` |
| `Clear` | `[]string` |

## `RenameRouteProps`

RenameRouteProps describes one route to rename: Route as it is declared now, Name the name it takes on. Dir, when HasDir is set, is the folder under sandbox/internal/routeslist it moves to ("" is the top); without it the route stays in the folder it sits in. Name may be its current one when only the folder changes.

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Route` | `string` |
| `Name` | `string` |
| `Dir` | `string` |
| `HasDir` | `bool` |

## `RebalanceRoutesProps`

RebalanceRoutesProps describes one rebalance of the chain: every route is laid down again Step rungs apart, in the order it runs now, the first one on Step.

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Step` | `int` |

## `ExplainRouteProps`

ExplainRouteProps describes one request to run against the declared routes without a server: its Method and its Path (query string included), and the Headers and Cookies it carries as "key=value" entries.

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Method` | `string` |
| `RequestPath` | `string` |
| `Headers` | `[]string` |
| `Cookies` | `[]string` |

## `DatabaseFieldProps`

DatabaseFieldProps describes one field to add to a table of a database's specs.yaml. Table is the table it lands in and Parent the nested collection inside that table, "" for a field of the table itself. Type is one of key, string, int, float, link or database; Target names the table a link points at and belongs to a link alone.

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Database` | `string` |
| `Table` | `string` |
| `Parent` | `string` |
| `Name` | `string` |
| `Type` | `string` |
| `Required` | `bool` |
| `Target` | `string` |

## `DatabaseFieldEditProps`

DatabaseFieldEditProps describes the change set-table-field applies to one field a table already declares. Name is the field as it is declared now and Rename the spelling it takes on ("" leaves it alone); Type and Target overwrite what is there when they are given, and an empty one leaves it as it is. Clear is how a key is taken off again — "required" or "target" — because an empty string cannot say "unset this" and "leave it alone" at once.

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Database` | `string` |
| `Table` | `string` |
| `Parent` | `string` |
| `Name` | `string` |
| `Rename` | `string` |
| `Type` | `string` |
| `Required` | `bool` |
| `Target` | `string` |
| `Clear` | `[]string` |

## `RoutePathProps`

RoutePathProps describes one entry to add to a route's `paths`. Id is the Entries field the slice binds to; Start and End are the raw segment indexes typed on the command line ("" is 0 and -1, the whole path); Trigger is what the slice has to read as for the route to run and TriggerType how it is compared ("" is equal) — a path with no Trigger is a plain capture — and TriggerNegate / TriggerIgnoreCase the two switches on it. Type is what the slice converts to: "string" (the default), "integer", "number" or "uuid", anything but string reading one segment alone. Position is the index to insert at (< 0 appends).

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Route` | `string` |
| `Id` | `string` |
| `Start` | `string` |
| `End` | `string` |
| `Type` | `string` |
| `TriggerType` | `string` |
| `Trigger` | `string` |
| `TriggerNegate` | `bool` |
| `TriggerIgnoreCase` | `bool` |
| `Description` | `string` |
| `Position` | `int` |

## `RoutePathEditProps`

RoutePathEditProps describes the change set-path applies to one entry of a route's `paths`. Id is the entry as it is declared now and Rename the id it takes on ("" leaves it alone); every other key overwrites what is there when it is given. Clear takes "trigger", "trigger-negate", "trigger-ignore-case", "type" or "description" off again.

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Route` | `string` |
| `Id` | `string` |
| `Rename` | `string` |
| `Start` | `string` |
| `End` | `string` |
| `Type` | `string` |
| `TriggerType` | `string` |
| `Trigger` | `string` |
| `TriggerNegate` | `bool` |
| `TriggerIgnoreCase` | `bool` |
| `Description` | `string` |
| `Clear` | `[]string` |

## `RouteParameterProps`

RouteParameterProps describes one entry to add to a route's `parameters`. Name is the query key or header name it is read under — its Entries field is the exported spelling of it. Type is "string", "integer", "number", "boolean", "datetime", "string-array" or "integer-array"; Fonts are where it is read from, in order — "query", "header", "cookie" ([] is the query string alone). Default is the raw literal typed on the command line ("" means unset). Trigger and TriggerType are a condition on the value that puts the parameter into what the route matches on. Position is the index to insert at (< 0 appends).

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Route` | `string` |
| `Name` | `string` |
| `Type` | `string` |
| `Fonts` | `[]string` |
| `Required` | `bool` |
| `Default` | `string` |
| `TriggerType` | `string` |
| `Trigger` | `string` |
| `TriggerNegate` | `bool` |
| `TriggerIgnoreCase` | `bool` |
| `Description` | `string` |
| `Examples` | `[]string` |
| `Position` | `int` |

## `RouteParameterEditProps`

RouteParameterEditProps describes the change set-parameter applies to one entry of a route's `parameters`. Name is the key as it is declared now and Rename the key it takes on ("" leaves it alone); Fonts replace the whole list when any is given; every other key overwrites what is there when it is given. Clear takes "description", "examples", "default", "required", "trigger", "trigger-negate" or "trigger-ignore-case" off again.

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Route` | `string` |
| `Name` | `string` |
| `Rename` | `string` |
| `Type` | `string` |
| `Fonts` | `[]string` |
| `Required` | `bool` |
| `Default` | `string` |
| `TriggerType` | `string` |
| `Trigger` | `string` |
| `TriggerNegate` | `bool` |
| `TriggerIgnoreCase` | `bool` |
| `Description` | `string` |
| `Examples` | `[]string` |
| `Clear` | `[]string` |

## `RouteBodyProps`

RouteBodyProps describes the body envelope of one route — everything about the request body but its schema (the json-schema of a json body, the form-schema of a form one), which is grown property by property with AddBodyField. Type is "none", "raw", "text", "json" or "form"; turning json into form and back carries a flat schema along. Required and Optional are the two sides of one switch, as are the empty strings and MaxBytes < 0 that mean "leave as is". DropSchema deletes the declared schema.

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Route` | `string` |
| `Type` | `string` |
| `Required` | `bool` |
| `Optional` | `bool` |
| `MaxBytes` | `int` |
| `ContentType` | `string` |
| `DropSchema` | `bool` |

## `RouteBodyFieldProps`

RouteBodyFieldProps describes one property of a route's body schema — its json-schema, or the flat form-schema of a form body. Name is the dotted path it sits at ("address.city"), and every other field is one keyword of the supported subset: the raw literals typed on the command line, where "" means unset. Type is "string", "boolean", "int", "float" or "object", and Array wraps the whole of it in an array schema. AdditionalProperties and NoAdditionalProperties are the two sides of one switch.

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Route` | `string` |
| `Name` | `string` |
| `Type` | `string` |
| `Required` | `bool` |
| `Array` | `bool` |
| `Min` | `string` |
| `Max` | `string` |
| `ExclusiveMin` | `string` |
| `ExclusiveMax` | `string` |
| `Format` | `string` |
| `Pattern` | `string` |
| `Enum` | `[]string` |
| `Const` | `string` |
| `Nullable` | `bool` |
| `MinItems` | `string` |
| `MaxItems` | `string` |
| `UniqueItems` | `bool` |
| `AdditionalProperties` | `bool` |
| `NoAdditionalProperties` | `bool` |

## `RouteBodyFieldEditProps`

RouteBodyFieldEditProps describes the change set-body-field applies to one property a route's body json-schema already declares. Name is the dotted path it sits at and Rename the leaf spelling it takes on (it stays in the object it is declared in); every other key is one keyword of the supported subset, overwriting what is there when it is given. Clear names the keywords to take off instead — "required", "array", "min", "max", "exclusive-min", "exclusive-max", "format", "pattern", "enum", "const", "nullable", "min-items", "max-items", "unique-items" or "additional-properties" — which is the one thing an empty value cannot say.

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Route` | `string` |
| `Name` | `string` |
| `Rename` | `string` |
| `Type` | `string` |
| `Required` | `bool` |
| `Array` | `bool` |
| `Min` | `string` |
| `Max` | `string` |
| `ExclusiveMin` | `string` |
| `ExclusiveMax` | `string` |
| `Format` | `string` |
| `Pattern` | `string` |
| `Enum` | `[]string` |
| `Const` | `string` |
| `Nullable` | `bool` |
| `MinItems` | `string` |
| `MaxItems` | `string` |
| `UniqueItems` | `bool` |
| `AdditionalProperties` | `bool` |
| `NoAdditionalProperties` | `bool` |
| `Clear` | `[]string` |

## `RouteBodyImportProps`

RouteBodyImportProps describes one example payload to read a route's body json-schema off. Json is the document itself and File a path to read it from — exactly one of the two — and the inference walks it: an object becomes an object property, a list an array of whatever its first item is, and a scalar the type it is written as. Required lists every key the example carries in its object's required set, InferFormat reads an email, a uuid, a date-time or a uri back as the format it spells, and Replace drops the schema that is there instead of adding to it.

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Route` | `string` |
| `Json` | `string` |
| `File` | `string` |
| `Required` | `bool` |
| `Replace` | `bool` |
| `InferFormat` | `bool` |

## `PageProps`

PageProps describes one html page to scaffold: the project directory, the name the page carries — its path under assets/frontend/ without the .html, slashes allowed ("blog/post"), "index" the page "/" answers — and the <title> the scaffolded html carries ("" defaults to the name).

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Name` | `string` |
| `Title` | `string` |

## `DocProps`

DocProps describes one doc to create under docs/. Name is the doc's directory, optionally nested under its parent ("PublicApi/api.Actions"). Themes are the theme ids of <ProjectName>Config/themes.yaml the doc belongs to: required on a first-level doc, forbidden on a sub-doc.

| Field | Type |
| --- | --- |
| `Path` | `string` |
| `Name` | `string` |
| `Description` | `string` |
| `Themes` | `[]string` |

## `Actions`

Actions is the whole set of operations agnos performs on a project. Every field takes the project directory as its first input (`path`, or the Path of a props struct) and reports failure as an error; the ones that change the tree re-render it before returning, so a project is always left in a built state.

| Field | Type | Description |
| --- | --- | --- |
| `Build` | `func(props BuildProps) error` | Build re-renders every generated file of the project and hands the result to the runtime named by the props. |
| `Compile` | `func(props CompileProps) error` | Compile cross-compiles the project's cmd/ binaries into release/, one file per named target. |
| `Verify` | `func(path string) error` | Verify checks the project against the schema every generator assumes and writes nothing; it reports every violation at once. |
| `Start` | `func(props StartProps) error` | Start scaffolds a new project: the config directory, go.mod, the sandbox skeleton and a first build. |
| `EnableExtension` | `func(path string, name string) error` | EnableExtension turns one generation mechanic on in the project's extensions.yaml and rebuilds, so what that mechanic owns is rendered from here on. |
| `DisableExtension` | `func(path string, name string) error` | DisableExtension turns one generation mechanic off. Nothing is removed: agnos stops rendering what that mechanic owns and the files it wrote become the project's, to keep or to edit by hand. Deleting them is what the matching <x>-purge is for. |
| `ListExtensions` | `func(path string) ([]ExtensionInfo, error)` | ListExtensions returns one row per generation mechanic of the catalog, saying which ones this project turned on. |
| `DepsInit` | `func(path string) error` | DepsInit adds the dependency layer (sandbox/deps/ and adapters/availables/standard/) to a project that has none. |
| `DepsPurge` | `func(path string) error` | DepsPurge removes the dependency layer and every installed dep with it. |
| `AddDep` | `func(props AddDepProps) error` | AddDep installs one dep of the built-in list: its contract under sandbox/deps/, one adapter filling it under adapters/libs/ and that adapter's go.mod require. |
| `RemoveDep` | `func(props RemoveDepProps) error` | RemoveDep uninstalls one installed dep: its adapters, their requires, and then the contract itself. It refuses a dep that still has an adapter installed unless props.WithAdapters says to take those too. |
| `ListDeps` | `func(path string) ([]DepInfo, error)` | ListDeps returns one row per dep of the embedded catalog, saying which the project has installed and which adapters fill each one. |
| `SetDep` | `func(props SetDepProps) error` | SetDep re-copies one remote dep at another version of its module and regenerates the shim that converts it. |
| `AddAdapter` | `func(props AddAdapterProps) error` | AddAdapter installs one further implementation of a contract the project already has, and — when props.Available names one — switches that available to it. |
| `RemoveAdapter` | `func(path string, adapter string) error` | RemoveAdapter uninstalls one adapter, its require and its files. It refuses one that an available still binds, and one written by the generator as half of a remote dep. |
| `SetAdapter` | `func(props SetAdapterProps) error` | SetAdapter changes which adapter fills one dep's field in one available, the only place that choice is recorded. |
| `ListAdapters` | `func(path string) ([]AdapterInfo, error)` | ListAdapters returns one row per adapter, of the embedded catalog and of the project, with the availables that bind each one. |
| `AddAvailable` | `func(path string, available string) error` | AddAvailable creates one further available, seeded with the standard available's selection so it starts filling every field. |
| `RemoveAvailable` | `func(path string, available string) error` | RemoveAvailable deletes one available. The standard one is refused: it is what cmd/main/main.go imports. |
| `CliInit` | `func(path string) error` | CliInit adds the CLI layer (cmd/main, the dispatcher and the help and version commands) to a project that has none. |
| `CliPurge` | `func(path string) error` | CliPurge removes the CLI layer and every command declared in it. |
| `AddCommand` | `func(props AddCommandProps) error` | AddCommand declares a new command: its command.yaml, its generated new.go and entries.go, and an InternalPureHandler.go to fill in. |
| `RemoveCommand` | `func(path string, name string) error` | RemoveCommand deletes one command and unwires it from the dispatch. |
| `SetCommand` | `func(props SetCommandProps) error` | SetCommand rewrites the command-level keys of one command's command.yaml. |
| `RenameCommand` | `func(props RenameCommandProps) error` | RenameCommand moves one command to a new name: its package, and the verb its first arg answers to when that verb was its name. |
| `RebalanceCommands` | `func(props RebalanceCommandsProps) error` | RebalanceCommands lays every command down again, Step rungs apart, in the order the chain runs them now. |
| `ListCommands` | `func(path string) ([]string, error)` | ListCommands renders every declared command as one line, in the order the dispatch runs them. |
| `ShowCommand` | `func(path string, command string) ([]string, error)` | ShowCommand renders one command's whole declaration — its args, its flags and the middlewares in front of it — as the lines of a tree. |
| `ExplainCommand` | `func(props ExplainCommandProps) ([]string, error)` | ExplainCommand runs one command line against the declared commands without running any, and says, command by command, whether it runs. |
| `AddFlag` | `func(props FlagProps) error` | AddFlag declares one flag on a command. |
| `SetFlag` | `func(props FlagEditProps) error` | SetFlag rewrites one declared flag of a command. |
| `RemoveFlag` | `func(path string, command string, name string) error` | RemoveFlag deletes one declared flag from a command. |
| `AddArg` | `func(props ArgProps) error` | AddArg declares one arg — a slice of the segments — on a command. |
| `SetArg` | `func(props ArgEditProps) error` | SetArg rewrites one declared arg of a command. |
| `RemoveArg` | `func(path string, command string, name string) error` | RemoveArg deletes one declared arg from a command. |
| `ServerInit` | `func(path string) error` | ServerInit adds the http server layer (sandbox/internal/server, the OpinatedAgnosServer lib, the health route of sandbox/internal/routeslist and the start-server command) to a project that has none, installing the CLI layer first when it is missing. |
| `ServerPurge` | `func(path string) error` | ServerPurge removes the server layer and every route declared in it. |
| `AddRoute` | `func(props AddRouteProps) error` | AddRoute declares a new route: its route.yaml, its generated new.go and entries.go, and an InternalPureHandler.go to fill in. |
| `RemoveRoute` | `func(path string, name string) error` | RemoveRoute deletes one route and unwires it from the dispatch. |
| `SetRoute` | `func(props RouteProps) error` | SetRoute rewrites the route-level keys of one route's route.yaml. |
| `AddPath` | `func(props RoutePathProps) error` | AddPath declares one slice of the request path on a route: the segments it reads, and the trigger they have to match when it declares one. |
| `SetPath` | `func(props RoutePathEditProps) error` | SetPath rewrites one entry of a route's `paths`, named by its id. |
| `RemovePath` | `func(path string, route string, id string) error` | RemovePath deletes one entry of a route's `paths`, named by its id. |
| `AddParameter` | `func(props RouteParameterProps) error` | AddParameter declares one value a route reads off the query string or the headers. |
| `SetParameter` | `func(props RouteParameterEditProps) error` | SetParameter rewrites one entry of a route's `parameters`, named by its key. |
| `RemoveParameter` | `func(path string, route string, name string) error` | RemoveParameter deletes one entry of a route's `parameters`, named by its key. |
| `SetBody` | `func(props RouteBodyProps) error` | SetBody rewrites the body keys of one route's route.yaml. |
| `AddBodyField` | `func(props RouteBodyFieldProps) error` | AddBodyField declares one property of a route's body schema (json- or form-schema, by the body's type), at the dotted path props.Name. |
| `RemoveBodyField` | `func(path string, route string, name string) error` | RemoveBodyField deletes one property from a route's body schema. |
| `SetBodyField` | `func(props RouteBodyFieldEditProps) error` | SetBodyField rewrites one property of a route's body schema, at the dotted path props.Name. |
| `ImportBody` | `func(props RouteBodyImportProps) error` | ImportBody declares a route's body schema from an example payload, inferring one property per key the example carries. |
| `ShowRoute` | `func(path string, route string) ([]string, error)` | ShowRoute renders one route's whole declaration — its paths, its parameters and its body schema — as the lines of a tree, ready to print. |
| `ListRoutes` | `func(path string) ([]string, error)` | ListRoutes renders every declared route as one line, in the order the dispatch runs them. |
| `ExplainRoute` | `func(props ExplainRouteProps) ([]string, error)` | ExplainRoute runs one request against the declared routes without a server and renders, route by route, whether it runs and why not. |
| `RenameRoute` | `func(props RenameRouteProps) error` | RenameRoute moves one route package to a new name. |
| `RebalanceRoutes` | `func(props RebalanceRoutesProps) error` | RebalanceRoutes gives every route a rung of its own, props.Step apart, in the order the chain runs them now. |
| `DatabaseInit` | `func(path string) error` | DatabaseInit adds the database layer (the store contract, the OpinatedAgnosDatabase lib and sandbox/internal/databases) to a project that has none. |
| `DatabasePurge` | `func(path string) error` | DatabasePurge removes the database layer and every database declared in it. |
| `AddDatabase` | `func(path string, name string, prefix string) error` | AddDatabase declares a new database: its specs.yaml, from which its api.go, new.go and methods.go are generated. |
| `RemoveDatabase` | `func(path string, name string) error` | RemoveDatabase deletes one database package whole. It refuses one carrying a hand-written methods_custom.go. |
| `AddTable` | `func(path string, database string, table string) error` | AddTable declares one collection of records on a database. |
| `RemoveTable` | `func(path string, database string, table string) error` | RemoveTable deletes one collection from a database. It refuses a table another table still links to. |
| `AddTableField` | `func(props DatabaseFieldProps) error` | AddTableField declares one field on a table, or on a nested collection of it. |
| `SetTableField` | `func(props DatabaseFieldEditProps) error` | SetTableField rewrites one field a table already declares. |
| `RemoveTableField` | `func(props DatabaseFieldProps) error` | RemoveTableField deletes one declared field from a table. |
| `ShowDatabase` | `func(path string, database string) ([]string, error)` | ShowDatabase renders one database's whole declaration — its tables, their fields and the methods each table generates — as the lines of a tree, ready to print. |
| `FrontInit` | `func(path string) error` | FrontInit adds the front layer (the OpinatedAgnosFront lib, the route serving every file of assets/frontend and that tree's index.html) to a project that has none, installing the server layer first when it is missing. |
| `FrontPurge` | `func(path string) error` | FrontPurge removes the front layer and the frontend route, leaving assets/frontend/ untouched. |
| `BackofficeInit` | `func(path string) error` | BackofficeInit adds the admin backoffice to a project: login, backoffice users, API tokens and the /api/admin JSON api, over a database of its own. It installs the server, front and database layers first when any is missing, and the catalog deps it calls into. Every file it writes is the project's from then on. |
| `BackofficePurge` | `func(path string) error` | BackofficePurge removes everything BackofficeInit wrote, leaving the layers it stood on, the deps it installed and the store on disk. |
| `AddPage` | `func(props PageProps) error` | AddPage scaffolds a new html page, assets/frontend/<name>.html, which the frontend route serves as soon as it exists. |
| `RemovePage` | `func(path string, name string) error` | RemovePage deletes one page, assets/frontend/<name>.html. |
| `AddDoc` | `func(props DocProps) error` | AddDoc creates one doc directory under docs/, with its props.yaml and a doc.md to fill in. |
| `RemoveDoc` | `func(path string, name string) error` | RemoveDoc deletes one doc directory and everything under it. |
| `AddCliExample` | `func(path string, name string) error` | AddCliExample creates one example under examples/cli/, with an example.sh stub that already runs. |
| `RemoveCliExample` | `func(path string, name string) error` | RemoveCliExample deletes one example of examples/cli/ whole. |
| `AddLibExample` | `func(path string, name string) error` | AddLibExample creates one example under examples/lib/, with an example.go stub that already runs. |
| `RemoveLibExample` | `func(path string, name string) error` | RemoveLibExample deletes one example of examples/lib/ whole. |
| `ExecTest` | `func(props ExecTestProps) error` | ExecTest runs the project's examples and checks each one against its golden result.yaml, reporting every example that diverged. |
| `UpdateTest` | `func(path string, name string) error` | UpdateTest runs one example by name, both sides, and rewrites its golden result.yaml with what the run produced, printing the changes. |
| `Interview` | `func(path string) error` | Interview runs the interactive session over a project: it asks what is to be done, generates the questions from the declaration of the command that answers it, and runs that command with the answers bound onto it. It writes nothing of its own — every command it dispatches runs the action behind it, which persists and builds for itself. |

[every contract](doc.md)
