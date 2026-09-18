package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/databaseconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// CollectDatabases reads every sandbox/internal/databases/<db>/specs.yaml and
// returns one data map per database, for the generated api.go, new.go and
// methods.go of that package. It is the database layer's CollectRoutes: what
// the map holds is the declaration itself, turned into the records and the
// methods the three templates spell — nothing here is a Go rendering of the
// storage.
//
// A database is not a surface of sandbox/api/: its methods are typed by table,
// so there is no []Database standing where Cli.Commands stands. The maps
// therefore feed the per-unit generator alone, never a field of the Sandbox.
func CollectDatabases(sandbox *api.Sandbox, io *smartio.SmartIO) ([]map[string]any, error) {
	var databases []map[string]any

	for _, dir := range io.ListDirs(utils.DatabasesDir) {
		name := lastSegmentOf(sandbox, dir)
		if name == "" {
			continue
		}

		content, err := io.ReadFile(utils.DatabasesDir + "/" + name + "/" + utils.DatabaseSpecsFile)
		if err != nil {
			continue
		}

		conf, err := databaseconf.New(sandbox, string(content))
		if err != nil {
			return nil, sandbox.Deps.Std.Errorf("databases/%s/%s: %w", name, utils.DatabaseSpecsFile, err)
		}

		if err := checkDatabaseSchema(sandbox, name, conf); err != nil {
			return nil, err
		}

		databases = append(databases, databaseData(sandbox, name, conf))
	}

	return databases, nil
}

// checkDatabaseSchema refuses a declaration the generator cannot turn into
// compiling Go: a link with no table to point at, and a nested collection
// holding another one — the second level has no parent record to be reached
// through, so no method could be generated for it.
func checkDatabaseSchema(sandbox *api.Sandbox, name string, conf *databaseconf.DatabaseConf) error {
	for _, table := range conf.Tables {
		for _, field := range table.Fields {
			switch field.Type {
			case databaseconf.FieldLink:
				if utils.FindDatabaseTable(sandbox, conf.Tables, field.Target) < 0 {
					return sandbox.Deps.Std.Errorf("databases/%s/%s: the link %s.%s targets %q, which is not a table of this database",
						name, utils.DatabaseSpecsFile, table.Name, field.Name, field.Target)
				}
			case databaseconf.FieldDatabase:
				for _, nested := range field.Fields {
					if nested.Type == databaseconf.FieldDatabase {
						return sandbox.Deps.Std.Errorf("databases/%s/%s: %s.%s.%s nests a collection inside a collection; only one level is generated",
							name, utils.DatabaseSpecsFile, table.Name, field.Name, nested.Name)
					}
				}
			}
		}
	}
	return nil
}

// databaseData is one database as the three templates read it: the package it
// declares, the records it holds, the schema literal new.go builds, and the
// methods api.go declares, new.go wires and methods.go writes the bodies of.
func databaseData(sandbox *api.Sandbox, name string, conf *databaseconf.DatabaseConf) map[string]any {
	tables := make([]map[string]any, 0, len(conf.Tables))
	records := []map[string]any{}
	methods := []map[string]any{}

	for _, table := range conf.Tables {
		data := databaseTableData(sandbox, conf, table)
		tables = append(tables, data)

		records = append(records, recordData(sandbox, table.Name, utils.ExportedName(sandbox, table.Name), table.Fields, true))
		for _, field := range table.Fields {
			if field.Type != databaseconf.FieldDatabase {
				continue
			}
			records = append(records, recordData(sandbox, field.Name, utils.ExportedName(sandbox, field.Name), field.Fields, false))
		}

		for _, method := range utils.DatabaseMethods(sandbox, table) {
			methods = append(methods, databaseMethodData(method))
		}
	}

	return map[string]any{
		"Name":    utils.DatabaseIdentifier(sandbox, name),
		"Package": name,
		"Type":    utils.ExportedName(sandbox, name),
		"Prefix":  conf.Prefix,
		"Tables":  tables,
		"Records": records,
		"Methods": methods,
	}
}

