package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/commandconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// CommandIdentifier normalizes a user-typed command name into the CLI verb:
// lowercased, spaces and underscores turned into dashes
// ("My Feature" -> "my-feature").
func CommandIdentifier(sandbox *api.Sandbox, name string) string {
	out := sandbox.Deps.Stringsdeps.ToLower(sandbox.Deps.Stringsdeps.TrimSpace(name))
	out = sandbox.Deps.Stringsdeps.ReplaceAll(out, " ", "-")
	out = sandbox.Deps.Stringsdeps.ReplaceAll(out, "_", "-")
	return out
}

// ValidateCommandName reports whether a user-typed command name normalizes to
// a usable CLI verb. The identifier becomes a directory name and a Go package
// clause, so anything outside [a-z][a-z0-9-]* would be propagated straight
// into `package bad_name!` and break the whole project's build. The check runs
// before any file is written.
func ValidateCommandName(sandbox *api.Sandbox, name string) error {
	identifier := CommandIdentifier(sandbox, name)
	if identifier == "" {
		return sandbox.Deps.Std.Errorf("a command needs a name")
	}
	if identifier[0] < 'a' || identifier[0] > 'z' {
		return sandbox.Deps.Std.Errorf("invalid command name %q: a command name must start with a lowercase letter", name)
	}
	for _, letter := range identifier {
		valid := (letter >= 'a' && letter <= 'z') ||
			(letter >= '0' && letter <= '9') ||
			letter == '-'
		if !valid {
			return sandbox.Deps.Std.Errorf(
				"invalid command name %q: only letters, digits, spaces, dashes and underscores are allowed (it becomes the directory %s of sandbox/internal/commands and a Go package name)",
				name, CommandPackage(sandbox, name))
		}
	}
	return nil
}

// NoteNormalizedCommandName logs the rewriting a command name went through,
// so a name that is silently lowercased or re-punctuated ("MyCmd" -> "mycmd")
// is never a surprise later.
func NoteNormalizedCommandName(sandbox *api.Sandbox, name string) {
	NoteNormalizedName(sandbox, "command", name)
}

// NoteNormalizedName is NoteNormalizedCommandName for any unit named the same
// way — a route, a database, a table, a field — kind naming which.
func NoteNormalizedName(sandbox *api.Sandbox, kind string, name string) {
	identifier := CommandIdentifier(sandbox, name)
	if identifier != name {
		sandbox.Deps.Std.Log("note: %s name %q normalized to %q \n", kind, name, identifier)
	}
}

// CommandPackage is the Go package / directory name for a command: the
// identifier with dashes turned into underscores ("my-feature" -> "my_feature").
func CommandPackage(sandbox *api.Sandbox, name string) string {
	return sandbox.Deps.Stringsdeps.ReplaceAll(CommandIdentifier(sandbox, name), "-", "_")
}

// CommandsDir is the tree the commands are declared in: every directory at
// any depth holding a command.yaml is a command, and every other one a folder
// grouping them.
const CommandsDir = "sandbox/internal/commands"

// CommandConfFile is the declaration of one command, beside its generated
// new.go and entries.go and its hand-written InternalPureHandler.go, and what
// makes its directory a command.
const CommandConfFile = "command.yaml"

// CommandHandlerFile is the one hand-written file of a command.
const CommandHandlerFile = "InternalPureHandler.go"

// CommandDirs is every command declared under CommandsDir, at any depth.
func CommandDirs(sandbox *api.Sandbox, io *smartio.SmartIO) []UnitDir {
	return FindUnitDirs(sandbox, io, CommandsDir, CommandConfFile)
}

// CommandDir is the project-relative directory holding the command named —
// by its package name or one of its verbs — in whatever folder it sits; one
// no command.yaml declares yet lands at the top of CommandsDir.
func CommandDir(sandbox *api.Sandbox, io *smartio.SmartIO, name string) string {
	pkg := CommandPackage(sandbox, ResolveCommandName(sandbox, io, name))
	if dir, found := FindUnitDir(sandbox, io, CommandsDir, CommandConfFile, pkg); found {
		return dir
	}
	return CommandsDir + "/" + pkg
}

// CommandConfPath is the project-relative path of a command's command.yaml.
func CommandConfPath(sandbox *api.Sandbox, io *smartio.SmartIO, name string) string {
	return CommandDir(sandbox, io, name) + "/" + CommandConfFile
}

