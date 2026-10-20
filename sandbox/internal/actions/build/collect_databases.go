package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/databaseconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/migrationconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// CollectDatabases reads every sandbox/internal/databases/<db>/database.yaml and
// returns one data map per database, for the generated.api.go, generated.new.go and
// generated.methods.go of that package. It is the database layer's CollectRoutes: what
// the map holds is the declaration itself, turned into the records and the
// methods the three templates spell — nothing here is a Go rendering of the
// storage.
//
// A database is not a surface of sandbox/api/: its methods are typed by table,
// so there is no []Database standing where Cli.Commands stands. The maps
// therefore feed the per-unit generator alone, never a field of the Sandbox.
func CollectDatabases(sandbox *api.Sandbox, io *stagedfs.StagedFS) ([]map[string]any, error) {
	var databases []map[string]any

	for _, dir := range io.ListDirs(utils.DatabasesDir) {
		name := lastSegmentOf(sandbox, dir)
		if name == "" {
			continue
		}

		content, err := io.ReadFile(utils.DatabasesDir + "/" + name + "/" + utils.DatabaseConfFile)
		if err != nil {
			continue
		}

		conf, err := databaseconf.New(sandbox, string(content))
		if err != nil {
			return nil, sandbox.Deps.StdDeps.Errorf("databases/%s/%s: %w", name, utils.DatabaseConfFile, err)
		}

		if err := checkDatabaseSchema(sandbox, name, conf); err != nil {
			return nil, err
		}

		data := databaseData(sandbox, name, conf)
		versions, has_steps, err := databaseVersionsData(sandbox, io, name, conf)
		if err != nil {
			return nil, err
		}
		data["Version"] = conf.Version
		data["Versions"] = versions
		data["HasSteps"] = has_steps

		databases = append(databases, data)
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
					return sandbox.Deps.StdDeps.Errorf("databases/%s/%s: the link %s.%s targets %q, which is not a table of this database",
						name, utils.DatabaseConfFile, table.Name, field.Name, field.Target)
				}
			case databaseconf.FieldObject:
				for _, nested := range field.Fields {
					if nested.Type == databaseconf.FieldObject {
						return sandbox.Deps.StdDeps.Errorf("databases/%s/%s: %s.%s.%s nests a collection inside a collection; only one level is generated",
							name, utils.DatabaseConfFile, table.Name, field.Name, nested.Name)
					}
				}
			}
		}
	}
	return nil
}

// databaseData is one database as the three templates read it: the package it
// declares, the records it holds, the schema literal generated.new.go builds, and the
// methods generated.api.go declares, generated.new.go wires and generated.methods.go writes the bodies of.
func databaseData(sandbox *api.Sandbox, name string, conf *databaseconf.DatabaseConf) map[string]any {
	tables := make([]map[string]any, 0, len(conf.Tables))
	records := []map[string]any{}
	methods := []map[string]any{}

	for _, table := range conf.Tables {
		data := databaseTableData(sandbox, conf, table)
		tables = append(tables, data)

		records = append(records, recordData(sandbox, table.Name, utils.GoIdentifier(sandbox, table.Name), table.Fields, true))
		for _, field := range table.Fields {
			if field.Type != databaseconf.FieldObject {
				continue
			}
			records = append(records, recordData(sandbox, field.Name, utils.GoIdentifier(sandbox, table.Name)+utils.GoIdentifier(sandbox, field.Name), field.Fields, false))
		}

		for _, method := range utils.DatabaseMethods(sandbox, table) {
			methods = append(methods, databaseMethodData(method))
		}
	}

	return map[string]any{
		"DatabaseName": utils.DatabaseName(sandbox, name),
		"Package":      name,
		"Type":         utils.GoIdentifier(sandbox, name),
		"KeyPrefix":    conf.KeyPrefix,
		"Tables":       tables,
		"Records":      records,
		"Methods":      methods,
	}
}

// databaseTableData is one table as the schema literal of generated.new.go reads it.
func databaseTableData(sandbox *api.Sandbox, conf *databaseconf.DatabaseConf, table databaseconf.Table) map[string]any {
	items := make([]map[string]any, 0, len(table.Fields))
	for _, field := range table.Fields {
		items = append(items, databaseItemData(sandbox, field))
	}

	return map[string]any{
		"Name":  table.Name,
		"Type":  utils.GoIdentifier(sandbox, table.Name),
		"Items": items,
	}
}

// databaseItemData is one field as the databasedeps.Field literal reads it, its
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

// databaseItemConst is the databasedeps.Field type constant one declared type is
// written as.
func databaseItemConst(kind string) string {
	switch kind {
	case databaseconf.FieldKey:
		return "databasedeps.Key"
	case databaseconf.FieldInteger:
		return "databasedeps.Int"
	case databaseconf.FieldNumber:
		return "databasedeps.Float"
	case databaseconf.FieldBytes:
		return "databasedeps.Bytes"
	case databaseconf.FieldLink:
		return "databasedeps.Link"
	case databaseconf.FieldObject:
		return "databasedeps.Nested"
	}
	return "databasedeps.String"
}

