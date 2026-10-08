package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/commandconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// CommandDocField is one flag or one arg as docs/Commands prints it: the
// spellings a user types, a single label carrying type, arity and bounds, the
// default, the declared description, and — for a flag a middleware adds — the
// page of that middleware.
type CommandDocField struct {
	Id          string
	Keys        string
	Type        string
	Default     string
	Description string
	From        string
	FromPage    string
}

// CommandDocReach is one command a middleware runs in front of — or, on a
// command's page, one middleware that runs in front of it: its page, and the
// condition it runs on ("may run", "only when --x …"), "" when it always does.
type CommandDocReach struct {
	Name      string
	Page      string
	Condition string
}

// CommandDoc is one command's page of docs/Commands, rendered from its
// command.yaml and the middlewares crossed with it: nothing here is written by
// hand on the page.
type CommandDoc struct {
	Name        string
	Identifier  string
	Page        string
	Aliases     string
	Summary     string
	Description string
	Usage       string
	Pattern     string
	Priority    int
	Middleware  bool
	Flags       []CommandDocField
	Args        []CommandDocField
	Examples    []string
	// Middlewares are the middlewares that run in front of the command.
	Middlewares []CommandDocReach
	// RunsBefore are, on a middleware's page, the commands it runs in front
	// of.
	RunsBefore []CommandDocReach
}

// CommandDocGroup is one category section of docs/Commands, holding the
// commands that declare that category.
type CommandDocGroup struct {
	Category string
	Commands []CommandDoc
}

// CommandDocs is docs/Commands: the commands grouped by category, and the
// middlewares apart, since nobody types one.
type CommandDocs struct {
	Groups      []CommandDocGroup
	Middlewares []CommandDoc
}

// commandDocOther is the category a command with no declared one falls into,
// matching what the generated help screen prints for it.
const commandDocOther = "Other"

// commandDocMayRun is how a page words a middleware whose trigger cannot be
// crossed with a command without a command line.
const commandDocMayRun = "may run — `explain-command` gives the exact answer"

// commandDocEntry is one declared command, read once for every crossing.
type commandDocEntry struct {
	Name string
	Conf *commandconf.CommandConf
}

// CollectCommandDocs renders every command.yaml under sandbox/internal/commands
// into the pages docs/Commands prints: the commands grouped by category in
// first-seen order — the same grouping the generated help screen uses — and the
// middlewares on a section of their own. Hidden commands are skipped, exactly
// as they are in help. Every command is crossed with every middleware whose
// trigger holds on it, so its page lists the flags they add, and a
// middleware's page the commands it runs in front of.
//
// The declaration is the only source: a command, a flag or an example reaches
// the page by being declared with `add-command`, `add-flag`, `add-arg` or
// `set-command`, never by the page being edited.
func CollectCommandDocs(sandbox *api.Sandbox, io *stagedfs.StagedFS) (CommandDocs, error) {
	docs := CommandDocs{}
	entries := []commandDocEntry{}

	for _, unit := range utils.CommandDirs(sandbox, io) {
		content, err := io.ReadFile(unit.Dir + "/" + utils.CommandConfFile)
		if err != nil {
			continue
		}
		conf, err := commandconf.New(sandbox, string(content))
		if err != nil {
			return docs, sandbox.Deps.StdDeps.Errorf("%s/%s: %w", unit.Dir, utils.CommandConfFile, err)
		}
		if conf.Hidden {
			continue
		}
		entries = append(entries, commandDocEntry{Name: unit.Name, Conf: conf})
	}

	index := map[string]int{}
	for _, entry := range entries {
		doc := commandDoc(sandbox, entry, entries)
		if doc.Middleware {
			docs.Middlewares = append(docs.Middlewares, doc)
			continue
		}

		category := entry.Conf.Category
		if category == "" {
			category = commandDocOther
		}
		position, seen := index[category]
		if !seen {
			position = len(docs.Groups)
			index[category] = position
			docs.Groups = append(docs.Groups, CommandDocGroup{Category: category})
		}
		docs.Groups[position].Commands = append(docs.Groups[position].Commands, doc)
	}

	return docs, nil
}