// ResolveCommandName is the package name of the command a user named: its
// package name, or one of the verbs it answers to ("add-flag", "help"). It is
// returned as given when nothing answers to it, so the caller's own "not
// found" names what was typed.
func ResolveCommandName(sandbox *api.Sandbox, io *smartio.SmartIO, name string) string {
	units := CommandDirs(sandbox, io)
	pkg := CommandPackage(sandbox, name)
	for _, unit := range units {
		if unit.Name == pkg {
			return name
		}
	}
	for _, unit := range units {
		content, err := io.ReadFile(unit.Dir + "/" + CommandConfFile)
		if err != nil {
			continue
		}
		conf, err := commandconf.New(sandbox, string(content))
		if err != nil {
			continue
		}
		for _, identifier := range conf.Identifiers() {
			if identifier == sandbox.Deps.Stringsdeps.TrimSpace(name) {
				return unit.Name
			}
		}
	}
	return name
}

// LoadCommandConf reads and parses the command.yaml of the command named —
// by its package name or one of its verbs — wherever under
// sandbox/internal/commands it sits.
func LoadCommandConf(sandbox *api.Sandbox, io *smartio.SmartIO, name string) (*commandconf.CommandConf, error) {
	name = ResolveCommandName(sandbox, io, name)
	if err := ValidateCommandName(sandbox, name); err != nil {
		return nil, err
	}
	path := CommandConfPath(sandbox, io, name)
	content, err := io.ReadFile(path)
	if err != nil {
		return nil, sandbox.Deps.Std.Errorf("command %q not found in %s", CommandIdentifier(sandbox, name), CommandsDir)
	}
	conf, err := commandconf.New(sandbox, string(content))
	if err != nil {
		return nil, sandbox.Deps.Std.Errorf("%s: %w", path, err)
	}
	return conf, nil
}

// SaveCommandConf renders conf back over the command.yaml of the command
// named.
func SaveCommandConf(sandbox *api.Sandbox, io *smartio.SmartIO, name string, conf *commandconf.CommandConf) error {
	return io.WriteFileOverwrite(CommandConfPath(sandbox, io, name), []byte(conf.Render()))
}

// CommandReservedIds are the Entries fields the generated entries.go spells
// itself, so no arg or flag may take them.
var CommandReservedIds = []string{"FullCommand"}

// CommandEntryId turns an arg or flag name typed on the command line into the
// exported Go name its Entries field carries: "out-file" -> "OutFile".
func CommandEntryId(sandbox *api.Sandbox, raw string) string {
	return RouteEntryId(sandbox, raw)
}

// CommandIdTaken reports whether id already names an arg or a flag of conf, or
// is reserved.
func CommandIdTaken(conf *commandconf.CommandConf, id string) bool {
	if contains(CommandReservedIds, id) {
		return true
	}
	return FindCommandArg(conf, id) >= 0 || FindCommandFlag(conf, id) >= 0
}

// FindCommandArg is the index of the arg of conf with that id, -1 when none.
func FindCommandArg(conf *commandconf.CommandConf, id string) int {
	for index, arg := range conf.Args {
		if arg.Id == id {
			return index
		}
	}
	return -1
}

// FindCommandFlag is the index of the flag of conf with that id — or answering
// to that key ("--out", "-o") — -1 when none.
func FindCommandFlag(conf *commandconf.CommandConf, id string) int {
	for index, flag := range conf.Flags {
		if flag.Id == id {
			return index
		}
		for _, key := range flag.Keys {
			if key == id {
				return index
			}
		}
	}
	return -1
}

// CommandKeyTaken is the flag of conf answering to one of keys, "" when none
// does.
func CommandKeyTaken(conf *commandconf.CommandConf, keys []string) string {
	for _, flag := range conf.Flags {
		for _, own := range flag.Keys {
			if contains(keys, own) {
				return own
			}
		}
	}
	return ""
}

// NextCommandSegment is the first segment no arg of conf reads yet — where
// `add-arg` puts an arg given no --start. It is -1 when an arg reads to the
// last segment, since nothing may come after it.
func NextCommandSegment(conf *commandconf.CommandConf) int {
	next := 0
	for _, arg := range conf.Args {
		if arg.End == commandconf.LastSegment {
			return -1
		}
		if arg.End+1 > next {
			next = arg.End + 1
		}
	}
	return next
}

