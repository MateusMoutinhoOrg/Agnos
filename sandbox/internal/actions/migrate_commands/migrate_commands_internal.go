package migrate_commands

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/commandconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/triggerconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// commandsDir holds one command per sub-directory.
const commandsDir = "sandbox/internal/commands"

// generatedCommands are the commands the build writes itself: their old files
// are dropped, and the build renders the new ones.
var generatedCommands = []string{"help", "version"}

// droppedFlags are the flags every command declared and the `project`
// middleware now declares once, in front of all of them.
var droppedFlags = []string{"path", "quiet"}

// renamedIds are the ids that would collide with the verb, which every
// migrated command reads as Entries.Command.
var renamedIds = map[string]string{"Command": "Target"}

// MigrateCommandsInternal rewrites every command under sandbox/internal/commands
// on an open SmartIO. A command already migrated — no entries.yaml — is left
// alone, so running it twice changes nothing.
func MigrateCommandsInternal(sandbox *api.Sandbox, io *smartio.SmartIO) error {
	for _, dir := range io.ListDirs(commandsDir) {
		parts := sandbox.Deps.Stringsdeps.Split(dir, "/")
		name := parts[len(parts)-1]
		base := commandsDir + "/" + name
		if name == "" || !io.IsFile(base+"/entries.yaml") {
			continue
		}

		if contains(generatedCommands, name) {
			io.RemoveDir(base + "/entries.yaml")
			io.RemoveDir(base + "/handler.go")
			continue
		}

		content, err := io.ReadFile(base + "/entries.yaml")
		if err != nil {
			return err
		}
		legacy, err := readLegacy(sandbox, string(content))
		if err != nil {
			return sandbox.Deps.Std.Errorf("commands/%s/entries.yaml: %w", name, err)
		}

		conf, sentinels := convert(sandbox, legacy)
		if err := io.WriteFileOverwrite(base+"/command.yaml", []byte(conf.Render())); err != nil {
			return err
		}

		handler, err := io.ReadFile(base + "/handler.go")
		if err != nil {
			return err
		}
		module_conf, err := utils.LoadModuleConf(sandbox, io)
		if err != nil {
			return err
		}
		rewritten := rewriteHandler(sandbox, string(handler), legacy, sentinels, module_conf.Module)
		if err := io.WriteFileOverwrite(base+"/InternalPureHandler.go", []byte(rewritten)); err != nil {
			return err
		}

		io.RemoveDir(base + "/entries.yaml")
		io.RemoveDir(base + "/handler.go")
		sandbox.Deps.Std.Log("migrate-commands %s \n", name)
	}
	return nil
}

// convert is one entries.yaml as a command.yaml: the verbs become the trigger
// of an arg on segment 0, each positional arg the segment after it, and each
// flag keeps its spellings under `keys`. It also reports the int flags that
// had no default and were read for their presence — they get -1, which the
// rewritten handler reads as "not given".
func convert(sandbox *api.Sandbox, legacy *legacyConf) (*commandconf.CommandConf, []string) {
	conf := commandconf.NewEmpty(sandbox)
	conf.Category = legacy.Category
	conf.Help = legacy.Help
	conf.LongDescription = legacy.LongDescription
	conf.Examples = legacy.Examples
	conf.Hidden = legacy.Hidden

	verbs := []string{}
	for _, identifier := range legacy.Identifiers {
		if !sandbox.Deps.Stringsdeps.HasPrefix(identifier, "-") {
			verbs = append(verbs, identifier)
		}
	}
	verb := commandconf.Arg{Id: "Command", Start: 0, End: 0, Type: commandconf.DefaultArgType}
	verb.Trigger = triggerconf.Trigger{Exists: true, Type: "equal", Values: []string{}}
	if len(verbs) == 1 {
		verb.Trigger.Value = verbs[0]
	} else {
		verb.Trigger.Type = triggerconf.OneOf
		verb.Trigger.Values = verbs
	}
	conf.Args = append(conf.Args, verb)

	for index, field := range legacy.Args {
		arg := commandconf.Arg{
			Id:          entryId(sandbox, field.Key),
			Start:       index + 1,
			End:         index + 1,
			Type:        argType(field.Type),
			Required:    field.Required,
			Default:     field.Default,
			HasDefault:  field.HasDefault,
			Description: field.Description,
		}
		if field.Array {
			arg.End = commandconf.LastSegment
			arg.Type = commandconf.DefaultArgType
		}
		conf.Args = append(conf.Args, arg)
	}

	sentinels := []string{}
	for _, field := range legacy.Flags {
		if contains(droppedFlags, field.Key) {
			continue
		}
		id := entryId(sandbox, field.Key)
		flag := commandconf.Flag{
			Id:          id,
			Keys:        field.Identifiers,
			Type:        flagType(field.Type, field.Array),
			Required:    field.Required,
			Default:     field.Default,
			HasDefault:  field.HasDefault,
			Min:         field.Min,
			HasMin:      field.HasMin,
			Max:         field.Max,
			HasMax:      field.HasMax,
			Enum:        []string{},
			Description: field.Description,
		}
		flag.HasKeys = !(len(field.Identifiers) == 1 && field.Identifiers[0] == commandconf.DefaultKey(sandbox, id))
		if len(field.Identifiers) == 0 {
			flag.Keys = []string{commandconf.DefaultKey(sandbox, id)}
			flag.HasKeys = false
		}
		if flag.Type == "integer" && !flag.HasDefault && !flag.Required {
			flag.Default, flag.HasDefault = "-1", true
			sentinels = append(sentinels, field.Key)
		}
		conf.Flags = append(conf.Flags, flag)
	}

	return conf, sentinels
}