// recordData is one Go record the package declares: the <X>Record every read
// hands back, the <X>Input every insert takes, and — for a table, which is the
// only thing List<T> ranges over — the <X>Filter that narrows it.
func recordData(sandbox *api.Sandbox, name string, kind string, fields []databaseconf.Field, filter bool) map[string]any {
	plain := []map[string]any{}
	for _, field := range fields {
		if field.Type == databaseconf.FieldObject {
			continue
		}
		plain = append(plain, databaseFieldData(sandbox, field))
	}

	return map[string]any{
		"Name":      name,
		"Type":      kind,
		"Fields":    plain,
		"HasFilter": filter,
	}
}

// databaseFieldData is one plain field as the record structs and the readers
// spell it: the Go name, the Go type, the reader that converts it back and
// which shape of filter it takes.
func databaseFieldData(sandbox *api.Sandbox, field databaseconf.Field) map[string]any {
	return map[string]any{
		"Name":    field.Name,
		"Go":      utils.GoIdentifier(sandbox, field.Name),
		"Kind":    field.Type,
		"GoType":  utils.DatabaseGoType(field.Type),
		"Reader":  databaseReader(field.Type),
		"IsText":  utils.DatabaseGoType(field.Type) == "string",
		"IsInt":   utils.DatabaseGoType(field.Type) == "int64",
		"IsFloat": utils.DatabaseGoType(field.Type) == "float64",
		"IsBytes": utils.DatabaseGoType(field.Type) == "[]byte",
	}
}

// databaseMethodData is one derived method as the three templates read it.
// The signature comes across whole, so generated.api.go, generated.new.go and generated.methods.go spell one
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

// databaseReader is the OpinionatedAgnosDatabase reader one declared type is read back
// through.
func databaseReader(kind string) string {
	switch utils.DatabaseGoType(kind) {
	case "int64":
		return "ReadInt"
	case "float64":
		return "ReadFloat"
	case "[]byte":
		return "ReadBytes"
	}
	return "ReadString"
}

// databaseVersionsData is every version of one database as its
// generated.migrations.go spells it, the first one first: the schema of each —
// frozen in migrations/v<N>.yaml for a past one, the declaration itself for
// the last — the operations of every one past the first, and whether a
// hand-written migrations/v<N>.go goes with it. The second result reports a
// database with at least one such step, which is what imports its migrations
// package.
func databaseVersionsData(sandbox *api.Sandbox, io *stagedfs.StagedFS, name string, conf *databaseconf.DatabaseConf) ([]map[string]any, bool, error) {
	versions := []map[string]any{}
	has_steps := false

	for version := 1; version <= conf.Version; version++ {
		tables := conf.Tables
		if version < conf.Version {
			frozen, err := loadFrozenSchema(sandbox, io, name, version)
			if err != nil {
				return nil, false, err
			}
			tables = frozen.Tables
		}

		data := map[string]any{
			"Version":  version,
			"Tables":   databaseTablesData(sandbox, conf, tables),
			"Ops":      []map[string]any{},
			"HasStep":  false,
			"Function": utils.DatabaseStepFunction(sandbox, version),
		}

		if version > 1 {
			ops, err := loadVersionOps(sandbox, io, name, version)
			if err != nil {
				return nil, false, err
			}
			data["Ops"] = ops
			if io.IsFile(utils.DatabaseStepPath(sandbox, name, version)) {
				data["HasStep"] = true
				has_steps = true
			}
		}

		versions = append(versions, data)
	}

	return versions, has_steps, nil
}

// loadFrozenSchema reads the declaration one past version froze.
func loadFrozenSchema(sandbox *api.Sandbox, io *stagedfs.StagedFS, name string, version int) (*databaseconf.DatabaseConf, error) {
	path := utils.DatabaseSchemaPath(sandbox, name, version)
	content, err := io.ReadFile(path)
	if err != nil {
		return nil, sandbox.Deps.StdDeps.Errorf("%s is missing: version %d of the database's data has no schema to migrate it from", path, version)
	}
	frozen, err := databaseconf.New(sandbox, string(content))
	if err != nil {
		return nil, sandbox.Deps.StdDeps.Errorf("%s: %w", path, err)
	}
	return frozen, nil
}

// loadVersionOps reads the operations of one version, as the template spells
// them.
func loadVersionOps(sandbox *api.Sandbox, io *stagedfs.StagedFS, name string, version int) ([]map[string]any, error) {
	path := utils.DatabaseOpsPath(sandbox, name, version)
	content, err := io.ReadFile(path)
	if err != nil {
		return nil, sandbox.Deps.StdDeps.Errorf("%s is missing: nothing says how the data reaches version %d", path, version)
	}
	conf, err := migrationconf.New(sandbox, string(content))
	if err != nil {
		return nil, sandbox.Deps.StdDeps.Errorf("%s: %w", path, err)
	}

	ops := make([]map[string]any, 0, len(conf.Ops))
	for _, op := range conf.Ops {
		ops = append(ops, map[string]any{
			"Kind":    op.Kind,
			"Table":   op.Table,
			"Parent":  op.Parent,
			"Field":   op.Field,
			"Type":    op.Type,
			"To":      op.To,
			"Default": op.Default,
		})
	}
	return ops, nil
}

// databaseTablesData is a list of tables as a schema literal reads them.
func databaseTablesData(sandbox *api.Sandbox, conf *databaseconf.DatabaseConf, tables []databaseconf.Table) []map[string]any {
	data := make([]map[string]any, 0, len(tables))
	for _, table := range tables {
		data = append(data, databaseTableData(sandbox, conf, table))
	}
	return data
}