// CommandFlagType maps the type spellings accepted on the command line onto
// the canonical command.yaml set; "" defaults to string, and array asks for
// the repeatable form of the type.
func CommandFlagType(sandbox *api.Sandbox, raw string, array bool) (string, bool) {
	kind := ""
	switch sandbox.Deps.Stringsdeps.ToLower(sandbox.Deps.Stringsdeps.TrimSpace(raw)) {
	case "", "string", "str":
		kind = "string"
	case "bool", "boolean":
		kind = "boolean"
	case "int", "integer":
		kind = "integer"
	case "float", "double", "number":
		kind = "number"
	case "string-array", "strings":
		kind = "string-array"
	case "integer-array", "int-array", "ints":
		kind = "integer-array"
	default:
		return "", false
	}
	if array {
		switch kind {
		case "string":
			kind = "string-array"
		case "integer":
			kind = "integer-array"
		case "string-array", "integer-array":
		default:
			return "", false
		}
	}
	return kind, true
}

// CommandArgType maps an arg type typed on the command line onto the
// canonical set; "" defaults to string.
func CommandArgType(sandbox *api.Sandbox, raw string) (string, bool) {
	switch sandbox.Deps.Stringsdeps.ToLower(sandbox.Deps.Stringsdeps.TrimSpace(raw)) {
	case "", "string", "str":
		return "string", true
	case "int", "integer":
		return "integer", true
	case "float", "double", "number":
		return "number", true
	case "uuid":
		return "uuid", true
	}
	return "", false
}

// CheckCommandLiteral reports a default that does not read as its type.
func CheckCommandLiteral(sandbox *api.Sandbox, kind string, label string, raw string) error {
	switch kind {
	case "boolean":
		if raw != "true" && raw != "false" {
			return sandbox.Deps.Std.Errorf("%s for a boolean must be true or false, got %q", label, raw)
		}
	case "integer", "integer-array":
		if _, err := sandbox.Deps.Stringsdeps.ParseInt(raw, 10, 64); err != nil {
			return sandbox.Deps.Std.Errorf("%s must be an integer, got %q", label, raw)
		}
	case "number":
		if value, err := sandbox.Deps.Stringsdeps.ParseFloat(raw, 64); err != nil || value-value != 0 {
			return sandbox.Deps.Std.Errorf("%s must be a number, got %q", label, raw)
		}
	case "uuid":
		matched, err := sandbox.Deps.Stringsdeps.MatchPattern(routeUuidPattern, raw)
		if err != nil || !matched {
			return sandbox.Deps.Std.Errorf("%s must be a uuid, got %q", label, raw)
		}
	}
	return nil
}

// ParseCommandBound reads a --min or --max: a number, for a numeric flag only.
func ParseCommandBound(sandbox *api.Sandbox, kind string, label string, raw string) (float64, error) {
	if kind != "integer" && kind != "number" && kind != "integer-array" {
		return 0, sandbox.Deps.Std.Errorf("--%s only applies to an integer or a number flag", label)
	}
	value, err := sandbox.Deps.Stringsdeps.ParseFloat(sandbox.Deps.Stringsdeps.TrimSpace(raw), 64)
	if err != nil {
		return 0, sandbox.Deps.Std.Errorf("--%s must be a number, got %q", label, raw)
	}
	return value, nil
}

// AppendPosition is the --position value meaning "put it last". It is the
// flag's default, so "not given" and "at the end" are the same request.
const AppendPosition = -1

// CheckPosition validates a --position against a list of size entries:
// AppendPosition means the end, and anything else must land inside 0..size.
func CheckPosition(sandbox *api.Sandbox, kind string, position int, size int) (int, error) {
	if position == AppendPosition {
		return size, nil
	}
	if position < 0 {
		return 0, sandbox.Deps.Std.Errorf("--position %d is negative: use an index from 0 to %d, or leave it out to append", position, size)
	}
	if position > size {
		return 0, sandbox.Deps.Std.Errorf("--position %d is out of range: this command has %d %s(s), so the accepted range is 0 to %d", position, size, kind, size)
	}
	return position, nil
}

// InsertCommandArg places arg at position inside args.
func InsertCommandArg(args []commandconf.Arg, arg commandconf.Arg, position int) []commandconf.Arg {
	if position < 0 || position >= len(args) {
		return append(args, arg)
	}
	out := make([]commandconf.Arg, 0, len(args)+1)
	out = append(out, args[:position]...)
	out = append(out, arg)
	return append(out, args[position:]...)
}

// InsertCommandFlag places flag at position inside flags.
func InsertCommandFlag(flags []commandconf.Flag, flag commandconf.Flag, position int) []commandconf.Flag {
	if position < 0 || position >= len(flags) {
		return append(flags, flag)
	}
	out := make([]commandconf.Flag, 0, len(flags)+1)
	out = append(out, flags[:position]...)
	out = append(out, flag)
	return append(out, flags[position:]...)
}