// commandDoc turns one parsed declaration into its page, crossed with every
// other command of the project.
func commandDoc(sandbox *api.Sandbox, entry commandDocEntry, entries []commandDocEntry) CommandDoc {
	conf := entry.Conf
	identifiers := conf.Identifiers()
	identifier := utils.CommandName(sandbox, entry.Name)
	if len(identifiers) > 0 && conf.Strict {
		identifier = identifiers[0]
	}

	doc := CommandDoc{
		Name:        entry.Name,
		Identifier:  identifier,
		Page:        commandDocPage(sandbox, entry),
		Aliases:     identifierList(sandbox, aliasesOf(identifiers)),
		Summary:     docCell(sandbox, conf.Summary),
		Description: docText(sandbox, conf.Description),
		Pattern:     conf.Pattern(),
		Priority:    conf.Priority,
		Middleware:  !conf.Strict,
		Examples:    conf.Examples,
	}

	for _, arg := range conf.Args {
		if arg.Trigger.Set {
			continue
		}
		doc.Args = append(doc.Args, CommandDocField{
			Id:          arg.Id,
			Type:        argTypeLabel(arg),
			Default:     defaultCell(arg.Default, arg.HasDefault),
			Description: docCell(sandbox, arg.Description),
		})
	}
	for _, flag := range conf.Flags {
		doc.Flags = append(doc.Flags, commandDocFlag(sandbox, flag, nil))
	}

	inherited := []commandconf.Flag{}
	for _, other := range entries {
		if other.Name == entry.Name {
			continue
		}
		if doc.Middleware {
			if other.Conf.Strict {
				if reach, condition := utils.CommandMiddlewareReach(sandbox, conf, other.Conf); reach != utils.NoReach {
					doc.RunsBefore = append(doc.RunsBefore, CommandDocReach{
						Name:      commandDocTitle(sandbox, other),
						Page:      commandDocPage(sandbox, other),
						Condition: reachCondition(reach, condition),
					})
				}
			}
			continue
		}

		reach, condition := utils.CommandMiddlewareReach(sandbox, other.Conf, conf)
		if reach == utils.NoReach {
			continue
		}
		doc.Middlewares = append(doc.Middlewares, CommandDocReach{
			Name:      commandDocTitle(sandbox, other),
			Page:      commandDocPage(sandbox, other),
			Condition: reachCondition(reach, condition),
		})
		for _, flag := range other.Conf.Flags {
			if utils.CommandKeyTaken(conf, flag.Keys) != "" {
				continue
			}
			field := commandDocFlag(sandbox, flag, &other)
			field.FromPage = commandDocPage(sandbox, other)
			doc.Flags = append(doc.Flags, field)
			inherited = append(inherited, flag)
		}
	}

	doc.Usage = commandUsage(sandbox, conf, inherited)
	return doc
}

// commandDocTitle is how a page names another command: its verb, or the
// package name of a middleware, which has none.
func commandDocTitle(sandbox *api.Sandbox, entry commandDocEntry) string {
	identifiers := entry.Conf.Identifiers()
	if len(identifiers) > 0 && entry.Conf.Strict {
		return identifiers[0]
	}
	return utils.CommandName(sandbox, entry.Name)
}

// commandDocPage is the file one command's page is written to, relative to
// docs/Commands.
func commandDocPage(sandbox *api.Sandbox, entry commandDocEntry) string {
	return sandbox.Deps.StringsDeps.ReplaceAll(commandDocTitle(sandbox, entry), " ", "-") + ".md"
}

// reachCondition is the cell a crossing is worded with: "" when the
// middleware always runs.
func reachCondition(reach utils.MiddlewareReach, condition string) string {
	switch {
	case reach == utils.MayRun && condition != "":
		return commandDocMayRun + "; " + condition
	case reach == utils.MayRun:
		return commandDocMayRun
	}
	return condition
}

// aliasesOf is every verb of a command but the first — the one the page
// titles the section with.
func aliasesOf(identifiers []string) []string {
	if len(identifiers) < 2 {
		return nil
	}
	return identifiers[1:]
}

// commandDocFlag renders one flag as its table row; from is the middleware
// declaring it, nil for one of the command's own.
func commandDocFlag(sandbox *api.Sandbox, flag commandconf.Flag, from *commandDocEntry) CommandDocField {
	field := CommandDocField{
		Id:          flag.Id,
		Keys:        identifierList(sandbox, flag.Keys),
		Type:        flagTypeLabel(sandbox, flag),
		Default:     defaultCell(flag.Default, flag.HasDefault),
		Description: docCell(sandbox, flag.Description),
	}
	if from != nil {
		field.From = commandDocTitle(sandbox, *from)
	}
	return field
}

