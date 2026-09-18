package show_database

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/databaseconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// The indents the tree is drawn with. A database is three lists deep at most —
// its tables, their fields and the fields of a nested collection — so the
// depth is carried as a prefix rather than drawn with box characters, which
// keeps a line the same whether it is copied into a terminal or a doc.
const (
	branch = "  "
	level  = "  "
)

// ShowDatabaseInternal reads one specs.yaml and renders it as the lines of a
// tree: the package it declares and where its keys are written, then each
// table, the fields it holds, and the methods those fields generate.
//
// It is the one command of the database surface that writes nothing. A
// database is declared by five editors, each of which prints the file it
// changed and nothing else, so the declaration as a whole was only ever
// readable as yaml. This is that declaration said the way the editors talk
// about it — and the methods beside it, which the yaml never showed at all.
func ShowDatabaseInternal(sandbox *api.Sandbox, io *smartio.SmartIO, database string) ([]string, error) {
	conf, err := utils.LoadDatabaseConf(sandbox, io, database)
	if err != nil {
		return nil, err
	}

	lines := []string{
		sandbox.Deps.Std.Sprintf("%s  %s", utils.ExportedName(sandbox, conf.Name), utils.DatabaseDir(sandbox, database)),
		sandbox.Deps.Std.Sprintf("%skeys under  %s", branch, conf.Prefix),
	}

	if len(conf.Tables) == 0 {
		return append(lines, branch+"no table declared yet"), nil
	}

	for _, table := range conf.Tables {
		lines = append(lines, "", sandbox.Deps.Std.Sprintf("%s%s", branch, table.Name))
		lines = append(lines, fieldLines(sandbox, table.Fields, branch+level)...)

		lines = append(lines, sandbox.Deps.Std.Sprintf("%s%smethods", branch, level))
		for _, method := range utils.DatabaseMethods(sandbox, table) {
			lines = append(lines, sandbox.Deps.Std.Sprintf("%s%s%s%s(%s) %s",
				branch, level, level, method.Name, method.Params, method.Results))
		}
	}

	return lines, nil
}

// fieldLines is one list of fields, the nested collection of a `database`
// field written one level in under it.
func fieldLines(sandbox *api.Sandbox, fields []databaseconf.Field, prefix string) []string {
	lines := []string{}

	for _, field := range fields {
		label := field.Type
		if field.Target != "" {
			label += " -> " + field.Target
		}
		if field.Required {
			label += ", required"
		}
		lines = append(lines, sandbox.Deps.Std.Sprintf("%s%s  %s", prefix, field.Name, label))
		lines = append(lines, fieldLines(sandbox, field.Fields, prefix+level)...)
	}

	return lines
}