// RemoveCommandArg drops the arg at index.
func RemoveCommandArg(args []commandconf.Arg, index int) []commandconf.Arg {
	out := make([]commandconf.Arg, 0, len(args)-1)
	out = append(out, args[:index]...)
	return append(out, args[index+1:]...)
}

// RemoveCommandFlag drops the flag at index.
func RemoveCommandFlag(flags []commandconf.Flag, index int) []commandconf.Flag {
	out := make([]commandconf.Flag, 0, len(flags)-1)
	out = append(out, flags[:index]...)
	return append(out, flags[index+1:]...)
}

// AppendUnique appends each value of extra to values, skipping duplicates.
func AppendUnique(values []string, extra []string) []string {
	for _, candidate := range extra {
		found := false
		for _, existing := range values {
			if existing == candidate {
				found = true
				break
			}
		}
		if !found {
			values = append(values, candidate)
		}
	}
	return values
}

// NewCommandFlag builds one flag from what was typed on the command line: its
// Entries id, its keys (--<name> when none is given, each starting with "-"),
// its type, and the default, bounds, enum, pattern and trigger it declares.
func NewCommandFlag(sandbox *api.Sandbox, props api.FlagProps) (commandconf.Flag, error) {
	strs := sandbox.Deps.Stringsdeps
	flag := commandconf.Flag{
		Id:          CommandEntryId(sandbox, props.Name),
		Required:    props.Required,
		Enum:        []string{},
		Pattern:     strs.TrimSpace(props.Pattern),
		Description: strs.TrimSpace(props.Description),
	}
	if err := ValidateEntryId(sandbox, "flag", props.Name, flag.Id); err != nil {
		return flag, err
	}

	kind, ok := CommandFlagType(sandbox, props.Type, props.Array)
	if !ok {
		return flag, sandbox.Deps.Std.Errorf("unknown flag type %q (use %s)", props.Type, strs.Join(commandconf.FlagTypes, ", "))
	}
	flag.Type = kind

	flag.Keys = []string{}
	for _, key := range props.Keys {
		key = strs.TrimSpace(key)
		if !strs.HasPrefix(key, "-") {
			return flag, sandbox.Deps.Std.Errorf("flag key %q must start with - or --", key)
		}
		flag.Keys = AppendUnique(flag.Keys, []string{key})
	}
	default_key := commandconf.DefaultKey(sandbox, flag.Id)
	flag.HasKeys = !(len(flag.Keys) == 0 || (len(flag.Keys) == 1 && flag.Keys[0] == default_key))
	if len(flag.Keys) == 0 {
		flag.Keys = []string{default_key}
	}

	if flag.Required && kind == "boolean" {
		return flag, sandbox.Deps.Std.Errorf("a boolean flag cannot be required (its absence already means false)")
	}
	if props.Default != "" {
		if flag.Required {
			return flag, sandbox.Deps.Std.Errorf("a flag cannot be both required and carry a default (the default already covers its absence)")
		}
		if err := CheckCommandLiteral(sandbox, kind, "default", props.Default); err != nil {
			return flag, err
		}
		flag.Default, flag.HasDefault = props.Default, true
	}
	if props.Min != "" {
		value, err := ParseCommandBound(sandbox, kind, "min", props.Min)
		if err != nil {
			return flag, err
		}
		flag.Min, flag.HasMin = value, true
	}
	if props.Max != "" {
		value, err := ParseCommandBound(sandbox, kind, "max", props.Max)
		if err != nil {
			return flag, err
		}
		flag.Max, flag.HasMax = value, true
	}
	if flag.HasMin && flag.HasMax && flag.Min > flag.Max {
		return flag, sandbox.Deps.Std.Errorf("min (%s) is greater than max (%s)", props.Min, props.Max)
	}
	for _, value := range props.Enum {
		if value = strs.TrimSpace(value); value != "" {
			flag.Enum = AppendUnique(flag.Enum, []string{value})
		}
	}
	if len(flag.Enum) > 0 && kind == "boolean" {
		return flag, sandbox.Deps.Std.Errorf("--enum does not apply to a boolean flag")
	}
	// A default outside the enum fails every run of the command, even one that
	// never passes the flag: the dispatch validates the default like a value.
	if flag.HasDefault && len(flag.Enum) > 0 && !contains(flag.Enum, flag.Default) {
		return flag, sandbox.Deps.Std.Errorf("default %q is not one of the --enum values (%s)", flag.Default, strs.Join(flag.Enum, ", "))
	}
	if flag.Pattern != "" {
		if _, err := strs.MatchPattern(flag.Pattern, ""); err != nil {
			return flag, sandbox.Deps.Std.Errorf("invalid --pattern %q: %s", flag.Pattern, err.Error())
		}
	}

	trigger, err := NewTrigger(sandbox, TriggerProps{
		Type:       props.TriggerType,
		Value:      props.Trigger,
		Negate:     props.TriggerNegate,
		IgnoreCase: props.TriggerIgnoreCase,
	})
	if err != nil {
		return flag, err
	}
	flag.Trigger = trigger
	return flag, nil
}