// defaultCell is a declared default as a table cell.
func defaultCell(value string, has bool) string {
	if !has {
		return ""
	}
	return "`" + value + "`"
}

// argTypeLabel is the cell carrying an arg's type: a range reads as a list.
func argTypeLabel(arg commandconf.Arg) string {
	label := arg.Type
	if label == "" {
		label = commandconf.DefaultArgType
	}
	if arg.End != arg.Start {
		label += ", repeatable"
	}
	if arg.Required {
		label += ", required"
	}
	return label
}

// flagTypeLabel is the one cell carrying everything the type of a flag
// implies: its kind, whether it must be given, its bounds and its enum.
func flagTypeLabel(sandbox *api.Sandbox, flag commandconf.Flag) string {
	label := flag.Type
	if label == "" {
		label = commandconf.DefaultFlagType
	}
	if flag.Required {
		label += ", required"
	}
	switch {
	case flag.HasMin && flag.HasMax:
		label += ", " + numberLabel(sandbox, flag.Type, flag.Min, true) + ".." + numberLabel(sandbox, flag.Type, flag.Max, true)
	case flag.HasMin:
		label += ", >= " + numberLabel(sandbox, flag.Type, flag.Min, true)
	case flag.HasMax:
		label += ", <= " + numberLabel(sandbox, flag.Type, flag.Max, true)
	}
	if len(flag.Enum) > 0 {
		label += ", one of " + sandbox.Deps.StringsDeps.Join(flag.Enum, "/")
	}
	return label
}

// commandUsage builds the one usage line of a command: its pattern, then every
// flag it declares in declaration order and every flag a middleware in front
// of it adds, each bracketed when optional.
func commandUsage(sandbox *api.Sandbox, conf *commandconf.CommandConf, inherited []commandconf.Flag) string {
	parts := []string{conf.Pattern()}
	for _, flag := range append(append([]commandconf.Flag{}, conf.Flags...), inherited...) {
		token := flagToken(sandbox, flag)
		if !flag.Required {
			token = "[" + token + "]"
		}
		parts = append(parts, token)
	}
	return sandbox.Deps.StringsDeps.Join(parts, " ")
}

// flagToken is one flag as the usage line spells it: its first key, plus a
// value placeholder for everything that is not a boolean switch.
func flagToken(sandbox *api.Sandbox, flag commandconf.Flag) string {
	key := flag.Keys[0]
	if flag.Type == "boolean" {
		return key
	}
	token := key + " <" + sandbox.Deps.StringsDeps.TrimLeft(key, "-") + ">"
	if flag.Type == "string-array" || flag.Type == "integer-array" {
		token += "..."
	}
	return token
}

// identifierList renders spellings as the inline code list a table cell carries, ""
// when there are none.
func identifierList(sandbox *api.Sandbox, keys []string) string {
	quoted := make([]string, 0, len(keys))
	for _, key := range keys {
		quoted = append(quoted, "`"+key+"`")
	}
	return sandbox.Deps.StringsDeps.Join(quoted, ", ")
}

// docCell flattens a declared description into one markdown table cell:
// no line breaks, and no bar to close the cell early.
func docCell(sandbox *api.Sandbox, raw string) string {
	flat := sandbox.Deps.StringsDeps.Join(sandbox.Deps.StringsDeps.Fields(raw), " ")
	return sandbox.Deps.StringsDeps.ReplaceAll(flat, "|", `\|`)
}

// docText normalizes a long description into paragraphs: the declaration
// hard-wraps its lines, which markdown would keep, so each blank-line-separated
// block is joined back into one line.
func docText(sandbox *api.Sandbox, raw string) string {
	blocks := sandbox.Deps.StringsDeps.Split(sandbox.Deps.StringsDeps.TrimSpace(raw), "\n\n")
	joined := make([]string, 0, len(blocks))
	for _, block := range blocks {
		if text := sandbox.Deps.StringsDeps.Join(sandbox.Deps.StringsDeps.Fields(block), " "); text != "" {
			joined = append(joined, text)
		}
	}
	return sandbox.Deps.StringsDeps.Join(joined, "\n\n")
}