// databaseTableData is one table as the schema literal of new.go reads it.
func databaseTableData(sandbox *api.Sandbox, conf *databaseconf.DatabaseConf, table databaseconf.Table) map[string]any {
	items := make([]map[string]any, 0, len(table.Fields))
	for _, field := range table.Fields {
		items = append(items, databaseItemData(sandbox, field))
	}

	return map[string]any{
		"Name":  table.Name,
		"Type":  utils.ExportedName(sandbox, table.Name),
		"Items": items,
	}
}

// databaseItemData is one field as the database.Item literal reads it, its
// nested collection included.
func databaseItemData(sandbox *api.Sandbox, field databaseconf.Field) map[string]any {
	nested := make([]map[string]any, 0, len(field.Fields))
	for _, child := range field.Fields {
		nested = append(nested, databaseItemData(sandbox, child))
	}

	return map[string]any{
		"Name":     field.Name,
		"Const":    databaseItemConst(field.Type),
		"Required": field.Required,
		"Target":   field.Target,
		"Nested":   nested,
	}
}

// databaseItemConst is the database.Item type constant one declared type is
// written as.
func databaseItemConst(kind string) string {
	switch kind {
	case databaseconf.FieldKey:
		return "database.Key"
	case databaseconf.FieldInt:
		return "database.Int"
	case databaseconf.FieldFloat:
		return "database.Float"
	case databaseconf.FieldLink:
		return "database.Link"
	case databaseconf.FieldDatabase:
		return "database.Database"
	}
	return "database.String"
}

// recordData is one Go record the package declares: the <X>Item every read
// hands back, the <X>New every insert takes, and — for a table, which is the
// only thing List<T> ranges over — the <X>Filtrage that narrows it.
func recordData(sandbox *api.Sandbox, name string, kind string, fields []databaseconf.Field, filtrage bool) map[string]any {
	plain := []map[string]any{}
	for _, field := range fields {
		if field.Type == databaseconf.FieldDatabase {
			continue
		}
		plain = append(plain, databaseFieldData(sandbox, field))
	}

	return map[string]any{
		"Name":        name,
		"Type":        kind,
		"Fields":      plain,
		"HasFiltrage": filtrage,
	}
}

// databaseFieldData is one plain field as the record structs and the readers
// spell it: the Go name, the Go type, the reader that converts it back and
// which shape of filter it takes.
func databaseFieldData(sandbox *api.Sandbox, field databaseconf.Field) map[string]any {
	return map[string]any{
		"Name":    field.Name,
		"Go":      utils.ExportedName(sandbox, field.Name),
		"Kind":    field.Type,
		"GoType":  utils.DatabaseGoType(field.Type),
		"Reader":  databaseReader(field.Type),
		"IsText":  utils.DatabaseGoType(field.Type) == "string",
		"IsInt":   utils.DatabaseGoType(field.Type) == "int64",
		"IsFloat": utils.DatabaseGoType(field.Type) == "float64",
	}
}

// databaseMethodData is one derived method as the three templates read it.
// The signature comes across whole, so api.go, new.go and methods.go spell one
// method the same way.
func databaseMethodData(method utils.DatabaseMethod) map[string]any {
	return map[string]any{
		"Kind":    method.Kind,
		"Name":    method.Name,
		"Table":   method.Table,
		"Type":    method.Type,
		"Record":  method.Record,
		"Field":   method.Field,
		"Params":  method.Params,
		"Results": method.Results,
		"Args":    method.Args,
		"Help":    method.Help,
	}
}

// databaseReader is the databaseio reader one declared type is read back
// through.
func databaseReader(kind string) string {
	switch utils.DatabaseGoType(kind) {
	case "int64":
		return "ReadInt"
	case "float64":
		return "ReadFloat"
	}
	return "ReadString"
}