// entryId is one entries.yaml field name as the Entries field it becomes:
// "max-bytes" -> MaxBytes, and "command", which the verb takes, -> Target.
func entryId(sandbox *api.Sandbox, key string) string {
	id := utils.RouteEntryId(sandbox, key)
	if renamed, is := renamedIds[id]; is {
		return renamed
	}
	return id
}

// argType is an entries.yaml arg type as a command.yaml one.
func argType(kind string) string {
	switch kind {
	case "int":
		return "integer"
	case "float":
		return "number"
	}
	return commandconf.DefaultArgType
}

// flagType is an entries.yaml flag type — with its array switch — as a
// command.yaml one.
func flagType(kind string, array bool) string {
	switch {
	case kind == "boolean":
		return "boolean"
	case kind == "int" && array:
		return "integer-array"
	case kind == "int":
		return "integer"
	case kind == "float":
		return "number"
	case array:
		return "string-array"
	}
	return commandconf.DefaultFlagType
}

// readers are the Get* of the old api.Command, each read back as a plain
// field of Entries.
var readers = []string{"GetStrings", "GetString", "GetBool", "GetInts", "GetInt", "GetFloats", "GetFloat"}

// rewriteHandler turns a handler.go into an InternalPureHandler.go: the new
// signature, every Get*("<name>") read off Entries (and --path off props),
// stdout through the response, and the ExitFailure pair a handler reported an
// error with through cliio.Fail. It is a text rewrite of a file every command
// wrote the same way; whatever it cannot rewrite is left for the compiler to
// point at.
func rewriteHandler(sandbox *api.Sandbox, source string, legacy *legacyConf, sentinels []string, module string) string {
	strs := sandbox.Deps.Stringsdeps

	source = strs.ReplaceAll(source,
		"func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {",
		"func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {")

	for _, key := range sentinels {
		source = strs.ReplaceAll(source, "command.GetItem(\""+key+"\") != nil", "entries."+entryId(sandbox, key)+" >= 0")
	}
	for _, reader := range readers {
		source = strs.ReplaceAll(source, "command."+reader+"(\"path\")", "props.Path")
		source = strs.ReplaceAll(source, "command."+reader+"(\"quiet\")", "false")
		for _, field := range append(append([]legacyField{}, legacy.Flags...), legacy.Args...) {
			source = strs.ReplaceAll(source, "command."+reader+"(\""+field.Key+"\")", "entries."+entryId(sandbox, field.Key))
		}
	}

	// The error-and-exit pair every handler reports a failed action with.
	lines := strs.Split(source, "\n")
	out := []string{}
	uses_cliio := false
	for index := 0; index < len(lines); index++ {
		line := lines[index]
		trimmed := strs.TrimSpace(line)
		indent := line[:len(line)-len(strs.TrimLeft(line, "\t"))]

		if strs.HasPrefix(trimmed, "sandbox.Deps.Std.Error(\"%s\\n\", ") && index+1 < len(lines) &&
			strs.TrimSpace(lines[index+1]) == "return api.ExitFailure" {
			message := strs.TrimSuffix(strs.TrimPrefix(trimmed, "sandbox.Deps.Std.Error(\"%s\\n\", "), ")")
			out = append(out, indent+"return cliio.Fail(sandbox, api.ExitFailure, \"\", "+message+")")
			uses_cliio = true
			index++
			continue
		}
		switch trimmed {
		case "return api.ExitOk":
			out = append(out, indent+"response.SetStatus(api.ExitOk)", indent+"return nil")
			continue
		case "return api.ExitFailure":
			out = append(out, indent+"return cliio.Fail(sandbox, api.ExitFailure, \"\", \"\")")
			uses_cliio = true
			continue
		case "return api.ExitUsage":
			out = append(out, indent+"return cliio.Fail(sandbox, api.ExitUsage, \"\", \"\")")
			uses_cliio = true
			continue
		}
		out = append(out, line)
	}
	source = strs.Join(out, "\n")

	source = strs.ReplaceAll(source, "sandbox.Deps.Std.Printf(", "response.Printf(")
	source = strs.ReplaceAll(source, "sandbox.Deps.Std.Error(", "response.Error(")

	if uses_cliio {
		parts := strs.Split(source, "import (\n")
		if len(parts) > 1 {
			source = parts[0] + "import (\n\t\"" + module + "/sandbox/internal/generated/cliio\"\n" + strs.Join(parts[1:], "import (\n")
		}
	}
	return source
}

func contains(values []string, value string) bool {
	for _, candidate := range values {
		if candidate == value {
			return true
		}
	}
	return false
}