// NewCommandArg builds one arg from what was typed on the command line: its
// Entries id, the segments it reads — from next, the first no arg reads yet,
// when no start is given, to its start when no end is — its type, and the
// default and trigger it declares.
func NewCommandArg(sandbox *api.Sandbox, props api.ArgProps, next int) (commandconf.Arg, error) {
	strs := sandbox.Deps.Stringsdeps
	arg := commandconf.Arg{
		Id:          CommandEntryId(sandbox, props.Name),
		Required:    props.Required,
		Description: strs.TrimSpace(props.Description),
	}
	if err := ValidateEntryId(sandbox, "arg", props.Name, arg.Id); err != nil {
		return arg, err
	}

	arg.Start = next
	if raw := strs.TrimSpace(props.Start); raw != "" {
		value, err := strs.Atoi(raw)
		if err != nil {
			return arg, sandbox.Deps.Std.Errorf("--start %q is not a segment index", raw)
		}
		arg.Start = value
	}
	if arg.Start < 0 {
		return arg, sandbox.Deps.Std.Errorf("the command already reads to its last segment: give the arg a --start before it")
	}
	arg.End = arg.Start
	if raw := strs.TrimSpace(props.End); raw != "" {
		value, err := strs.Atoi(raw)
		if err != nil {
			return arg, sandbox.Deps.Std.Errorf("--end %q is not a segment index (use -1 for the last one)", raw)
		}
		arg.End = value
	}
	if arg.End != commandconf.LastSegment && arg.End < arg.Start {
		return arg, sandbox.Deps.Std.Errorf("--end %d is before --start %d", arg.End, arg.Start)
	}

	kind, ok := CommandArgType(sandbox, props.Type)
	if !ok {
		return arg, sandbox.Deps.Std.Errorf("unknown arg type %q (use %s)", props.Type, strs.Join(commandconf.ArgTypes, ", "))
	}
	if kind != commandconf.DefaultArgType && arg.End != arg.Start {
		return arg, sandbox.Deps.Std.Errorf("a %s arg reads one segment: it takes no range", kind)
	}
	arg.Type = kind

	if props.Default != "" {
		if arg.Required {
			return arg, sandbox.Deps.Std.Errorf("an arg cannot be both required and carry a default (the default already covers its absence)")
		}
		if err := CheckCommandLiteral(sandbox, kind, "default", props.Default); err != nil {
			return arg, err
		}
		arg.Default, arg.HasDefault = props.Default, true
	}

	trigger, err := NewTrigger(sandbox, TriggerProps{
		Type:       props.TriggerType,
		Value:      props.Trigger,
		Negate:     props.TriggerNegate,
		IgnoreCase: props.TriggerIgnoreCase,
	})
	if err != nil {
		return arg, err
	}
	arg.Trigger = trigger
	return arg, nil
}

// GeneratedCommands are the commands the build writes itself — the cli
// layer's own help, version and help-flag — which no verb declares, renames
// or removes.
var GeneratedCommands = []string{"help", "version", "help_flag"}

// IsGeneratedCommand reports whether a command name is one of
// GeneratedCommands.
func IsGeneratedCommand(sandbox *api.Sandbox, name string) bool {
	return contains(GeneratedCommands, CommandPackage(sandbox, name))
}

// FindCommandFlagNamed is the index of the flag a user named: by its Entries
// id ("out-file" -> OutFile), by one of its keys ("--out"), or by a key
// without its dashes ("command" for --command). -1 when none answers to it.
func FindCommandFlagNamed(sandbox *api.Sandbox, conf *commandconf.CommandConf, name string) int {
	name = sandbox.Deps.Stringsdeps.TrimSpace(name)
	for _, candidate := range []string{CommandEntryId(sandbox, name), name, "--" + name, "-" + name} {
		if index := FindCommandFlag(conf, candidate); index >= 0 {
			return index
		}
	}
	return -1
}
