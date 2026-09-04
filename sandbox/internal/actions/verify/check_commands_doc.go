package verify

import (
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/commandconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// commandsDir holds one declared command per sub-directory.
const commandsDir = "sandbox/internal/commands"

// commandsDoc is the reference page for the command surface. It is written by
// hand — its columns carry judgement the declarations do not hold — so this is
// the file the command declarations can silently drift away from.
const commandsDoc = "docs/Commands/doc.md"

// CheckCommandsDoc enforces that the hand-written command reference still
// covers the declared surface: every visible command is named there, and every
// flag it declares is reachable in the page under at least one of its
// identifiers. It is the guard that makes `add-command` and `add-flag` safe to
// run without the page going stale.
//
// A project without that page, or without commands, has nothing to check.
func CheckCommandsDoc(deps *deps.Deps, io *smartio.SmartIO) []string {
	var violations []string

	if !io.IsDir(commandsDir) || !io.IsFile(commandsDoc) {
		return violations
	}

	content, err := io.ReadFile(commandsDoc)
	if err != nil {
		return []string{commandsDoc + " could not be read"}
	}
	page := string(content)

	for _, dir := range io.ListDirs(commandsDir) {
		name := lastSegment(dir)

		raw, err := io.ReadFile(dir + "/entries.yaml")
		if err != nil {
			continue
		}

		conf, err := commandconf.New(deps, string(raw))
		if err != nil {
			violations = append(violations, dir+"/entries.yaml is not parsable: "+err.Error())
			continue
		}

		if conf.Hidden {
			continue
		}

		violations = append(violations, checkCommandDocumented(page, name, conf)...)
	}

	return violations
}

// checkCommandDocumented reports one command's gaps in the reference page.
func checkCommandDocumented(page string, name string, conf *commandconf.CommandConf) []string {
	if !mentionsAny(page, conf.Identifiers) {
		return []string{commandsDoc + " documents no row for command " + name +
			" (every visible command belongs to one of its category tables)"}
	}

	var violations []string

	for _, flag := range conf.Flags {
		if mentionsAny(page, flag.Identifiers) {
			continue
		}
		violations = append(violations, commandsDoc+" never mentions flag "+flag.Key+
			" of command "+name+" (document it under one of its identifiers)")
	}

	return violations
}

// mentionsAny reports whether the page carries at least one of the
// identifiers. One spelling is enough: a page showing --target does not owe
// the reader -t as well.
func mentionsAny(page string, identifiers []string) bool {
	for _, identifier := range identifiers {
		if identifier != "" && strings.Contains(page, identifier) {
			return true
		}
	}
	return false
}
