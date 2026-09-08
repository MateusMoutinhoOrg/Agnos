package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/commandconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// CommandDocField is one flag or one positional argument as docs/Commands
// prints it: the identifiers a user types, a single label carrying type,
// arity and bounds, the default, and the declared description.
type CommandDocField struct {
	Key         string
	Identifiers string
	Type        string
	Default     string
	Description string
}

// CommandDoc is one command's section of docs/Commands, rendered from its
// entries.yaml alone: nothing here is written by hand on the page.
type CommandDoc struct {
	Name            string
	Identifier      string
	Aliases         string
	Help            string
	LongDescription string
	Usage           string
	Flags           []CommandDocField
	Args            []CommandDocField
	Examples        []string
}

// CommandDocGroup is one category section of docs/Commands, holding the
// commands that declare that category.
type CommandDocGroup struct {
	Category string
	Commands []CommandDoc
}

// commandsDocDir holds one declared command per sub-directory.
const commandsDocDir = "sandbox/internal/commands"

// commandDocOther is the category a command with no declared one falls into,
// matching what the generated help screen prints for it.
const commandDocOther = "Other"

// CollectCommandDocs renders every sandbox/internal/commands/<name>/entries.yaml
// into the sections docs/Commands prints, grouped by category in first-seen
// order — the same grouping the generated help screen uses. Hidden commands
// are skipped, exactly as they are in help.
//
// The declaration is the only source: a command, a flag or an example reaches
// the page by being declared with `add-command`, `add-flag`, `add-arg` or
// `set-command`, never by the page being edited.
func CollectCommandDocs(deps *deps.Deps, io *smartio.SmartIO) ([]CommandDocGroup, error) {
	var groups []CommandDocGroup
	index := map[string]int{}

	for _, dir := range io.ListDirs(commandsDocDir) {
		name := lastSegmentOf(deps, dir)
		if name == "" {
			continue
		}

		content, err := io.ReadFile(commandsDocDir + "/" + name + "/entries.yaml")
		if err != nil {
			continue
		}

		conf, err := commandconf.New(deps, string(content))
		if err != nil {
			return nil, deps.Std.Errorf("commands/%s/entries.yaml: %w", name, err)
		}

		if conf.Hidden {
			continue
		}

		category := conf.Category
		if category == "" {
			category = commandDocOther
		}

		position, seen := index[category]
		if !seen {
			position = len(groups)
			index[category] = position
			groups = append(groups, CommandDocGroup{Category: category})
		}

		groups[position].Commands = append(groups[position].Commands, commandDoc(deps, name, conf))
	}

	return groups, nil
}

// commandDoc turns one parsed declaration into its page section.
func commandDoc(deps *deps.Deps, name string, conf *commandconf.CommandConf) CommandDoc {
	identifier := name
	if len(conf.Identifiers) > 0 {
		identifier = conf.Identifiers[0]
	}

	doc := CommandDoc{
		Name:            name,
		Identifier:      identifier,
		Aliases:         identifierList(deps, aliasesOf(conf.Identifiers)),
		Help:            docCell(deps, conf.Help),
		LongDescription: docText(deps, conf.LongDescription),
		Examples:        conf.Examples,
	}

	for _, flag := range conf.Flags {
		doc.Flags = append(doc.Flags, commandDocField(deps, flag))
	}
	for _, arg := range conf.Args {
		doc.Args = append(doc.Args, commandDocField(deps, arg))
	}

	doc.Usage = commandUsage(deps, identifier, conf)

	return doc
}

// aliasesOf is every identifier of a command but the first — the verb the
// page titles the section with.
func aliasesOf(identifiers []string) []string {
	if len(identifiers) < 2 {
		return nil
	}
	return identifiers[1:]
}

// commandDocField renders one flag or positional as its table row.
func commandDocField(deps *deps.Deps, field commandconf.Field) CommandDocField {
	value := ""
	if field.HasDefault {
		value = "`" + field.Default + "`"
	}

	return CommandDocField{
		Key:         field.Key,
		Identifiers: identifierList(deps, field.Identifiers),
		Type:        fieldTypeLabel(deps, field),
		Default:     value,
		Description: docCell(deps, field.Description),
	}
}

// fieldTypeLabel is the one cell carrying everything the type of a field
// implies: its kind, whether it repeats, whether it must be given, and the
// bounds a numeric field declares.
func fieldTypeLabel(deps *deps.Deps, field commandconf.Field) string {
	label := field.Type
	if label == "" {
		label = "string"
	}
	if field.Array {
		label += ", repeatable"
	}
	if field.Required {
		label += ", required"
	}
	if bounds := fieldBounds(deps, field); bounds != "" {
		label += ", " + bounds
	}
	return label
}

// fieldBounds spells the min/max a numeric field declares, "" when it
// declares neither.
func fieldBounds(deps *deps.Deps, field commandconf.Field) string {
	switch {
	case field.HasMin && field.HasMax:
		return numberLabel(deps, field, field.Min, true) + ".." + numberLabel(deps, field, field.Max, true)
	case field.HasMin:
		return ">= " + numberLabel(deps, field, field.Min, true)
	case field.HasMax:
		return "<= " + numberLabel(deps, field, field.Max, true)
	}
	return ""
}

// commandUsage builds the one usage line of a command: its verb, every flag
// in declaration order (bracketed when optional) and every positional after
// them.
func commandUsage(deps *deps.Deps, identifier string, conf *commandconf.CommandConf) string {
	parts := []string{identifier}

	for _, flag := range conf.Flags {
		token := flagToken(flag)
		if !flag.Required {
			token = "[" + token + "]"
		}
		parts = append(parts, token)
	}

	for _, arg := range conf.Args {
		token := "<" + arg.Key + ">"
		if arg.Array {
			token += "..."
		}
		if !arg.Required {
			token = "[" + token + "]"
		}
		parts = append(parts, token)
	}

	return deps.Stringsdeps.Join(parts, " ")
}

// flagToken is one flag as the usage line spells it: its first identifier,
// plus a value placeholder for everything that is not a boolean switch.
func flagToken(flag commandconf.Field) string {
	identifier := "--" + flag.Key
	if len(flag.Identifiers) > 0 {
		identifier = flag.Identifiers[0]
	}
	if flag.Type == "boolean" {
		return identifier
	}
	token := identifier + " <" + flag.Key + ">"
	if flag.Array {
		token += "..."
	}
	return token
}

// identifierList renders identifiers as the inline code list a table cell
// carries, "" when there are none (a positional argument).
func identifierList(deps *deps.Deps, identifiers []string) string {
	quoted := make([]string, 0, len(identifiers))
	for _, identifier := range identifiers {
		quoted = append(quoted, "`"+identifier+"`")
	}
	return deps.Stringsdeps.Join(quoted, ", ")
}

// docCell flattens a declared description into one markdown table cell:
// no line breaks, and no bar to close the cell early.
func docCell(deps *deps.Deps, raw string) string {
	flat := deps.Stringsdeps.Join(deps.Stringsdeps.Fields(raw), " ")
	return deps.Stringsdeps.ReplaceAll(flat, "|", `\|`)
}

// docText normalizes a long description into paragraphs: the declaration
// hard-wraps its lines, which markdown would keep, so each blank-line-separated
// block is joined back into one line.
func docText(deps *deps.Deps, raw string) string {
	blocks := deps.Stringsdeps.Split(deps.Stringsdeps.TrimSpace(raw), "\n\n")
	joined := make([]string, 0, len(blocks))
	for _, block := range blocks {
		if text := deps.Stringsdeps.Join(deps.Stringsdeps.Fields(block), " "); text != "" {
			joined = append(joined, text)
		}
	}
	return deps.Stringsdeps.Join(joined, "\n\n")
}
